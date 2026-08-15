package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// Test_Example_ValidateSSHKeyInSecret demonstrates how to validate if a Kubernetes secret
// contains a valid SSH private key by attempting to SSH into a target host.
func Test_Example_ValidateSSHKeyInSecret(t *testing.T) {
	ctx := contextutils.WithVerbose(context.TODO())

	// Generate an SSH key pair for testing
	keyPair, err := sshutils.GenerateKeyPair("", nil)
	require.NoError(t, err)

	err = keyPair.Validate(ctx)
	require.NoError(t, err)

	privateKey, err := keyPair.GetPrivateKey()
	require.NoError(t, err)

	clusterName := testutils.GetKindClusterNameForTest(t)
	_, err = kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	kubernetes, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	const namespaceName = "example-ssh-validation"
	const secretName = "example-ssh-key"
	const secretKey = "id_ed25519"

	namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)
	defer namespace.Delete(ctx)

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

	options := &kubernetesparameteroptions.ValidateSshKeyInSecretOptions{
		Namespace:             namespaceName,
		SecretName:            secretName,
		SecretKey:             secretKey,
		TargetHost:            "example-host.com",
		TargetUser:            "username",
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
