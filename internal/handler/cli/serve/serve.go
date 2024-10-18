package serve

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/apenella/ransidble/internal/configuration"
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/service/executor"
	projectService "github.com/apenella/ransidble/internal/domain/core/service/project"
	taskService "github.com/apenella/ransidble/internal/domain/core/service/task"
	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
	server "github.com/apenella/ransidble/internal/handler/http"
	projectHandler "github.com/apenella/ransidble/internal/handler/http/project"
	taskHandler "github.com/apenella/ransidble/internal/handler/http/task"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch"
	localprojectpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository"
	taskpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
	"github.com/apenella/ransidble/internal/infrastructure/unpack"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	createTaskAnsiblePlaybookPath = "/tasks/ansible-playbook/:project_id"
	getHealthPath                 = "/health"
	getProjectPath                = "/projects/:id"
	getProjectsPath               = "/projects"
	getTaskPath                   = "/tasks/:id"
	getTasksPath                  = "/tasks"
)

var (
	// ErrStartDispatcher represents an error when starting the dispatcher
	ErrStartDispatcher = fmt.Errorf("error starting dispatcher")
	// ErrLoadProjects represents an error when loading projects
	ErrLoadProjects = fmt.Errorf("error loading projects")
)

// NewCommand returns a new cobra.Command to serve a Ransidble server
func NewCommand(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve is a command to start a Ransidble server",
		Long:  "Serve is a command to start a Ransidble server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			ctx := cmd.Context()
			log := logger.NewLogger()
			filesystem := afero.NewOsFs()

			// At the moment, the project repository loads the projects from the local storage. In the future, the plan is to have a database where you need to create a project before running it.
			projectsRepository := localprojectpersistence.NewLocalProjectRepository(
				filesystem,
				config.Server.Project.LocalStoragePath,
				log,
			)
			errLoadProjects := projectsRepository.LoadProjects()
			if errLoadProjects != nil {
				err = fmt.Errorf("%w: %s", ErrLoadProjects, errLoadProjects)
				return
			}

			fetchFactory := fetch.NewFactory()
			fetchFactory.Register(entity.ProjectTypeLocal, fetch.NewLocalStorage(
				filesystem,
				log,
			))

			unpackFactory := unpack.NewFactory()
			unpackFactory.Register(entity.ProjectFormatPlain, unpack.NewPlainFormat(
				filesystem,
				log,
			))

			workspaceBuilder := workspace.NewBuilder(
				filesystem,
				fetchFactory,
				unpackFactory,
				projectsRepository,
				log,
			)

			taskRepository := taskpersistence.NewMemoryTaskRepository()
			dispatcher := executor.NewDispatch(
				config.Server.WorkerPoolSize,
				workspaceBuilder,
				log,
			)

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

			createTaskAnsiblePlaybookService := taskService.NewCreateTaskAnsiblePlaybookService(
				dispatcher,
				taskRepository,
				projectsRepository,
				log,
			)
			createTaskAnsiblePlaybookHandler := taskHandler.NewCreateTaskAnsiblePlaybookHandler(createTaskAnsiblePlaybookService, log)
			router.POST(createTaskAnsiblePlaybookPath, createTaskAnsiblePlaybookHandler.Handle)

			getTaskService := taskService.NewGetTaskService(taskRepository, log)
			getTaskHandler := taskHandler.NewGetTaskHandler(getTaskService, log)
			router.GET(getTaskPath, getTaskHandler.Handle)

			getProjectService := projectService.NewGetProjectService(projectsRepository, log)
			getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
			router.GET(getProjectPath, getProjectHandler.Handle)
			getProjecListtHandler := projectHandler.NewGetProjecListtHandler(getProjectService, log)
			router.GET(getProjectsPath, getProjecListtHandler.Handle)

			// Wait for interrupt signal to gracefully shutdown the server
			quitCh := make(chan os.Signal, 1)
			signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
			errListenAndServeCh := make(chan error)

			srv := server.NewServer(config.Server.HTTPListenAddress, router, log)

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

			return
		},
	}

	return cmd
}
