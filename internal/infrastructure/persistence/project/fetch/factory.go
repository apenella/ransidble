package fetch

import "github.com/apenella/ransidble/internal/domain/ports/repository"

// Factory represents the factory for fetching source code components
type Factory struct {
	factory map[string]repository.SourceCodeFetcher
}

// NewFactory creates a new Factory for fetching source code
func NewFactory() *Factory {
	return &Factory{
		factory: make(map[string]repository.SourceCodeFetcher),
	}
}

// Register registers a new source code fetcher
func (f *Factory) Register(name string, fetcher repository.SourceCodeFetcher) {
	f.factory[name] = fetcher
}

// Get gets an source code fetcher by name
func (f *Factory) Get(name string) repository.SourceCodeFetcher {
	return f.factory[name]
}
