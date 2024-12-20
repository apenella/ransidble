package filesystem

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMkdirAll(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	path := "new-dir"
	perm := os.ModePerm

	err := fs.MkdirAll(path, perm)
	assert.Nil(t, err)

	_, err = fs.fs.Stat(path)
	assert.Nil(t, err)
}

func TestStat(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	path := "file"
	_, err := fs.Stat(path)
	assert.Nil(t, err)
}

func TestRemoveAll(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	path := "deleteme"

	fs.fs.Mkdir(path, os.ModePerm)

	_, err := fs.fs.Stat(path)
	assert.Nil(t, err)

	err = fs.RemoveAll(path)
	assert.Nil(t, err)

	_, err = fs.fs.Stat(path)
	assert.EqualError(t, err, "stat ../../../test/fixtures/filesystem/deleteme: no such file or directory")
}

func TestTempDir(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	prefix := "temp-dir"

	path, err := fs.TempDir("", prefix)
	assert.Nil(t, err)

	_, err = fs.fs.Stat(path)
	assert.Nil(t, err)
}

func TestDirExists(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	path := "dir"

	exists, err := fs.DirExists(path)
	assert.Nil(t, err)
	assert.True(t, exists)
}

func TestOpen(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	path := "file"

	file, err := fs.Open(path)
	assert.Nil(t, err)
	assert.NotNil(t, file)
}

func TestCreate(t *testing.T) {

	fsBase := afero.NewCopyOnWriteFs(
		afero.NewReadOnlyFs(
			afero.NewBasePathFs(afero.NewOsFs(), "../../../test/fixtures/filesystem"),
		),
		afero.NewMemMapFs(),
	)

	fs := NewFilesystem(fsBase)
	path := "new-file"

	file, err := fs.Create(path)
	assert.Nil(t, err)
	assert.NotNil(t, file)
}
