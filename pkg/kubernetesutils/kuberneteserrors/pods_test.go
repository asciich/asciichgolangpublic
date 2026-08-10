package kuberneteserrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kuberneteserrors"
)

func TestIsPodInFailedState(t *testing.T) {
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
			name:     "ErrPodInFailedState returns true",
			err:      kuberneteserrors.ErrPodInFailedState,
			expected: true,
		},
		{
			name:     "wrapped ErrPodInFailedState returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrPodInFailedState),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrPodInFailedState returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrPodInFailedState)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsPodInFailedState(tt.err)
			if result != tt.expected {
				t.Errorf("IsPodInFailedState(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsPodNotFound(t *testing.T) {
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
			name:     "ErrPodNotFound returns true",
			err:      kuberneteserrors.ErrPodNotFound,
			expected: true,
		},
		{
			name:     "wrapped ErrPodNotFound returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrPodNotFound),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrPodNotFound returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrPodNotFound)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsPodNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsPodNotFound(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsPodDeleted(t *testing.T) {
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
			name:     "ErrPodDeleted returns true",
			err:      kuberneteserrors.ErrPodDeleted,
			expected: true,
		},
		{
			name:     "wrapped ErrPodDeleted returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrPodDeleted),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrPodDeleted returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrPodDeleted)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsPodDeleted(tt.err)
			if result != tt.expected {
				t.Errorf("IsPodDeleted(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsPodNotRunning(t *testing.T) {
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
			name:     "ErrPodNotRunning returns true",
			err:      kuberneteserrors.ErrPodNotRunning,
			expected: true,
		},
		{
			name:     "wrapped ErrPodNotRunning returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrPodNotRunning),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrPodNotRunning returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrPodNotRunning)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsPodNotRunning(tt.err)
			if result != tt.expected {
				t.Errorf("IsPodNotRunning(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsPodNotSucceeded(t *testing.T) {
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
			name:     "ErrPodNotSucceeded returns true",
			err:      kuberneteserrors.ErrPodNotSucceeded,
			expected: true,
		},
		{
			name:     "wrapped ErrPodNotSucceeded returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrPodNotSucceeded),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrPodNotSucceeded returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrPodNotSucceeded)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsPodNotSucceeded(tt.err)
			if result != tt.expected {
				t.Errorf("IsPodNotSucceeded(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsPodTimeout(t *testing.T) {
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
			name:     "ErrPodTimeout returns true",
			err:      kuberneteserrors.ErrPodTimeout,
			expected: true,
		},
		{
			name:     "wrapped ErrPodTimeout returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrPodTimeout),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrPodTimeout returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrPodTimeout)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsPodTimeout(tt.err)
			if result != tt.expected {
				t.Errorf("IsPodTimeout(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsContainerNotFound(t *testing.T) {
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
			name:     "ErrContainerNotFound returns true",
			err:      kuberneteserrors.ErrContainerNotFound,
			expected: true,
		},
		{
			name:     "wrapped ErrContainerNotFound returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrContainerNotFound),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrContainerNotFound returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrContainerNotFound)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsContainerNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("IsContainerNotFound(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestIsUnableToRetrieveLogs(t *testing.T) {
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
			name:     "ErrUnableToRetrieveLogs returns true",
			err:      kuberneteserrors.ErrUnableToRetrieveLogs,
			expected: true,
		},
		{
			name:     "wrapped ErrUnableToRetrieveLogs returns true",
			err:      fmt.Errorf("some context: %w", kuberneteserrors.ErrUnableToRetrieveLogs),
			expected: true,
		},
		{
			name:     "unrelated error returns false",
			err:      errors.New("some other error"),
			expected: false,
		},
		{
			name:     "double wrapped ErrUnableToRetrieveLogs returns true",
			err:      fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", kuberneteserrors.ErrUnableToRetrieveLogs)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kuberneteserrors.IsUnableToRetrieveLogs(tt.err)
			if result != tt.expected {
				t.Errorf("IsUnableToRetrieveLogs(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}
