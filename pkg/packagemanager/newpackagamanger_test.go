package packagemanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanagergeneric"
)

func Test_NewPackageManagerGeneric_autoDetect(t *testing.T) {
	for _, image := range getImagesTotest() {
		t.Run("auto detect", func(t *testing.T) {
			ctx := getCtx()
			container := getTestContainer(ctx, t, image)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
			require.NoError(t, err)
			require.NotNil(t, packageManager)

			// The auto-detected package manager must be functional:
			packageName := getAlreadyInstalledPackage(t, image)
			isInstalled, err := packageManager.IsPackageInstalled(ctx, packageName)
			require.NoError(t, err)
			require.True(t, isInstalled)
		})
	}
}
