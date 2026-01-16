package dockerutils_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/commandexecutordocker"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/commandexecutorlinuxuserutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/linuxuseroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

const username = "testuser"

func getNativeContainer(ctx context.Context, t *testing.T, imageName string) containerinterfaces.Container {
	const name = "test-runasuser"

	container, err := nativedocker.GetContainerByName(name)
	require.NoError(t, err)

	require.NoError(t, container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true}))

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: imageName,
		Command:   []string{"sleep", "1m"},
	})
	require.NoError(t, err)

	err = commandexecutorlinuxuserutils.Create(ctx, container, &linuxuseroptions.CreateOptions{
		UserName: "testuser",
	})
	require.NoError(t, err)

	return container
}

func getCommandExecutorContainer(ctx context.Context, t *testing.T, imageName string) containerinterfaces.Container {
	const name = "test-runasuser"

	docker, err := commandexecutordocker.GetLocalCommandExecutorDocker()
	require.NoError(t, err)

	container, err := docker.GetContainerByName(name)
	require.NoError(t, err)

	require.NoError(t, container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true}))

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: imageName,
		Command:   []string{"sleep", "1m"},
	})
	require.NoError(t, err)

	return container
}

func TestRunAsUser(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		name     string
		executor commandexecutorinterfaces.CommandExecutor
		useSudo  bool
	}{
		{name: "exec", executor: commandexecutorexecoo.NewExec(), useSudo: true},
		{name: "bash", executor: commandexecutorbashoo.Bash(), useSudo: true},
		{name: "native docker archlinux", executor: getNativeContainer(ctx, t, "archlinux"), useSudo: false},
		{name: "native docker ubuntu", executor: getNativeContainer(ctx, t, "ubuntu"), useSudo: false},
		{name: "commandexecutor docker archlinux", executor: getCommandExecutorContainer(ctx, t, "archlinux"), useSudo: false},
		{name: "commandexecutor docker ubuntu", executor: getCommandExecutorContainer(ctx, t, "ubuntu"), useSudo: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := commandexecutorlinuxuserutils.Create(ctx, tt.executor, &linuxuseroptions.CreateOptions{
				UserName: username,
				UseSudo:  tt.useSudo,
			})
			require.NoError(t, err)
			defer func() {
				err := commandexecutorlinuxuserutils.Delete(ctx, tt.executor, &linuxuseroptions.DeleteOptions{
					UserName: username,
					UseSudo:  tt.useSudo,
				})
				require.NoError(t, err)
			}()

			output, err := tt.executor.RunCommandAndGetStdoutAsString(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command:            []string{"whoami"},
					RunAsUser:          username,
					UseSudoToRunAsUser: tt.useSudo,
				},
			)
			require.NoError(t, err)
			require.EqualValues(t, "testuser", strings.TrimSpace(output))
		})
	}
}
