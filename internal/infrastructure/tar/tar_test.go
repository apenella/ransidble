package tar

import (
	"archive/tar"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {

	/*
		The test/fixtures/tar/extract-file.tar file content is:
			.
			├── dir
			│   └── file-2.txt
			└── file-1.txt
	*/

	sourceBase := filepath.Join("fixtures", "tar")
	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test"),
		),
		afero.NewMemMapFs(),
	)
	testFile := filepath.Join(sourceBase, "extract-file.tar")
	workingDir := filepath.Join("working-dir")
	sourceCodeFileReader, _ := fs.Open(testFile)

	tests := []struct {
		desc        string
		tar         *Tar
		reader      io.Reader
		destination string
		err         error
		arrangeFunc func(t *testing.T, tar *Tar)
		assertFunc  func(t *testing.T, tar *Tar)
	}{
		{
			desc:        "Testing extracting content from a tar file",
			tar:         NewTar(fs, logger.NewFakeLogger()),
			reader:      sourceCodeFileReader,
			destination: workingDir,
			err:         errors.New(""),
			arrangeFunc: func(t *testing.T, tar *Tar) {
				fs.RemoveAll(workingDir)
			},
			assertFunc: func(t *testing.T, tar *Tar) {
				_, err := fs.Stat(filepath.Join(workingDir, "file-1.txt"))
				assert.Nil(t, err)
				_, err = fs.Stat(filepath.Join(workingDir, "dir", "file-2.txt"))
				assert.Nil(t, err)

			},
		},
		{
			desc:        "Testing error extracting content from a tar file when reader is not provided",
			tar:         NewTar(fs, logger.NewFakeLogger()),
			reader:      nil,
			destination: workingDir,
			err:         ErrReaderNotProvided,
			arrangeFunc: func(t *testing.T, tar *Tar) {},
			assertFunc:  func(t *testing.T, tar *Tar) {},
		},
		{
			desc:        "Testing error extracting content from a tar file when destination is not provided",
			tar:         NewTar(fs, logger.NewFakeLogger()),
			reader:      sourceCodeFileReader,
			destination: "",
			err:         ErrDestinationNotProvided,
			arrangeFunc: func(t *testing.T, tar *Tar) {},
			assertFunc:  func(t *testing.T, tar *Tar) {},
		},
		{
			desc:        "Testing error extracting content from a tar file when filesystem is not provided",
			tar:         NewTar(nil, logger.NewFakeLogger()),
			reader:      sourceCodeFileReader,
			destination: workingDir,
			err:         ErrFilesystemNotProvided,
			arrangeFunc: func(t *testing.T, tar *Tar) {},
			assertFunc:  func(t *testing.T, tar *Tar) {},
		},
		{
			desc:        "Testing error extracting content from a tar file when a directory already exists on the working directory",
			tar:         NewTar(fs, logger.NewFakeLogger()),
			reader:      sourceCodeFileReader,
			destination: workingDir,
			err:         ErrCreatingFileFromTar,
			arrangeFunc: func(t *testing.T, tar *Tar) {
				fs.MkdirAll(filepath.Join(workingDir, "dir"), 0755)
			},
			assertFunc: func(t *testing.T, tar *Tar) {},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			// t.Parallel()
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.tar)
			}

			err := test.tar.Extract(test.reader, test.destination)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)

				test.assertFunc(t, test.tar)
			}
		})
	}
}

func TestExtractRegularFile(t *testing.T) {
	// The test only validates the input arguments. The file extraction is validated in TestExtract

	tests := []struct {
		desc        string
		tar         *Tar
		tr          *tar.Reader
		header      *tar.Header
		destination string
		err         error
	}{
		{
			desc:   "Testing error extracting a tar file when the header is not provided",
			tar:    NewTar(afero.NewMemMapFs(), logger.NewFakeLogger()),
			header: nil,
			err:    ErrTarFileHeaderNotProvided,
		},
		{
			desc:        "Testing error extracting a tar file when destination is not provided",
			tar:         NewTar(afero.NewMemMapFs(), logger.NewFakeLogger()),
			header:      &tar.Header{},
			destination: "",
			err:         ErrDestinationNotProvided,
		},
		{
			desc:        "Testing error extracting a tar file when tar reader is not provided",
			tar:         NewTar(afero.NewMemMapFs(), logger.NewFakeLogger()),
			tr:          nil,
			header:      &tar.Header{},
			destination: "destination",
			err:         ErrReaderNotProvided,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			err := test.tar.extractRegularFile(test.tr, test.header, test.destination)
			assert.NotNil(t, err)
			assert.Equal(t, err, test.err)
		})
	}
}
