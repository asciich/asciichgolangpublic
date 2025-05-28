package kubernetesutils_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_Example_SecretByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Prepare start ...
	const clusterName = "kind"

	// Ensure a local kind cluster is available for testing:
	_, err := commandexecutor.Bash().RunOneLiner(ctx, fmt.Sprintf("kind create cluster -n '%s' || true", clusterName))
	require.NoError(t, err)

	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, clusterName)
	require.NoError(t, err)

	// Create an example secret. This implicitly generates the namespace if it does not exist.
	const namespaceName = "testnamespace"
	const secretName = "example-secret"
	_, err = cluster.CreateSecret(ctx, namespaceName, secretName, &kubernetesutils.CreateSecretOptions{
		SecretData: map[string][]byte{"my-secret": []byte("very-secret")},
	})
	require.NoError(t, err)
	// ... prepare finished.

	// Our created secret exists:
	exists, err := cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.True(t, exists)

	// The same secret name in the default namespace does not exist:
	exists, err = cluster.SecretByNameExists(ctx, "default", secretName)
	require.NoError(t, err)
	require.False(t, exists)

	// This secret is expected to be in the same namespace but does not exist:
	exists, err = cluster.SecretByNameExists(ctx, namespaceName, "secret-does-not-exist")
	require.NoError(t, err)
	require.False(t, exists)

	// If we delete our secret again...
	err = cluster.DeleteSecretByName(ctx, namespaceName, secretName)
	require.NoError(t, err)

	// ... our secret becomes absent:
	exists, err = cluster.SecretByNameExists(ctx, namespaceName, secretName)
	require.NoError(t, err)
	require.False(t, exists)
}
