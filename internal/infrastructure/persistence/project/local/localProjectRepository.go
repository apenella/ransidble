package local

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

const (
	// ErrProjectNotFound represents a project not found error
	ErrProjectNotFound = "project not found"
	// ErrProjectAlreadyExists represents a project already exists error
	ErrProjectAlreadyExists = "project already exists"
)

type LocalProjectRepository struct {
	// Filesystem path where projects are stored
	Fs afero.Fs
	// Path represents the local storage path
	Path string
	// Projects represents the projects
	Projects map[string]*entity.Project

	mutex  sync.Mutex
	logger repository.Logger
}

// NewLocalProjectRepository creates a new local project repository
func NewLocalProjectRepository(fs afero.Fs, path string, logger repository.Logger) *LocalProjectRepository {
	return &LocalProjectRepository{
		Fs:       fs,
		Path:     path,
		Projects: make(map[string]*entity.Project),
		logger:   logger,
	}
}

// LoadProjects loads projects from local storage
func (r *LocalProjectRepository) LoadProjects() error {

	var err error

	_, err = afero.IsDir(r.Fs, r.Path)
	if err != nil {
		return fmt.Errorf("error checking if path %s is a directory: %w", r.Path, err)
	}

	projects, err := afero.ReadDir(r.Fs, r.Path)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", r.Path, err)
	}

	for _, project := range projects {
		if project.IsDir() {
			projectPath := filepath.Join(r.Path, project.Name())
			projectEntity := entity.NewProject(project.Name(), projectPath, entity.ProjectTypeLocal)

			r.logger.Debug(fmt.Sprintf("Loading project %s from %s", project.Name(), projectPath))

			fmt.Println(">>>>>", project.Name(), projectPath, projectEntity)
			err = r.SafeStore(project.Name(), projectEntity)
			if err != nil {
				r.logger.Error(fmt.Sprintf("Error loading project %s: %s", project.Name(), err))
			}
		} else {
			r.logger.Warn(fmt.Sprintf("Path %s is not a directory, ignoring", project.Name()))
		}
	}

	return nil
}

// Find returns a project by id
func (r *LocalProjectRepository) Find(id string) (*entity.Project, error) {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	project, ok := r.Projects[id]
	if !ok {
		return nil, fmt.Errorf(ErrProjectNotFound)
	}

	return project, nil
}

// FindAll returns all projects
func (r *LocalProjectRepository) FindAll() ([]*entity.Project, error) {

	// TODO: return sorted projects

	projects := []*entity.Project{}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, project := range r.Projects {
		projects = append(projects, project)
	}

	return projects, nil
}

// Remove removes a project by id
func (r *LocalProjectRepository) Remove(id string) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.Projects[id]; !ok {
		return fmt.Errorf(ErrProjectNotFound)
	}

	delete(r.Projects, id)

	return nil
}

// SafeStore stores a project by id
func (r *LocalProjectRepository) SafeStore(id string, project *entity.Project) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.Projects[id]; ok {
		return fmt.Errorf(ErrProjectAlreadyExists)
	}

	r.Projects[id] = project

	return nil
}

// Store stores a project by id
func (r *LocalProjectRepository) Store(id string, project *entity.Project) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.Projects[id] = project

	return nil
}

// Update updates a project by id
func (r *LocalProjectRepository) Update(id string, project *entity.Project) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.Projects[id]; !ok {
		return fmt.Errorf(ErrProjectNotFound)
	}

	r.Projects[id] = project

	return nil
}
