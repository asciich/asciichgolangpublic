package kindutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getClusterName() (clusterName string) {
	return "kind-ci-test"
}

func getKindByImplementationName(implementationName string) (kind kindutils.Kind) {
	if implementationName == "commandExecutorKind" {
		return kindutils.MustGetLocalCommandExecutorKind()
	}

	logging.LogFatalWithTracef(
		"Unknown implmentation name '%s'",
		implementationName,
	)

	return nil
}

func TestKind_CreateAndDeleteCluster(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKind"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				clusterName := getClusterName()

				kind := getKindByImplementationName(tt.implementationName)

				err := kind.DeleteClusterByName(ctx, clusterName)
				require.NoError(t, err)

				require.False(t, mustutils.Must(kind.ClusterByNameExists(ctx, clusterName)))

				for i := 0; i < 2; i++ {
					_, err := kind.CreateClusterByName(ctx, clusterName)
					require.NoError(t, err)

					require.True(t, mustutils.Must(kind.ClusterByNameExists(ctx, clusterName)))
				}

				for i := 0; i < 2; i++ {
					err = kind.DeleteClusterByName(ctx, clusterName)
					require.NoError(t, err)

					require.False(t, mustutils.Must(kind.ClusterByNameExists(ctx, clusterName)))
				}
			},
		)
	}
}

func TestKind_CreateNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKind"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true
				ctx := getCtx()
				clusterName := getClusterName()

				kind := getKindByImplementationName(tt.implementationName)

				cluster, err := kind.CreateClusterByName(ctx, clusterName)
				require.NoError(t, err)

				namespaceName := "test-namespace"

				err = cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				for i := 0; i < 2; i++ {
					_, err := cluster.CreateNamespaceByName(ctx, namespaceName)
					require.NoError(t, err)
					require.True(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))
				}

				for i := 0; i < 2; i++ {
					err := cluster.DeleteNamespaceByName(ctx, namespaceName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))
				}
			},
		)
	}
}
