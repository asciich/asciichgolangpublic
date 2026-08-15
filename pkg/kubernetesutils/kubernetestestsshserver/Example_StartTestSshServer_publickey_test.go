package kubernetestestsshserver_test

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetestestsshserver"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
)

// Test_Example_StartTestSshServer_publickey demonstrates how to start a test SSH server
// with public key authentication in a Kubernetes cluster and verify connectivity.
func Test_Example_StartTestSshServer_publickey(t *testing.T) {
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

	// Generate an SSH key pair for testing using sshutils
	keyPair, err := sshutils.GenerateKeyPair(sshutils.SSH_KEY_TYPE_ED25519, nil)
	require.NoError(t, err)

	// Define the namespace and pod name for the SSH server
	const namespaceName = "test-ssh-publickey"
	const podName = "test-ssh-server-publickey"
	const sshUsername = "keyuser"

	// Use the public key material directly (it already contains the key type prefix)
	publicKeyString := keyPair.PublicKey.KeyMaterial

	// Start the test SSH server with public key authentication
	sshServerPod, err := kubernetestestsshserver.StartTestSshServerInCluster(ctx, kubernetesCluster, &kubernetestestsshserver.StartTestSshServerOptions{
		KubernetesNamespace: namespaceName,
		PodName:             podName,
		SSHUsername:         sshUsername,
		SSHPublicKey:        publicKeyString,
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
	testClientPodName := "test-ssh-client-publickey"
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

	// Wait for the test client pod and SSH server to be fully ready
	time.Sleep(10 * time.Second)

	// Base64 encode the private key to avoid shell escaping issues
	privateKeyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPair.PrivateKey.KeyMaterial))
	sshServerHost := podName + "." + namespaceName + ".svc.cluster.local"

	// Create a script that decodes the key from base64 and connects via SSH
	testCommand := []string{"sh", "-c",
		"apk add --no-cache openssh-client && " +
			"mkdir -p ~/.ssh && " +
			"echo '" + privateKeyBase64 + "' | base64 -d > ~/.ssh/id_ed25519 && " +
			"chmod 600 ~/.ssh/id_ed25519 && " +
			"ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 -i ~/.ssh/id_ed25519 " + sshUsername + "@" + sshServerHost + " 'echo SSH_CONNECTION_SUCCESS'"}

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
