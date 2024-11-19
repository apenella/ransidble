package executor

import (
	"context"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/core/service/workspace"
	"github.com/apenella/ransidble/internal/infrastructure/executor"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestDispatchTaskExecution(t *testing.T) {
	// This test ensures that the dispatcher is able to execute a task

	t.Run("Testing the execution dispatcher flow", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing the execution dispatcher flow")

		// arrange workspace mocks for testing the dispatcher
		mockWorkspace := &workspace.MockWorkspace{}
		mockWorkspace.On("Prepare").Return(nil)
		mockWorkspace.On("GetWorkingDir").Return("/tmp", nil)
		mockWorkspace.On("Cleanup").Return(nil)
		// arrange ansible playbook executor mocks for testing the dispatcher
		ansiblePlaybookExecutor := executor.NewMockAnsiblePlaybook()
		ansiblePlaybookExecutor.On("Run", context.TODO(), "/tmp", &entity.AnsiblePlaybookParameters{}).Return(nil)

		workspaceBuilder := &workspace.MockBuilder{
			Workspace: mockWorkspace,
		}

		dispatch := NewDispatch(
			1,
			workspaceBuilder,
			ansiblePlaybookExecutor,
			logger.NewFakeLogger(),
		)

		// Initialize the dispatcher
		err := dispatch.Start(context.TODO())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		task := &entity.Task{
			ID:         "task-id",
			Status:     "PENDING",
			Parameters: &entity.AnsiblePlaybookParameters{},
			Command:    "ansible-playbook",
			ProjectID:  "project-id",
		}

		err = dispatch.Execute(task)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Wait for the task to be processed
		time.Sleep(1 * time.Second)

		// Assertions
		assert.Equal(
			t,
			entity.SUCCESS,
			task.Status,
		)

		mockWorkspace.AssertExpectations(t)
		ansiblePlaybookExecutor.AssertExpectations(t)

		// teardown the dispatcher by stopping it
		dispatch.Stop()
	})
}

func TestDispatchTaskExecutionContextCancelled(t *testing.T) {

	t.Run("Testing the execution dispatcher flow when context is cancelled", func(t *testing.T) {
		// This test ensures that dispatcher is stopped when context is cancelled

		t.Parallel()
		t.Log("Testing the execution dispatcher flow when context is cancelled")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockWorkspace := &workspace.MockWorkspace{}
		ansiblePlaybookExecutor := executor.NewMockAnsiblePlaybook()
		workspaceBuilder := &workspace.MockBuilder{
			Workspace: mockWorkspace,
		}

		dispatch := NewDispatch(
			1,
			workspaceBuilder,
			ansiblePlaybookExecutor,
			logger.NewFakeLogger(),
		)

		err := dispatch.Start(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// cancel function forces the context cancellation
		cancel()

		// Wait for the dispatcher to be stopped
		time.Sleep(1 * time.Second)

		// When the dispatcher is stopped due to context cancellation, a panic is expected when trying to execute a task. That ensures that the dispatcher is stopped after the context is cancelled
		assert.Panics(t, func() {
			dispatch.Execute(&entity.Task{})
		})

	})

}
