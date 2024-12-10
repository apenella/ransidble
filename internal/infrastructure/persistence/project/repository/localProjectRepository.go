package repository

import (
	"fmt"
	"path/filepath"
	"strings"
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
	// ErrLocalProjectRepositoryPathNotExists represents a path not exists error
	ErrLocalProjectRepositoryPathNotExists = "path not exists"
	// ErrLocalProjectRepositoryPathMustBeDirectory represents a path must be a directory error
	ErrLocalProjectRepositoryPathMustBeDirectory = "path must be a directory"
	// ErrStoringProjectToLocalProjectRepository represents an error storing a project to local project repository
	ErrStoringProjectToLocalProjectRepository = "error storing project to local project repository"

	// ExtensionTarGz represents the tar.gz extension
	ExtensionTarGz = ".tar.gz"
)

// LocalProjectRepository represents a repository on local storage
type LocalProjectRepository struct {
	// Filesystem path where projects are stored
	fs afero.Fs
	// Path represents the local storage path
	path string
	// Projects represents the projects
	projects map[string]*entity.Project

	mutex  sync.Mutex
	logger repository.Logger
}

// NewLocalProjectRepository creates a new local project repository
func NewLocalProjectRepository(fs afero.Fs, path string, logger repository.Logger) *LocalProjectRepository {
	return &LocalProjectRepository{
		fs:       fs,
		path:     path,
		projects: make(map[string]*entity.Project),
		logger:   logger,
	}
}

// LoadProjects loads projects from local storage
func (r *LocalProjectRepository) LoadProjects() error {

	var err error
	var pathIsDir bool
	var projectEntity *entity.Project
	var projectFormat string
	var projectName string
	var projectPath string

	_, err = r.fs.Stat(r.path)
	if err != nil {
		r.logger.Error(
			ErrLocalProjectRepositoryPathNotExists,
			map[string]interface{}{
				"component": "LocalProjectRepository.LoadProjects",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"path":      r.path,
			})

		return fmt.Errorf(ErrLocalProjectRepositoryPathNotExists)
	}

	pathIsDir, err = afero.IsDir(r.fs, r.path)
	if err != nil {
		// This block handles an unexpected error returned by afero.IsDir
		errMsg := "An error occurred checking if path for a local project repository is a directory"
		r.logger.Error(
			errMsg,
			map[string]interface{}{
				"component": "LocalProjectRepository.LoadProjects",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"path":      r.path,
			})

		return fmt.Errorf(errMsg, err)
	}
	if !pathIsDir {
		r.logger.Error(
			ErrLocalProjectRepositoryPathMustBeDirectory,
			map[string]interface{}{
				"component": "LocalProjectRepository.LoadProjects",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"path":      r.path,
			})

		return fmt.Errorf(ErrLocalProjectRepositoryPathMustBeDirectory)
	}

	projects, err := afero.ReadDir(r.fs, r.path)
	if err != nil {
		// This block handles an unexpected error returned by afero.ReadDir
		errMsg := "An error occurred reading directory for a local project repository"
		r.logger.Error(
			errMsg,
			map[string]interface{}{
				"component": "LocalProjectRepository.LoadProjects",
				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"path":      r.path,
			})

		return fmt.Errorf(errMsg, err)
	}

	for _, project := range projects {

		projectName = project.Name()

		if project.Mode().IsRegular() {
			// When is found a regular file and project name ends with .tar.gz we consider it as a tar.gz file, otherwise we skip the file
			if strings.HasSuffix(project.Name(), ExtensionTarGz) {
				projectFormat = entity.ProjectFormatTarGz
				projectName = strings.TrimSuffix(project.Name(), ExtensionTarGz)
			} else {
				continue
			}
		}

		if project.IsDir() {
			projectFormat = entity.ProjectFormatPlain
		}

		projectPath = filepath.Join(r.path, project.Name())
		projectEntity = entity.NewProject(projectName, projectPath, projectFormat, entity.ProjectTypeLocal)

		r.logger.Debug(fmt.Sprintf("Loading project %s from %s", project.Name(), projectPath))

		err = r.SafeStore(projectName, projectEntity)
		if err != nil {
			r.logger.Error(
				ErrStoringProjectToLocalProjectRepository,
				map[string]interface{}{
					"component":  "LocalProjectRepository.LoadProjects",
					"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
					"project_id": project.Name(),
				})

			return fmt.Errorf("%s. %w", ErrStoringProjectToLocalProjectRepository, err)
		}
	}

	return nil
}

// Find returns a project by id
func (r *LocalProjectRepository) Find(id string) (*entity.Project, error) {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	project, ok := r.projects[id]
	if !ok {
		r.logger.Error(
			ErrProjectNotFound,
			map[string]interface{}{
				"component":  "LocalProjectRepository.Find",
				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"project_id": id,
			})

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

	for _, project := range r.projects {
		projects = append(projects, project)
	}

	return projects, nil
}

// Remove removes a project by id
func (r *LocalProjectRepository) Remove(id string) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.projects[id]; !ok {
		r.logger.Error(
			ErrProjectNotFound,
			map[string]interface{}{
				"component":  "LocalProjectRepository.Remove",
				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"project_id": id,
			})

		return fmt.Errorf(ErrProjectNotFound)
	}

	delete(r.projects, id)

	return nil
}

// SafeStore stores a project by id
func (r *LocalProjectRepository) SafeStore(id string, project *entity.Project) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.projects[id]; ok {
		r.logger.Error(
			ErrProjectAlreadyExists,
			map[string]interface{}{
				"component":  "LocalProjectRepository.SafeStore",
				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"project_id": id,
			})

		return fmt.Errorf(ErrProjectAlreadyExists)
	}

	r.projects[id] = project

	return nil
}

// Store stores a project by id
func (r *LocalProjectRepository) Store(id string, project *entity.Project) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.projects[id] = project

	return nil
}

// Update updates a project by id
func (r *LocalProjectRepository) Update(id string, project *entity.Project) error {

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.projects[id]; !ok {
		r.logger.Error(
			ErrProjectNotFound,
			map[string]interface{}{
				"component":  "LocalProjectRepository.Update",
				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
				"project_id": id,
			})

		return fmt.Errorf(ErrProjectNotFound)
	}

	r.projects[id] = project

	return nil
}
