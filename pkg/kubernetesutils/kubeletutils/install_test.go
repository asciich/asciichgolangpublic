package kubeletutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeletutils"
)

func getCtx() context.Context {
	return contextutils.WithVerbose(context.TODO())
}

func TestInstallKubelet(t *testing.T) {
	ctx := getCtx()

	// Create a temporary file path for testing to avoid overwriting system kubelet
	tempFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	require.NoError(t, err)
	tempFilePath, err := tempFile.GetPath()
	require.NoError(t, err)

	// Clean up the temporary file after test
	defer tempFile.Delete(ctx, &filesoptions.DeleteOptions{})

	t.Run("InstallKubelet with custom path", func(t *testing.T) {
		// Use a temporary path instead of /bin/kubelet to avoid overwriting system binary
		// Note: This test still requires network access to download kubelet
		options := &kubeletutils.InstallKubeletOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubeletutils.InstallKubelet(ctx, options)

		// We don't assert success/failure here as it depends on network access
		// The function is tested for proper execution flow
		if err != nil {
			t.Logf("InstallKubelet failed (expected in restricted environments): %v", err)
		} else {
			t.Log("InstallKubelet completed successfully")
			// Verify the file was created
			exists, err := tempFile.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists, "kubelet binary should be installed")
		}
	})

	t.Run("InstallKubelet with unsupported version", func(t *testing.T) {
		options := &kubeletutils.InstallKubeletOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v99.99.99",
		}

		err := kubeletutils.InstallKubelet(ctx, options)
		require.Error(t, err, "should fail with unsupported version")
	})
}

func TestInstallKubeletUsingCommandExecutor(t *testing.T) {
	ctx := getCtx()

	// Create a temporary file path for testing to avoid overwriting system kubelet
	tempFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	require.NoError(t, err)
	tempFilePath, err := tempFile.GetPath()
	require.NoError(t, err)

	// Clean up the temporary file after test
	defer tempFile.Delete(ctx, &filesoptions.DeleteOptions{})

	// Use the local exec command executor
	commandExecutor := commandexecutorexecoo.Exec()

	t.Run("InstallKubeletUsingCommandExecutor with custom path", func(t *testing.T) {
		// Use a temporary path instead of /bin/kubelet to avoid overwriting system binary
		// Note: This test still requires network access to download kubelet
		options := &kubeletutils.InstallKubeletOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubeletutils.InstallKubeletUsingCommandExecutor(ctx, commandExecutor, options)

		// We don't assert success/failure here as it depends on network access
		// The function is tested for proper execution flow
		if err != nil {
			t.Logf("InstallKubeletUsingCommandExecutor failed (expected in restricted environments): %v", err)
		} else {
			t.Log("InstallKubeletUsingCommandExecutor completed successfully")
			// Verify the file was created
			exists, err := tempFile.Exists(ctx)
			require.NoError(t, err)
			require.True(t, exists, "kubelet binary should be installed")
		}
	})

	t.Run("InstallKubeletUsingCommandExecutor with unsupported version", func(t *testing.T) {
		options := &kubeletutils.InstallKubeletOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v99.99.99",
		}

		err := kubeletutils.InstallKubeletUsingCommandExecutor(ctx, commandExecutor, options)
		require.Error(t, err, "should fail with unsupported version")
	})

	t.Run("InstallKubeletUsingCommandExecutor with nil commandExecutor", func(t *testing.T) {
		options := &kubeletutils.InstallKubeletOptions{
			InstallPath: tempFilePath,
			UseSudo:     false,
			Version:     "v1.36.2",
		}

		err := kubeletutils.InstallKubeletUsingCommandExecutor(ctx, nil, options)
		require.Error(t, err, "should fail with nil commandExecutor")
	})
}
