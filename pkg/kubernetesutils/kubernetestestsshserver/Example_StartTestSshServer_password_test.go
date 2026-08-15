package kubernetestestsshserver_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetestestsshserver"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// Test_Example_StartTestSshServer_password demonstrates how to start a test SSH server
// with password authentication in a Kubernetes cluster and verify connectivity.
func Test_Example_StartTestSshServer_password(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	kubernetesCluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	if err != nil {
		require.NoError(t, err)
	}
	// ... prepare test environment finished.
	// -----

	// Define the namespace and pod name for the SSH server
	const namespaceName = "test-ssh-password"
	const podName = "test-ssh-server-password"
	const sshUsername = "testuser"
	const sshPassword = "testpassword123"

	// Start the test SSH server with password authentication
	sshServerPod, err := kubernetestestsshserver.StartTestSshServerInCluster(ctx, kubernetesCluster, &kubernetestestsshserver.StartTestSshServerOptions{
		KubernetesNamespace: namespaceName,
		PodName:             podName,
		SSHUsername:         sshUsername,
		SSHPassword:         sshPassword,
	})
	require.NoError(t, err)
	require.NotNil(t, sshServerPod)
	defer sshServerPod.Delete(ctx)

	// Verify the SSH server pod is running
	exists, err := sshServerPod.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Get the namespace
	namespace, err := sshServerPod.GetNamespace()
	require.NoError(t, err)

	// Deploy a test client pod in the same namespace to verify SSH connectivity
	testClientPodName := "test-ssh-client-password"
	testClientPod, err := namespace.CreatePod(ctx, &kubernetesparameteroptions.KubernetesRunCommandOptions{
		PodName: testClientPodName,
		Image:   "alpine:latest",
		RunCommandOptions: &parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", "sleep 300"},
		},
		DeleteAlreadyExistingPod: true,
	})
	require.NoError(t, err)
	require.NotNil(t, testClientPod)
	defer testClientPod.Delete(ctx)

	// Wait for the test client pod to be ready
	time.Sleep(5 * time.Second)

	// Execute SSH connection test from the client pod to the SSH server pod
	// The SSH server pod name resolves via Kubernetes DNS as: <pod-name>.<namespace>.svc.cluster.local
	sshServerHost := podName + "." + namespaceName + ".svc.cluster.local"
	testCommand := []string{"sh", "-c", "apk add --no-cache openssh-client sshpass && sshpass -p '" + sshPassword + "' ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 " + sshUsername + "@" + sshServerHost + " 'echo SSH_CONNECTION_SUCCESS'"}

	output, err := testClientPod.RunCommandInContainer(ctx, &kubernetesparameteroptions.KubernetesRunCommandOptions{
		ContainerName: testClientPodName,
		RunCommandOptions: &parameteroptions.RunCommandOptions{
			Command: testCommand,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, output.Stdout)
	require.Contains(t, string(*output.Stdout), "SSH_CONNECTION_SUCCESS")

	// Clean up: Delete the test client pod
	err = testClientPod.Delete(ctx)
	require.NoError(t, err)

	// Verify the test client pod is deleted
	exists, err = testClientPod.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Clean up: Delete the SSH server pod
	err = sshServerPod.Delete(ctx)
	require.NoError(t, err)

	// Verify the SSH server pod is deleted
	exists, err = sshServerPod.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)
}
