package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/archive"
	"github.com/spf13/afero"
)

const (
	// DefaultWorkerPoolSize represents the default size of the worker pool
	DefaultWorkerPoolSize = 1
)

var (
	// ErrDispatcherStartingWorker represents an error when starting a worker
	ErrDispatcherStartingWorker = "error starting worker"
)

// Dispatcher represents a dispatcher to run tasks
type Dispatcher struct {
	// archiver factory to get the archiver
	archiverFactory *archive.ArchiveFactory
	// fs is the filesystem
	fs afero.Fs
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
	// workingDir is the working directory
	workingDir string
}

// NewDispatcher creates a new dispatcher
func NewDispatcher(workers int, fs afero.Fs, archiveFactory *archive.ArchiveFactory, workingDir string, logger repository.Logger) *Dispatcher {

	if workers == 0 {
		workers = DefaultWorkerPoolSize
	}

	// if logger == nil {
	// 	logger = zap.NewNop()
	// }

	return &Dispatcher{
		archiverFactory: archiveFactory,
		fs:              fs,
		logger:          logger,
		queue:           make(chan *entity.Task, workers),
		stopCh:          make(chan struct{}),
		workerPool:      make(chan chan *entity.Task, workers),
		workers:         make([]*Worker, 0, workers),
		workingDir:      workingDir,
	}
}

// Start starts the dispatcher
func (d *Dispatcher) Start(ctx context.Context) (err error) {

	d.onceStart.Do(func() {
		for i := 0; i < cap(d.queue); i++ {
			worker := NewWorker(d.workerPool, d.fs, d.archiverFactory, d.workingDir, d.logger)
			d.workers = append(d.workers, worker)
			workerStartErr := worker.Start(ctx)
			if workerStartErr != nil {
				msg := fmt.Sprintf("%s: %v", ErrDispatcherStartingWorker, workerStartErr)
				d.logger.Error(msg)
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
					d.logger.Info("Dispatcher stopped")
					return
				}
			}
		}()
	})

	return
}

func (d *Dispatcher) Stop() {
	d.logger.Info("Stopping dispatcher...")

	d.onceStop.Do(func() {
		close(d.stopCh)
	})
}

func (d *Dispatcher) Execute(task *entity.Task) error {
	d.queue <- task
	return nil
}
