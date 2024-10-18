package workspace

import (
	"fmt"
	"path/filepath"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/spf13/afero"
)

var (
	// ErrSourceCodeFetcherNotProvided represents an error when the source code fetcher is not provided
	ErrSourceCodeFetcherNotProvided = fmt.Errorf("source code fetcher not provided")
	// ErrSourceCodeUnpackerNotProvided represents an error when the source code unpacker is not provided
	ErrSourceCodeUnpackerNotProvided = fmt.Errorf("source code unpacker not provided")
	// ErrProjectRepositoryNotProvided represents an error when the project repository is not provided
	ErrProjectRepositoryNotProvided = fmt.Errorf("project repository not provided")
	// ErrTaskNotProvided represents an error when the task is not provided
	ErrTaskNotProvided = fmt.Errorf("task not provided")
	// ErrProjectNotProvided represents an error when the project is not provided
	ErrProjectNotProvided = fmt.Errorf("project not provided")
	// ErrFindingProject represents an error when the project is not found
	ErrFindingProject = fmt.Errorf("error finding project")
	// ErrWorkingDirNotDefined represents an error when the workspace path is not defined
	ErrWorkingDirNotDefined = fmt.Errorf("workspace path not defined")
	// ErrFetchingProject represents an error when fetching a project
	ErrFetchingProject = fmt.Errorf("error fetching project")
	// ErrUnpackingProject represents an error when unpacking a project
	ErrUnpackingProject = fmt.Errorf("error unpacking project")
	// ErrProjectFetcherNotAvailable represents an error when the project fetcher is not available
	ErrProjectFetcherNotAvailable = fmt.Errorf("project fetcher not available")
	// ErrProjectUnpackerNotAvailable represents an error when the project unpacker is not available
	ErrProjectUnpackerNotAvailable = fmt.Errorf("project unpacker not available")
	// ErrWorkingDirAlreadyExists represents an error when the working directory already exists
	ErrWorkingDirAlreadyExists = fmt.Errorf("working directory already exists")
	// ErrCreatingWorkingDirFolder represents an error when the working directory cannot be created
	ErrCreatingWorkingDirFolder = fmt.Errorf("error creating working directory folder")
)

// FuncOptions represents the function to set the options
type FuncOptions func(*Workspace)

// Workspace represents the location where the project is stored before being executed
type Workspace struct {
	// fetchFactory returns the fetcher to get the project from the catalog
	fetchFactory repository.SourceCodeFetchFactory
	// fs is the filesystem
	fs afero.Fs
	// workingDir is the working directory path, where the project source code is stored
	workingDir string
	// repository is the repository to get the project from the catalog
	repository repository.ProjectRepository
	// task is the task to be executed
	task *entity.Task
	// unpackFactory returns the unpacker to unpack the project
	unpackFactory repository.SourceCodeUnpackFactory
	// logger is the logger
	logger repository.Logger
}

// NewWorkspace creates a new workspace
func NewWorkspace(options ...FuncOptions) *Workspace {
	w := &Workspace{}

	w.Options(options...)
	return w
}

// Options sets the options for the workspace
func (w *Workspace) Options(options ...FuncOptions) *Workspace {
	for _, option := range options {
		option(w)
	}
	return w
}

