package fetch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFetchFileFromLocalFilesystem(t *testing.T) {
	sourceBase := filepath.Join("fixtures", "persistence-project-fetch")
	sourceProject2 := filepath.Join("project-2.tar.gz")
	source := filepath.Join(sourceBase, sourceProject2)
	workingDir := filepath.Join("working-dir")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		fetch       *LocalFetchFile
		source      string
		workingDir  string
		err         error
		arrangeFunc func(*testing.T, *LocalFetchFile)
		assertFunc  func(*testing.T, *LocalFetchFile)
	}{
		{
			desc:       "Testing fetch file from local filesystem",
			fetch:      NewLocalFetchFile(fs, logger.NewFakeLogger()),
			source:     source,
			workingDir: workingDir,
			err:        nil,
			arrangeFunc: func(t *testing.T, fetch *LocalFetchFile) {
				fetch.fs.MkdirAll(workingDir, os.ModePerm)
			},
			assertFunc: func(t *testing.T, fetch *LocalFetchFile) {
				_, err := fetch.fs.Stat(filepath.Join(workingDir, sourceProject2))
				assert.Nil(t, err)
			},
		},
		{
			desc:       "Testing error fetching file from local filesystem when source file does exists",
			fetch:      NewLocalFetchFile(fs, logger.NewFakeLogger()),
			source:     "not-exists",
			workingDir: workingDir,
			err:        ErrSourceCodeNotExists,
			arrangeFunc: func(t *testing.T, fetch *LocalFetchFile) {
				fetch.fs.MkdirAll(workingDir, os.ModePerm)
			},
		},
		{
			desc:       "Testing error fetching file from local filesystem when working directory does exists",
			fetch:      NewLocalFetchFile(fs, logger.NewFakeLogger()),
			source:     source,
			workingDir: "not-exists",
			err:        ErrWorkingDirNotExists,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.fetch)
			}

			err := test.fetch.Fetch(test.source, test.workingDir)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				test.assertFunc(t, test.fetch)
			}
		})
	}
}
