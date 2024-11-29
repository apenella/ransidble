package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectNotFoundError(t *testing.T) {
	tests := []struct {
		desc     string
		err      error
		expected string
	}{
		{
			desc:     "Testing project not found error",
			err:      NewProjectNotFoundError(fmt.Errorf("project not found")),
			expected: "project not found",
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
