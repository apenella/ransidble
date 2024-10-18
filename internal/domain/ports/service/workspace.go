package service

import (
	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
)

// WorkspaceBuilder interface to build a workspace
type WorkspaceBuilder interface {
	WithTask(task *entity.Task) *workspace.Builder
	Build() *workspace.Workspace
}
