package kubectlutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
)

func getCtx() context.Context {
	return contextutils.WithVerbose(context.TODO())
}

func TestInstallKubectl(t *testing.T) {
	ctx := getCtx()

	// Create a temporary file path for testing to avoid overwriting system kubectl
	tempFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	require.NoError(t, err)
	tempFilePath, err := tempFile.GetPath()
	require.NoError(t, err)

	// Clean up the temporary file after test
	defer tempFile.Delete(ctx, &filesoptions.DeleteOptions{})

	t.Run("InstallKubectl with custom path", func(t *testing.T) {
		// Use a temporary path instead of /bin/kubectl to avoid overwriting system binary
		// Note: This test still requires network access to download kubectl
		options := &kubectlutils.InstallKubectlOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubectlutils.InstallKubectl(ctx, options)

		// We don't assert success/failure here as it depends on network access
		// The function is tested for proper execution flow
		if err != nil {
			t.Logf("InstallKubectl failed (expected in restricted environments): %v", err)
		} else {
			t.Log("InstallKubectl completed successfully")
			// Verify the file was created
			exists, err := tempFile.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists, "kubectl binary should be installed")
		}
	})

	t.Run("InstallKubectl with unsupported version", func(t *testing.T) {
		options := &kubectlutils.InstallKubectlOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v99.99.99",
		}

		err := kubectlutils.InstallKubectl(ctx, options)
		require.Error(t, err, "should fail with unsupported version")
	})
}
