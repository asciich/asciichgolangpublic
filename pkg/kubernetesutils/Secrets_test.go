package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kuberneteserrors"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_SecretByNameExists(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const secretName = "secretname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = namespace.DeleteSecretByName(ctx, secretName)
				require.NoError(t, err)

				exists, err := namespace.SecretByNameExists(ctx, secretName)
				require.NoError(t, err)
				require.False(t, exists)

				secret, err := namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{SecretData: map[string][]byte{}})
				require.NoError(t, err)

				exists, err = secret.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				exists, err = namespace.SecretByNameExists(ctx, secretName)
				require.NoError(t, err)
				require.True(t, exists)

				for i := 0; i < 2; i++ {
					err = namespace.DeleteSecretByName(ctx, secretName)
					require.NoError(t, err)

					exists, err := namespace.SecretByNameExists(ctx, secretName)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}

func Test_GetSecret_ErrorIfNotExist(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const secretName = "secretname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = namespace.DeleteSecretByName(ctx, secretName)
				require.NoError(t, err)

				secret, err := namespace.GetSecretByName(secretName)
				require.NoError(t, err)

				got, err := secret.Read(ctx)
				require.ErrorIs(t, err, kuberneteserrors.ErrSecretNotFound)
				require.True(t, kuberneteserrors.IsSecretNotFoundError(err))
				require.Nil(t, got)
			},
		)
	}
}

func Test_CreateSecretInNonExistentNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const secretName = "secretname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				// ensure namespace is absent:
				err := kubernetes.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				exists, err := kubernetes.NamespaceByNameExists(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, exists)

				// create the secret in the absent namespace:
				_, err = kubernetes.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{SecretData: map[string][]byte{"my-secret": []byte("value")}})
				require.NoError(t, err)

				// Namespace is implicitly generated:
				exists, err = kubernetes.NamespaceByNameExists(ctx, namespaceName)
				require.NoError(t, err)
				require.True(t, exists)

				// and secret is generated as well:
				exists, err = kubernetes.SecretByNameExists(ctx, namespaceName, secretName)
				require.NoError(t, err)
				require.True(t, exists)
			},
		)
	}
}

func Test_SecretReadWriteUpdate(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const secretName = "secretname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Ensure secret is absent before starting:
				err = namespace.DeleteSecretByName(ctx, secretName)
				require.NoError(t, err)

				// --- Write: create secret with initial data ---
				initialData := map[string][]byte{
					"username": []byte("admin"),
					"password": []byte("initial-password"),
				}

				secret, err := namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: initialData,
				})
				require.NoError(t, err)

				// --- Read: verify initial data is stored correctly ---
				gotData, err := secret.Read(ctx)
				require.NoError(t, err)
				require.Equal(t, initialData, gotData)

				// --- Update: overwrite secret with new data ---
				updatedData := map[string][]byte{
					"username": []byte("admin"),
					"password": []byte("updated-password"),
				}

				secret, err = namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: updatedData,
				})
				require.NoError(t, err)

				// --- Read: verify updated data is stored correctly ---
				gotData, err = secret.Read(ctx)
				require.NoError(t, err)
				require.Equal(t, updatedData, gotData)

				// --- Read: verify old data is no longer present ---
				require.NotEqual(t, initialData, gotData)

				// --- Update: idempotency check, applying same data again should not fail ---
				secret, err = namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: updatedData,
				})
				require.NoError(t, err)

				gotData, err = secret.Read(ctx)
				require.NoError(t, err)
				require.Equal(t, updatedData, gotData)

				// Cleanup:
				err = namespace.DeleteSecretByName(ctx, secretName)
				require.NoError(t, err)

				exists, err := namespace.SecretByNameExists(ctx, secretName)
				require.NoError(t, err)
				require.False(t, exists)
			},
		)
	}
}

