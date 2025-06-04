package repository

// import (
// 	"fmt"
// 	"path/filepath"
// 	"strings"
// 	"sync"

// 	"github.com/apenella/ransidble/internal/domain/core/entity"
// 	"github.com/apenella/ransidble/internal/domain/ports/repository"
// 	"github.com/spf13/afero"
// )

// const (
// 	// ErrInitializingLocalProjectRepository represents an error initializing a local project repository
// 	ErrInitializingLocalProjectRepository = "error initializing local project repository"
// 	// ErrLocalProjectRepositoryCheckPathIsDirectory represents an error checking if path for a local project repository is a directory
// 	ErrLocalProjectRepositoryCheckPathIsDirectory = "an error occurred checking if path for a local project repository is a directory"
// 	// ErrLocalProjectRepositoryPathMustBeDirectory represents a path must be a directory error
// 	ErrLocalProjectRepositoryPathMustBeDirectory = "path must be a directory"
// 	// ErrLocalProjectRepositoryPathNotExists represents a path not exists error
// 	ErrLocalProjectRepositoryPathNotExists = "path not exists"
// 	// ErrProjectAlreadyExists represents a project already exists error
// 	ErrProjectAlreadyExists = "project already exists"
// 	// ErrProjectNotFound represents a project not found error
// 	ErrProjectNotFound = "project not found"
// 	// ErrProjectNotInitializedStorage is returned when the storage is not initialized
// 	ErrProjectNotInitializedStorage = "project storage not initialized"
// 	// ErrReadingLocalProjectRepositoryPath represents an error reading directory for a local project repository
// 	ErrReadingLocalProjectRepositoryPath = "an error occurred reading directory for a local project repository"
// 	// ErrStoringProjectToLocalProjectRepository represents an error storing a project to local project repository
// 	ErrStoringProjectToLocalProjectRepository = "error storing project to local project repository"
// 	// ErrLocalProjectRepositoryPathNotProvided represents an error when path is not provided
// 	ErrLocalProjectRepositoryPathNotProvided = "path not provided"
// )

// // LocalProjectRepository represents a repository on local storage
// type LocalProjectRepository struct {
// 	// fs path where projects are stored
// 	fs afero.Fs
// 	// path represents the local storage path
// 	path string
// 	// store represents the projects repository
// 	store map[string]*entity.Project

// 	// backendDriverer
// 	driver backendDriverer

// 	mutex  sync.Mutex
// 	logger repository.Logger
// }

// // NewLocalProjectRepository creates a new local project repository
// func NewLocalProjectRepository(fs afero.Fs, path string, driver backendDriverer, logger repository.Logger) *LocalProjectRepository {
// 	return &LocalProjectRepository{
// 		fs:     fs,
// 		path:   path,
// 		driver: driver,
// 		store:  make(map[string]*entity.Project),
// 		logger: logger,
// 	}
// }

// // Initialize initializes the local project repository
// func (r *LocalProjectRepository) Initialize() error {
// 	var err error

// 	if r.fs == nil {
// 		r.logger.Error(
// 			ErrProjectNotInitializedStorage,
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.Initialize",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 			})

// 		return fmt.Errorf(ErrProjectNotInitializedStorage)
// 	}

// 	if r.path == "" {
// 		r.logger.Error(
// 			ErrLocalProjectRepositoryPathNotProvided,
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.Initialize",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 			})

// 		return fmt.Errorf(ErrLocalProjectRepositoryPathNotProvided)
// 	}

// 	_, err = r.fs.Stat(r.path)
// 	if err != nil {
// 		r.logger.Info(
// 			"Creating local storage repository",
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.Initialize",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"path":      r.path,
// 			})

// 		err = r.fs.MkdirAll(r.path, 0755)
// 		if err != nil {
// 			r.logger.Error(
// 				fmt.Sprintf("%s: %s", ErrInitializingLocalProjectRepository, err.Error()),
// 				map[string]interface{}{
// 					"component": "LocalProjectRepository.Initialize",
// 					"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 					"path":      r.path,
// 				})

// 			return fmt.Errorf("%s: %w", ErrInitializingLocalProjectRepository, err)
// 		}
// 	}

// 	return nil
// }

// // LoadProjects loads projects from local storage
// func (r *LocalProjectRepository) LoadProjects() error {

// 	var err error
// 	var pathIsDir bool
// 	var projectEntity *entity.Project
// 	var projectFormat string
// 	var projectName string
// 	var projectPath string

// 	if r.store == nil {
// 		r.logger.Error(
// 			ErrProjectNotInitializedStorage,
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.LoadProjects",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 			})

// 		return fmt.Errorf(ErrProjectNotInitializedStorage)
// 	}

// 	if r.path == "" {
// 		r.logger.Error(
// 			ErrLocalProjectRepositoryPathNotProvided,
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.Initialize",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 			})

// 		return fmt.Errorf(ErrLocalProjectRepositoryPathNotProvided)
// 	}

// 	_, err = r.fs.Stat(r.path)
// 	if err != nil {
// 		r.logger.Error(
// 			fmt.Sprintf("%s: %s", ErrLocalProjectRepositoryPathNotExists, err.Error()),
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.LoadProjects",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"path":      r.path,
// 			})

// 		return fmt.Errorf("%s: %w", ErrLocalProjectRepositoryPathNotExists, err)
// 	}

// 	pathIsDir, err = afero.IsDir(r.fs, r.path)
// 	if err != nil {
// 		// This block handles an unexpected error returned by afero.IsDir. Since received error is not expected, it is logged and returned to avoid unexpected behavior
// 		r.logger.Error(
// 			fmt.Sprintf("%s: %s", ErrLocalProjectRepositoryCheckPathIsDirectory, err.Error()),
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.LoadProjects",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"path":      r.path,
// 			})

