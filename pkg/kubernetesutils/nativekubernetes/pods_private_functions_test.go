package nativekubernetes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractExitCodeFromError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		expectedOk   bool
	}{
		{
			name:         "nil error",
			err:          nil,
			expectedCode: 0,
			expectedOk:   false,
		},
		{
			name:         "exit code 1",
			err:          errors.New("command terminated with exit code 1"),
			expectedCode: 1,
			expectedOk:   true,
		},
		{
			name:         "exit code 0",
			err:          errors.New("command terminated with exit code 0"),
			expectedCode: 0,
			expectedOk:   true,
		},
		{
			name:         "exit code 127",
			err:          errors.New("command terminated with exit code 127"),
			expectedCode: 127,
			expectedOk:   true,
		},
		{
			name:         "exit code 255",
			err:          errors.New("command terminated with exit code 255"),
			expectedCode: 255,
			expectedOk:   true,
		},
		{
			name:         "error without exit code",
			err:          errors.New("connection refused"),
			expectedCode: 0,
			expectedOk:   false,
		},
		{
			name:         "error with exit code in different format",
			err:          errors.New("process exited with status 1"),
			expectedCode: 0,
			expectedOk:   false,
		},
		{
			name:         "exit code embedded in longer message",
			err:          errors.New("Error executing command: command terminated with exit code 2: some additional context"),
			expectedCode: 2,
			expectedOk:   true,
		},
		{
			name:         "exit code embedded in longer message",
			err:          errors.New("Error executing command: command terminated with exit code 1"),
			expectedCode: 1,
			expectedOk:   true,
		},
		{
			name:         "empty error message",
			err:          errors.New(""),
			expectedCode: 0,
			expectedOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, ok := extractExitCodeFromError(tt.err)
			require.Equal(t, tt.expectedCode, code)
			require.Equal(t, tt.expectedOk, ok)
		})
	}
}
