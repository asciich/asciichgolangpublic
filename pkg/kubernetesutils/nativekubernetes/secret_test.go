package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_CreateAndDeleteSecret(t *testing.T) {
	ctx := getCtx()

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

	t.Run("happy path", func(t *testing.T) {
		const secretName = "testsecret"
		const namespaceName = "default"

		secretData := map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret123"),
		}

		err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
		require.NoError(t, err)

		secretNames, err := nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
		require.NoError(t, err)
		require.NotContains(t, secretNames, secretName)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreateSecret(ctx, clientset, namespaceName, secretName, secretData)
			require.NoError(t, err)

			secretNames, err = nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
			require.NoError(t, err)
			require.Contains(t, secretNames, secretName)

			readData, err := nativekubernetes.ReadSecret(ctx, clientset, namespaceName, secretName)
			require.NoError(t, err)
			require.Equal(t, secretData, readData)

			err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
			require.NoError(t, err)

			secretNames, err = nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
			require.NoError(t, err)
			require.NotContains(t, secretNames, secretName)
		}
	})
}

func Test_DeleteSecretAlreadyAbsent(t *testing.T) {
	ctx := getCtx()

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

	t.Run("already deleted secret", func(t *testing.T) {
		const secretName = "testsecret-absent"
		const namespaceName = "default"

		// Ensure secret is absent
		err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
		require.NoError(t, err)

		// Deleting an already absent secret should not error:
		err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
		require.NoError(t, err)
	})
}

func Test_ReadSecretNotFound(t *testing.T) {
	ctx := getCtx()

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

	t.Run("read non-existent secret returns error", func(t *testing.T) {
		const secretName = "nonexistent-secret"
		const namespaceName = "default"

		// Ensure secret does not exist
		err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, secretName)
		require.NoError(t, err)

		_, err = nativekubernetes.ReadSecret(ctx, clientset, namespaceName, secretName)
		require.Error(t, err)
	})
}

func Test_ListSecrets(t *testing.T) {
	ctx := getCtx()

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

	t.Run("create and delete secrets with list in between", func(t *testing.T) {
		const namespaceName = "default"

		secretNames := []string{"listsecret-1", "listsecret-2", "listsecret-3"}
		secretData := map[string][]byte{"key": []byte("value")}

		// Ensure all test secrets are absent before starting
		for _, name := range secretNames {
			err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, name)
			require.NoError(t, err)
		}

		// Create secrets one by one and verify list grows
		for i, name := range secretNames {
			err = nativekubernetes.CreateSecret(ctx, clientset, namespaceName, name, secretData)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, created := range secretNames[:i+1] {
				require.Contains(t, listed, created)
			}
			for _, notYetCreated := range secretNames[i+1:] {
				require.NotContains(t, listed, notYetCreated)
			}
		}

		// Delete secrets one by one and verify list shrinks
		for i, name := range secretNames {
			err = nativekubernetes.DeleteSecret(ctx, clientset, namespaceName, name)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListSecrets(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, deleted := range secretNames[:i+1] {
				require.NotContains(t, listed, deleted)
			}
			for _, stillPresent := range secretNames[i+1:] {
				require.Contains(t, listed, stillPresent)
			}
		}
	})
}
