package fetch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "persistence-project-fetch")
	workingDir := filepath.Join("working-dir")

	project1Name := "project-1"
	project2Name := "project-2"
	project2File := strings.Join([]string{project2Name, "tar", "gz"}, ".")

	project1ExpectedFile := filepath.Join(workingDir, "site.yaml")
	project2ExpectedFile := filepath.Join(workingDir, project2File)
	sourceProject1 := filepath.Join(project1Name)
	sourceProject2 := filepath.Join(project2File)

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		storage     *LocalStorage
		project     *entity.Project
		workingDir  string
		err         error
		arrangeFunc func(*testing.T, *LocalStorage)
		assertFunc  func(*testing.T, *LocalStorage)
	}{
		{
			desc:    "Testing fetch a project in plain format from local storage",
			storage: NewLocalStorage(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      project1Name,
				Reference: filepath.Join(sourceBase, sourceProject1),
				Format:    "plain",
				Storage:   "local",
			},
			workingDir: workingDir,
			err:        nil,
			arrangeFunc: func(t *testing.T, storage *LocalStorage) {
				storage.fs.MkdirAll(workingDir, os.ModePerm)
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Stat(project1ExpectedFile)
				assert.Nil(t, err)
			},
		},
		{
			desc:    "Testing fetch a project in tar.gz format from local storage",
			storage: NewLocalStorage(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      project2Name,
				Reference: filepath.Join(sourceBase, sourceProject2),
				Format:    "targz",
				Storage:   "local",
			},
			workingDir: workingDir,
			err:        nil,
			arrangeFunc: func(t *testing.T, storage *LocalStorage) {
				storage.fs.MkdirAll(workingDir, os.ModePerm)
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Stat(project2ExpectedFile)
				assert.Nil(t, err)
			},
		},
		{
			desc:       "Testing error fetching a project from local storage when project is not provided",
			storage:    NewLocalStorage(fs, logger.NewFakeLogger()),
			project:    nil,
			workingDir: workingDir,
			err:        ErrProjectNotProvided,
		},
		{
			desc:    "Testing error fetching a project from local storage when working directory is not provided",
			storage: NewLocalStorage(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      project1Name,
				Reference: filepath.Join(sourceBase, sourceProject1),
				Format:    "plain",
				Storage:   "local",
			},
			workingDir: "",
			err:        ErrWorkingDirNotProvided,
		},
		{
			desc: "Testing error fetching a project from local storage when filesystem is not initialized",
			storage: &LocalStorage{
				fs:     nil,
				logger: logger.NewFakeLogger(),
			},
			project: &entity.Project{
				Name:      project1Name,
				Reference: filepath.Join(sourceBase, sourceProject1),
				Format:    "plain",
				Storage:   "local",
			},
			workingDir: workingDir,
			err:        ErrFileSystemNotInitialized,
		},
		{
			desc:    "Testing error fetching a project from local storage when working directory does not exists",
			storage: NewLocalStorage(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      project1Name,
				Reference: filepath.Join(sourceBase, sourceProject1),
				Format:    "plain",
				Storage:   "local",
			},
			workingDir: "not-exists",
			err:        ErrWorkingDirNotExists,
		},
		{
			desc:    "Testing error fetching a project from local storage when project reference is not provided",
			storage: NewLocalStorage(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      project1Name,
				Reference: "",
				Format:    "plain",
				Storage:   "local",
			},
			workingDir: workingDir,
			err:        ErrProjectReferenceNotProvided,
		},
		{
			desc:    "Testing error fetching a project from local storage when project reference is invalid",
			storage: NewLocalStorage(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      project1Name,
				Reference: "not-exists",
				Format:    "plain",
				Storage:   "local",
			},
			workingDir: workingDir,
			err:        fmt.Errorf("%s: %w", ErrInvalidProjectReference, errors.New("stat ../../../../../test/not-exists: no such file or directory")),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.storage)
			}

			err := test.storage.Fetch(test.project, test.workingDir)
			if test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)

				if test.assertFunc != nil {
					test.assertFunc(t, test.storage)
				}
			}
		})
	}
}
