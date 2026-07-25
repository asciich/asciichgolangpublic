package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestNamespace_CreateAndDeleteNamespace(t *testing.T) {
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

				cluster := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				err := cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				for i := 0; i < 2; i++ {
					namespace, err := cluster.CreateNamespaceByName(ctx, namespaceName)
					require.NoError(t, err)

					exists, err := namespace.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					exists, err = cluster.NamespaceByNameExists(ctx, namespaceName)
					require.NoError(t, err)
					require.True(t, exists)
				}

				for i := 0; i < 2; i++ {
					err := cluster.DeleteNamespaceByName(ctx, namespaceName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

					namespace, err := cluster.GetNamespaceByName(namespaceName)
					require.NoError(t, err)

					exists, err := namespace.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
