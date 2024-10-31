package executor

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestWorkerGenerateID(t *testing.T) {
	t.Parallel()

	id := genereteID()
	// ensure the id is a valid uuid

	t.Run("ID is not empty", func(t *testing.T) {
		t.Log("Testing generateID for a worker is not empty")
		if id == "" {
			assert.NotEmpty(t, id, "ID must not be empty")
		}
	})

	t.Run("ID has 36 characters", func(t *testing.T) {
		t.Log("Testing generateID for a worker has 36 characters")
		assert.Len(t, id, 36, "ID must be a valid uuid with 36 characters")
	})

	t.Run("ID is a valid uuid", func(t *testing.T) {
		t.Log("Testing generateID for a worker is a valid uuid")
		uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
		assert.True(t, uuidRegex.MatchString(id), "ID is not a valid uuid")
	})
}

func TestCreateWorkspace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc    string
		worker  *Worker
		task    *entity.Task
		err     error
		arrange func(*testing.T, *Worker) error
	}{
		{
			desc: "Testing create a workspace for a task",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&workspace.MockBuilder{
					Workspace: &workspace.MockWorkspace{},
				},
				nil,
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*workspace.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*workspace.MockBuilder).Workspace.On("Prepare").Return(nil)

				return nil
			},
			err: nil,
		},
		{
			desc: "Testing error creating a workspace for a task",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&workspace.MockBuilder{
					Workspace: &workspace.MockWorkspace{},
				},
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*workspace.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*workspace.MockBuilder).Workspace.On("Prepare").Return(fmt.Errorf("Error preparing workspace"))

				return nil
			},
			err: fmt.Errorf("%s: %s", ErrPreparingWorkspace, "Error preparing workspace"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.arrange != nil {
				err := test.arrange(t, test.worker)
				if err != nil {
					t.Error(err)
				}
			}

			workspace, err := test.worker.createWorkspace(test.task)
			if err != nil {
				assert.Equal(t, test.err, err, "Error must be the expected")
			} else {
				assert.NotNil(t, workspace, "Workspace must not be nil")
			}
		})
	}
}
