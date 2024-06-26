package persistence

import (
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

type MemoryTaskRepository struct {
	store map[string]*entity.Task
	mutex sync.Mutex
}

// NewMemoryRepository creates a new memory repository
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		store: make(map[string]*entity.Task),
	}
}

// Find returns a task by id
func (m *MemoryTaskRepository) Find(id string) (*entity.Task, error) {

	if m.store == nil || m == nil {
		return nil, entity.ErrNotInitializedStorage
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	task, ok := m.store[id]
	if !ok {
		return nil, entity.ErrTaskNotFound
	}

	return task, nil
}

// FindAll returns all tasks
func (m *MemoryTaskRepository) FindAll() ([]*entity.Task, error) {
	tasks := []*entity.Task{}

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
		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if !ok {
		return entity.ErrTaskNotFound
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.store, id)

	return nil
}

// SafeStore stores a task and return an error if the task already exists
func (m *MemoryTaskRepository) SafeStore(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.store[id]
	if ok {
		return entity.ErrTaskAlreadyExists
	}

	m.store[id] = task

	return nil
}

// Store stores a task
func (m *MemoryTaskRepository) Store(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store[id] = task

	return nil
}

// Update updates a task
func (m *MemoryTaskRepository) Update(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if !ok {
		return entity.ErrTaskNotFound
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store[id] = task

	return nil
}
