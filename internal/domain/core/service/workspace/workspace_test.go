package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	domainerror "github.com/apenella/ransidble/internal/domain/core/error"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWorkingDir(t *testing.T) {

	tests := []struct {
		desc      string
		workspace *Workspace
		expected  string
		err       error
	}{
		{
			desc: "Testing error getting the working directory from a workspace when task is not provided",
			workspace: &Workspace{
				workingDir: "",
				logger:     logger.NewFakeLogger(),
			},
			expected: "",
			err:      ErrTaskNotProvided,
		},
		{
			desc: "Testing error getting the working directory from a workspace when the working directory is not set",
			workspace: &Workspace{
				workingDir: "",
				logger:     logger.NewFakeLogger(),
				task: &entity.Task{
					ID: "task-id",
				},
			},
			expected: "",
			err:      ErrWorkingDirNotDefined,
		},
		{
			desc: "Testing the GetWorkingDir from a workspace",
			workspace: &Workspace{
				workingDir: "/tmp",
				logger:     logger.NewFakeLogger(),
				task: &entity.Task{
					ID: "task-id",
				},
			},
			expected: "/tmp",
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			workingDir, err := test.workspace.GetWorkingDir()
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, test.expected, workingDir)
			}
		})
	}
}

func TestGenerateWorkingDirPath(t *testing.T) {
	tests := []struct {
		desc        string
		workspace   *Workspace
		projectID   string
		taskID      string
		err         error
		arrangeFunc func(*testing.T, *Workspace)
	}{
		{
			desc: "Testing error generating the working directory path from a workspace when project id is not provided",
			workspace: &Workspace{
				workingDir: "",
				logger:     logger.NewFakeLogger(),
			},
			projectID: "",
			taskID:    "task-id",
			err:       ErrProjectNotProvided,
		},
		{
			desc: "Testing error generating the working directory path from a workspace when task id is not provided",
			workspace: &Workspace{
				workingDir: "",
				logger:     logger.NewFakeLogger(),
			},
			projectID: "project-id",
			taskID:    "",
			err:       ErrTaskNotProvided,
		},
		{
			desc: "Testing error generating the temporary directory path from a workspace when filesystem is not provided",
			workspace: &Workspace{
				workingDir: "/tmp",
				logger:     logger.NewFakeLogger(),
			},
			projectID: "project-id",
			taskID:    "task-id",
			err:       ErrFilesystemNotProvided,
		},
		{
			desc: "Testing the GenerateWorkingDirPath from a workspace",
			workspace: &Workspace{
				workingDir: "/tmp",
				logger:     logger.NewFakeLogger(),
				fs:         repository.NewMockFilesystemer(),
			},
			projectID: "project-id",
			taskID:    "task-id",
			arrangeFunc: func(t *testing.T, w *Workspace) {
				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return("/tmp/ransidble", nil)
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.workspace)
			}

			workingDir, err := test.workspace.generateWorkingDirPath(test.projectID, test.taskID)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Equal(t, filepath.Join("/tmp/ransidble", test.projectID, test.taskID), workingDir)
			}
		})
	}
}

