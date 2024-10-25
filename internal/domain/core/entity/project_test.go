package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectValidate(t *testing.T) {
	type fields struct {
		Format    string
		Name      string
		Reference string
		Type      string
	}
	tests := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a project",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "reference",
				Type:      "local",
			},
			wantErr: false,
		},
		{
			desc: "Validating a project with empty format",
			fields: fields{
				Format:    "",
				Name:      "project",
				Reference: "reference",
				Type:      "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project with empty name",
			fields: fields{
				Format:    "plain",
				Name:      "",
				Reference: "reference",
				Type:      "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project with empty reference",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "",
				Type:      "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project with empty type",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "reference",
				Type:      "",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project with invalid type",
			fields: fields{
				Format:    "plain",
				Name:      "project",
				Reference: "reference",
				Type:      "invalid-type",
			},
			wantErr: true,
		},
		{
			desc: "Validating a project with invalid format",
			fields: fields{
				Format:    "invalid-format",
				Name:      "project",
				Reference: "reference",
				Type:      "local",
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
				Type:      test.fields.Type,
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
