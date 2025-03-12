package store

import (
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLocalStorage_Initialize(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "persistence-project-store")
	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		storage     *LocalStorage
		arrangeFunc func(*testing.T, *LocalStorage)
		assertFunc  func(*testing.T, *LocalStorage)
		err         error
	}{
		{
			desc: "Testing initialize local storage when path already exists",
			storage: NewLocalStorage(
				fs,
				sourceBase,
				logger.NewFakeLogger(),
			),
			assertFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Stat(sourceBase)
				assert.Nil(t, err, fmt.Sprintf("error checking file %s", sourceBase))
			},
			err: nil,
		},
		{
			desc: "Testing initialize local storage when path does not exists",
			storage: NewLocalStorage(
				fs,
				filepath.Join(sourceBase, "storage"),
				logger.NewFakeLogger(),
			),
			assertFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Stat(filepath.Join(sourceBase, "storage"))
				assert.Nil(t, err)
			},
			err: nil,
		},
		{
			desc: "Testing error initializing local storage when filesystem is not initialized",
			storage: NewLocalStorage(
				nil,
				sourceBase,
				logger.NewFakeLogger(),
			),
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			err:        fmt.Errorf(ErrStorageHandlerNotInitialized),
		},
		{
			desc: "Testing error initializing local storage when path is not provided",
			storage: NewLocalStorage(
				fs,
				"",
				logger.NewFakeLogger(),
			),
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			err:        fmt.Errorf(ErrStoragePathNotProvided),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.storage)
			}

			err := test.storage.Initialize()
			if err != nil && test.err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Nil(t, test.err)

				if test.assertFunc != nil {
					test.assertFunc(t, test.storage)
				}
			}
		})
	}
}

func TestLocalStorage_Store(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "persistence-project-store")
	sourceProjectFile := "project-1.tar.gz"
	localStoragePath := filepath.Join("local-storage")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	err := fs.MkdirAll(localStoragePath, 0755)
	if err != nil {
		t.Fatalf("error creating local storage directory: %s", err)
	}

	srcFile, err := fs.Open(filepath.Join(sourceBase, sourceProjectFile))
	if err != nil {
		t.Fatalf("error opening source file: %s", err)
	}

	tests := []struct {
		desc        string
		storage     *LocalStorage
		project     *entity.Project
		srcFile     io.Reader
		arrangeFunc func(*testing.T, *LocalStorage)
		assertFunc  func(*testing.T, *LocalStorage)
		err         error
	}{
		{
			desc:    "Testing store a project in local storage",
			storage: NewLocalStorage(fs, localStoragePath, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-1",
				Reference: sourceProjectFile,
				Format:    "targz",
				Storage:   "local",
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Stat(filepath.Join(localStoragePath, sourceProjectFile))
				assert.Nil(t, err, fmt.Sprintf("error checking file %s", filepath.Join(localStoragePath, sourceProjectFile)))
			},
			srcFile: srcFile,
			err:     nil,
		},
		{
			desc:    "Testing store a project overwriting a project in local storage",
			storage: NewLocalStorage(fs, localStoragePath, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-1",
				Reference: sourceProjectFile,
				Format:    "targz",
				Storage:   "local",
			},
			arrangeFunc: func(t *testing.T, storage *LocalStorage) {
				storage.Store(&entity.Project{
					Name:      "project-1",
					Reference: sourceProjectFile,
					Format:    "targz",
					Storage:   "local",
				}, srcFile)
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Stat(filepath.Join(localStoragePath, sourceProjectFile))
				assert.Nil(t, err, fmt.Sprintf("error checking file %s", filepath.Join(localStoragePath, sourceProjectFile)))
			},
			srcFile: srcFile,
			err:     nil,
		},
		{
			desc:       "Testing error storing a project in local storage when project is not provided",
			storage:    NewLocalStorage(fs, localStoragePath, logger.NewFakeLogger()),
			project:    nil,
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    srcFile,
			err:        fmt.Errorf(ErrProjectNotProvided),
		},
		{
			desc:       "Testing error storing a project in local storage when file is not provided",
			storage:    NewLocalStorage(fs, localStoragePath, logger.NewFakeLogger()),
			project:    &entity.Project{},
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    nil,
			err:        fmt.Errorf(ErrProjectFileNotProvided),
		},
		{
			desc:       "Testing error storing a project in local storage when project reference is not provided",
			storage:    NewLocalStorage(fs, localStoragePath, logger.NewFakeLogger()),
			project:    &entity.Project{},
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    srcFile,
			err:        fmt.Errorf(ErrProjectReferenceNotProvided),
		},
		{
			desc:    "Testing error storing a project in local storage when storage filesystem is not initialized",
			storage: NewLocalStorage(nil, "", logger.NewFakeLogger()),
			project: &entity.Project{
				Reference: sourceProjectFile,
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    srcFile,
			err:        fmt.Errorf(ErrStorageHandlerNotInitialized),
		},
		{
			desc:    "Testing error storing a project in local storage when storage path is not provided",
			storage: NewLocalStorage(fs, "", logger.NewFakeLogger()),
			project: &entity.Project{
				Reference: sourceProjectFile,
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    srcFile,
			err:        fmt.Errorf(ErrStoragePathNotProvided),
		},
		{
			desc:    "Testing error storing a project in local storage when storage path does not exists",
			storage: NewLocalStorage(fs, "unexisting", logger.NewFakeLogger()),
			project: &entity.Project{
				Reference: sourceProjectFile,
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    srcFile,
			err:        fmt.Errorf("%s: %s", ErrStoragePathNotExists, "stat ../../../../../test/unexisting: no such file or directory"),
		},
		{
			desc:    "Testing error storing a project in local storage when storage path is not a directory",
			storage: NewLocalStorage(fs, "invalid-storage-path", logger.NewFakeLogger()),
			project: &entity.Project{
				Reference: sourceProjectFile,
			},
			arrangeFunc: func(t *testing.T, storage *LocalStorage) {
				_, err := storage.fs.Create("invalid-storage-path")
				assert.Nil(t, err, fmt.Sprintf("error creating file %s", "invalid-storage-path"))
			},
			srcFile: srcFile,
			err:     fmt.Errorf(ErrStoragePathNotDirectory),
		},
		{
			desc:    "Testing error storing a project in local storage when opening destination file fails",
			storage: NewLocalStorage(fs, localStoragePath, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-1",
				Reference: filepath.Join("unexisting", sourceProjectFile),
				Format:    "targz",
				Storage:   "local",
			},
			assertFunc: func(t *testing.T, storage *LocalStorage) {},
			srcFile:    srcFile,
			err:        fmt.Errorf("%s: %s", ErrOpeningDestinationFileInLocalStorage, "open local-storage/unexisting: file does not exist"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.storage)
			}

			err := test.storage.Store(test.project, test.srcFile)
			if err != nil && test.err != nil {
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
