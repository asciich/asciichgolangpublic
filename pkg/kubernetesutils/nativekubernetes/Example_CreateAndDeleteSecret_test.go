package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

// This example shows how to create and delete a Secret.
func Test_Example_CreateAndDeleteSecret(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Define the name of the Secret we use for testing
	const secretName = "example-secret"
	const namespaceName = "default"

	secretData := map[string][]byte{
		"username": []byte("admin"),
		"password": []byte("secret123"),
	}

	// Ensure the Secret is absent:
	err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
	require.NoError(t, err)

	secretNames, err := nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, secretNames, secretName)

	// Create the Secret:
	err = nativekubernetes.CreateSecret(ctx, clientset, namespaceName, secretName, secretData)
	require.NoError(t, err)

	secretNames, err = nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.Contains(t, secretNames, secretName)

	// Read the Secret and verify the data:
	readData, err := nativekubernetes.ReadSecret(ctx, clientset, namespaceName, secretName)
	require.NoError(t, err)
	require.Equal(t, secretData, readData)

	// The create function is idempotent:
	err = nativekubernetes.CreateSecret(ctx, clientset, namespaceName, secretName, secretData)
	require.NoError(t, err)

	secretNames, err = nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.Contains(t, secretNames, secretName)

	// Delete the Secret again so we have no leftovers:
	err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
	require.NoError(t, err)

	secretNames, err = nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, secretNames, secretName)

	// The delete function is idempotent as well.
	// If the Secret is already absent no error will be raised.
	err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
	require.NoError(t, err)

	secretNames, err = nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, secretNames, secretName)
}
