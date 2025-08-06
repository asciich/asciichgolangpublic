package commandexecutor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getCommandExecutorByImplementationName(implementationName string) (commandExecutor commandexecutorinterfaces.CommandExecutor) {
	if implementationName == "Bash" {
		return commandexecutorbashoo.Bash()
	}

	if implementationName == "Exec" {
		return commandexecutorexecoo.Exec()
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				commandExecutor := getCommandExecutorByImplementationName(tt.implementationName)

				copy, err := commandexecutor.GetDeepCopyOfCommandExecutor(commandExecutor)
				require.NoError(t, err)

				expectedHostDescription, err := commandExecutor.GetHostDescription()
				require.NoError(t, err)

				hostDescription, err := copy.GetHostDescription()
				require.NoError(t, err)

				require.EqualValues(
					t,
					expectedHostDescription,
					hostDescription,
				)
			},
		)
	}
}
