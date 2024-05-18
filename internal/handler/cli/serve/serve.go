package serve

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/apenella/ransidble/internal/configuration"
	ansibleplaybookservice "github.com/apenella/ransidble/internal/domain/core/service/command"
	taskService "github.com/apenella/ransidble/internal/domain/core/service/task"
	server "github.com/apenella/ransidble/internal/handler/http"
	ansibleplaybookhandler "github.com/apenella/ransidble/internal/handler/http/command/ansible-playbook"
	taskHandler "github.com/apenella/ransidble/internal/handler/http/task"
	"github.com/apenella/ransidble/internal/infrastructure/executor"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	taskpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var (
	// ErrStartDispatcher represents an error when starting the dispatcher
	ErrStartDispatcher = fmt.Errorf("error starting dispatcher")
)

func NewCommand(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve is a command to start a Ransidble server",
		Long:  "Serve is a command to start a Ransidble server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			ctx := cmd.Context()
			log := logger.NewLogger()

			tasksStore := taskpersistence.NewMemoryPersistence()
			dispatcher := executor.NewDispatcher(config.WorkerPoolSize, log)
			go func() {
				errStartDispatcher := dispatcher.Start(ctx)
				if errStartDispatcher != nil {
					err = fmt.Errorf("%w: %s", ErrStartDispatcher, errStartDispatcher)
					return
				}
			}()

			router := echo.New()
			router.Use(middleware.Logger())
			router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
				Level: 5,
			}))

			ansiblePlaybookService := ansibleplaybookservice.NewAnsiblePlaybookService(dispatcher, tasksStore, log)
			ansiblePlaybookHandler := ansibleplaybookhandler.NewAnsiblePlaybookHandler(ansiblePlaybookService, tasksStore, log)
			router.POST("/command/ansible-playbook", ansiblePlaybookHandler.Handle)

			getTaskService := taskService.NewGetTaskService(tasksStore, log)
			getTaskHandler := taskHandler.NewGetTaskHandler(getTaskService, log)
			router.GET("/task/:id", getTaskHandler.Handle)

			// Wait for interrupt signal to gracefully shutdown the server
			quitCh := make(chan os.Signal, 1)
			signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
			errListenAndServeCh := make(chan error)

			srv := server.NewServer(config.HTTPListenAddress, router, log)

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				errListenAndServe := srv.Start(ctx)
				if errListenAndServe != nil {
					errListenAndServeCh <- errListenAndServe
				}
				wg.Done()
			}()

			select {
			case errListenAndServe := <-errListenAndServeCh:
				if errListenAndServe != nil {
					cmd.Println("Ransidble server stopped due to an error:", errListenAndServe)
					err = fmt.Errorf("%w: %s", server.ErrServerStarting, errListenAndServe)
				}
			case <-quitCh:
				cmd.Println("Received signal to stop the Ransidble server...")
				srv.Stop()
				dispatcher.Stop()
			}

			wg.Wait()

			return nil
		},
	}

	return cmd
}
