package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetestestsshserver"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
)

// Test_Example_ValidateSSHKeyInSecret demonstrates how to validate if a Kubernetes secret
// contains a valid SSH private key by attempting to SSH into a target host.
func Test_Example_ValidateSSHKeyInSecret(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Generate an SSH key pair for testing
	keyPair, err := sshutils.GenerateKeyPair(sshutils.SSH_KEY_TYPE_ED25519, nil)
	require.NoError(t, err)

	err = keyPair.Validate(ctx)
	require.NoError(t, err)

	privateKey, err := keyPair.GetPrivateKey()
	require.NoError(t, err)

	// Use the shared Kind cluster
	kubernetes, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	const namespaceName = "example-ssh-validation"
	const secretName = "example-ssh-key"
	const secretKey = "id_ed25519"
	const sshServerPodName = "ssh-server-validate"
	const sshUsername = "testuser"

	namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)
	defer namespace.Delete(ctx)

	// Deploy a test SSH server with public key authentication in the cluster
	publicKeyString := keyPair.PublicKey.KeyMaterial

	sshServerPod, err := kubernetestestsshserver.StartTestSshServerInCluster(ctx, kubernetes, &kubernetestestsshserver.StartTestSshServerOptions{
		KubernetesNamespace: namespaceName,
		PodName:             sshServerPodName,
		SSHUsername:         sshUsername,
		SSHPublicKey:        publicKeyString,
	})
	require.NoError(t, err)
	require.NotNil(t, sshServerPod)
	defer sshServerPod.Delete(ctx)

	// Create a secret containing the private SSH key
	err = namespace.DeleteSecretByName(ctx, secretName)
	require.NoError(t, err)

	_, err = namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{
			secretKey: []byte(privateKey.KeyMaterial),
		},
	})
	require.NoError(t, err)

	defer func() {
		_ = namespace.DeleteSecretByName(ctx, secretName)
	}()

	exists, err := namespace.SecretByNameExists(ctx, secretName)
	require.NoError(t, err)
	require.True(t, exists)

	// Validate the SSH key in the secret by connecting to the test SSH server
	// The SSH server is reachable via its Service DNS: <podName>.<namespace>.svc.cluster.local
	sshServerHost := sshServerPodName + "." + namespaceName + ".svc.cluster.local"

	options := &kubernetesparameteroptions.ValidateSshKeyInSecretOptions{
		Namespace:             namespaceName,
		SecretName:            secretName,
		SecretKey:             secretKey,
		TargetHost:            sshServerHost,
		TargetUser:            sshUsername,
		TargetPort:            22,
		SkipHostKeyValidation: true,
		ConnectionTimeout:     "10 seconds",
		ConnectionAttempts:    1,
	}
	success, err := kubernetes.ValidateSSHKeyInSecret(ctx, options)
	require.NoError(t, err)
	require.True(t, success)

	err = kubernetes.DeleteNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)
}
