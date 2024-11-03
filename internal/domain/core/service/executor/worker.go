package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
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
	// ErrAnsiblePlaybookExecutorDefined represents an error when the ansible playbook executor is not found
	ErrAnsiblePlaybookExecutorDefined = fmt.Errorf("ansible playbook executor not defined")
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

	// ansiblePlaybookExecutor is the ansible playbook executor
	ansiblePlaybookExecutor AnsiblePlaybookExecutor
}

func genereteID() string {
	return uuid.New().String()
}

// NewWorker creates a new worker
func NewWorker(
	workerPool chan chan *entity.Task,
	workspaceBuilder service.WorkspaceBuilder,
	ansiblePlaybookExecutor AnsiblePlaybookExecutor,
	logger repository.Logger,
) *Worker {

	id := genereteID()

	return &Worker{
		ansiblePlaybookExecutor: ansiblePlaybookExecutor,
		id:                      id, // set random alphanumeric
		logger:                  logger,
		stopCh:                  make(chan struct{}),
		taskChan:                make(chan *entity.Task),
		workerPool:              workerPool,
		workspaceBuilder:        workspaceBuilder,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) (err error) {

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
				case task, ok := <-w.taskChan:
					if !ok {
						w.Stop()
					}
					err = w.handleTask(ctx, task)
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

// handleTask handles a task to run by the worker
func (w *Worker) handleTask(ctx context.Context, task *entity.Task) error {
	var workingDir string
	var workspace service.Workspacer
	var err error

	task.Accepted()

	workspace, err = w.createWorkspace(task)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s", ErrPreparingWorkspace, err.Error())
		task.Failed(errMsg)
		w.logger.Error(errMsg, map[string]interface{}{
			"component": "Worker.handleTask",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
			"task_id":   task.ID,
		})
		return err
	}
	defer func() {
		err = workspace.Cleanup()
	}()

	workingDir, err = workspace.GetWorkingDir()
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s", ErrGettingWorkingDir, err.Error())
		task.Failed(errMsg)
		w.logger.Error(errMsg, map[string]interface{}{
			"component": "Worker.handleTask",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
			"task_id":   task.ID,
		})
		err = fmt.Errorf("%s", errMsg)
		return err
	}

	switch task.Command {
	case entity.AnsiblePlaybookCommand:
		// w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Running a playbook"))
		// _, ok := t.Parameters.(*entity.AnsiblePlaybookParameters)
		// if !ok {
		// 	errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, ErrTaskInvalidParameters)
		// 	t.Failed(errorMsg)
		// 	w.logger.Error(errorMsg, map[string]interface{}{
		// 		"component": "Worker.handleTask",
		// 		"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
		// 		"task_id":   t.ID,
		// 	})

		// 	return fmt.Errorf("%s", errorMsg)
		// }

		// t.Running()
		// ansibleplaybook := executor.NewAnsiblePlaybook()
		// errRunAnsiblePlaybook := ansibleplaybook.Run(ctx, workingDir, t.Parameters.(*entity.AnsiblePlaybookParameters))
		// if errRunAnsiblePlaybook != nil {
		// 	errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, errRunAnsiblePlaybook)
		// 	t.Failed(errorMsg)
		// 	w.logger.Error(errorMsg, map[string]interface{}{
		// 		"component": "Worker.handleTask",
		// 		"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
		// 		"task_id":   t.ID,
		// 	})

		// 	return fmt.Errorf("%s", errorMsg)
		// }

		_, ok := task.Parameters.(*entity.AnsiblePlaybookParameters)
		if !ok {
			errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, ErrTaskInvalidParameters)
			w.logger.Error(errorMsg, map[string]interface{}{
				"component": "Worker.handleTask",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
				"task_id":   task.ID,
			})

			return fmt.Errorf("%s", errorMsg)
		}

		task.Running()
		err = w.handleAnsiblePlaybookTask(ctx, task, workingDir)
		if err != nil {
			errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, err)
			task.Failed(errorMsg)
			w.logger.Error(errorMsg, map[string]interface{}{
				"component": "Worker.handleTask",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
				"task_id":   task.ID,
			})

			return fmt.Errorf("%s", errorMsg)
		}

		task.Success()
		w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, "Playbook successfully executed"), map[string]interface{}{
			"component": "Worker.handleTask",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			"task_id":   task.ID,
		})

	default:
		errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, "Task with an unknown command")
		task.Failed(errorMsg)
		w.logger.Error(errorMsg, map[string]interface{}{
			"component": "Worker.handleTask",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			"task_id":   task.ID,
		})

		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}

// createWorkspace creates a workspace for a task
func (w *Worker) createWorkspace(task *entity.Task) (service.Workspacer, error) {

	// the wsp and err are defined in the return statement to be able to handle the error in the defer function
	wsp := w.workspaceBuilder.WithTask(task).Build()

	err := wsp.Prepare()
	if err != nil {
		errMssg := fmt.Sprintf("%s: %s", ErrPreparingWorkspace, err.Error())
		w.logger.Error(errMssg, map[string]interface{}{
			"component": "Worker.createWorkspace",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			"task_id":   task.ID,
		})
		err = fmt.Errorf("%s", errMssg)
		return nil, err
	}

	return wsp, nil
}

// handleAnsiblePlaybookTask runs an ansible-playbook task
func (w *Worker) handleAnsiblePlaybookTask(ctx context.Context, task *entity.Task, workingDir string) error {

	if w.ansiblePlaybookExecutor == nil {
		errMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, ErrAnsiblePlaybookExecutorDefined)
		w.logger.Error(errMsg, map[string]interface{}{
			"component": "Worker.handleAnsiblePlaybookTask",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			"task_id":   task.ID,
		})

		return fmt.Errorf("%s", errMsg)
	}

	w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, "Running a playbook"), map[string]interface{}{
		"component": "Worker.handleAnsiblePlaybookTask",
		"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
		"task_id":   task.ID,
	})

	// ansibleplaybook := executor.NewAnsiblePlaybook()
	errRunAnsiblePlaybook := w.ansiblePlaybookExecutor.Run(ctx, workingDir, task.Parameters.(*entity.AnsiblePlaybookParameters))
	if errRunAnsiblePlaybook != nil {
		errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, task.ID, errRunAnsiblePlaybook)
		w.logger.Error(errorMsg, map[string]interface{}{
			"component": "Worker.handleAnsiblePlaybookTask",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/executor",
			"task_id":   task.ID,
		})

		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
