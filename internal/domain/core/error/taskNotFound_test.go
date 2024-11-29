package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskNotFound(t *testing.T) {
	tests := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc:     "Testing task not found error",
			err:      NewTaskNotFoundError(fmt.Errorf("task not found")),
			expected: "task not found",
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
