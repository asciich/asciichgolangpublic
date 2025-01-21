package commandlineinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestCommandLineInterface_IsLinePromptOnly(t *testing.T) {
	tests := []struct {
		line                 string
		expectedIsPromptOnly bool
	}{
		{"", false},       // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"s", false},      // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"sh", false},     // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"sh-", false},    // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"sh-5", false},   // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"sh-5.", false},  // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"sh-5.2", false}, // Every prompt must end with a "$" to avoid halfways loaded prompts are detected as full prompt.
		{"sh-5.2$", true},
		{"sh-5.2$ ", true},
		{"sh-5.2$ \n", true},
		{"\nsh-5.2$ \n", true},
		{"sh-5.2$ l", false},
		{"sh-5.2$ ls", false},
		{"sh-5.2$ ls\n", false},
		{"-sh-5.2$", true},
		{"-sh-5.2$ ", true},
		{"-sh-5.2$ \n", true},
		{"\n-sh-5.2$ \n", true},
		{"-sh-5.2$ l", false},
		{"-sh-5.2$ ls", false},
		{"-sh-5.2$ ls\n", false},
		{"bash-5.1$", true},
		{"bash-5.1$ ", true},
		{"bash-5.1$ \n", true},
		{"\nbash-5.1$ \n", true},
		{"bash-5.1$ l", false},
		{"bash-5.1$ ls", false},
		{"bash-5.1$ ls\n", false},
		{"-bash-5.1$", true},
		{"-bash-5.1$ ", true},
		{"-bash-5.1$ \n", true},
		{"\n-bash-5.1$ \n", true},
		{"-bash-5.1$ l", false},
		{"-bash-5.1$ ls", false},
		{"-bash-5.1$ ls\n", false},
		{"[user@host directory]$", true},
		{"[user@host directory]$ ", true},
		{"[user@host directory]$ \n", true},
		{"\n[user@host directory]$ \n", true},
		{"[user@host directory]$ l", false},
		{"[user@host directory]$ ls", false},
		{"[user@host directory]$ ls\n", false},
		{"[user@host\ndirectory]$\n", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsPromptOnly,
					IsLinePromptOnly(tt.line),
				)
			},
		)
	}
}
