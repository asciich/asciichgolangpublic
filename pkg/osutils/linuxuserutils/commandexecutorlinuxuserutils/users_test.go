package commandexecutorlinuxuserutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/commandexecutorlinuxuserutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/linuxuseroptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_CreateAndDeleteUser(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		name      string
		container containerinterfaces.Container
	}{
		{
			"archlinux",
			mustutils.Must(nativedocker.RunContainer(
				ctx,
				&dockeroptions.DockerRunContainerOptions{
					ImageName: "archlinux",
					Name:      "test-create-and-delete-user-archlinux",
					Command:   []string{"sleep", "1m"},
				},
			)),
		},
		{
			"ubuntu",
			mustutils.Must(nativedocker.RunContainer(
				ctx,
				&dockeroptions.DockerRunContainerOptions{
					ImageName: "ubuntu",
					Name:      "test-create-and-delete-user-ubuntu",
					Command:   []string{"sleep", "1m"},
				},
			)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const userName = "testuser"

			for range 2 {
				err := commandexecutorlinuxuserutils.Delete(ctx, tt.container, &linuxuseroptions.DeleteOptions{UserName: userName, Force: true})
				require.NoError(t, err)

				exists, err := commandexecutorlinuxuserutils.Exists(ctx, tt.container, userName)
				require.NoError(t, err)
				require.False(t, exists)
			}

			for range 2 {
				err := commandexecutorlinuxuserutils.Create(
					ctx,
					tt.container,
					&linuxuseroptions.CreateOptions{
						UserName: userName,
					})
				require.NoError(t, err)

				exists, err := commandexecutorlinuxuserutils.Exists(ctx, tt.container, userName)
				require.NoError(t, err)
				require.True(t, exists)
			}

			for range 2 {
				err := commandexecutorlinuxuserutils.Delete(ctx, tt.container, &linuxuseroptions.DeleteOptions{UserName: userName, Force: true})
				require.NoError(t, err)

				exists, err := commandexecutorlinuxuserutils.Exists(ctx, tt.container, userName)
				require.NoError(t, err)
				require.False(t, exists)
			}
		})
	}
}
