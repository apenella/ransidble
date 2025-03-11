package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectAlreadyExistsError(t *testing.T) {
	tests := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc:     "Testing project already exists error",
			err:      NewProjectAlreadyExistsError(fmt.Errorf("project already exists")),
			expected: "project already exists",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			t.Log(test.desc)
			assert.Equal(t, test.expected, test.err.Error())
		})
	}
}
