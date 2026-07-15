package errorutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUserAbort(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		// nil
		{
			name:     "nil error returns false",
			err:      nil,
			expected: false,
		},

		// abort
		{
			name:     "lowercase 'abort' returns true",
			err:      errors.New("abort"),
			expected: true,
		},
		{
			name:     "uppercase 'ABORT' returns true",
			err:      errors.New("ABORT"),
			expected: true,
		},
		{
			name:     "mixed case 'Abort' returns true",
			err:      errors.New("Abort"),
			expected: true,
		},
		{
			name:     "'abort' in sentence returns true",
			err:      errors.New("user aborted the operation"),
			expected: true,
		},

		// interrupt
		{
			name:     "lowercase 'interrupt' returns true",
			err:      errors.New("interrupt"),
			expected: true,
		},
		{
			name:     "uppercase 'INTERRUPT' returns true",
			err:      errors.New("INTERRUPT"),
			expected: true,
		},
		{
			name:     "mixed case 'Interrupt' returns true",
			err:      errors.New("Interrupt"),
			expected: true,
		},
		{
			name:     "'interrupt' in sentence returns true",
			err:      errors.New("signal: interrupt"),
			expected: true,
		},

		// no item
		{
			name:     "lowercase 'no item' returns true",
			err:      errors.New("no item"),
			expected: true,
		},
		{
			name:     "uppercase 'NO ITEM' returns true",
			err:      errors.New("NO ITEM"),
			expected: true,
		},
		{
			name:     "mixed case 'No Item' returns true",
			err:      errors.New("No Item"),
			expected: true,
		},
		{
			name:     "'no item' in sentence returns true",
			err:      errors.New("fzf: no item selected"),
			expected: true,
		},

		// unrelated errors
		{
			name:     "unrelated error returns false",
			err:      errors.New("something went wrong"),
			expected: false,
		},
		{
			name:     "empty error message returns false",
			err:      errors.New(""),
			expected: false,
		},
		{
			name:     "timeout error returns false",
			err:      errors.New("context deadline exceeded"),
			expected: false,
		},
		{
			name:     "permission error returns false",
			err:      errors.New("permission denied"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsUserAbort(tt.err))
		})
	}
}
