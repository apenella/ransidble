package archive

import "github.com/apenella/ransidble/internal/domain/ports/repository"

type ArchiveFactory struct {
	factory map[string]repository.Archiver
}

func NewArchiveFactory() *ArchiveFactory {
	return &ArchiveFactory{
		factory: make(map[string]repository.Archiver),
	}
}

func (f *ArchiveFactory) Register(name string, archiver repository.Archiver) {
	f.factory[name] = archiver
}

func (f *ArchiveFactory) Get(name string) repository.Archiver {
	return f.factory[name]
}
