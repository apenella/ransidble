package serve

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/apenella/ransidble/internal/configuration"
	projectService "github.com/apenella/ransidble/internal/domain/core/service/project"
	taskService "github.com/apenella/ransidble/internal/domain/core/service/task"
	server "github.com/apenella/ransidble/internal/handler/http"
	projectHandler "github.com/apenella/ransidble/internal/handler/http/project"
	taskHandler "github.com/apenella/ransidble/internal/handler/http/task"
	"github.com/apenella/ransidble/internal/infrastructure/archive"
	"github.com/apenella/ransidble/internal/infrastructure/executor"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	localprojectpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/project/local"
	taskpersistence "github.com/apenella/ransidble/internal/infrastructure/persistence/task"
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

			workingDir, errCreateWorkingdir := afero.TempDir(filesystem, "", "ransidble")
			if err != nil {
				err = fmt.Errorf("%w: %s", archive.ErrCreatingWorkingDirFolder, errCreateWorkingdir)
				return
			}

			projectsRepository := localprojectpersistence.NewLocalProjectRepository(filesystem, config.Server.Project.LocalStoragePath, log)
			errLoadProjects := projectsRepository.LoadProjects()
			if errLoadProjects != nil {
				err = fmt.Errorf("%w: %s", ErrLoadProjects, errLoadProjects)
				return
			}

			archiveFactory := archive.NewArchiveFactory()
			archiveFactory.Register("local", archive.NewLocalStorageArchive(filesystem, log))

			taskRepository := taskpersistence.NewMemoryTaskRepository()
			dispatcher := executor.NewDispatcher(config.Server.WorkerPoolSize, filesystem, archiveFactory, workingDir, log)
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

			createTaskAnsiblePlaybookService := taskService.NewCreateTaskAnsiblePlaybookService(dispatcher, taskRepository, projectsRepository, log)
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