func TestPrepare(t *testing.T) {
	tests := []struct {
		desc        string
		workspace   *Workspace
		err         error
		arrangeFunc func(*testing.T, *Workspace)
	}{
		{
			desc: "Testing error preparing the workspace when fetchFactory is not provided",
			workspace: &Workspace{
				logger: logger.NewFakeLogger(),
			},
			err: ErrSourceCodeFetcherNotProvided,
		},
		{
			desc: "Testing error preparing the workspace when unpackFactory is not provided",
			workspace: &Workspace{
				logger:       logger.NewFakeLogger(),
				fetchFactory: &repository.MockProjectSourceCodeFetchFactory{},
			},
			err: ErrSourceCodeUnpackerNotProvided,
		},
		{
			desc: "Testing error preparing the workspace when repository is not provided",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
			},
			err: ErrProjectRepositoryNotProvided,
		},
		{
			desc: "Testing error preparing the workspace when task is not provided",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
			},
			err: ErrTaskNotProvided,
		},
		{
			desc: "Testing error preparing the workspace when filesystem is not provided",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task:          &entity.Task{},
			},
			err: ErrFilesystemNotProvided,
		},
		{
			desc: "Testing error preparing the workspace when project id is not provided in the task",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task:          &entity.Task{},
				fs:            repository.NewMockFilesystemer(),
			},
			err: ErrProjectNotProvided,
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when finding the project",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: domainerror.NewProjectNotFoundError(
				fmt.Errorf("%s %s: %w", ErrFindingProject, "project-id", errors.New("error finding project")),
			),
			arrangeFunc: func(t *testing.T, w *Workspace) {
				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(nil, errors.New("error finding project"))
			},
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when generating the working directory path",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: fmt.Errorf("%s: %w", "error generating workspace path",
				fmt.Errorf("%s. project: project-id. task: task-id. %w", "temporal directory cannot be created",
					errors.New("error generating working directory path")),
			),
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)
				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return("", errors.New("error generating working directory path"))
			},
		},
		{
			desc: "Testing error preparing the workspace when the working directory already exists",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: ErrWorkingDirAlreadyExists,
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				workingDir := "/tmp"

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)
				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil, nil)
			},
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when creating the working directory",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: fmt.Errorf(fmt.Sprintf("%s: %s", ErrCreatingWorkingDirFolder, errors.New("error creating working directory"))),
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				workingDir := "/tmp"

				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil, os.ErrNotExist)
				w.fs.(*repository.MockFilesystemer).On("MkdirAll", filepath.Join(workingDir, w.task.ProjectID, w.task.ID), mock.Anything).Return(errors.New("error creating working directory"))

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)
			},
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when getting the fetcher",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: ErrProjectFetcherNotAvailable,
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				workingDir := "/tmp"

				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil, os.ErrNotExist)
				w.fs.(*repository.MockFilesystemer).On("MkdirAll", filepath.Join(workingDir, w.task.ProjectID, w.task.ID), mock.Anything).Return(nil)

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)

				w.fetchFactory.(*repository.MockProjectSourceCodeFetchFactory).On("Get", "local").Return(nil)
			},
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when fetching the project source code",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: fmt.Errorf("%s: %w", ErrFetchingProject.Error(), errors.New("error fetching project source code")),
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				workingDir := "/tmp"

				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil, os.ErrNotExist)
				w.fs.(*repository.MockFilesystemer).On("MkdirAll", filepath.Join(workingDir, w.task.ProjectID, w.task.ID), mock.Anything).Return(nil)

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)

				fetcher := &repository.MockProjectSourceCodeFetcher{}
				fetcher.On("Fetch", project, filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(errors.New("error fetching project source code"))

				w.fetchFactory.(*repository.MockProjectSourceCodeFetchFactory).On("Get", "local").Return(fetcher)
			},
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when getting the unpacker",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: ErrProjectUnpackerNotAvailable,
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				workingDir := "/tmp"

				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil, os.ErrNotExist)
				w.fs.(*repository.MockFilesystemer).On("MkdirAll", filepath.Join(workingDir, w.task.ProjectID, w.task.ID), mock.Anything).Return(nil)

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)

				fetcher := &repository.MockProjectSourceCodeFetcher{}
				fetcher.On("Fetch", project, filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil)

				w.fetchFactory.(*repository.MockProjectSourceCodeFetchFactory).On("Get", "local").Return(fetcher)

				w.unpackFactory.(*repository.MockProjectSourceCodeUnpackFactory).On("Get", "plain").Return(nil)
			},
		},
		{
			desc: "Testing error preparing the workspace when an error occurs when unpacking the project source code",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ProjectID: "project-id",
					ID:        "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: fmt.Errorf("%s: %w", ErrUnpackingProject.Error(), errors.New("error unpacking project source code")),
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				workingDir := "/tmp"

				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil, os.ErrNotExist)
				w.fs.(*repository.MockFilesystemer).On("MkdirAll", filepath.Join(workingDir, w.task.ProjectID, w.task.ID), mock.Anything).Return(nil)

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)

				fetcher := &repository.MockProjectSourceCodeFetcher{}
				fetcher.On("Fetch", project, filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(nil)

				w.fetchFactory.(*repository.MockProjectSourceCodeFetchFactory).On("Get", "local").Return(fetcher)

				unpacker := &repository.MockProjectSourceCodeUnpacker{}
				unpacker.On("Unpack", project, filepath.Join(workingDir, w.task.ProjectID, w.task.ID)).Return(errors.New("error unpacking project source code"))

				w.unpackFactory.(*repository.MockProjectSourceCodeUnpackFactory).On("Get", "plain").Return(unpacker)
			},
		},
		{
			desc: "Testing preparing the workspace",
			workspace: &Workspace{
				logger:        logger.NewFakeLogger(),
				fetchFactory:  &repository.MockProjectSourceCodeFetchFactory{},
				unpackFactory: &repository.MockProjectSourceCodeUnpackFactory{},
				repository:    &repository.MockProjectRepository{},
				task: &entity.Task{
					ID:        "task-id",
					ProjectID: "project-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: nil,
			arrangeFunc: func(t *testing.T, w *Workspace) {
				project := &entity.Project{
					Format:    "plain",
					Name:      "project-id",
					Reference: "project-id",
					Storage:   "local",
				}
				task := &entity.Task{
					ID:        "task-id",
					ProjectID: "project-id",
				}

				workingDir := "/tmp"

				w.fs.(*repository.MockFilesystemer).On("TempDir", "", "ransidble").Return(workingDir, nil)
				w.fs.(*repository.MockFilesystemer).On("Stat", filepath.Join(workingDir, task.ProjectID, task.ID)).Return(nil, os.ErrNotExist)
				w.fs.(*repository.MockFilesystemer).On("MkdirAll", filepath.Join(workingDir, task.ProjectID, task.ID), mock.Anything).Return(nil)

				w.repository.(*repository.MockProjectRepository).On("Find", "project-id").Return(project, nil)

				fetcher := &repository.MockProjectSourceCodeFetcher{}
				fetcher.On("Fetch", project, filepath.Join(workingDir, task.ProjectID, task.ID)).Return(nil)

				w.fetchFactory.(*repository.MockProjectSourceCodeFetchFactory).On("Get", "local").Return(fetcher)

				unpacker := &repository.MockProjectSourceCodeUnpacker{}
				unpacker.On("Unpack", project, filepath.Join(workingDir, task.ProjectID, task.ID)).Return(nil)

				w.unpackFactory.(*repository.MockProjectSourceCodeUnpackFactory).On("Get", "plain").Return(unpacker)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.workspace)
			}

			err := test.workspace.Prepare()
			assert.Equal(t, test.err, err)
		})
	}
}

