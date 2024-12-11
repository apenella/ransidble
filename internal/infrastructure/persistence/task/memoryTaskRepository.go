package persistence

import (
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
)

// MemoryTaskRepository struct to store tasks in memory
type MemoryTaskRepository struct {
	store  map[string]*entity.Task
	mutex  sync.Mutex
	logger repository.Logger
}

// NewMemoryTaskRepository creates a new MemoryTaskRepository
func NewMemoryTaskRepository(logger repository.Logger) *MemoryTaskRepository {
	return &MemoryTaskRepository{
		store:  make(map[string]*entity.Task),
		logger: logger,
	}
}

// Find returns a task by id
func (m *MemoryTaskRepository) Find(id string) (*entity.Task, error) {

	if m.store == nil || m == nil {
		m.logger.Error(
			entity.ErrNotInitializedStorage.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Find",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return nil, entity.ErrNotInitializedStorage
	}

	m.logger.Debug(
		"Finding task",
		map[string]interface{}{
			"component": "MemoryTaskRepository.Find",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
			"task_id":   id,
		},
	)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	task, ok := m.store[id]
	if !ok {
		m.logger.Error(
			entity.ErrTaskNotFound.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Find",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return nil, entity.ErrTaskNotFound
	}

	return task, nil
}

// FindAll returns all tasks
func (m *MemoryTaskRepository) FindAll() ([]*entity.Task, error) {
	tasks := []*entity.Task{}

	if m.store == nil || m == nil {
		m.logger.Error(
			entity.ErrNotInitializedStorage.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.FindAll",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
			},
		)

		return nil, entity.ErrNotInitializedStorage
	}

	m.logger.Debug(
		"Finding all tasks",
		map[string]interface{}{
			"component": "MemoryTaskRepository.FindAll",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
		},
	)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, task := range m.store {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Remove removes a task by id
func (m *MemoryTaskRepository) Remove(id string) error {

	if m.store == nil || m == nil {
		m.logger.Error(
			entity.ErrNotInitializedStorage.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Remove",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if !ok {
		m.logger.Error(
			entity.ErrTaskNotFound.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Remove",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrTaskNotFound
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.store, id)

	m.logger.Debug(
		"Task removed",
		map[string]interface{}{
			"component": "MemoryTaskRepository.Remove",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
			"task_id":   id,
		},
	)

	return nil
}

// SafeStore stores a task and return an error if the task already exists
func (m *MemoryTaskRepository) SafeStore(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		m.logger.Error(
			entity.ErrNotInitializedStorage.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.SafeStore",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrNotInitializedStorage
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.store[id]
	if ok {
		m.logger.Error(
			entity.ErrTaskAlreadyExists.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.SafeStore",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrTaskAlreadyExists
	}

	m.store[id] = task

	m.logger.Debug(
		"Task stored",
		map[string]interface{}{
			"component": "MemoryTaskRepository.SafeStore",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
			"task_id":   id,
		},
	)

	return nil
}

// Store stores a task
func (m *MemoryTaskRepository) Store(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		m.logger.Error(
			entity.ErrNotInitializedStorage.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Store",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrNotInitializedStorage
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store[id] = task

	m.logger.Debug(
		"Task stored",
		map[string]interface{}{
			"component": "MemoryTaskRepository.Store",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
			"task_id":   id,
		},
	)

	return nil
}

// Update updates a task
func (m *MemoryTaskRepository) Update(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		m.logger.Error(
			entity.ErrNotInitializedStorage.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Update",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if !ok {
		m.logger.Error(
			entity.ErrTaskNotFound.Error(),
			map[string]interface{}{
				"component": "MemoryTaskRepository.Update",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
				"task_id":   id,
			},
		)

		return entity.ErrTaskNotFound
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store[id] = task

	m.logger.Debug(
		"Task updated",
		map[string]interface{}{
			"component": "MemoryTaskRepository.Update",
			"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/task",
			"task_id":   id,
		},
	)

	return nil
}
