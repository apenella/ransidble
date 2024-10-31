package service

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
)

// Workspacer interface to manage a workspace
type Workspacer interface {
	Prepare() error
	Cleanup() error
	GetWorkingDir() (string, error)
}

// WorkspaceBuilder interface to build a workspace
type WorkspaceBuilder interface {
	WithTask(task *entity.Task) WorkspaceBuilder
	Build() Workspacer
}
