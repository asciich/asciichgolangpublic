package kindutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils/kindparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
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
				clusterName := continuousintegration.GetDefaultKindClusterName()

				kind := getKindByImplementationName(tt.implementationName)

				err := kind.DeleteClusterByName(ctx, clusterName)
				require.NoError(t, err)

				require.False(t, mustutils.Must(kind.ClusterByNameExists(ctx, clusterName)))

				for i := 0; i < 2; i++ {
					kindCluster, err := kind.CreateCluster(ctx, &kindparameteroptions.CreateClusterOptions{Name: clusterName})
					require.NoError(t, err)

					require.True(t, mustutils.Must(kind.ClusterByNameExists(ctx, clusterName)))

					require.EqualValues(t, "kind-"+clusterName, mustutils.Must(kindCluster.GetName()))
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
				clusterName := continuousintegration.GetDefaultKindClusterName()
				defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)

				kind := getKindByImplementationName(tt.implementationName)

				cluster, err := kind.CreateCluster(ctx, &kindparameteroptions.CreateClusterOptions{Name: clusterName})
				require.NoError(t, err)
				defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)

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

				// cleanup
				err = kind.DeleteClusterByName(ctx, clusterName)
				require.NoError(t, err)
			},
		)
	}
}

func TestKind_CreateMultipleWorkers(t *testing.T) {
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
				clusterName := continuousintegration.GetDefaultKindClusterName()
				defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)

				kind := getKindByImplementationName(tt.implementationName)

				// Default is to create a cluster with only one node/ no workers.
				cluster, err := kind.CreateCluster(ctx, &kindparameteroptions.CreateClusterOptions{Name: clusterName})
				require.NoError(t, err)
				defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)

				nodeNames, err := cluster.ListNodeNames(ctx)
				require.NoError(t, err)
				require.Len(t, nodeNames, 1)

			},
		)
	}
}
