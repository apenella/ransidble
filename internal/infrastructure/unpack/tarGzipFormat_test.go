package unpack

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/tar"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTarGzipFormatUnpack(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "unpack", "project-targz")
	workingDir := sourceBase
	sourceProjectTargz := "project.tar.gz"
	sourceFile := filepath.Join("site.yml")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		unpack      *TarGzipFormat
		project     *entity.Project
		workingDir  string
		err         error
		arrangeFunc func(*testing.T, *TarGzipFormat)
		assertFunc  func(*testing.T, *TarGzipFormat)
	}{
		{
			desc:   "Testing unpack project in tar.gz format",
			unpack: NewTarGzipFormat(fs, tar.NewTar(fs, logger.NewFakeLogger()), logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-targz",
				Format:    "targz",
				Reference: sourceProjectTargz,
				Storage:   "local",
			},
			workingDir:  workingDir,
			err:         errors.New(""),
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {},
			assertFunc: func(t *testing.T, unpack *TarGzipFormat) {
				_, err := unpack.fs.Stat(filepath.Join(workingDir, sourceFile))
				assert.Nil(t, err)
			},
		},
		{
			desc:        "Testing error unpacking project in tar.gz format when project is not provided",
			unpack:      NewTarGzipFormat(fs, NewMockTarExtractor(), logger.NewFakeLogger()),
			project:     nil,
			workingDir:  workingDir,
			err:         ErrProjectNotProvided,
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {},
			assertFunc:  func(t *testing.T, unpack *TarGzipFormat) {},
		},
		{
			desc:        "Testing error unpacking project in tar.gz format when working directory is not provided",
			unpack:      NewTarGzipFormat(fs, NewMockTarExtractor(), logger.NewFakeLogger()),
			project:     &entity.Project{},
			workingDir:  "",
			err:         ErrWorkingDirNotProvided,
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {},
			assertFunc:  func(t *testing.T, unpack *TarGzipFormat) {},
		},
		{
			desc:        "Testing error unpacking project in tar.gz format when filesystem is not provided",
			unpack:      NewTarGzipFormat(nil, NewMockTarExtractor(), logger.NewFakeLogger()),
			project:     &entity.Project{},
			workingDir:  workingDir,
			err:         ErrFilesystemNotProvided,
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {},
			assertFunc:  func(t *testing.T, unpack *TarGzipFormat) {},
		},
		{
			desc: "Testing error unpacking project in tar.gz format when project tar extractor is not provided",
			unpack: &TarGzipFormat{
				fs:     fs,
				logger: logger.NewFakeLogger(),
			},
			project:     &entity.Project{},
			workingDir:  workingDir,
			err:         ErrTarExtractorNotProvided,
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {},
			assertFunc:  func(t *testing.T, unpack *TarGzipFormat) {},
		},
		{
			desc:   "Testing error unpacking project in tar.gz format when project reference is not provided",
			unpack: NewTarGzipFormat(fs, NewMockTarExtractor(), logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-targz",
				Format:    "targz",
				Reference: "",
				Storage:   "local",
			},
			workingDir:  workingDir,
			err:         ErrProjectReferenceNotProvided,
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {},
			assertFunc:  func(t *testing.T, unpack *TarGzipFormat) {},
		},
		{
			desc:   "Testing error unpacking project in tar.gz format when there is an error extracting tar file",
			unpack: NewTarGzipFormat(fs, NewMockTarExtractor(), logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-targz",
				Format:    "targz",
				Reference: sourceProjectTargz,
				Storage:   "local",
			},
			workingDir: workingDir,
			err:        fmt.Errorf("%s: %w", ErrExtractingSourceCodeFile, errors.New("error extracting tar file")),
			arrangeFunc: func(t *testing.T, unpack *TarGzipFormat) {
				unpack.extractor.(*MockTarExtractor).On("Extract", mock.Anything, workingDir).Return(errors.New("error extracting tar file"))
			},
			assertFunc: func(t *testing.T, unpack *TarGzipFormat) {},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrangeFunc != nil {
				test.arrangeFunc(t, test.unpack)
			}

			err := test.unpack.Unpack(test.project, test.workingDir)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				test.assertFunc(t, test.unpack)
			}
		})
	}
}
