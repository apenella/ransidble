package serve

import (
	"fmt"

	"github.com/apenella/ransidble/internal/configuration"
	ansibleplaybookservice "github.com/apenella/ransidble/internal/domain/core/service/ansibleplaybook"
	server "github.com/apenella/ransidble/internal/handler/http"
	ansibleplaybookhandler "github.com/apenella/ransidble/internal/handler/http/command/ansible-playbook"
	"github.com/apenella/ransidble/internal/infrastructure/executor"
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
			dispatcher := executor.NewDispatcher(config.WorkerPoolSize)

			go func() {
				errStartDispatcher := dispatcher.Start(ctx)
				if err != nil {
					err = fmt.Errorf("%w: %s", ErrStartDispatcher, errStartDispatcher)
					return
				}
			}()

			router := echo.New()
			router.Use(middleware.Logger())
			router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
				Level: 5,
			}))

			//ansiblePlaybookExecutor := ansibleplaybookexecutor.NewAnsiblePlaybookRun()

			ansiblePlaybookService := ansibleplaybookservice.NewAnsiblePlaybookService(dispatcher)
			ansiblePlaybookHandler := ansibleplaybookhandler.NewAnsiblePlaybookHandler(ansiblePlaybookService)

			router.POST("/command/ansible-playbook", ansiblePlaybookHandler.Handle)

			s := server.NewServer(config.HTTPListenAddress, router)

			err = s.Start()
			if err != nil {
				cmd.Println("Error starting server", err)
				return err
			}
			return nil
		},
	}

	return cmd
}
