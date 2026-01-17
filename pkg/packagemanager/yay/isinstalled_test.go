package yay_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func getArchLinuxContainer(t *testing.T) containerinterfaces.Container {
	container, err := nativedocker.NewDocker().RunContainer(getCtx(),
		&dockeroptions.DockerRunContainerOptions{
			Name:      "test-archlinux-yay",
			ImageName: "archlinux",
			Command:   []string{"sleep", "1m"},
		},
	)
	require.NoError(t, err)
	return container
}

func Test_YayInstalled(t *testing.T) {
	t.Run("default archlinux", func(t *testing.T) {
		ctx := getCtx()

		container := getArchLinuxContainer(t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := yay.IsInstalled(ctx, container)
		require.NoError(t, err)
		// Out of the box yay is not installed on archlinux:
		require.False(t, isInstalled)
	})
}
