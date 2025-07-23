package dockerutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/containerutils/containerinterfaces"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/dockerutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getDockerContainerToTest(t *testing.T, implementationName string, containerName string) (container containerinterfaces.Container) {
	if implementationName == "commandExectuorDockerContainer" {
		executor, err := dockerutils.GetLocalCommandExecutorDocker()
		require.NoError(t, err)

		container, err := executor.GetContainerByName(containerName)
		require.NoError(t, err)
		return container
	}

	logging.LogFatalWithTracef("Unkown implementaion name: '%s'", implementationName)

	return nil
}

func TestContainers_IsHostRunning(t *testing.T) {

	tests := []struct {
		implementationName string
	}{
		{"commandExectuorDockerContainer"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {

				container := getDockerContainerToTest(t, tt.implementationName, "thisContainerDoesNotRun")

				isRunning, err := container.IsRunning(getCtx())
				require.NoError(t, err)
				require.False(t, isRunning)
			},
		)
	}
}
