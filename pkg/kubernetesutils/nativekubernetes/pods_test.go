package nativekubernetes_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_CreateAndDeletePod(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	namespaceName := "test-createanddeletepod"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	t.Run("happy path", func(t *testing.T) {

		const podName = "testpod"

		err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
		require.NoError(t, err)

		exists, err := nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
		require.NoError(t, err)
		require.False(t, exists)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreatePod(
				ctx,
				clientset,
				namespaceName,
				&kubernetesparameteroptions.RunCommandOptions{
					PodName: podName,
					Image:   "ubunt",
					Command: []string{"bash", "-c", "sleep 1m"},
				})
			require.NoError(t, err)

			exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
			require.NoError(t, err)
			require.True(t, exists)

			err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
			require.NoError(t, err)

			exists, err = nativekubernetes.PodExists(ctx, clientset, podName, namespaceName)
			require.NoError(t, err)
			require.False(t, exists)
		}
	})
}

func Test_WaitForDeleted(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("already deleted pod", func(t *testing.T) {
		podName := "testpod"
		namespaceName := "test-createanddeletepod"

		// Ensure pod is absent
		err = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
		require.NoError(t, err)

		// Check there's no wait for an already deleted pod:
		err = nativekubernetes.WaitForPodDeleted(ctx, clientset, podName, namespaceName, time.Second*1)
		require.NoError(t, err)
	})
}

