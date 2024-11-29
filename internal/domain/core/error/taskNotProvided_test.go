package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskNotProvidedError(t *testing.T) {
	tests := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc:     "Testing task not provided error",
			err:      NewTaskNotProvidedError(fmt.Errorf("task not provided")),
			expected: "task not provided",
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
