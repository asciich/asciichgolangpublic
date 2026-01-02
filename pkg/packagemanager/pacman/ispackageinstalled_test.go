package pacman_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
)

func Test_IsPackageInstalled(t *testing.T) {
	t.Run("empty package name", func(t *testing.T) {
		ctx := getCtx()
		pacman, container := getPacmanInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := pacman.IsPackageInstalled(ctx, "")
		require.Error(t, err)
		require.False(t, isInstalled)
	})

	t.Run("Already installed package", func(t *testing.T) {
		ctx := getCtx()
		pacman, container := getPacmanInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := pacman.IsPackageInstalled(ctx, "archlinux-keyring")
		require.NoError(t, err)
		require.True(t, isInstalled)
	})

	t.Run("Not installed package", func(t *testing.T) {
		ctx := getCtx()
		pacman, container := getPacmanInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := pacman.IsPackageInstalled(ctx, "nvidia-utils")
		require.NoError(t, err)
		require.False(t, isInstalled)
	})
}
