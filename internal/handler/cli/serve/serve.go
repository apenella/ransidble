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
	ansibleexecutor "github.com/apenella/ransidble/internal/infrastructure/executor"
	"github.com/apenella/ransidble/internal/infrastructure/filesystem"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch"
	localprojectpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository"
	taskpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
	"github.com/apenella/ransidble/internal/infrastructure/tar"
	"github.com/apenella/ransidble/internal/infrastructure/unpack"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	CreateProjectPath             = "/projects"
	CreateTaskAnsiblePlaybookPath = "/tasks/ansible-playbook/:project_id"
	GetHealthPath                 = "/health"
	GetProjectPath                = "/projects/:id"
	GetProjectsPath               = "/projects"
	GetTaskPath                   = "/tasks/:id"
	GetTasksPath                  = "/tasks"
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

			log := logger.NewLogger()
			afs := afero.NewOsFs()
			fs := filesystem.NewFilesystem(afs)

			// At this moment, the project repository loads the projects from the local storage. In the future, the plan is to have a database where you need to create a project before running it.
			projectsRepository := localprojectpersistence.NewLocalProjectRepository(
				afs,
				config.Server.Project.LocalStoragePath,
				log,
			)

			errLoadProjects := projectsRepository.LoadProjects()
			if errLoadProjects != nil {
				errMsg := fmt.Sprintf("%s: %s", ErrLoadProjects, errLoadProjects)

				log.Error(
					errMsg,
					map[string]interface{}{
						"component": "Serve",
						"package":   "github.com/apenella/ransidble/internal/handler/cli/serve",
					})

				err = fmt.Errorf("%s", errMsg)
				return
			}

			fetchFactory := fetch.NewFactory()
			fetchFactory.Register(
				entity.ProjectTypeLocal,
				fetch.NewLocalStorage(
					afs,
					log,
				),
			)

			unpackFactory := unpack.NewFactory()
			unpackFactory.Register(entity.ProjectFormatPlain, unpack.NewPlainFormat(
				afs,
				log,
			))

			tarExtractor := tar.NewTar(afs, log)
			unpackFactory.Register(entity.ProjectFormatTarGz, unpack.NewTarGzipFormat(
				afs,
				tarExtractor,
				log,
			))

			workspaceBuilder := workspace.NewBuilder(
				fs,
				fetchFactory,
				unpackFactory,
				projectsRepository,
				log,
			)

			dispatcher := executor.NewDispatch(
				config.Server.WorkerPoolSize,
				workspaceBuilder,
				ansibleexecutor.NewAnsiblePlaybook(log),
				log,
			)

			taskRepository := taskpersistence.NewMemoryTaskRepository(log)
			createTaskAnsiblePlaybookService := taskService.NewCreateTaskAnsiblePlaybookService(
				dispatcher,
				taskRepository,
				projectsRepository,
				log,
			)

			createTaskAnsiblePlaybookHandler := taskHandler.NewCreateTaskAnsiblePlaybookHandler(createTaskAnsiblePlaybookService, log)

			getTaskService := taskService.NewGetTaskService(taskRepository, log)
			getTaskHandler := taskHandler.NewGetTaskHandler(getTaskService, log)

			getProjectService := projectService.NewGetProjectService(projectsRepository, log)
			getProjectHandler := projectHandler.NewGetProjectHandler(getProjectService, log)
			getProjectListHandler := projectHandler.NewGetProjectListHandler(getProjectService, log)

			createProjectHandler := projectHandler.NewCreateProjectHandler(log)

			go func() {
				errStartDispatcher := dispatcher.Start(cmd.Context())
				if errStartDispatcher != nil {
					errMsg := fmt.Sprintf("%s: %s", ErrStartDispatcher, errStartDispatcher)
					log.Error(
						errMsg,
						map[string]interface{}{
							"component": "Serve",
							"package":   "github.com/apenella/ransidble/internal/handler/cli/serve",
						})

					err = fmt.Errorf("%s", errMsg)
					return
				}
			}()

			router := echo.New()
			router.Use(middleware.Logger())
			router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
				Level: 5,
			}))

			router.POST(CreateProjectPath, createProjectHandler.Handle)
			router.POST(CreateTaskAnsiblePlaybookPath, createTaskAnsiblePlaybookHandler.Handle)
			router.GET(GetTaskPath, getTaskHandler.Handle)
			router.GET(GetProjectPath, getProjectHandler.Handle)
			router.GET(GetProjectsPath, getProjectListHandler.Handle)

			// Wait for interrupt signal to gracefully shutdown the server
			quitCh := make(chan os.Signal, 1)
			signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
			errListenAndServeCh := make(chan error)

			srv := server.NewServer(config.Server.HTTPListenAddress, router, log)

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				errListenAndServe := srv.Start(cmd.Context())
				if errListenAndServe != nil {
					errListenAndServeCh <- errListenAndServe
				}
				wg.Done()
			}()

			select {
			case errListenAndServe := <-errListenAndServeCh:
				if errListenAndServe != nil {
					errMsg := fmt.Sprintf("%s: %s", server.ErrServerStarting, errListenAndServe)
					log.Error(
						errMsg,
						map[string]interface{}{
							"component": "Serve",
							"package":   "github.com/apenella/ransidble/internal/handler/cli/serve",
						})
					err = fmt.Errorf("%s", errMsg)
				}
			case <-quitCh:
				log.Info(
					"Received signal to stop the Ransidble server",
					map[string]interface{}{
						"component": "Serve",
						"package":   "github.com/apenella/ransidble/internal/handler/cli/serve",
					})

				srv.Stop()
				dispatcher.Stop()
			}

			wg.Wait()

			cmd.Println("Ransidble server stopped")

			return
		},
	}

	return cmd
}
