package unpack

import "github.com/apenella/ransidble/internal/domain/ports/repository"

// Factory represents the factory for unpacking source code components
type Factory struct {
	factory map[string]repository.SourceCodeUnpacker
}

// NewFactory creates a new Factory for unpacking source code
func NewFactory() *Factory {
	return &Factory{
		factory: make(map[string]repository.SourceCodeUnpacker),
	}
}

// Register registers a new unpacker
func (f *Factory) Register(name string, unpacker repository.SourceCodeUnpacker) {
	f.factory[name] = unpacker
}

// Get gets an unpacker by name
func (f *Factory) Get(name string) repository.SourceCodeUnpacker {
	return f.factory[name]
}