func Test_ListSecretNames(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const secretNameA = "secretname-a"
				const secretNameB = "secretname-b"
				const secretNameC = "secretname-c"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Ensure all secrets are absent before starting:
				err = namespace.DeleteSecretByName(ctx, secretNameA)
				require.NoError(t, err)
				err = namespace.DeleteSecretByName(ctx, secretNameB)
				require.NoError(t, err)
				err = namespace.DeleteSecretByName(ctx, secretNameC)
				require.NoError(t, err)

				// --- List: empty namespace should return no secrets (or not contain our secrets) ---
				names, err := namespace.ListSecretNames(ctx)
				require.NoError(t, err)
				require.NotContains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.NotContains(t, names, secretNameC)

				// --- Create first secret and verify it appears in the list ---
				_, err = namespace.CreateSecret(ctx, secretNameA, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{"key": []byte("value-a")},
				})
				require.NoError(t, err)

				names, err = namespace.ListSecretNames(ctx)
				require.NoError(t, err)
				require.Contains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.NotContains(t, names, secretNameC)

				// --- Create second and third secret and verify all appear in the list ---
				_, err = namespace.CreateSecret(ctx, secretNameB, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{"key": []byte("value-b")},
				})
				require.NoError(t, err)

				_, err = namespace.CreateSecret(ctx, secretNameC, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{"key": []byte("value-c")},
				})
				require.NoError(t, err)

				names, err = namespace.ListSecretNames(ctx)
				require.NoError(t, err)
				require.Contains(t, names, secretNameA)
				require.Contains(t, names, secretNameB)
				require.Contains(t, names, secretNameC)

				// --- Delete one secret and verify it no longer appears in the list ---
				err = namespace.DeleteSecretByName(ctx, secretNameB)
				require.NoError(t, err)

				names, err = namespace.ListSecretNames(ctx)
				require.NoError(t, err)
				require.Contains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.Contains(t, names, secretNameC)

				// Cleanup:
				err = namespace.DeleteSecretByName(ctx, secretNameA)
				require.NoError(t, err)
				err = namespace.DeleteSecretByName(ctx, secretNameC)
				require.NoError(t, err)

				names, err = namespace.ListSecretNames(ctx)
				require.NoError(t, err)
				require.NotContains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.NotContains(t, names, secretNameC)
			},
		)
	}
}

func Test_ListSecrets(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const secretNameA = "secretname-a"
				const secretNameB = "secretname-b"
				const secretNameC = "secretname-c"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Ensure all secrets are absent before starting:
				err = namespace.DeleteSecretByName(ctx, secretNameA)
				require.NoError(t, err)
				err = namespace.DeleteSecretByName(ctx, secretNameB)
				require.NoError(t, err)
				err = namespace.DeleteSecretByName(ctx, secretNameC)
				require.NoError(t, err)

				getSecretNames := func(ctx context.Context) []string {
					secrets, err := namespace.ListSecrets(ctx)
					require.NoError(t, err)
					names := []string{}
					for _, secret := range secrets {
						name, err := secret.GetName()
						require.NoError(t, err)
						names = append(names, name)
					}
					return names
				}

				// --- List: empty namespace should return no secrets (or not contain our secrets) ---
				names := getSecretNames(ctx)
				require.NotContains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.NotContains(t, names, secretNameC)

				// --- Create first secret and verify it appears in the list ---
				_, err = namespace.CreateSecret(ctx, secretNameA, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{"key": []byte("value-a")},
				})
				require.NoError(t, err)

				names = getSecretNames(ctx)
				require.Contains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.NotContains(t, names, secretNameC)

				// --- Create second and third secret and verify all appear in the list ---
				_, err = namespace.CreateSecret(ctx, secretNameB, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{"key": []byte("value-b")},
				})
				require.NoError(t, err)

				_, err = namespace.CreateSecret(ctx, secretNameC, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{"key": []byte("value-c")},
				})
				require.NoError(t, err)

				names = getSecretNames(ctx)
				require.Contains(t, names, secretNameA)
				require.Contains(t, names, secretNameB)
				require.Contains(t, names, secretNameC)

				// --- Delete one secret and verify it no longer appears in the list ---
				err = namespace.DeleteSecretByName(ctx, secretNameB)
				require.NoError(t, err)

				names = getSecretNames(ctx)
				require.Contains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.Contains(t, names, secretNameC)

				// Cleanup:
				err = namespace.DeleteSecretByName(ctx, secretNameA)
				require.NoError(t, err)
				err = namespace.DeleteSecretByName(ctx, secretNameC)
				require.NoError(t, err)

				names = getSecretNames(ctx)
				require.NotContains(t, names, secretNameA)
				require.NotContains(t, names, secretNameB)
				require.NotContains(t, names, secretNameC)
			},
		)
	}
}
