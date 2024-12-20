package unpack

import (
	"io"

	"github.com/stretchr/testify/mock"
)

// MockTarExtractor interface used to extract tar files
type MockTarExtractor struct {
	mock.Mock
}

// NewMockTarExtractor method creates a new TarGzipFormat struct
func NewMockTarExtractor() *MockTarExtractor {
	return &MockTarExtractor{}
}

// Extract method extracts a tar file
func (m *MockTarExtractor) Extract(reader io.Reader, dest string) error {
	args := m.Called(reader, dest)
	return args.Error(0)
}
