package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			desc: "Validating a task",
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
			desc: "Validating a task with empty id",
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
			desc: "Validating a task with empty status",
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
			desc: "Validating a task with empty command",
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
			desc: "Validating a task with empty parameters",
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
			desc: "Validating a task with empty project id",
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
			desc: "Validating a task with invalid status",
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
			desc: "Validating a task with invalid command",
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