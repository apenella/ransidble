package persistence

import "github.com/apenella/ransidble/internal/domain/core/entity"

type MemoryPersistence struct {
	store map[string]*entity.Task
}

// NewMemoryRepository creates a new memory repository
func NewMemoryPersistence() *MemoryPersistence {
	return &MemoryPersistence{
		store: make(map[string]*entity.Task),
	}
}

// Find returns a task by id
func (m *MemoryPersistence) Find(id string) (*entity.Task, error) {

	if m.store == nil || m == nil {
		return nil, entity.ErrNotInitializedStorage
	}

	task, ok := m.store[id]
	if !ok {
		return nil, entity.ErrTaskNotFound
	}

	return task, nil
}

// FindAll returns all tasks
func (m *MemoryPersistence) FindAll() ([]*entity.Task, error) {
	tasks := []*entity.Task{}
	for _, task := range m.store {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Remove removes a task by id
func (m *MemoryPersistence) Remove(id string) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if !ok {
		return entity.ErrTaskNotFound
	}

	delete(m.store, id)
	return nil
}

// SafeStore stores a task and return an error if the task already exists
func (m *MemoryPersistence) SafeStore(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if ok {
		return entity.ErrTaskAlreadyExists
	}

	m.store[id] = task
	return nil
}

// Store stores a task
func (m *MemoryPersistence) Store(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	m.store[id] = task
	return nil
}

// Update updates a task
func (m *MemoryPersistence) Update(id string, task *entity.Task) error {

	if m.store == nil || m == nil {
		return entity.ErrNotInitializedStorage
	}

	_, ok := m.store[id]
	if !ok {
		return entity.ErrTaskNotFound
	}

	m.store[id] = task
	return nil
}
