package entity

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Task statuses

	// Task status when the task is accepted to be executed
	ACCEPTED = "ACCEPTED"
	// Task status when the task is failed
	FAILED = "FAILED"
	// Task status when the task is pending. This status is used when the task is not yet accepted to be executed
	PENDING = "PENDING"
	// Task status when the task starts running
	RUNNING = "RUNNING"
	// Task status when the task is successfully executed
	SUCCESS = "SUCCESS"

	// Kinds of tasks
	ANSIBLE_PLAYBOOK = "ansible-playbook"
)

var (
	// Errors

	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = fmt.Errorf("task not found")
	// ErrTaskAlreadyExists is returned when you try to store a task that already exists
	ErrTaskAlreadyExists = fmt.Errorf("task already exists")
	// ErrNotInitializedStorage is returned when the storage is not initialized
	ErrNotInitializedStorage = fmt.Errorf("storage not initialized")
)

// Task represents a task to be executed
type Task struct {
	Command      string      `json:"command" validate:"required"`
	CompletedAt  string      `json:"completed_at"`
	CreatedAt    string      `json:"created_at"`
	ErrorMessage string      `json:"error_message,omitempty"`
	ExecutedAt   string      `json:"executed_at"`
	ID           string      `json:"id" validate:"required"`
	Parameters   interface{} `json:"parameters" validate:"required"`
	Project      *Project    `json:"project"`
	Status       string      `json:"status" validate:"required"`

	statusMutex sync.Mutex
}

// NewTask creates a new task
func NewTask(id string, command string, parameters interface{}) *Task {
	return &Task{
		CreatedAt:  time.Now().Format(time.RFC3339),
		ID:         id,
		Command:    command,
		Parameters: parameters,
		Status:     PENDING,
	}
}

// Accepted sets the task status to ACCEPTED
func (t *Task) Accepted() {
	t.statusMutex.Lock()
	defer t.statusMutex.Unlock()
	t.Status = ACCEPTED
	t.CreatedAt = time.Now().Format(time.RFC3339)
}

// Failed sets the task status to FAILED
func (t *Task) Failed(errorMsg string) {
	t.statusMutex.Lock()
	defer t.statusMutex.Unlock()
	t.Status = FAILED
	t.ErrorMessage = errorMsg
	t.CompletedAt = time.Now().Format(time.RFC3339)
}

// Success sets the task status to SUCCESS
func (t *Task) Success() {
	t.statusMutex.Lock()
	defer t.statusMutex.Unlock()
	t.Status = SUCCESS
	t.CompletedAt = time.Now().Format(time.RFC3339)
}

// Running sets the task status to RUNNING
func (t *Task) Running() {
	t.statusMutex.Lock()
	defer t.statusMutex.Unlock()
	t.Status = RUNNING
	t.ExecutedAt = time.Now().Format(time.RFC3339)
}
