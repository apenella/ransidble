package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectNotProvidedError(t *testing.T) {
	tests := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc:     "Testing project not provided error",
			err:      NewProjectNotProvidedError(fmt.Errorf("project not provided")),
			expected: "project not provided",
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
