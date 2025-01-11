package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getCommandExecutorByImplementationName(implementationName string) (commandExecutor CommandExecutor) {
	if implementationName == "Bash" {
		return Bash()
	}

	if implementationName == "Exec" {
		return Exec()
	}

	LogFatalf("Unnown implementation name: '%s'", implementationName)

	return nil
}

func TestCommandExecutor_GetDeepCopyOfCommandExecutor(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"Bash"},
		{"Exec"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				commandExecutor := getCommandExecutorByImplementationName(tt.implementationName)

				copy := MustGetDeepCopyOfCommandExecutor(commandExecutor)

				assert.EqualValues(
					commandExecutor.MustGetHostDescription(),
					copy.MustGetHostDescription(),
				)
			},
		)
	}
}
