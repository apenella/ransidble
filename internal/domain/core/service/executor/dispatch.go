package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
)

const (
	// DefaultWorkerPoolSize represents the default size of the worker pool
	DefaultWorkerPoolSize = 1
)

var (
	// ErrDispatcherStartingWorker represents an error when starting a worker
	ErrDispatcherStartingWorker = "error starting worker"
)

// Dispatch represents a dispatcher to run tasks
type Dispatch struct {
	// logger is the logger of the dispatcher
	logger repository.Logger
	// onceStart is the sync.Once to stop the dispatcher
	onceStart sync.Once
	// onceStop is the sync.Once to stop the dispatcher
	onceStop sync.Once
	// queue is the queue of tasks to be executed
	queue chan *entity.Task
	// stopCh is the channel to stop the dispatcher
	stopCh chan struct{}
	// workerPool is the pool of workers
	workerPool chan chan *entity.Task
	// workers list of workers
	workers []*Worker
	// workspaceBuilder is the workspace builder
	workspaceBuilder service.WorkspaceBuilder
	// ansiblePlaybookExecutor is the ansible playbook executor
	ansiblePlaybookExecutor AnsiblePlaybookExecutor
}

// NewDispatch creates a new dispatcher to run tasks
func NewDispatch(
	workers int,
	workspaceBuilder service.WorkspaceBuilder,
	ansiblePlaybookExecutor AnsiblePlaybookExecutor,
	logger repository.Logger,
) *Dispatch {

	if workers == 0 {
		workers = DefaultWorkerPoolSize
	}

	return &Dispatch{
		ansiblePlaybookExecutor: ansiblePlaybookExecutor,
		logger:                  logger,
		queue:                   make(chan *entity.Task, workers),
		stopCh:                  make(chan struct{}),
		workerPool:              make(chan chan *entity.Task, workers),
		workers:                 make([]*Worker, 0, workers),
		workspaceBuilder:        workspaceBuilder,
	}
}

// Start starts the dispatcher
func (d *Dispatch) Start(ctx context.Context) (err error) {

	d.onceStart.Do(func() {

		for i := 0; i < cap(d.queue); i++ {
			worker := NewWorker(
				d.workerPool,
				d.workspaceBuilder,
				d.ansiblePlaybookExecutor,
				d.logger)
			d.workers = append(d.workers, worker)
			workerStartErr := worker.Start(ctx)

			if workerStartErr != nil {
				msg := fmt.Sprintf("%s: %v", ErrDispatcherStartingWorker, workerStartErr)
				d.logger.Error(msg, map[string]interface{}{
					"component": "Dispatch.Start",
					"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
				})
				err = fmt.Errorf(msg)
				return
			}
		}

		// main loop of the dispatcher must receive tasks from the queue. Then achieve the worker channel from the worker pool and send the task to the worker channel
		go func() {
			for {
				select {
				case task := <-d.queue:
					workerChannel := <-d.workerPool
					workerChannel <- task
				case <-ctx.Done():
					d.Stop()
				case <-d.stopCh:
					var wg sync.WaitGroup
					wg.Add(len(d.workers))
					for _, worker := range d.workers {
						worker.Stop()
						wg.Done()
					}
					wg.Wait()
					d.logger.Info("Dispatcher stopped", map[string]interface{}{
						"component": "Dispatch.Start",
						"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
					})
					close(d.queue)
					return
				}
			}
		}()
	})

	return
}

// Stop stops the dispatcher
func (d *Dispatch) Stop() {
	d.logger.Info("Stopping dispatcher", map[string]interface{}{
		"component": "Dispatch.Stop",
		"package":   "github.com/apenella/ransidble/internal/domain/core/service/task",
	})

	d.onceStop.Do(func() {
		close(d.stopCh)
	})
}

// Execute executes a task
func (d *Dispatch) Execute(task *entity.Task) error {
	d.queue <- task
	return nil
}
