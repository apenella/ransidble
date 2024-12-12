package entity

import (
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	// ACCEPTED status when the task is accepted to be executed
	ACCEPTED = "ACCEPTED"
	// FAILED status when the task is failed
	FAILED = "FAILED"
	// PENDING status when the task is pending. This status is used when the task is not yet accepted to be executed
	PENDING = "PENDING"
	// RUNNING status when the task starts running
	RUNNING = "RUNNING"
	// SUCCESS status when the task is successfully executed
	SUCCESS = "SUCCESS"

	// AnsiblePlaybookCommand identifies the task as an Ansible playbook task
	AnsiblePlaybookCommand = "ansible-playbook"
)

// Task entity represents a task to be executed
type Task struct {
	// Command represents the command type to be executed. This field is required and must be one of the following values: ansible-playbook
	Command string `json:"command" validate:"required,oneof=ansible-playbook"`
	// CompletedAt represents the time when the task is completed
	CompletedAt string `json:"completed_at"`
	// CreatedAt represents the time when the task is created
	CreatedAt string `json:"created_at"`
	// ErrorMessage represents the error message when the task is failed
	ErrorMessage string `json:"error_message,omitempty"`
	// ExecutedAt represents the time when the task is executed
	ExecutedAt string `json:"executed_at"`
	// ID represents the task ID. This field is required
	ID string `json:"id" validate:"required"`
	// Parameters represents the task parameters. This field is required
	Parameters interface{} `json:"parameters" validate:"required"`
	// ProjectID represents the project ID. This field is required when the command is ansible-playbook
	ProjectID string `json:"project_id" validate:"required_if=Command ansible-playbook"`
	// Status represents the task status. This field is required and must be one of the following values: ACCEPTED, FAILED, PENDING, RUNNING, SUCCESS
	Status string `json:"status" validate:"required,oneof=ACCEPTED FAILED PENDING RUNNING SUCCESS"`

	statusMutex sync.Mutex
}

// NewTask creates a new task
func NewTask(id string, projectID string, command string, parameters interface{}) *Task {
	return &Task{
		Command:    command,
		CreatedAt:  time.Now().Format(time.RFC3339),
		ID:         id,
		Parameters: parameters,
		ProjectID:  projectID,
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

// Validate validates the task entity
func (t *Task) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
