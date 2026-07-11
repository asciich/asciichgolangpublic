package kuberneteserrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kuberneteserrors"
)

func TestIsSecretNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error returns false",
			err:      nil,
			expected: false,
		},
		{
			name:     "ErrSecretNotFound returns true",
			err:      kuberneteserrors.ErrSecretNotFound,
			expected: true,
		},
		{
			name:     "wrapped ErrSecretNotFound returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrSecretNotFound),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrSecretNotFound returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrSecretNotFound)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsSecretNotFoundError(tt.err)
			if result != tt.expected {
				t.Errorf("IsSecretNotFoundError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}
