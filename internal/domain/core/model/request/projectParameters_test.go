package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectParametersValidate(t *testing.T) {
	type fields struct {
		Format string
		// Name      string
		// Reference string
		Storage string
	}
	test := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a ProjectParameters",
			fields: fields{
				Format: "targz",
				// Name:   "project",
				// Reference: "reference",
				Storage: "local",
			},
			wantErr: false,
		},
		{
			desc: "Validating a ProjectParameters with empty format",
			fields: fields{
				// Name: "project",
				// Reference: "reference",
				Storage: "local",
			},
			wantErr: true,
		},
		// {
		// 	desc: "Validating a ProjectParameters with empty name",
		// 	fields: fields{
		// 		Format: "targz",
		// 		// Reference: "reference",
		// 		Storage: "local",
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	desc: "Validating a ProjectParameters with empty reference",
		// 	fields: fields{
		// 		Format:  "targz",
		// 		Name:    "project",
		// 		Storage: "local",
		// 	},
		// 	wantErr: true,
		// },
		{
			desc: "Validating a ProjectParameters with empty storage",
			fields: fields{
				Format: "targz",
				// Name:   "project",
				// Reference: "reference",
			},
			wantErr: true,
		},
		{
			desc: "Validating a ProjectParameters with invalid format",
			fields: fields{
				Format: "invalid",
				// Name:   "project",
				// Reference: "reference",
				Storage: "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a ProjectParameters with invalid storage",
			fields: fields{
				Format: "targz",
				// Name:   "project",
				// Reference: "reference",
				Storage: "invalid",
			},
			wantErr: true,
		},
	}
	for _, test := range test {
		t.Run(test.desc, func(t *testing.T) {
			p := &ProjectParameters{
				Format: test.fields.Format,
				// Name:   test.fields.Name,
				// Reference: test.fields.Reference,
				Storage: test.fields.Storage,
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
