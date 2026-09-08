package netutilserrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
)

func TestIsConnectionRefusedError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "exact sentinel error",
			err:      netutilserrors.ErrConnectionRefused,
			expected: true,
		},
		{
			name:     "wrapped sentinel error",
			err:      fmt.Errorf("failed to connect: %w", netutilserrors.ErrConnectionRefused),
			expected: true,
		},
		{
			name:     "double wrapped sentinel error",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", netutilserrors.ErrConnectionRefused)),
			expected: true,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "unrelated error",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "different error with same message text",
			err:      errors.New("connection refused"),
			expected: false,
		},
		{
			name:     "wrapped unrelated error",
			err:      fmt.Errorf("context: %w", errors.New("timeout")),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := netutilserrors.IsConnectionRefusedError(tt.err)
			require.Equal(t, tt.expected, result)
		})
	}
}
