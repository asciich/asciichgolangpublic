package commandexecutor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCommandExecutorByImplementationName(implementationName string) (commandExecutor CommandExecutor) {
	if implementationName == "Bash" {
		return Bash()
	}

	if implementationName == "Exec" {
		return Exec()
	}

	logging.LogFatalf("Unnown implementation name: '%s'", implementationName)

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
			testutils.MustFormatAsTestname(tt),
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