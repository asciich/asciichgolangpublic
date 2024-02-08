package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecRunCommandAndGetStdoutAsString(t *testing.T) {
	tests := []struct {
		command        []string
		expectedOutput string
	}{
		{[]string{"echo", "hello"}, "hello\n"},
		{[]string{"echo", "hello world"}, "hello world\n"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				var exec CommandExecutor = Exec()
				output := exec.MustRunCommandAndGetStdoutAsString(
					&RunCommandOptions{
						Command: tt.command,
						Verbose: verbose,
					},
				)

				output2 := exec.MustRunCommandAndGetStdoutAsString(
					&RunCommandOptions{
						Command:            tt.command,
						Verbose:            verbose,
						LiveOutputOnStdout: true,
					},
				)

				assert.EqualValues(tt.expectedOutput, output)
				assert.EqualValues(tt.expectedOutput, output2)
			},
		)
	}
}
