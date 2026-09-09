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
		{
			name:     "other sentinel error (no route to host)",
			err:      netutilserrors.ErrNoRouteToHost,
			expected: false,
		},
		{
			name:     "wrapped other sentinel error (no route to host)",
			err:      fmt.Errorf("failed to connect: %w", netutilserrors.ErrNoRouteToHost),
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

func TestIsNoRouteToHostError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "exact sentinel error",
			err:      netutilserrors.ErrNoRouteToHost,
			expected: true,
		},
		{
			name:     "wrapped sentinel error",
			err:      fmt.Errorf("failed to connect: %w", netutilserrors.ErrNoRouteToHost),
			expected: true,
		},
		{
			name:     "double wrapped sentinel error",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", netutilserrors.ErrNoRouteToHost)),
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
			err:      errors.New("no route to host"),
			expected: false,
		},
		{
			name:     "wrapped unrelated error",
			err:      fmt.Errorf("context: %w", errors.New("timeout")),
			expected: false,
		},
		{
			name:     "other sentinel error (connection refused)",
			err:      netutilserrors.ErrConnectionRefused,
			expected: false,
		},
		{
			name:     "wrapped other sentinel error (connection refused)",
			err:      fmt.Errorf("failed to connect: %w", netutilserrors.ErrConnectionRefused),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := netutilserrors.IsNoRouteToHostError(tt.err)
			require.Equal(t, tt.expected, result)
		})
	}
}
