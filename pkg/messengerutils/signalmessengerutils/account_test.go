package signalmessengerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
)

func TestIsAccountNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// ✅ Valid E.164 numbers
		{
			name:     "valid Swiss mobile number",
			input:    "+41791234567",
			expected: true,
		},
		{
			name:     "valid US number",
			input:    "+12025550123",
			expected: true,
		},
		{
			name:     "valid UK number",
			input:    "+447911123456",
			expected: true,
		},
		{
			name:     "valid minimum length (7 digits after country code)",
			input:    "+11234567",
			expected: true,
		},
		{
			name:     "valid maximum length (14 digits after country code)",
			input:    "+123456789012345",
			expected: true,
		},
		{
			name:     "valid German number",
			input:    "+4915112345678",
			expected: true,
		},
		{
			name:     "valid single digit country code",
			input:    "+11234567",
			expected: true,
		},

		// ❌ Invalid — missing plus sign
		{
			name:     "missing plus sign",
			input:    "41791234567",
			expected: false,
		},
		{
			name:     "starts with 00 instead of plus",
			input:    "0041791234567",
			expected: false,
		},

		// ❌ Invalid — wrong country code
		{
			name:     "country code starts with zero",
			input:    "+0791234567",
			expected: false,
		},

		// ❌ Invalid — too short / too long
		{
			name:     "too short (5 digits after country code)",
			input:    "+41791",
			expected: false,
		},
		{
			name:     "too long (15 digits after country code)",
			input:    "+1234567890123456",
			expected: false,
		},
		{
			name:     "only plus sign",
			input:    "+",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},

		// ❌ Invalid — non-digit characters
		{
			name:     "contains spaces",
			input:    "+41 79 123 45 67",
			expected: false,
		},
		{
			name:     "contains dashes",
			input:    "+41-79-123-45-67",
			expected: false,
		},
		{
			name:     "contains parentheses",
			input:    "+1(202)5550123",
			expected: false,
		},
		{
			name:     "contains letters",
			input:    "+417912345AB",
			expected: false,
		},
		{
			name:     "alphanumeric string",
			input:    "abcdefgh",
			expected: false,
		},

		// ❌ Invalid — formatting issues
		{
			name:     "plus sign in the middle",
			input:    "417+91234567",
			expected: false,
		},
		{
			name:     "multiple plus signs",
			input:    "++41791234567",
			expected: false,
		},
		{
			name:     "plus sign at the end",
			input:    "41791234567+",
			expected: false,
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: false,
		},
		{
			name:     "newline character",
			input:    "+41791234567\n",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := signalmessengerutils.IsAccountNumber(tt.input)
			require.EqualValues(t, tt.expected, result)
		})
	}
}