func Test_ListPods(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	namespaceName := "test-listpods"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("create and delete pods with list in between", func(t *testing.T) {
		podNames := []string{"listpod-1", "listpod-2", "listpod-3"}

		// Ensure all test pods are absent before starting
		for _, name := range podNames {
			err = nativekubernetes.DeletePod(ctx, clientset, name, namespaceName)
			require.NoError(t, err)
		}

		// Create pods one by one and verify list grows
		for i, name := range podNames {
			err = nativekubernetes.CreatePod(
				ctx,
				clientset,
				namespaceName,
				&kubernetesparameteroptions.RunCommandOptions{
					PodName: name,
					Image:   "ubuntu",
					Command: []string{"bash", "-c", "sleep 1m"},
				})
			require.NoError(t, err)

			listed, err := nativekubernetes.ListPods(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, created := range podNames[:i+1] {
				require.Contains(t, listed, created)
			}
			for _, notYetCreated := range podNames[i+1:] {
				require.NotContains(t, listed, notYetCreated)
			}
		}

		// Delete pods one by one and verify list shrinks
		for i, name := range podNames {
			err = nativekubernetes.DeletePod(ctx, clientset, name, namespaceName)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListPods(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, deleted := range podNames[:i+1] {
				require.NotContains(t, listed, deleted)
			}
			for _, stillPresent := range podNames[i+1:] {
				require.Contains(t, listed, stillPresent)
			}
		}
	})
}

func Test_CopyFileToPod(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	namespaceName := "test-copyfiletopod"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	const podName = "test-copyfiletopod-pod"
	containerName := podName

	// Create a test pod
	err = nativekubernetes.CreatePod(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.RunCommandOptions{
			PodName:                  podName,
			Image:                    "ubuntu",
			Command:                  []string{"sh", "-c", "trap \"echo Caught SIGTERM, exiting...; exit 0\" TERM; while true; do sleep .1; done"},
			DeleteAlreadyExistingPod: true,
			WaitForPodRunning:        true,
		})
	require.NoError(t, err)

	// Ensure pod is cleaned up after test
	defer func() {
		_ = nativekubernetes.DeletePod(ctx, clientset, podName, namespaceName)
	}()
	// ... prepare test environment finished.
	// -----

	t.Run("copy file to pod", func(t *testing.T) {
		// Generate random test content
		testContent := "This is test file content for CopyFileToPod test.\nLine 2 of test content.\nLine 3 with special chars: !@#$%^&*()"

		// Create a temporary file with test content
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
		require.NoError(t, err)

		// Define destination path in pod
		destPath := "/tmp/testfile.txt"

		// Copy file to pod
		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
		require.NoError(t, err)

		// Verify file exists and has correct content
		commandOutput, err := nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"cat", destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)

		readContent, err := commandOutput.GetStdoutAsString()
		require.NoError(t, err)
		require.Equal(t, testContent, readContent)
	})

	t.Run("copy file to nested directory in pod", func(t *testing.T) {
		// Generate random test content
		testContent := "Test content for nested directory"

		// Create a temporary file with test content
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
		require.NoError(t, err)

		// Define destination path in nested directory
		destPath := "/tmp/nested/dir/testfile.txt"

		// Create the nested directory first
		_, err = nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"mkdir", "-p", "/tmp/nested/dir"},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)

		// Copy file to pod
		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
		require.NoError(t, err)

		// Verify file exists and has correct content
		commandOutput, err := nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"cat", destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)

		readContent, err := commandOutput.GetStdoutAsString()
		require.NoError(t, err)
		require.Equal(t, testContent, readContent)
	})

	t.Run("copy binary file to pod", func(t *testing.T) {
		// Generate binary test content
		testContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}

		// Create a temporary file with binary content
		localPath, err := tempfiles.CreateTemporaryFileFromContentBytes(ctx, testContent)
		require.NoError(t, err)

		// Define destination path
		destPath := "/tmp/binaryfile.bin"

		// Copy file to pod
		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
		require.NoError(t, err)

		// Verify file exists and has correct content using xxd or od
		commandOutput, err := nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"od", "-A", "x", "-t", "x1", destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)

		stdout, err := commandOutput.GetStdoutAsString()
		require.NoError(t, err)
		// Verify the binary content is present (od will show hex values)
		require.Contains(t, stdout, "00 01 02 03")
		require.Contains(t, stdout, "ff fe fd")
	})

	t.Run("copy file with special characters in name", func(t *testing.T) {
		testContent := "Test content for file with special name"

		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, testContent)
		require.NoError(t, err)

		// Use a destination with special characters
		destPath := "/tmp/test-file_with.special.chars.txt"

		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
		require.NoError(t, err)

		// Verify file exists and has correct content
		commandOutput, err := nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"cat", destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)

		readContent, err := commandOutput.GetStdoutAsString()
		require.NoError(t, err)
		require.Equal(t, testContent, readContent)
	})

	t.Run("error handling - nil config", func(t *testing.T) {
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "test")
		require.NoError(t, err)

		err = nativekubernetes.CopyFileToPod(ctx, nil, localPath, "/tmp/test.txt", podName, containerName, namespaceName)
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil")
	})

	t.Run("error handling - empty local file path", func(t *testing.T) {
		err = nativekubernetes.CopyFileToPod(ctx, config, "", "/tmp/test.txt", podName, containerName, namespaceName)
		require.Error(t, err)
		require.Contains(t, err.Error(), "localFile")
	})

	t.Run("error handling - empty destination path", func(t *testing.T) {
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "test")
		require.NoError(t, err)

		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, "", podName, containerName, namespaceName)
		require.Error(t, err)
		require.Contains(t, err.Error(), "destPath")
	})

	t.Run("error handling - empty pod name", func(t *testing.T) {
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "test")
		require.NoError(t, err)

		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, "/tmp/test.txt", "", containerName, namespaceName)
		require.Error(t, err)
		require.Contains(t, err.Error(), "podName")
	})

	t.Run("error handling - empty namespace", func(t *testing.T) {
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "test")
		require.NoError(t, err)

		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, "/tmp/test.txt", podName, containerName, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "namespaceName")
	})

	t.Run("error handling - non-existent pod", func(t *testing.T) {
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "test")
		require.NoError(t, err)

		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, "/tmp/test.txt", "non-existent-pod", containerName, namespaceName)
		require.Error(t, err)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		// First, create a file with initial content
		initialContent := "Initial content"
		destPath := "/tmp/overwrite-test.txt"

		// Create initial file using exec
		_, err := nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"sh", "-c", "echo '" + initialContent + "' > " + destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)

		// Verify initial content
		commandOutput, err := nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"cat", destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)
		content, err := commandOutput.GetStdoutAsString()
		require.NoError(t, err)
		require.Equal(t, initialContent+"\n", content)

		// Now overwrite with new content
		newContent := "New content after overwrite"
		localPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, newContent)
		require.NoError(t, err)

		err = nativekubernetes.CopyFileToPod(ctx, config, localPath, destPath, podName, containerName, namespaceName)
		require.NoError(t, err)

		// Verify new content
		commandOutput, err = nativekubernetes.Exec(
			ctx,
			config,
			namespaceName,
			&kubernetesparameteroptions.RunCommandOptions{
				PodName:       podName,
				Command:       []string{"cat", destPath},
				ContainerName: containerName,
			},
		)
		require.NoError(t, err)
		content, err = commandOutput.GetStdoutAsString()
		require.NoError(t, err)
		require.Equal(t, newContent, content)
	})
}