// Prepare fetches and unpacks a project
func (w *Workspace) Prepare() error {

	var err error
	var workingDir string

	if w.fetchFactory == nil {
		w.logger.Error(ErrSourceCodeFetcherNotProvided.Error(), map[string]interface{}{
			"component": "Workspace.Prepare",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
		})

		return ErrSourceCodeFetcherNotProvided
	}

	if w.unpackFactory == nil {
		w.logger.Error(ErrSourceCodeUnpackerNotProvided.Error(), map[string]interface{}{
			"component": "Workspace.Prepare",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
		})

		return ErrSourceCodeUnpackerNotProvided
	}

	if w.repository == nil {
		w.logger.Error(ErrProjectRepositoryNotProvided.Error(), map[string]interface{}{
			"component": "Workspace.Prepare",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
		})

		return ErrProjectRepositoryNotProvided
	}

	if w.task == nil {
		w.logger.Error(ErrTaskNotProvided.Error(), map[string]interface{}{
			"component": "Workspace.Prepare",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
		})

		return ErrTaskNotProvided
	}

	projectID := w.task.ProjectID
	if projectID == "" {
		w.logger.Error(ErrProjectNotProvided.Error(), map[string]interface{}{
			"component": "Workspace.Prepare",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
		})

		return ErrProjectNotProvided
	}

	project, err := w.repository.Find(projectID)
	if err != nil {
		w.logger.Error(fmt.Sprintf("%s: %s", ErrFindingProject, err.Error()), map[string]interface{}{
			"component":  "Workspace.Prepare",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id": projectID,
		})

		return domainerror.NewProjectNotFoundError(
			fmt.Errorf("%s %s: %w", ErrFindingProject, projectID, err),
		)
	}

	workingDir, err = w.generateWorkingDirPath(projectID, w.task.ID)
	if err != nil {
		w.logger.Error(fmt.Sprintf("%s: %s", "error generating workspace path", err.Error()), map[string]interface{}{
			"component":  "Workspace.generateWorkingDir",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id": projectID,
			"task_id":    w.task.ID,
		})

		return fmt.Errorf("%s: %w", "error generating workspace path", err)
	}
	w.workingDir = workingDir

	_, err = w.fs.Stat(workingDir)
	if err == nil {
		w.logger.Error(
			ErrWorkingDirAlreadyExists.Error(),
			map[string]interface{}{
				"component":   "Workspace.Prepare",
				"package":     "github.com/apenella/ransidble/internal/domain/core/service/workspace",
				"project_id":  projectID,
				"task_id":     w.task.ID,
				"working_dir": workingDir,
			})

		return ErrWorkingDirAlreadyExists
	}

	err = w.fs.MkdirAll(workingDir, 0755)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %s", ErrCreatingWorkingDirFolder, err)
		w.logger.Error(
			errorMsg,
			map[string]interface{}{
				"component":   "Workspace.Prepare",
				"package":     "github.com/apenella/ransidble/internal/domain/core/service/workspace",
				"project_id":  projectID,
				"task_id":     w.task.ID,
				"working_dir": workingDir,
			})
		return fmt.Errorf("%s", errorMsg)
	}

	fetcher := w.fetchFactory.Get(project.Type)
	if fetcher == nil {
		w.logger.Error(ErrProjectFetcherNotAvailable.Error(), map[string]interface{}{
			"component":    "Workspace.Prepare",
			"package":      "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id":   projectID,
			"task_id":      w.task.ID,
			"storage_type": project.Type,
		})

		return ErrProjectFetcherNotAvailable
	}

	err = fetcher.Fetch(project, workingDir)
	if err != nil {
		w.logger.Error(fmt.Sprintf("%s: %s", ErrFetchingProject.Error(), err.Error()), map[string]interface{}{
			"component":  "Workspace.Prepare",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id": projectID,
			"task_id":    w.task.ID,
		})

		return fmt.Errorf("%s: %w", ErrFetchingProject.Error(), err)
	}

	unpacker := w.unpackFactory.Get(project.Format)
	if unpacker == nil {
		w.logger.Error(ErrProjectUnpackerNotAvailable.Error(), map[string]interface{}{
			"component":      "Workspace.Prepare",
			"package":        "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id":     projectID,
			"task_id":        w.task.ID,
			"project_format": project.Format,
		})

		return ErrProjectUnpackerNotAvailable
	}

	err = unpacker.Unpack(project, workingDir)
	if err != nil {
		w.logger.Error(fmt.Sprintf("%s: %s", ErrUnpackingProject.Error(), err.Error()), map[string]interface{}{
			"component":  "Workspace.Prepare",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id": projectID,
			"task_id":    w.task.ID,
		})

		return fmt.Errorf("%s: %w", ErrUnpackingProject.Error(), err)
	}

	return nil
}

// generateWorkingDirPath generates the workspace path
func (w *Workspace) generateWorkingDirPath(projectID, taskID string) (string, error) {

	workspaceBaseDir, err := afero.TempDir(w.fs, "", "ransidble")
	if err != nil {
		w.logger.Error(fmt.Sprintf("%s %s %s: %s", "temporal directory cannot be created", projectID, taskID, err), map[string]interface{}{
			"component":  "Workspace.generateWorkingDirPath",
			"package":    "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"project_id": projectID,
			"task_id":    taskID,
		})

		err = fmt.Errorf("%s. project: %s. task: %s. %w", "temporal directory cannot be created", projectID, taskID, err)
		return "", err
	}

	return filepath.Join(workspaceBaseDir, projectID, taskID), nil
}

// GetWorkingDir returns the working directory
func (w *Workspace) GetWorkingDir() (string, error) {

	if w.workingDir == "" {
		w.logger.Error(ErrWorkingDirNotDefined.Error(), map[string]interface{}{
			"component": "Workspace.GetWorkingDir",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"task_id":   w.task.ID,
		})
		return "", ErrWorkingDirNotDefined
	}

	return w.workingDir, nil
}

// Cleanup cleans the workspace
func (w *Workspace) Cleanup() error {

	if w.workingDir == "" {
		w.logger.Error(ErrWorkingDirNotDefined.Error(), map[string]interface{}{
			"component": "Workspace.Cleanup",
			"package":   "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"task_id":   w.task.ID,
		})
		return ErrWorkingDirNotDefined
	}

	err := w.fs.RemoveAll(w.workingDir)
	if err != nil {
		w.logger.Error(fmt.Sprintf("%s: %s", "error removing working directory", err.Error()), map[string]interface{}{
			"component":   "Workspace.Cleanup",
			"package":     "github.com/apenella/ransidble/internal/domain/core/service/workspace",
			"task_id":     w.task.ID,
			"working_dir": w.workingDir,
		})
		return fmt.Errorf("%s: %w", "error removing working directory", err)
	}

	return nil
}
