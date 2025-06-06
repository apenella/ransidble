package repository

import (
	"io"

	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// ProjectRepository represents a repository to manage projects
type ProjectRepository interface {
	Find(id string) (*entity.Project, error)
	FindAll() ([]*entity.Project, error)
	// Remove(id string) error
	SafeStore(id string, project *entity.Project) error
	// Store(id string, project *entity.Project) error
	// Update(id string, project *entity.Project) error
}

// SourceCodeFetcher represents the component to fetch a project from a repository
type SourceCodeFetcher interface {
	Fetch(project *entity.Project, destination string) error
}

// SourceCodeFetchFactory represents the component to create a SourceCodeFetcher
type SourceCodeFetchFactory interface {
	Get(projectType string) SourceCodeFetcher
}

// SourceCodeStorer represents the component to save a project in a storage
type SourceCodeStorer interface {
	Store(project *entity.Project, file io.Reader) error
}

// SourceCodeStorageFactory represents the component to create a SourceCodeStorer
type SourceCodeStorageFactory interface {
	Get(projectType string) SourceCodeStorer
}

// Unpacker represents the component to archive and unarchive projects before executing tasks
type Unpacker interface {
	Unpack(project *entity.Project, workingDir string) error
}

// SourceCodeUnpacker represents the component to unpack a project
type SourceCodeUnpacker interface {
	Unpack(project *entity.Project, destination string) error
}

// SourceCodeUnpackFactory represents the component to create a SourceCodeUnpacker
type SourceCodeUnpackFactory interface {
	Get(projectType string) SourceCodeUnpacker
}
