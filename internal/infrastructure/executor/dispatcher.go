package executor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

const (
	// DefaultWorkerPoolSize represents the default size of the worker pool
	DefaultWorkerPoolSize = 1
)

var (
	// ErrDispatcherStartingWorker represents an error when starting a worker
	ErrDispatcherStartingWorker = errors.New("error starting worker")
)

// Dispatcher represents a dispatcher to run tasks
type Dispatcher struct {
	workerPool chan chan *entity.Task
	// queue is the queue of tasks to be executed
	queue chan *entity.Task
	// stopCh is the channel to stop the dispatcher
	stopCh chan struct{}
	// once is the sync.Once to stop the dispatcher
	once sync.Once
}

// NewDispatcher creates a new dispatcher
func NewDispatcher(workers int) *Dispatcher {

	if workers == 0 {
		workers = DefaultWorkerPoolSize
	}

	return &Dispatcher{
		queue:      make(chan *entity.Task, workers),
		stopCh:     make(chan struct{}),
		workerPool: make(chan chan *entity.Task, workers),
	}
}

// Start starts the dispatcher
func (d *Dispatcher) Start(ctx context.Context) (err error) {

	d.once.Do(func() {
		for i := 0; i < cap(d.queue); i++ {
			worker := NewWorker(d.workerPool)
			workerStartErr := worker.Start(ctx)
			if err != nil {
				err = fmt.Errorf("%w: %v", ErrDispatcherStartingWorker, workerStartErr)
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
					for i := 0; i < cap(d.workerPool); i++ {
						workerChannel := <-d.workerPool
						close(workerChannel)
					}
					return
				}
			}
		}()
	})

	return nil
}

func (d *Dispatcher) Stop() {
	d.once.Do(func() {
		close(d.stopCh)
	})
}

func (d *Dispatcher) Execute(task *entity.Task) error {
	d.queue <- task
	return nil
}
