package repository

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestNewLocalProjectRepository tests the NewLocalProjectRepository method
func TestNewLocalProjectRepository(t *testing.T) {
	t.Parallel()
	t.Log("Testing NewLocalProjectRepository")

	fs := afero.NewMemMapFs()
	persistence := NewLocalProjectRepository(fs, "/tmp", nil)
	expected := &LocalProjectRepository{
		fs:     fs,
		logger: nil,
		path:   "/tmp",
		store:  make(map[string]*entity.Project),
	}

	assert.Equal(t, persistence, expected)
}

func TestLocalProjectRepository_LoadProjects(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "persistence-project-repository")
	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc       string
		repository *LocalProjectRepository
		// path     string
		// fs       afero.Fs
		expected map[string]*entity.Project
		err      error
	}{
		{
			desc: "Testing load projects from local storage",
			repository: NewLocalProjectRepository(
				fs,
				sourceBase,
				logger.NewFakeLogger(),
			),
			expected: map[string]*entity.Project{
				"project-1": {
					Name:      "project-1",
					Format:    "plain",
					Reference: filepath.Join(sourceBase, "project-1"),
					Storage:   "local",
				},
				"project-2": {
					Name:      "project-2",
					Format:    "targz",
					Reference: filepath.Join(sourceBase, "project-2.tar.gz"),
					Storage:   "local",
				},
			},
			err: nil,
		},
		{
			desc: "Testing error when loading projects from local storage and path does not exists",
			repository: NewLocalProjectRepository(
				fs,
				"not-exists",
				logger.NewFakeLogger(),
			),
			expected: map[string]*entity.Project{},
			err:      ErrLocalProjectRepositoryPathNotExists,
		},
		{
			desc: "Testing error when loading projects from local storage and path is not a directory",
			repository: NewLocalProjectRepository(
				fs,
				filepath.Join(sourceBase, "project-2.tar.gz"),
				logger.NewFakeLogger(),
			),
			expected: map[string]*entity.Project{},
			err:      ErrLocalProjectRepositoryPathMustBeDirectory,
		},
		{
			desc: "Testing error when loading projects from local storage and error occurs storing project",
			repository: &LocalProjectRepository{
				fs:   fs,
				path: sourceBase,
				store: map[string]*entity.Project{
					"project-1": {},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{},
			err:      fmt.Errorf("%s. %w", ErrStoringProjectToLocalProjectRepository, errors.New("project already exists")),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			err := test.repository.LoadProjects()
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Equal(t, test.expected, test.repository.store)
			}
		})
	}
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
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: &entity.Project{Name: "project1"},
			err:      nil,
		},
		{
			desc: "Testing finding a project error when project does not exist",
			name: "project2",
			persistence: &LocalProjectRepository{
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: nil,
			err:      ErrProjectNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
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
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
					"project2": {Name: "project2"},
				},
				logger: logger.NewFakeLogger(),
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
				store:  map[string]*entity.Project{},
				logger: logger.NewFakeLogger(),
			},
			expected: []*entity.Project{},
			err:      nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			projects, err := test.persistence.FindAll()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.ElementsMatch(t, test.expected, projects)
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
		expected    map[string]*entity.Project
		err         error
	}{
		{
			desc: "Testing remove a project in local persistence",
			name: "project1",
			persistence: &LocalProjectRepository{
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{},
			err:      nil,
		},
		{
			desc: "Testing remove a project in local persistence when project does not exist",
			name: "project2",
			persistence: &LocalProjectRepository{
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{
				"project1": {Name: "project1"},
			},
			err: ErrProjectNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			err := test.persistence.Remove(test.name)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence.store)
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
		expected    map[string]*entity.Project
		err         error
	}{
		{
			desc: "Testing safe store a project in local persistence",
			name: "project1",
			project: &entity.Project{
				Name: "project1",
			},
			persistence: &LocalProjectRepository{
				store:  map[string]*entity.Project{},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{
				"project1": {Name: "project1"},
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
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{
				"project1": {Name: "project1"},
			},
			err: ErrProjectAlreadyExists,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			err := test.persistence.SafeStore(test.name, test.project)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence.store)
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
		expected    map[string]*entity.Project
		err         error
	}{
		{
			desc: "Testing store a project in local persistence",
			name: "project1",
			project: &entity.Project{
				Name: "project1",
			},
			persistence: &LocalProjectRepository{
				store:  map[string]*entity.Project{},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{
				"project1": {Name: "project1"},
			},
			err: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			err := test.persistence.Store(test.name, test.project)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence.store)
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
		expected    map[string]*entity.Project
		err         error
	}{
		{
			desc: "Testing update a project in local persistence",
			name: "project1",
			project: &entity.Project{
				Name: "project1_new",
			},
			persistence: &LocalProjectRepository{
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{
				"project1": {Name: "project1_new"},
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
				store: map[string]*entity.Project{
					"project1": {Name: "project1"},
				},
				logger: logger.NewFakeLogger(),
			},
			expected: map[string]*entity.Project{
				"project1": {Name: "project1"},
			},
			err: ErrProjectNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			err := test.persistence.Update(test.name, test.project)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, test.persistence.store)
			}
		})
	}
}
