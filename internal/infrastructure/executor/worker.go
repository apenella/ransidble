package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	request "github.com/apenella/ransidble/internal/domain/core/model/request/ansible-playbook"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	executor "github.com/apenella/ransidble/internal/infrastructure/executor/ansible-playbook"
	"github.com/google/uuid"
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
}

// NewWorker creates a new worker
func NewWorker(workerPool chan chan *entity.Task, logger repository.Logger) *Worker {
	return &Worker{
		// set random alphanumeric
		id:         uuid.New().String(),
		stopCh:     make(chan struct{}),
		taskChan:   make(chan *entity.Task),
		workerPool: workerPool,
		logger:     logger,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) error {

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
						w.logger.Debug(fmt.Sprintf("Worker %s: Running a playbook %s", w.id, t.ID))
						_, ok = t.Parameters.(*request.AnsiblePlaybookParameters)
						if !ok {
							errorMsg := fmt.Sprintf("Worker %s: Task '%s' created with an invalid parameters", w.id, t.ID)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
							continue
						}
						t.Running()

						ansibleplaybook := executor.NewAnsiblePlaybook()
						errRunAnsiblePlaybook := ansibleplaybook.Run(ctx, t.Parameters.(*request.AnsiblePlaybookParameters))
						if errRunAnsiblePlaybook != nil {
							errorMsg := fmt.Sprintf("Worker %s: Task '%s' failed: %s", w.id, t.ID, errRunAnsiblePlaybook)
							t.Failed(errorMsg)
							w.logger.Error(errorMsg)
						} else {
							t.Success()
							w.logger.Debug(fmt.Sprintf("Worker %s: Task '%s' successfully executed", w.id, t.ID))
						}

					default:
						errorMsg := fmt.Sprintf("Worker %s: Task '%s' created with an unknown command", w.id, t.ID)
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
		close(w.stopCh)
	})
}