// 		return fmt.Errorf("%s: %w", ErrLocalProjectRepositoryCheckPathIsDirectory, err)
// 	}
// 	if !pathIsDir {
// 		r.logger.Error(
// 			ErrLocalProjectRepositoryPathMustBeDirectory,
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.LoadProjects",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"path":      r.path,
// 			})

// 		return fmt.Errorf(ErrLocalProjectRepositoryPathMustBeDirectory)
// 	}

// 	projects, err := afero.ReadDir(r.fs, r.path)
// 	if err != nil {
// 		// This block handles an unexpected error returned by afero.ReadDir. Since received error is not expected, it is logged and returned to avoid unexpected behavior
// 		r.logger.Error(
// 			fmt.Sprintf("%s: %s", ErrReadingLocalProjectRepositoryPath, err.Error()),
// 			map[string]interface{}{
// 				"component": "LocalProjectRepository.LoadProjects",
// 				"package":   "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"path":      r.path,
// 			})

// 		return fmt.Errorf("%s: %w", ErrReadingLocalProjectRepositoryPath, err)
// 	}

// 	for _, project := range projects {

// 		projectName = project.Name()

// 		if project.Mode().IsRegular() {
// 			// When is found a regular file and project name ends with .tar.gz we consider it as a tar.gz file, otherwise we skip the file
// 			if strings.HasSuffix(project.Name(), entity.ExtensionTarGz) {
// 				projectFormat = entity.ProjectFormatTarGz
// 				projectName = strings.TrimSuffix(project.Name(), entity.ExtensionTarGz)
// 			} else {
// 				continue
// 			}
// 		}

// 		if project.IsDir() {
// 			projectFormat = entity.ProjectFormatPlain
// 		}

// 		projectPath = filepath.Join(r.path, project.Name())
// 		projectEntity = entity.NewProject(projectName, projectPath, projectFormat, entity.ProjectTypeLocal)

// 		r.logger.Debug(
// 			"Loading project",
// 			map[string]interface{}{
// 				"component":    "LocalProjectRepository.LoadProjects",
// 				"package":      "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"project_id":   project.Name(),
// 				"project_path": projectPath,
// 			})

// 		err = r.SafeStore(projectName, projectEntity)
// 		if err != nil {
// 			r.logger.Error(
// 				fmt.Sprintf("%s: %s", ErrStoringProjectToLocalProjectRepository, err.Error()),
// 				map[string]interface{}{
// 					"component":  "LocalProjectRepository.LoadProjects",
// 					"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 					"project_id": project.Name(),
// 				})

// 			return fmt.Errorf("%s. %w", ErrStoringProjectToLocalProjectRepository, err)
// 		}
// 	}

// 	return nil
// }

// // Find returns a project by id
// func (r *LocalProjectRepository) Find(id string) (*entity.Project, error) {

// 	if r.store == nil {
// 		r.logger.Error(
// 			ErrProjectNotInitializedStorage,
// 			map[string]interface{}{
// 				"component":  "LocalProjectRepository.Find",
// 				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"project_id": id,
// 			})

// 		return nil, fmt.Errorf(ErrProjectNotInitializedStorage)
// 	}

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	project, ok := r.store[id]
// 	if !ok {
// 		r.logger.Error(
// 			ErrProjectNotFound,
// 			map[string]interface{}{
// 				"component":  "LocalProjectRepository.Find",
// 				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"project_id": id,
// 			})

// 		return nil, fmt.Errorf(ErrProjectNotFound)
// 	}

// 	return project, nil
// }

// // FindAll returns all projects
// func (r *LocalProjectRepository) FindAll() ([]*entity.Project, error) {

// 	// TODO: return sorted projects

// 	projects := []*entity.Project{}

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	for _, project := range r.store {
// 		projects = append(projects, project)
// 	}

// 	return projects, nil
// }

// // Remove removes a project by id
// func (r *LocalProjectRepository) Remove(id string) error {

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	if _, ok := r.store[id]; !ok {
// 		r.logger.Error(
// 			ErrProjectNotFound,
// 			map[string]interface{}{
// 				"component":  "LocalProjectRepository.Remove",
// 				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"project_id": id,
// 			})

// 		return fmt.Errorf(ErrProjectNotFound)
// 	}

// 	delete(r.store, id)

// 	return nil
// }

// // SafeStore stores a project by id
// func (r *LocalProjectRepository) SafeStore(id string, project *entity.Project) error {

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	if _, ok := r.store[id]; ok {
// 		r.logger.Error(
// 			ErrProjectAlreadyExists,
// 			map[string]interface{}{
// 				"component":  "LocalProjectRepository.SafeStore",
// 				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"project_id": id,
// 			})

// 		return fmt.Errorf(ErrProjectAlreadyExists)
// 	}

// 	r.store[id] = project

// 	return nil
// }

// // Store stores a project by id
// func (r *LocalProjectRepository) Store(id string, project *entity.Project) error {

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	r.store[id] = project

// 	return nil
// }

// // Update updates a project by id
// func (r *LocalProjectRepository) Update(id string, project *entity.Project) error {

// 	r.mutex.Lock()
// 	defer r.mutex.Unlock()

// 	if _, ok := r.store[id]; !ok {
// 		r.logger.Error(
// 			ErrProjectNotFound,
// 			map[string]interface{}{
// 				"component":  "LocalProjectRepository.Update",
// 				"package":    "github.com/apenella/ransidble/internal/infrastructure/persistence/project/repository",
// 				"project_id": id,
// 			})

// 		return fmt.Errorf(ErrProjectNotFound)
// 	}

// 	r.store[id] = project

// 	return nil
// }
