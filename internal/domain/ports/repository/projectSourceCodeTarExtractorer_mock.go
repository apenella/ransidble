package repository

import (
	"io"

	"github.com/stretchr/testify/mock"
)

// MockProjectSourceCodeTarExtractorer interface used to extract tar files
type MockProjectSourceCodeTarExtractorer struct {
	mock.Mock
}

// Ensure MockProjectSourceCodeTarExtractorer implements the SourceCodeTarExtractorer interface
var _ SourceCodeTarExtractorer = (*MockProjectSourceCodeTarExtractorer)(nil)

// NewMockProjectSourceCodeTarExtractorer method creates a new TarGzipFormat struct
func NewMockProjectSourceCodeTarExtractorer() *MockProjectSourceCodeTarExtractorer {
	return &MockProjectSourceCodeTarExtractorer{}
}

// Extract method extracts a tar file
func (m *MockProjectSourceCodeTarExtractorer) Extract(reader io.Reader, dest string) error {
	args := m.Called(reader, dest)
	return args.Error(0)
}
