package commandexecutor_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
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
		return commandexecutor.Bash()
	}

	if implementationName == "Exec" {
		return commandexecutor.Exec()
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

func Test_WithLiveOutputOnStdout(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, commandexecutor.IsLiveOutputOnStdoutEnabled(nil))
	})

	t.Run("with nil", func(t *testing.T) {
		ctx := commandexecutor.WithLiveOutputOnStdout(nil)
		require.True(t, commandexecutor.IsLiveOutputOnStdoutEnabled(ctx))
	})

	t.Run("with silent", func(t *testing.T) {
		ctx := commandexecutor.WithLiveOutputOnStdout(contextutils.ContextSilent())
		require.True(t, commandexecutor.IsLiveOutputOnStdoutEnabled(ctx))
	})

	t.Run("with silent", func(t *testing.T) {
		ctx := commandexecutor.WithLiveOutputOnStdout(contextutils.ContextVerbose())
		require.True(t, commandexecutor.IsLiveOutputOnStdoutEnabled(ctx))
	})

	t.Run("silent context", func(t *testing.T) {
		require.False(t, commandexecutor.IsLiveOutputOnStdoutEnabled(contextutils.ContextSilent()))
	})

	t.Run("verbose context", func(t *testing.T) {
		require.False(t, commandexecutor.IsLiveOutputOnStdoutEnabled(contextutils.ContextVerbose()))
	})

	t.Run("if verbose", func(t *testing.T) {
		require.False(
			t,
			commandexecutor.IsLiveOutputOnStdoutEnabled(
				commandexecutor.WithLiveOutputOnStdoutIfVerbose(contextutils.ContextSilent()),
			),
		)

		require.True(
			t,
			commandexecutor.IsLiveOutputOnStdoutEnabled(
				commandexecutor.WithLiveOutputOnStdoutIfVerbose(contextutils.ContextVerbose()),
			),
		)
	})
}
