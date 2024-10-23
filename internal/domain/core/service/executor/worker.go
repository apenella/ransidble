package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/apenella/ransidble/internal/infrastructure/executor"
	"github.com/google/uuid"
)

const (
	// WorkerTaskMessagePrefix represents the prefix message for worker errors
	WorkerTaskMessagePrefix = "Worker %s: Task '%s'. %s"
)

var (
	// ErrProjectNotFound represents an error when the project is not found
	ErrProjectNotFound = fmt.Errorf("project not found")
	// ErrTaskInvalidParameters represents an error when the task has invalid parameters
	ErrTaskInvalidParameters = fmt.Errorf("task has invalid parameters")
	// ErrUnachiverNotFound represents an error when the unarchiver is not found
	ErrUnachiverNotFound = fmt.Errorf("unarchiver not found")
	// ErrRemovingWorkingDirFolder represents an error removing working directory folder
	ErrRemovingWorkingDirFolder = fmt.Errorf("error removing working directory folder")
	// ErrPreparingWorkspace represents an error when preparing the workspace
	ErrPreparingWorkspace = fmt.Errorf("error preparing workspace")
	// ErrGettingWorkingDir represents an error when getting the working directory
	ErrGettingWorkingDir = fmt.Errorf("error getting working directory")
)

// Worker represents a worker to run tasks
type Worker struct {
	// id is the id of the worker
	id string
	// logger is the logger of the worker
	logger repository.Logger
	// onceStart is the sync.Once to start the worker
	onceStart sync.Once
	// onceStop is the sync.Once to stop the worker
	onceStop sync.Once
	// stopCh is the channel to stop the worker
	stopCh chan struct{}
	// taskChan is the channel to receive tasks
	taskChan chan *entity.Task
	// workerPool is the pool of workers to synchronize to the dispatcher
	workerPool chan chan *entity.Task
	// workspaceBuilder is the workspace builder
	workspaceBuilder service.WorkspaceBuilder
}

func genereteID() string {
	return uuid.New().String()
}

// NewWorker creates a new worker
func NewWorker(workerPool chan chan *entity.Task, workspaceBuilder service.WorkspaceBuilder, logger repository.Logger) *Worker {

	id := genereteID()

	return &Worker{
		id:               id, // set random alphanumeric
		logger:           logger,
		stopCh:           make(chan struct{}),
		taskChan:         make(chan *entity.Task),
		workerPool:       workerPool,
		workspaceBuilder: workspaceBuilder,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) (err error) {
	var workingDir string
	var workspace *workspace.Workspace

	w.onceStart.Do(func() {
		go func() {
			w.logger.Info(fmt.Sprintf("Starting worker %s", w.id), map[string]interface{}{
				"component": "Worker.Start",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
				"worker_id": w.id,
			})

			for {
				w.workerPool <- w.taskChan

				select {
				case t, ok := <-w.taskChan:
					if !ok {
						w.Stop()
					}
					t.Accepted()

					workspace = w.workspaceBuilder.WithTask(t).Build()

					defer func() {
						err = workspace.Cleanup()
					}()

					err = workspace.Prepare()
					if err != nil {
						errMssg := fmt.Sprintf("%s: %s", ErrPreparingWorkspace, err.Error())
						t.Failed(errMssg)
						w.logger.Error(errMssg, map[string]interface{}{
							"component": "Worker.Start",
							"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
							"task_id":   t.ID,
						})
						err = fmt.Errorf("%s", errMssg)
						continue
					}

					workingDir, err = workspace.GetWorkingDir()
					if err != nil {
						errMsg := fmt.Sprintf("%s: %s", ErrGettingWorkingDir, err.Error())
						t.Failed(errMsg)
						w.logger.Error(errMsg, map[string]interface{}{
							"component": "CreateTaskAnsiblePlaybookService.Run",
							"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
							"task_id":   t.ID,
						})
						err = fmt.Errorf("%s", errMsg)
						continue
					}

					switch t.Command {
					case entity.ANSIBLE_PLAYBOOK:
						w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Running a playbook"))
						_, ok = t.Parameters.(*entity.AnsiblePlaybookParameters)
						if !ok {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, ErrTaskInvalidParameters)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg, map[string]interface{}{
								"component": "Worker.Start",
								"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
								"task_id":   t.ID,
							})

							continue
						}

						t.Running()
						ansibleplaybook := executor.NewAnsiblePlaybook()
						errRunAnsiblePlaybook := ansibleplaybook.Run(ctx, workingDir, t.Parameters.(*entity.AnsiblePlaybookParameters))
						if errRunAnsiblePlaybook != nil {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, errRunAnsiblePlaybook)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg, map[string]interface{}{
								"component": "Worker.Start",
								"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
								"task_id":   t.ID,
							})
						} else {
							t.Success()
							w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Playbook successfully executed"), map[string]interface{}{
								"component": "Worker.Start",
								"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
								"task_id":   t.ID,
							})
						}

					default:
						errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Task with an unknown command")
						t.Failed(errorMsg)
						w.logger.Error(errorMsg, map[string]interface{}{
							"component": "Worker.Start",
							"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
							"task_id":   t.ID,
						})
					}

				case <-ctx.Done():
					w.Stop()
				case <-w.stopCh:
					w.logger.Info(fmt.Sprintf("Worker %s stopped", w.id), map[string]interface{}{
						"component": "Worker.Start",
						"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
						"worker_id": w.id,
					})
					close(w.taskChan)
					return
				}
			}
		}()
	})
	return nil
}

// Stop stops the worker
func (w *Worker) Stop() {
	w.logger.Info(fmt.Sprintf("Stopping worker %s", w.id), map[string]interface{}{
		"component": "Worker.Stop",
		"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
		"worker_id": w.id,
	})

	w.onceStop.Do(func() {
		close(w.stopCh)
	})
}
