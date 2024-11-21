package workspace

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/domain/ports/service"
	"github.com/spf13/afero"
)

// Builder represents the location where the project is stored before being executed
type Builder struct {
	options []FuncOptions
}

// NewBuilder creates a new Builder
func NewBuilder(
	fs afero.Fs,
	fetchFactory repository.SourceCodeFetchFactory,
	unpackFactory repository.SourceCodeUnpackFactory,
	repository repository.ProjectRepository,
	logger repository.Logger,
) *Builder {
	// The builder receives all the common dependencies required to create a workspace. It simplifies the creation of a workspace by providing only the custom dependencies required to create a workspace.
	return &Builder{
		options: []FuncOptions{
			func(w *Workspace) {
				w.unpackFactory = unpackFactory
			},
			func(w *Workspace) {
				w.fetchFactory = fetchFactory
			},
			func(w *Workspace) {
				w.repository = repository
			},
			func(w *Workspace) {
				w.logger = logger
			},
			func(w *Workspace) {
				w.fs = fs
			},
		},
	}
}

// WithTask sets the project
func (w *Builder) WithTask(task *entity.Task) service.WorkspaceBuilder {
	w.options = append(w.options, func(w *Workspace) {
		w.task = task
	})
	return w
}

// Build creates a new workspace
func (w *Builder) Build() service.Workspacer {
	workspace := &Workspace{}
	for _, option := range w.options {
		option(workspace)
	}
	return workspace
}
