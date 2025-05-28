package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/testutils"
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

				secret, err := namespace.CreateSecret(ctx, secretName, &kubernetesutils.CreateSecretOptions{SecretData: map[string][]byte{}})
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
