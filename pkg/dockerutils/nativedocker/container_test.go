package nativedocker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func Test_NewContainer(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		container, err := nativedocker.NewContainer("")
		require.Error(t, err)
		require.Nil(t, container)
	})

	t.Run("simplename", func(t *testing.T) {
		container, err := nativedocker.NewContainer("simplename")
		require.NoError(t, err)
		require.NotNil(t, container)

		name, err := container.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "simplename", name)
	})
}

func Test_ContainerIsACommandExecutor(t *testing.T) {
	// It should be possible to run commands in a container in the same way as other CommandExectuors.
	// This test validates a container fullfils the CommandExecutor interface:

	var commandExecutor commandexecutorinterfaces.CommandExecutor
	var err error

	commandExecutor, err = nativedocker.NewContainer("containername")
	require.NoError(t, err)
	require.NotNil(t, commandExecutor)
}

func Test_RunCommand(t *testing.T) {
	var tests = []struct {
		cmd string
	}{
		{cmd: "printenv"},
		{cmd: "/usr/bin/printenv"},
	}

	for _, tt := range tests {
		t.Run("path variable using "+tt.cmd, func(t *testing.T) {
			ctx := getCtx()
			const containerName = "test-nativedocker-path-variable"
			docker := nativedocker.NewDocker()
			require.NoError(t, docker.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true}))

			container, err := docker.RunContainer(ctx, &dockeroptions.DockerRunContainerOptions{
				ImageName: "ubuntu",
				Name:      containerName,
				Command:   []string{"sleep", "1m"},
			})
			require.NoError(t, err)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			stdout, err := container.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
				Command: []string{tt.cmd, "PATH"},
			})
			require.NoError(t, err)
			require.Contains(t, stdout, "/bin")
		})
	}
}
