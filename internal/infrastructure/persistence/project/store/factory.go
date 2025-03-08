package store

import "github.com/apenella/ransidble/internal/domain/ports/repository"

// Factory represents the factory for storing source code components
type Factory struct {
	factory map[string]repository.SourceCodeStorer
}

// NewFactory creates a new Factory for storing source code
func NewFactory() *Factory {
	return &Factory{
		factory: make(map[string]repository.SourceCodeStorer),
	}
}

// Register registers a new source code storer
func (f *Factory) Register(name string, storer repository.SourceCodeStorer) {
	f.factory[name] = storer
}

// Get gets an source code storer by name
func (f *Factory) Get(name string) repository.SourceCodeStorer {
	return f.factory[name]
}
