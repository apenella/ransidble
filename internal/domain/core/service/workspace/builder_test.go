package workspace

import (
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/apenella/ransidble/internal/infrastructure/persistence/project/fetch"
	"github.com/apenella/ransidble/internal/infrastructure/unpack"
	"github.com/stretchr/testify/assert"
)

func TestBuildWorkspace(t *testing.T) {
	t.Parallel()
	t.Log("Testing the BuildWorkspace function")

	// fs := afero.NewMemMapFs()
	fs := repository.NewMockFilesystemer()
	fetchFactory := fetch.NewFactory()
	unpackFactory := unpack.NewFactory()
	repository := repository.NewMockProjectRepository()
	logger := logger.NewFakeLogger()

	tasks := &entity.Task{
		ID:         "task-id",
		Status:     "PENDING",
		Parameters: &entity.AnsiblePlaybookParameters{},
		Command:    "ansible-playbook",
		ProjectID:  "project-id",
	}

	expected := &Workspace{
		fetchFactory:  fetchFactory,
		fs:            fs,
		logger:        logger,
		repository:    repository,
		task:          tasks,
		unpackFactory: unpackFactory,
	}

	builder := NewBuilder(
		fs,
		fetchFactory,
		unpackFactory,
		repository,
		logger,
	)

	workspace := builder.WithTask(tasks).Build()

	assert.Equal(t, expected, workspace)
}
