package filesystem

import (
	"os"

	"github.com/spf13/afero"
)

// Filesystem represents the filesystem
type Filesystem struct {
	fs afero.Fs
}

// NewFilesystem creates a new filesystem
func NewFilesystem(fs afero.Fs) *Filesystem {
	return &Filesystem{
		fs: fs,
	}
}

// MkdirAll creates a directory
func (f *Filesystem) MkdirAll(path string, perm os.FileMode) error {
	return f.fs.MkdirAll(path, perm)
}

// Stat returns the file information
func (f *Filesystem) Stat(path string) (os.FileInfo, error) {
	return f.fs.Stat(path)
}

// RemoveAll removes a directory
func (f *Filesystem) RemoveAll(path string) error {
	return f.fs.RemoveAll(path)
}

// TempDir creates a temporary directory
func (f *Filesystem) TempDir(dir, prefix string) (string, error) {
	return afero.TempDir(f.fs, dir, prefix)
}

// DirExists checks if a directory exists
func (f *Filesystem) DirExists(path string) (bool, error) {
	return afero.DirExists(f.fs, path)
}

// Open opens a file
func (f *Filesystem) Open(name string) (afero.File, error) {
	return f.fs.Open(name)
}

// Create creates a file
func (f *Filesystem) Create(name string) (afero.File, error) {
	return f.fs.Create(name)
}
