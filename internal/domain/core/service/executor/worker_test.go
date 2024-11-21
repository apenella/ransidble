package executor

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/apenella/ransidble/internal/domain/core/entity"
	"github.com/apenella/ransidble/internal/domain/ports/repository"
	"github.com/apenella/ransidble/internal/infrastructure/executor"
	"github.com/apenella/ransidble/internal/infrastructure/logger"
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestWorkerGenerateID(t *testing.T) {
	id := genereteID()
	// ensure the id is a valid uuid

	t.Run("ID is not empty", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing generateID for a worker is not empty")
		if id == "" {
			assert.NotEmpty(t, id, "ID must not be empty")
		}
	})

	t.Run("ID has 36 characters", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing generateID for a worker has 36 characters")
		assert.Len(t, id, 36, "ID must be a valid uuid with 36 characters")
	})

	t.Run("ID is a valid uuid", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing generateID for a worker is a valid uuid")
		uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
		assert.True(t, uuidRegex.MatchString(id), "ID is not a valid uuid")
	})
}

func TestCreateWorkspace(t *testing.T) {
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
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				executor.NewAnsiblePlaybook(),
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
				// The arrange function is used to mock the workspace builder to return a workspace without errors

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				return nil
			},
			err: nil,
		},
		{
			desc: "Testing error creating a workspace for a task",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				executor.NewAnsiblePlaybook(),
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
				// The arrange function is used to mock the workspace builder to return an error when preparing the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(fmt.Errorf("Error preparing workspace"))

				return nil
			},
			err: fmt.Errorf("%s: %s", ErrPreparingWorkspace, "Error preparing workspace"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrange != nil {
				err := test.arrange(t, test.worker)
				if err != nil {
					t.Error(err)
				}
			}

			wsp, err := test.worker.createWorkspace(test.task)
			if err != nil {
				assert.Equal(t, test.err, err, "Error must be the expected")
			} else {
				assert.NotNil(t, wsp, "Workspace must not be nil")
				test.worker.workspaceBuilder.(*repository.MockBuilder).Workspace.AssertExpectations(t)
			}
		})
	}
}

func TestHandleAnsiblePlaybookTask(t *testing.T) {

	tests := []struct {
		desc       string
		worker     *Worker
		task       *entity.Task
		workingDir string
		err        error
		arrange    func(*testing.T, *Worker) error
	}{
		{
			desc: "Testing handle an ansible-playbook task",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			workingDir: "/tmp",
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the ansible playbook executor to return a success when running the ansible playbook

				if w.ansiblePlaybookExecutor == nil {
					return fmt.Errorf("Ansible playbook executor must not be nil")
				}

				_, ok := w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor)
				if !ok {
					return fmt.Errorf("Ansible playbook executor must have expectations")
				}

				w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor).On("Run", context.TODO(), "/tmp", &entity.AnsiblePlaybookParameters{}).Return(nil)

				return nil
			},
		},
		{
			desc: "Testing error handling an ansible-playbook task when ansible playbook executor is nil",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				nil,
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			workingDir: "/tmp",
			err:        fmt.Errorf(ErrAnsiblePlaybookExecutorDefined.Error()),
		},
		{
			desc: "Testing error handling an ansible-playbook when ansible playbook executor returns an error",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			workingDir: "/tmp",
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the ansible playbook executor to return an error when running the ansible playbook

				if w.ansiblePlaybookExecutor == nil {
					return fmt.Errorf("Ansible playbook executor must not be nil")
				}

				_, ok := w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor)
				if !ok {
					return fmt.Errorf("Ansible playbook executor must have expectations")
				}

				w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor).On("Run", context.TODO(), "/tmp", &entity.AnsiblePlaybookParameters{}).Return(fmt.Errorf("error running ansible playbook"))

				return nil
			},
			err: fmt.Errorf("error running ansible playbook"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrange != nil {
				err := test.arrange(t, test.worker)
				if err != nil {
					t.Error(err)
				}
			}

			err := test.worker.handleAnsiblePlaybookTask(context.TODO(), test.task, test.workingDir)
			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error(), "Error must be the expected")
			} else {
				test.worker.workspaceBuilder.(*repository.MockBuilder).Workspace.AssertExpectations(t)
			}
		})
	}
}

