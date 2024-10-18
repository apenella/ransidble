package local

import (
	"fmt"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestNewLocalProjectRepository tests the NewLocalProjectRepository method
func TestNewLocalProjectRepository(t *testing.T) {
	fs := afero.NewMemMapFs()
	persistence := NewLocalProjectRepository(fs, "/tmp", nil)
	expected := &LocalProjectRepository{
		Fs:       fs,
		logger:   nil,
		Path:     "/tmp",
		Projects: make(map[string]*entity.Project),
	}

	assert.Equal(t, persistence, expected)
}

// TestLocalProjectRepository_Find tests the Find method
func TestLocalProjectRepository_Find(t *testing.T) {
	tests := []struct {
		desc        string
		name        string
		persistence *LocalProjectRepository
		expected    *entity.Project
		err         error
	}{
		{
			desc: "Testing find a project in local persistence",
			name: "project1",
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: &entity.Project{Name: "project1"},
			err:      nil,
		},
		{
			desc: "Testing finding a project error when project does not exist",
			name: "project2",
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: nil,
			err:      fmt.Errorf(ErrProjectNotFound),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			project, err := test.persistence.Find(test.name)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, project)
			}
		})
	}
}

// TestLocalProjectRepository_FindAll tests the FindAll method
func TestLocalProjectRepository_FindAll(t *testing.T) {
	tests := []struct {
		desc        string
		persistence *LocalProjectRepository
		expected    []*entity.Project
		err         error
	}{
		{
			desc: "Testing find all projects in local persistence",
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
					"project2": {Name: "project2"},
				},
			},
			expected: []*entity.Project{
				{Name: "project1"},
				{Name: "project2"},
			},
			err: nil,
		},
		{
			desc: "Testing find all projects in local persistence with empty projects",
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{},
			},
			expected: []*entity.Project{},
			err:      nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			projects, err := test.persistence.FindAll()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, projects)
			}
		})
	}
}

// TestLocalProjectRepository_Remove tests the Remove method
func TestLocalProjectRepository_Remove(t *testing.T) {
	tests := []struct {
		desc        string
		name        string
		persistence *LocalProjectRepository
		expected    *LocalProjectRepository
		err         error
	}{
		{
			desc: "Testing remove a project in local persistence",
			name: "project1",
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{},
			},
			err: nil,
		},
		{
			desc: "Testing remove a project in local persistence when project does not exist",
			name: "project2",
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			err: fmt.Errorf(ErrProjectNotFound),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.persistence.Remove(test.name)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}

// TestLocalProjectRepository_SafeStore tests the SafeStore method
func TestLocalProjectRepository_SafeStore(t *testing.T) {
	tests := []struct {
		desc        string
		name        string
		project     *entity.Project
		persistence *LocalProjectRepository
		expected    *LocalProjectRepository
		err         error
	}{
		{
			desc: "Testing safe store a project in local persistence",
			name: "project1",
			project: &entity.Project{
				Name: "project1",
			},
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			err: nil,
		},
		{
			desc: "Testing safe store a project in local persistence when project already exists",
			name: "project1",
			project: &entity.Project{
				Name: "project1",
			},
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			err: fmt.Errorf(ErrProjectAlreadyExists),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.persistence.SafeStore(test.name, test.project)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}

// TestLocalProjectRepository_Store tests the Store method
func TestLocalProjectRepository_Store(t *testing.T) {
	tests := []struct {
		desc        string
		name        string
		project     *entity.Project
		persistence *LocalProjectRepository
		expected    *LocalProjectRepository
		err         error
	}{
		{
			desc: "Testing store a project in local persistence",
			name: "project1",
			project: &entity.Project{
				Name: "project1",
			},
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.persistence.Store(test.name, test.project)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}

// TestLocalProjectRepository_Update tests the Update method
func TestLocalProjectRepository_Update(t *testing.T) {
	tests := []struct {
		desc        string
		name        string
		project     *entity.Project
		persistence *LocalProjectRepository
		expected    *LocalProjectRepository
		err         error
	}{
		{
			desc: "Testing update a project in local persistence",
			name: "project1",
			project: &entity.Project{
				Name: "project1_new",
			},
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1_new"},
				},
			},
			err: nil,
		},
		{
			desc: "Testing update a project in local persistence when project does not exist",
			name: "project2",
			project: &entity.Project{
				Name: "project2",
			},
			persistence: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			expected: &LocalProjectRepository{
				Projects: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
			},
			err: fmt.Errorf(ErrProjectNotFound),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			err := test.persistence.Update(test.name, test.project)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence)
			}
		})
	}
}
