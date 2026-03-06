package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectIDNotProvidedError(t *testing.T) {
	tests := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc:     "Testing project id not provided error",
			err:      NewProjectIDNotProvidedError(fmt.Errorf("project id not provided")),
			expected: "project id not provided",
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
