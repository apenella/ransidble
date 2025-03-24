package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProject(t *testing.T) {
	t.Log("Testing project entity creation")
	t.Parallel()

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
			t.Parallel()

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

func TestProjectSourceCodeExtension(t *testing.T) {
	type fields struct {
		Format    string
		Name      string
		Reference string
		Storage   string
	}
	tests := []struct {
		desc     string
		fields   fields
		expected string
		err      error
	}{
		{
			desc: "Testing get source code extension when format is targz",
			fields: fields{
				Format:    ProjectFormatTarGz,
				Name:      "project",
				Reference: "reference",
				Storage:   "local",
			},
			expected: ExtensionTarGz,
			err:      nil,
		},
		{
			desc: "Testing get source code extension when format is plain",
			fields: fields{
				Format:    ProjectFormatPlain,
				Name:      "project",
				Reference: "reference",
				Storage:   "local",
			},
			expected: "",
			err:      nil,
		},
		{
			desc: "Testing get source code extension with invalid format",
			fields: fields{
				Format:    "invalid-format",
				Name:      "project",
				Reference: "reference",
				Storage:   "local",
			},
			expected: "",
			err:      fmt.Errorf("format invalid-format not supported"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			p := &Project{
				Format:    test.fields.Format,
				Name:      test.fields.Name,
				Reference: test.fields.Reference,
				Storage:   test.fields.Storage,
			}
			got, err := p.ProjectSourceCodeExtension()

			if err != nil && test.err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.Nil(t, err, "got an unexpected error")
				assert.Nil(t, test.err, "an expected error not received")
				assert.Equal(t, test.expected, got)
			}
		})
	}
}

func TestValidateProjectFormat(t *testing.T) {
	tests := []struct {
		desc   string
		format string
		err    error
	}{
		{
			desc:   "Testing validate project format with plain format",
			format: "plain",
			err:    nil,
		},
		{
			desc:   "Testing validate project format with targz format",
			format: "targz",
			err:    nil,
		},
		{
			desc:   "Testing validate project format with invalid format",
			format: "invalid-format",
			err:    fmt.Errorf("invalid format: invalid-format"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			err := ValidateProjectFormat(test.format)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, err, "got an unexpected error")
				assert.Nil(t, test.err, "an expected error not received")
			}
		})
	}
}

func TestValidateProjectStorage(t *testing.T) {
	tests := []struct {
		desc    string
		storage string
		err     error
	}{
		{
			desc:    "Testing validate project storage with local storage",
			storage: "local",
			err:     nil,
		},
		{
			desc:    "Testing validate project storage with invalid storage",
			storage: "invalid-storage",
			err:     fmt.Errorf("invalid storage type: invalid-storage"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			err := ValidateProjectStorage(test.storage)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, test.err, "an expected error not received")
				assert.Nil(t, err, "got an unexpected error")
			}
		})
	}
}

func TestValidateProjectFileExtension(t *testing.T) {
	tests := []struct {
		desc string
		file string
		err  error
	}{
		{
			desc: "Testing validate project file extension with valid file extension",
			file: "file.tar.gz",
			err:  nil,
		},
		{
			desc: "Testing validate project file extension with invalid file extension",
			file: "file.zip",
			err:  fmt.Errorf("file file.zip extension not supported"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)
			t.Parallel()

			err := ValidateProjectFileExtension(test.file)

			if err != nil {
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.Nil(t, test.err, "an expected error not received")
				assert.Nil(t, err, "got an unexpected error")
			}
		})
	}
}
