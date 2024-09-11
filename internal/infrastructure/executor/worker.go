package executor

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/archive"
	executor "github.com/apenella/ransidble/internal/infrastructure/executor/ansible-playbook"
	"github.com/google/uuid"
	"github.com/spf13/afero"
)

const (
	// WorkerTaskMessagePrefix represents the prefix message for worker errors
	WorkerTaskMessagePrefix = "Worker %s: Task '%s'. %s"
)

var (
	// ErrCreateWorkingDirFolder represents an error creating working directory folder
	ErrCreateWorkingDirFolder = fmt.Errorf("error creating working directory folder")
	// ErrProjectNotFound represents an error when the project is not found
	ErrProjectNotFound = fmt.Errorf("project not found")
	// ErrTaskInvalidParameters represents an error when the task has invalid parameters
	ErrTaskInvalidParameters = fmt.Errorf("task has invalid parameters")
	// ErrUnachiverNotFound represents an error when the unarchiver is not found
	ErrUnachiverNotFound = fmt.Errorf("unarchiver not found")
	// ErrRemovingWorkingDirFolder represents an error removing working directory folder
	ErrRemovingWorkingDirFolder = fmt.Errorf("error removing working directory folder")
)

// Worker represents a worker to run tasks
type Worker struct {
	// archiver factory to get the archiver
	archiverFactory *archive.ArchiveFactory
	// fs is the filesystem
	fs afero.Fs
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
	// workingDir is the working directory
	workingDir string
}

func genereteID() string {
	return uuid.New().String()
}

// NewWorker creates a new worker
func NewWorker(workerPool chan chan *entity.Task, fs afero.Fs, archiveFactory *archive.ArchiveFactory, workingDir string, logger repository.Logger) *Worker {

	id := genereteID()
	workingDir = filepath.Join(workingDir, id)

	return &Worker{
		archiverFactory: archiveFactory,
		fs:              fs,
		id:              id, // set random alphanumeric
		logger:          logger,
		stopCh:          make(chan struct{}),
		taskChan:        make(chan *entity.Task),
		workerPool:      workerPool,
		workingDir:      workingDir,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) error {
	var err error

	err = w.fs.MkdirAll(w.workingDir, 0755)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateWorkingDirFolder, err)
	}

	w.onceStart.Do(func() {
		go func() {
			w.logger.Info(fmt.Sprintf("Starting worker %s", w.id))

			for {
				w.workerPool <- w.taskChan

				select {
				case t, ok := <-w.taskChan:
					if !ok {
						w.Stop()
					}
					t.Accepted()

					switch t.Command {
					case entity.ANSIBLE_PLAYBOOK:
						w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Running a playbook"))
						_, ok = t.Parameters.(*entity.AnsiblePlaybookParameters)
						if !ok {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, ErrTaskInvalidParameters)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
							continue
						}

						if t.Project == nil {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, ErrProjectNotFound)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
							continue
						}

						archiver := w.archiverFactory.Get(t.Project.Type)
						if archiver == nil {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, ErrUnachiverNotFound)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
							continue
						}

						projectTaskWorkingDir := filepath.Join(w.workingDir, t.Project.Name, t.ID)
						err = archiver.Unarchive(t.Project, projectTaskWorkingDir)
						if err != nil {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, fmt.Errorf("error unarchiving project: %w", err))
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
							continue
						}

						t.Running()
						ansibleplaybook := executor.NewAnsiblePlaybook()
						errRunAnsiblePlaybook := ansibleplaybook.Run(ctx, projectTaskWorkingDir, t.Parameters.(*entity.AnsiblePlaybookParameters))
						if errRunAnsiblePlaybook != nil {
							errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, errRunAnsiblePlaybook)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
						} else {
							t.Success()
							w.logger.Debug(fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Playbook successfully executed"))
						}

						err := w.fs.RemoveAll(projectTaskWorkingDir)
						if err != nil {
							errorMsg := fmt.Sprintf("%s: %s", ErrRemovingWorkingDirFolder, err)
							w.logger.Error(errorMsg)
						}
					default:
						errorMsg := fmt.Sprintf(WorkerTaskMessagePrefix, w.id, t.ID, "Task with an unknown command")
						t.Failed(errorMsg)
						w.logger.Error(errorMsg)
					}

				case <-ctx.Done():
					w.Stop()
				case <-w.stopCh:
					w.logger.Info(fmt.Sprintf("Worker %s stopped", w.id))
					close(w.taskChan)
					return
				}
			}
		}()
	})
	return nil
}

func (w *Worker) Stop() {
	w.logger.Info(fmt.Sprintf("Stopping worker %s", w.id))

	w.onceStop.Do(func() {

		err := w.fs.RemoveAll(w.workingDir)
		if err != nil {
			errorMsg := fmt.Sprintf("%s: %s", ErrRemovingWorkingDirFolder, err)
			w.logger.Error(errorMsg)
		}

		close(w.stopCh)
	})
}
