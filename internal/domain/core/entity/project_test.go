package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProject(t *testing.T) {
	project := NewProject("project", "reference", "plain", "local")

	assert.Equal(t, "plain", project.Format)
	assert.Equal(t, "project", project.Name)
	assert.Equal(t, "reference", project.Reference)
	assert.Equal(t, "local", project.Storage)
}

func TestProjectValidate(t *testing.T) {
	type fields struct {
		Format    string
		Name      string
		Reference string
		Storage   string
	}
	tests := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a project entity",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "reference",
				Storage:   "local",
			},
			wantErr: false,
		},
		{
			desc: "Validating a project entity with empty format",
			fields: fields{
				Format:    "",
				Name:      "project",
				Reference: "reference",
				Storage:   "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project entity with empty name",
			fields: fields{
				Format:    "plain",
				Name:      "",
				Reference: "reference",
				Storage:   "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project entity with empty reference",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "",
				Storage:   "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project entity with empty type",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "reference",
				Storage:   "",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project entity with invalid type",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "reference",
				Storage:   "invalid-type",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project entity with invalid format",
			fields: fields{
				Format:    "invalid-format",
				Name:      "project",
				Reference: "reference",
				Storage:   "local",
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			p := &Project{
				Format:    test.fields.Format,
				Name:      test.fields.Name,
				Reference: test.fields.Reference,
				Storage:   test.fields.Storage,
			}
			err := p.Validate()

			if err != nil {
				assert.Equal(t, test.wantErr, true, err.Error())
			}

			if test.wantErr {
				assert.NotNil(t, err)
			}

		})
	}
}
