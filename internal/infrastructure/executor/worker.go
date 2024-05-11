package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/google/uuid"
)

// Worker represents a worker to run tasks
type Worker struct {
	id         string
	once       sync.Once
	stopCh     chan struct{}
	taskChan   chan *entity.Task
	workerPool chan chan *entity.Task
}

// NewWorker creates a new worker
func NewWorker(workerPool chan chan *entity.Task) *Worker {
	return &Worker{
		// set random alphanumeric
		id:         uuid.New().String(),
		stopCh:     make(chan struct{}),
		taskChan:   make(chan *entity.Task),
		workerPool: workerPool,
	}
}

// Start starts the worker
func (w *Worker) Start(ctx context.Context) error {
	go func() {

		fmt.Printf("Starting worker %s\n", w.id)

		for {
			w.workerPool <- w.taskChan

			select {
			case task, ok := <-w.taskChan:
				if !ok {
					w.Stop()
				}

				fmt.Printf("Worker %s: Running task %s\n", w.id, task.ID)

			case <-ctx.Done():
				w.Stop()
			case <-w.stopCh:
				fmt.Printf("Stoped worker %s\n", w.id)
				return
			}

		}
	}()
	return nil
}

func (w *Worker) Stop() {
	fmt.Printf("Stopping worker %s\n", w.id)
	w.once.Do(func() {
		close(w.stopCh)
	})
}
