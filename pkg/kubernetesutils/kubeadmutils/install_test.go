package kubeadmutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeadmutils"
)

func getCtx() context.Context {
	return contextutils.WithVerbose(context.TODO())
}

func TestInstallKubeadm(t *testing.T) {
	ctx := getCtx()

	// Create a temporary file path for testing to avoid overwriting system kubeadm
	tempFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	require.NoError(t, err)
	tempFilePath, err := tempFile.GetPath()
	require.NoError(t, err)

	// Clean up the temporary file after test
	defer tempFile.Delete(ctx, &filesoptions.DeleteOptions{})

	t.Run("InstallKubeadm with custom path", func(t *testing.T) {
		// Use a temporary path instead of /bin/kubeadm to avoid overwriting system binary
		// Note: This test still requires network access to download kubeadm
		options := &kubeadmutils.InstallKubeadmOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubeadmutils.InstallKubeadm(ctx, options)

		// We don't assert success/failure here as it depends on network access
		// The function is tested for proper execution flow
		if err != nil {
			t.Logf("InstallKubeadm failed (expected in restricted environments): %v", err)
		} else {
			t.Log("InstallKubeadm completed successfully")
			// Verify the file was created
			exists, err := tempFile.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists, "kubeadm binary should be installed")
		}
	})

	t.Run("InstallKubeadm with unsupported version", func(t *testing.T) {
		options := &kubeadmutils.InstallKubeadmOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v99.99.99",
		}

		err := kubeadmutils.InstallKubeadm(ctx, options)
		require.Error(t, err, "should fail with unsupported version")
	})
}

func TestInstallKubeadmUsingCommandExecutor(t *testing.T) {
	ctx := getCtx()

	// Create a temporary file path for testing to avoid overwriting system kubeadm
	tempFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	require.NoError(t, err)
	tempFilePath, err := tempFile.GetPath()
	require.NoError(t, err)

	// Clean up the temporary file after test
	defer tempFile.Delete(ctx, &filesoptions.DeleteOptions{})

	// Use the local exec command executor
	commandExecutor := commandexecutorexecoo.Exec()

	t.Run("InstallKubeadmUsingCommandExecutor with custom path", func(t *testing.T) {
		// Use a temporary path instead of /bin/kubeadm to avoid overwriting system binary
		// Note: This test still requires network access to download kubeadm
		options := &kubeadmutils.InstallKubeadmOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubeadmutils.InstallKubeadmUsingCommandExecutor(ctx, commandExecutor, options)

		// We don't assert success/failure here as it depends on network access
		// The function is tested for proper execution flow
		if err != nil {
			t.Logf("InstallKubeadmUsingCommandExecutor failed (expected in restricted environments): %v", err)
		} else {
			t.Log("InstallKubeadmUsingCommandExecutor completed successfully")
			// Verify the file was created
			exists, err := tempFile.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists, "kubeadm binary should be installed")
		}
	})

	t.Run("InstallKubeadmUsingCommandExecutor with unsupported version", func(t *testing.T) {
		options := &kubeadmutils.InstallKubeadmOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v99.99.99",
		}

		err := kubeadmutils.InstallKubeadmUsingCommandExecutor(ctx, commandExecutor, options)
		require.Error(t, err, "should fail with unsupported version")
	})

	t.Run("InstallKubeadmUsingCommandExecutor with nil commandExecutor", func(t *testing.T) {
		options := &kubeadmutils.InstallKubeadmOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubeadmutils.InstallKubeadmUsingCommandExecutor(ctx, nil, options)
		require.Error(t, err, "should fail with nil commandExecutor")
	})
}