func TestCleanup(t *testing.T) {
	tests := []struct {
		desc        string
		workspace   *Workspace
		err         error
		arrangeFunc func(*testing.T, *Workspace)
	}{
		{
			desc: "Testing error cleaning up the workspace when filesystem is not provided",
			workspace: &Workspace{
				logger: logger.NewFakeLogger(),
				task: &entity.Task{
					ID: "task-id",
				},
			},
			err: ErrFilesystemNotProvided,
		},
		{
			desc: "Testing error cleaning up the workspace when working directory is not defined",
			workspace: &Workspace{
				logger: logger.NewFakeLogger(),
				task: &entity.Task{
					ID: "task-id",
				},
				fs: repository.NewMockFilesystemer(),
			},
			err: ErrWorkingDirNotDefined,
		},
		{
			desc: "Testing error cleaning up the workspace when an error occurs when removing the working directory",
			workspace: &Workspace{
				logger: logger.NewFakeLogger(),
				task: &entity.Task{
					ID: "task-id",
				},
				fs:         repository.NewMockFilesystemer(),
				workingDir: "/tmp",
			},
			err: fmt.Errorf("%s: %w", ErrRemoveWorkingDir, errors.New("error removing working directory")),
			arrangeFunc: func(t *testing.T, w *Workspace) {
				w.fs.(*repository.MockFilesystemer).On("RemoveAll", w.workingDir).Return(errors.New("error removing working directory"))
			},
		},
		{
			desc: "Testing cleaning up the workspace",
			workspace: &Workspace{
				logger: logger.NewFakeLogger(),
				task: &entity.Task{
					ID: "task-id",
				},
				fs:         repository.NewMockFilesystemer(),
				workingDir: "/tmp",
			},
			err: nil,
			arrangeFunc: func(t *testing.T, w *Workspace) {
				w.fs.(*repository.MockFilesystemer).On("RemoveAll", w.workingDir).Return(nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.workspace)
			}

			err := test.workspace.Cleanup()
			assert.Equal(t, test.err, err)
		})
	}
}
