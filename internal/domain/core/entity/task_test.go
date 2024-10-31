package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	task := NewTask("id", "project-id", "command", map[string]interface{}{})

	assert.Equal(t, "id", task.ID)
	assert.Equal(t, "project-id", task.ProjectID)
	assert.Equal(t, "command", task.Command)
	assert.Equal(t, map[string]interface{}{}, task.Parameters)
	assert.Equal(t, PENDING, task.Status)
}

func TestTaskValidate(t *testing.T) {
	type fields struct {
		Command    string
		ID         string
		Parameters interface{}
		ProjectID  string
		Status     string
	}
	tests := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a task entity ",
			fields: fields{
				Command:    "ansible-playbook",
				ID:         "task-id",
				Parameters: map[string]interface{}{},
				ProjectID:  "project-id",
				Status:     "ACCEPTED",
			},
			wantErr: false,
		},
		{
			desc: "Validating a task entity with empty id",
			fields: fields{
				ID:         "",
				Status:     "ACCEPTED",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			wantErr: true,
		},
		{
			desc: "Validating a task entity with empty status",
			fields: fields{
				ID:         "task-id",
				Status:     "",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			wantErr: true,
		},
		{
			desc: "Validating a task entity with empty command",
			fields: fields{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: map[string]interface{}{},
				Command:    "",
				ProjectID:  "project-id",
			},
			wantErr: true,
		},
		{
			desc: "Validating a task entity with empty parameters",
			fields: fields{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: nil,
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			wantErr: true,
		},
		{
			desc: "Validating a task entity with empty project id",
			fields: fields{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "",
			},
			wantErr: true,
		},
		{
			desc: "Validating a task entity with invalid status",
			fields: fields{
				ID:         "task-id",
				Status:     "invalid-status",
				Parameters: map[string]interface{}{},
				Command:    "ansible-playbook",
				ProjectID:  "project-id",
			},
			wantErr: true,
		},
		{
			desc: "Validating a task entity with invalid command",
			fields: fields{
				ID:         "task-id",
				Status:     "ACCEPTED",
				Parameters: map[string]interface{}{},
				Command:    "invalid-command",
				ProjectID:  "project-id",
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			task := &Task{
				Command:    test.fields.Command,
				ID:         test.fields.ID,
				Parameters: test.fields.Parameters,
				ProjectID:  test.fields.ProjectID,
				Status:     test.fields.Status,
			}

			err := task.Validate()
			if err != nil {
				assert.Equal(t, test.wantErr, true, err.Error())
			}

			if test.wantErr {
				assert.NotNil(t, err)
			}

		})
	}
}

func TestAccepted(t *testing.T) {
	t.Log("Testing task entity accepted method")

	task := NewTask("id", "project-id", "command", map[string]interface{}{})
	task.Accepted()

	assert.Equal(t, ACCEPTED, task.Status)
}

func TestRunning(t *testing.T) {
	t.Log("Testing task entity running method")

	task := NewTask("id", "project-id", "command", map[string]interface{}{})
	task.Running()

	assert.Equal(t, RUNNING, task.Status)
}

func TestFailed(t *testing.T) {
	t.Log("Testing task entity failed method")

	task := NewTask("id", "project-id", "command", map[string]interface{}{})
	task.Failed("error message")

	assert.Equal(t, FAILED, task.Status)
	assert.Equal(t, "error message", task.ErrorMessage)
}

func TestSuccess(t *testing.T) {
	t.Log("Testing task entity success method")

	task := NewTask("id", "project-id", "command", map[string]interface{}{})
	task.Success()

	assert.Equal(t, SUCCESS, task.Status)
}
