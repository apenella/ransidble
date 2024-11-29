package repository

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// MockFilesystemer represents a mock filesystemer
type MockFilesystemer struct {
	mock.Mock
}

// NewMockFilesystemer creates a new mock filesystemer
func NewMockFilesystemer() *MockFilesystemer {
	return &MockFilesystemer{}
}

// MkdirAll creates a directory
func (m *MockFilesystemer) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

// RemoveAll removes a directory
func (m *MockFilesystemer) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

// Stat returns the file information
func (m *MockFilesystemer) Stat(path string) (os.FileInfo, error) {
	var fi os.FileInfo
	var err error

	args := m.Called(path)

	if args.Get(0) != nil {
		fi = args.Get(0).(os.FileInfo)
	}

	if args.Error(1) != nil {
		err = args.Error(1)
	}

	return fi, err
}

// TempDir creates a temporary directory
func (m *MockFilesystemer) TempDir(dir, prefix string) (string, error) {
	args := m.Called(dir, prefix)
	return args.String(0), args.Error(1)
}
