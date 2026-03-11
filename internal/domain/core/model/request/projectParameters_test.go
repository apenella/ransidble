package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectParametersValidate(t *testing.T) {
	type fields struct {
		Format  string
		Storage string
		Version string
	}
	test := []struct {
		desc    string
		fields  fields
		wantErr bool
	}{
		{
			desc: "Validating a ProjectParameters",
			fields: fields{
				Format:  "targz",
				Storage: "local",
			},
			wantErr: false,
		},
		{
			desc: "Validating a ProjectParameters with empty format",
			fields: fields{
				Storage: "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a ProjectParameters with empty storage",
			fields: fields{
				Format: "targz",
			},
			wantErr: true,
		},
		{
			desc: "Validating a ProjectParameters with invalid format",
			fields: fields{
				Format:  "invalid",
				Storage: "local",
			},
			wantErr: true,
		},
		{
			desc: "Validating a ProjectParameters with invalid storage",
			fields: fields{
				Format:  "targz",
				Storage: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "Validating a ProjectParameters with a version",
			fields: fields{
				Format:  "targz",
				Storage: "local",
				Version: "v1.0.0",
			},
			wantErr: false,
		},
	}
	for _, test := range test {
		t.Run(test.desc, func(t *testing.T) {
			p := &ProjectParameters{
				Format:  test.fields.Format,
				Storage: test.fields.Storage,
				Version: test.fields.Version,
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
