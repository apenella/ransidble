package persistence

import (
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryTaskRepository(t *testing.T) {
	t.Parallel()
	t.Log("Testing NewMemoryTaskRepository")

	persistence := NewMemoryTaskRepository()
	expected := &MemoryTaskRepository{
		store: make(map[string]*entity.Task),
	}

	assert.Equal(t, persistence, expected)
}

// TestMemoryTaskRepository_Find tests the Find method
func TestMemoryTaskRepository_Find(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		persistence *MemoryTaskRepository
		expected    *entity.Task
		err         error
	}{
		{
			desc: "Testing find a task in memory persistence",
			id:   "task1",
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &entity.Task{ID: "task1"},
			err:      nil,
		},
		{
			desc: "Testing finding a task error when store is not initialized",
			id:   "task2",
			persistence: &MemoryTaskRepository{
				store: nil,
			},
			expected: nil,
			err:      entity.ErrNotInitializedStorage,
		},
		{
			desc: "Testing finding a task error when task does not exist",
			id:   "task3",
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: nil,
			err:      entity.ErrTaskNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			task, err := test.persistence.Find(test.id)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, task)
			}
		})
	}
}

// TestMemoryTaskRepository_FindAll tests the FindAll method
func TestMemoryTaskRepository_FindAll(t *testing.T) {
	tests := []struct {
		desc        string
		persistence *MemoryTaskRepository
		expected    []*entity.Task
		err         error
	}{
		{
			desc: "Testing find all tasks in memory persistence",
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
					"task2": {ID: "task2"},
				},
			},
			expected: []*entity.Task{{ID: "task1"}, {ID: "task2"}},
			err:      nil,
		},
		{
			desc: "Testing finding all tasks error when store is not initialized",
			persistence: &MemoryTaskRepository{
				store: nil,
			},
			expected: []*entity.Task{},
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			tasks, err := test.persistence.FindAll()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.ElementsMatch(t, test.expected, tasks)
			}
		})
	}
}

// TestMemoryTaskRepository_Remove tests the Remove method
func TestMemoryTaskRepository_Remove(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		persistence *MemoryTaskRepository
		expected    *MemoryTaskRepository
		err         error
	}{
		{
			desc: "Testing remove a task in memory persistence",
			id:   "task1",
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &MemoryTaskRepository{
				store: make(map[string]*entity.Task),
			},
			err: nil,
		},
		{
			desc: "Testing removing a task error when store is not initialized",
			id:   "task2",
			persistence: &MemoryTaskRepository{
				store: nil,
			},
			expected: &MemoryTaskRepository{},
			err:      entity.ErrNotInitializedStorage,
		},
		{
			desc: "Testing removing a task error when task does not exist",
			id:   "task3",
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &MemoryTaskRepository{},
			err:      entity.ErrTaskNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			err := test.persistence.Remove(test.id)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}

// TestMemoryTaskRepository_SafeStore tests the SafeStore method
func TestMemoryTaskRepository_SafeStore(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		task        *entity.Task
		persistence *MemoryTaskRepository
		expected    *MemoryTaskRepository
		err         error
	}{
		{
			desc: "Testing safe store a task in memory persistence",
			id:   "task1",
			task: &entity.Task{ID: "task1"},
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{},
			},
			expected: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			err: nil,
		},
		{
			desc: "Testing safe store a task error when task already exists",
			id:   "task1",
			task: &entity.Task{ID: "task1"},
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &MemoryTaskRepository{},
			err:      entity.ErrTaskAlreadyExists,
		},
		{
			desc: "Testing safe store a task error when store is not initialized",
			id:   "task2",
			task: &entity.Task{ID: "task2"},
			persistence: &MemoryTaskRepository{
				store: nil,
			},
			expected: &MemoryTaskRepository{},
			err:      entity.ErrNotInitializedStorage,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			err := test.persistence.SafeStore(test.id, test.task)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}

// TestMemoryTaskRepository_Store tests the Store method
func TestMemoryTaskRepository_Store(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		task        *entity.Task
		persistence *MemoryTaskRepository
		expected    *MemoryTaskRepository
		err         error
	}{
		{
			desc: "Testing store a task in memory persistence",
			id:   "task1",
			task: &entity.Task{ID: "task1"},
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{},
			},
			expected: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			err: nil,
		},
		{
			desc: "Testing store a task overwriting a task in memory persistence",
			id:   "task1",
			task: &entity.Task{ID: "task1_new"},
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1_new"},
				},
			},
			err: nil,
		},
		{
			desc: "Testing store a task error when store is not initialized",
			id:   "task2",
			task: &entity.Task{ID: "task2"},
			persistence: &MemoryTaskRepository{
				store: nil,
			},
			expected: &MemoryTaskRepository{},
			err:      entity.ErrNotInitializedStorage,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			err := test.persistence.Store(test.id, test.task)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}

// TestMemoryTaskRepository_Update tests the Update method
func TestMemoryTaskRepository_Update(t *testing.T) {
	tests := []struct {
		desc        string
		id          string
		task        *entity.Task
		persistence *MemoryTaskRepository
		expected    *MemoryTaskRepository
		err         error
	}{
		{
			desc: "Testing update a task in memory persistence",
			id:   "task1",
			task: &entity.Task{ID: "task1"},
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			err: nil,
		},
		{
			desc: "Testing update a task error when task does not exist",
			id:   "task2",
			task: &entity.Task{ID: "task2"},
			persistence: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			expected: &MemoryTaskRepository{
				store: map[string]*entity.Task{
					"task1": {ID: "task1"},
				},
			},
			err: entity.ErrTaskNotFound,
		},
		{
			desc: "Testing update a task error when store is not initialized",
			id:   "task3",
			task: &entity.Task{ID: "task3"},
			persistence: &MemoryTaskRepository{
				store: nil,
			},
			expected: &MemoryTaskRepository{},
			err:      entity.ErrNotInitializedStorage,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			err := test.persistence.Update(test.id, test.task)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}
