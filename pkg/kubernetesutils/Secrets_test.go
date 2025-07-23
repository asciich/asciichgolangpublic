package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_SecretByNameExists(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		// {"commandExecutorKubernetes"},
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

func Test_CreateSecretInNonExistentNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		// {"commandExecutorKubernetes"},
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
