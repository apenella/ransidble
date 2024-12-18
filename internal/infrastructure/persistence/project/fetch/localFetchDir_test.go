package fetch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLocalFetchDirFromLocalFilesystem(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "persistence-project-fetch")
	sourceProject1 := filepath.Join("project-1")
	source := filepath.Join(sourceBase, sourceProject1)
	sourceFile := filepath.Join("site.yaml")
	workingDir := filepath.Join("working-dir")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		fetch       *LocalFetchDir
		source      string
		workingDir  string
		err         error
		arrangeFunc func(*testing.T, *LocalFetchDir)
		assertFunc  func(*testing.T, *LocalFetchDir)
	}{
		{
			desc:       "Testing fetch directory from local filesystem",
			fetch:      NewLocalFetchDir(fs, logger.NewFakeLogger()),
			source:     source,
			workingDir: workingDir,
			err:        nil,
			arrangeFunc: func(t *testing.T, fetch *LocalFetchDir) {
				fetch.fs.MkdirAll(workingDir, os.ModePerm)
			},
			assertFunc: func(t *testing.T, fetch *LocalFetchDir) {
				_, err := fetch.fs.Stat(filepath.Join(workingDir, sourceFile))
				assert.Nil(t, err)
			},
		},
		{
			desc:       "Testing error fetching directory from local filesystem when filesystem is not initialized",
			fetch:      NewLocalFetchDir(nil, logger.NewFakeLogger()),
			source:     source,
			workingDir: workingDir,
			err:        ErrFileSystemNotInitialized,
		},
		{
			desc:       "Testing error fetching directory from local filesystem when source directory does exists",
			fetch:      NewLocalFetchDir(fs, logger.NewFakeLogger()),
			source:     "not-exists",
			workingDir: workingDir,
			err:        ErrSourceCodeNotExists,
			arrangeFunc: func(t *testing.T, fetch *LocalFetchDir) {
				fetch.fs.MkdirAll(workingDir, os.ModePerm)
			},
		},
		{
			desc:       "Testing error fetching directory from local filesystem when working directory does not exists",
			fetch:      NewLocalFetchDir(fs, logger.NewFakeLogger()),
			source:     source,
			workingDir: "not-exists",
			err:        ErrWorkingDirNotExists,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.fetch)
			}

			err := test.fetch.Fetch(test.source, test.workingDir)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Nil(t, test.err)

				if test.assertFunc != nil {
					test.assertFunc(t, test.fetch)
				}
			}
		})
	}
}
