package aptget_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_IsPackageInstalledAptGet(t *testing.T) {
	t.Run("empty package name", func(t *testing.T) {
		ctx := getCtx()
		aptGet, container := getAptGetInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := aptGet.IsPackageInstalled(ctx, "")
		require.Error(t, err)
		require.False(t, isInstalled)
	})

	t.Run("Already installed package", func(t *testing.T) {
		ctx := getCtx()
		aptGet, container := getAptGetInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := aptGet.IsPackageInstalled(ctx, "base-files")
		require.NoError(t, err)
		require.True(t, isInstalled)
	})

	t.Run("Not installed package", func(t *testing.T) {
		ctx := getCtx()
		aptGet, container := getAptGetInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		isInstalled, err := aptGet.IsPackageInstalled(ctx, "nvidia-utils-535")
		require.NoError(t, err)
		require.False(t, isInstalled)
	})
}
