package unpack

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestPlainFormatUnpack(t *testing.T) {

	sourceBase := filepath.Join("fixtures", "unpack")
	sourceProjectPlain := filepath.Join("project-plain")
	workingDir := filepath.Join(sourceBase, sourceProjectPlain)
	sourceFile := filepath.Join("site.yaml")

	fs := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test"),
		),
		afero.NewMemMapFs(),
	)

	tests := []struct {
		desc        string
		unpack      *PlainFormat
		project     *entity.Project
		workingDir  string
		err         error
		arrangeFunc func(*testing.T, *PlainFormat)
		assertFunc  func(*testing.T, *PlainFormat)
	}{
		{
			desc:   "Testing unpack project in plain format",
			unpack: NewPlainFormat(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-plain",
				Format:    "plain",
				Reference: sourceProjectPlain,
				Storage:   "local",
			},
			workingDir:  workingDir,
			err:         nil,
			arrangeFunc: func(t *testing.T, unpack *PlainFormat) {},
			assertFunc: func(t *testing.T, unpack *PlainFormat) {
				_, err := unpack.fs.Stat(filepath.Join(workingDir, sourceFile))
				assert.Nil(t, err)
			},
		},
		{
			desc:        "Testing error unpacking project in plain format when project is not provided",
			unpack:      NewPlainFormat(fs, logger.NewFakeLogger()),
			project:     nil,
			workingDir:  workingDir,
			err:         ErrProjectNotProvided,
			arrangeFunc: func(t *testing.T, unpack *PlainFormat) {},
			assertFunc:  func(t *testing.T, unpack *PlainFormat) {},
		},
		{
			desc:   "Testing error unpacking project in plain format when working directory is not provided",
			unpack: NewPlainFormat(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-plain",
				Format:    "plain",
				Reference: sourceProjectPlain,
				Storage:   "local",
			},
			workingDir:  "",
			err:         ErrWorkingDirNotProvided,
			arrangeFunc: func(t *testing.T, unpack *PlainFormat) {},
			assertFunc:  func(t *testing.T, unpack *PlainFormat) {},
		},
		{
			desc:   "Testing error unpacking project in plain format when filesystem is not provided",
			unpack: NewPlainFormat(nil, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-plain",
				Format:    "plain",
				Reference: sourceProjectPlain,
				Storage:   "local",
			},
			workingDir:  workingDir,
			err:         ErrFilesystemNotProvided,
			arrangeFunc: func(t *testing.T, unpack *PlainFormat) {},
			assertFunc:  func(t *testing.T, unpack *PlainFormat) {},
		},
		{
			desc:   "Testing error unpacking project in plain format when working directory does not exists",
			unpack: NewPlainFormat(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-plain",
				Format:    "plain",
				Reference: sourceProjectPlain,
				Storage:   "local",
			},
			workingDir:  "not-exists",
			err:         fmt.Errorf("%s: %w", ErrWorkingDirNotExists, errors.New("stat ../../../test/not-exists: no such file or directory")),
			arrangeFunc: func(t *testing.T, unpack *PlainFormat) {},
			assertFunc:  func(t *testing.T, unpack *PlainFormat) {},
		},
		{
			desc:   "Testing error unpacking project in plain format when working directory is a regular file",
			unpack: NewPlainFormat(fs, logger.NewFakeLogger()),
			project: &entity.Project{
				Name:      "project-plain",
				Format:    "plain",
				Reference: sourceProjectPlain,
				Storage:   "local",
			},
			workingDir: filepath.Join(sourceBase, "file"),
			err:        ErrWorkingDirIsNotDirectory,
			arrangeFunc: func(t *testing.T, unpack *PlainFormat) {
				unpack.fs.OpenFile(filepath.Join(sourceBase, "file"), os.O_CREATE, os.ModePerm)
			},
			assertFunc: func(t *testing.T, unpack *PlainFormat) {},
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