func TestHandleTask(t *testing.T) {

	tests := []struct {
		desc         string
		worker       *Worker
		task         *entity.Task
		expectedTask *entity.Task
		err          error
		arrange      func(*testing.T, *Worker) error
	}{
		{
			desc: "Testing error handling a task when there is an error creating a workspace",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the workspace builder to return an error when preparing the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(fmt.Errorf("error preparing workspace"))

				return nil
			},
			expectedTask: &entity.Task{
				ID:         "task-id",
				Status:     "FAILED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			err: fmt.Errorf("%s: %s", ErrPreparingWorkspace, "error preparing workspace"),
		},
		{
			desc: "Testing error handling a task when there is an error getting the working directory from the workspace",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the workspace builder to return an error when getting the working directory. It is also required to mock the cleanup method to avoid an error when the worker tries to cleanup the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("GetWorkingDir").Return("", fmt.Errorf("error getting working directory"))

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Cleanup").Return(nil)

				return nil
			},
			expectedTask: &entity.Task{
				ID:         "task-id",
				Status:     "FAILED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			err: fmt.Errorf("%s: %s", ErrGettingWorkingDir, "error getting working directory"),
		},
		{
			desc: "Testing error handling a task when the task command is not an ansible-playbook command",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "not-ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the workspace builder to return an error when getting the working directory. It is also required to mock the cleanup method to avoid an error when the worker tries to cleanup the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("GetWorkingDir").Return("/tmp", nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Cleanup").Return(nil)

				return nil
			},
			expectedTask: &entity.Task{
				ID:         "task-id",
				Status:     "FAILED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "not-ansible-playbook",
				ProjectID:  "project-id",
			},
			err: fmt.Errorf(ErrUnknownCommandType.Error()),
		},
		{
			desc: "Testing error handling a task when the task parameters are not an ansible-playbook parameters",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the workspace builder to return an error when getting the working directory. It is also required to mock the cleanup method to avoid an error when the worker tries to cleanup the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("GetWorkingDir").Return("/tmp", nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Cleanup").Return(nil)

				return nil
			},
			expectedTask: &entity.Task{
				ID:         "task-id",
				Status:     "FAILED",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			err: fmt.Errorf(ErrAnsiblePlaybookTaskInvalidParameters.Error()),
		},
		{
			desc: "Testing error handling a task when ansible playbook executor returns an error",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the workspace builder to return an error when getting the working directory. It is also required to mock the cleanup method to avoid an error when the worker tries to cleanup the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("GetWorkingDir").Return("/tmp", nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Cleanup").Return(nil)

				if w.ansiblePlaybookExecutor == nil {
					return fmt.Errorf("Ansible playbook executor must not be nil")
				}

				_, ok = w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor)
				if !ok {
					return fmt.Errorf("Ansible playbook executor must have expectations")
				}

				w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor).On("Run", context.TODO(), "/tmp", &entity.AnsiblePlaybookParameters{}).Return(fmt.Errorf("error running ansible playbook"))

				return nil
			},
			expectedTask: &entity.Task{
				ID:         "task-id",
				Status:     "FAILED",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			err: fmt.Errorf("%s: %s", ErrAnsiblePlaybookTaskFailed, "error running ansible playbook"),
		},
		{
			desc: "Testing handle an ansible-playbook task",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			task: &entity.Task{
				ID:         "task-id",
				Status:     "PENDING",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			arrange: func(t *testing.T, w *Worker) error {
				// The arrange function is used to mock the workspace builder to return an error when getting the working directory. It is also required to mock the cleanup method to avoid an error when the worker tries to cleanup the workspace

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("GetWorkingDir").Return("/tmp", nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Cleanup").Return(nil)

				if w.ansiblePlaybookExecutor == nil {
					return fmt.Errorf("Ansible playbook executor must not be nil")
				}

				_, ok = w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor)
				if !ok {
					return fmt.Errorf("Ansible playbook executor must have expectations")
				}

				w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor).On("Run", context.TODO(), "/tmp", &entity.AnsiblePlaybookParameters{}).Return(nil)

				return nil
			},
			expectedTask: &entity.Task{
				ID:         "task-id",
				Status:     "SUCCESS",
				Parameters: &entity.AnsiblePlaybookParameters{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			err: &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			if test.arrange != nil {
				err := test.arrange(t, test.worker)
				if err != nil {
					t.Error(err)
				}
			}

			err := test.worker.handleTask(context.TODO(), test.task)
			if err != nil {
				t.Log(err)
				assert.Equal(t, test.err.Error(), err.Error(), "Error must be the expected")
			} else {
				test.worker.workspaceBuilder.(*repository.MockBuilder).Workspace.AssertExpectations(t)
			}
			assert.Equal(t, test.expectedTask.Status, test.task.Status, "Task must be the expected")
		})
	}
}
func TestWorkerStop(t *testing.T) {

	tests := []struct {
		desc   string
		worker *Worker
	}{
		{
			desc: "Testing stopping a worker",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			// Start the worker
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()
			err := test.worker.Start(ctx)
			assert.NoError(t, err, "Starting worker should not return an error")

			// Stop the worker
			test.worker.Stop()

			// Ensure the worker is stopped with a deadline
			select {
			case <-test.worker.stopCh:
				t.Log("Worker stopped successfully")
			case <-time.After(2 * time.Second):
				t.Error("Worker did not stop as expected within the deadline")
			}
		})
	}
}
func TestWorkerStart(t *testing.T) {
	tests := []struct {
		desc        string
		worker      *Worker
		arrangeFunc func(*testing.T, *Worker) error
		actionFunc  func(*testing.T, *Worker) error
		verifyFunc  func(*testing.T, *Worker) error
		cancelFunc  func(*testing.T, *Worker) error
	}{
		{
			desc: "Testing executor worker starts successfully",
			worker: NewWorker(
				make(chan chan *entity.Task),
				&repository.MockBuilder{
					Workspace: &repository.MockWorkspace{},
				},
				NewMockAnsiblePlaybookExecutor(),
				logger.NewFakeLogger(),
			),
			arrangeFunc: func(t *testing.T, w *Worker) error {

				if w.workspaceBuilder == nil {
					return fmt.Errorf("Workspace builder must not be nil")
				}

				_, ok := w.workspaceBuilder.(*repository.MockBuilder)
				if !ok {
					return fmt.Errorf("Workspace builder must have expectations")
				}

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Prepare").Return(nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("GetWorkingDir").Return("/tmp", nil)

				w.workspaceBuilder.(*repository.MockBuilder).Workspace.On("Cleanup").Return(nil)

				if w.ansiblePlaybookExecutor == nil {
					return fmt.Errorf("Ansible playbook executor must not be nil")
				}

				_, ok = w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor)
				if !ok {
					return fmt.Errorf("Ansible playbook executor must have expectations")
				}

				w.ansiblePlaybookExecutor.(*MockAnsiblePlaybookExecutor).On("Run", context.TODO(), "/tmp", &entity.AnsiblePlaybookParameters{}).Return(nil)

				return nil
			},
			actionFunc: func(t *testing.T, w *Worker) error {

				if w.taskChan == nil {
					return fmt.Errorf("Task channel must not be nil")
				}

				if w.workerPool == nil {
					return fmt.Errorf("Worker pool must not be nil")
				}

				taskChan := <-w.workerPool

				task := &entity.Task{
					ID:         "task-id",
					Status:     "PENDING",
					Parameters: &entity.AnsiblePlaybookParameters{},
					Command:    "ansible-playbook",
					ProjectID:  "project-id",
				}

				taskChan <- task

				return nil
			},
			verifyFunc: func(t *testing.T, w *Worker) error {
				// Verify that the worker has started
				return nil
			},
			cancelFunc: func(t *testing.T, w *Worker) error {
				w.Stop()
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)

			// Setup the worker
			if test.arrangeFunc != nil {
				errSetup := test.arrangeFunc(t, test.worker)
				if errSetup != nil {
					t.Error(errSetup)
				}
			}

			// Start the worker
			err := test.worker.Start(context.TODO())
			assert.NoError(t, err, "Starting worker should not return an error")

			if test.actionFunc != nil {
				errAction := test.actionFunc(t, test.worker)
				if errAction != nil {
					t.Error(errAction)
				}
			}

			// Cancel the context or close stopCh to stop the worker
			if test.cancelFunc != nil {
				errCancel := test.cancelFunc(t, test.worker)
				if errCancel != nil {
					t.Error(errCancel)
				}
			}

			// Verify the worker behavior
			if test.verifyFunc != nil {
				errVerify := test.verifyFunc(t, test.worker)
				if errVerify != nil {
					t.Error(errVerify)
				}
			}
		})
	}
}
