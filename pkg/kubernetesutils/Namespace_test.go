package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/commandexecutorkubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getKubernetesByImplementationName(ctx context.Context, implementationName string) kubernetesinterfaces.KubernetesCluster {
	clusterName := "kubernetesutils" // We use one kind cluster for all the tests here.

	if implementationName == "commandExecutorKubernetes" {
		// Ensure a local kind cluster is available for testing:
		mustutils.Must(kindutils.CreateCluster(ctx, clusterName))

		return mustutils.Must(commandexecutorkubernetes.GetClusterByName("kind-" + clusterName))
	}

	if implementationName == "nativeKubernetes" {
		// Ensure a local kind cluster is available for testing:
		mustutils.Must(kindutils.CreateCluster(ctx, clusterName))

		return mustutils.Must(nativekubernetes.GetClusterByName(getCtx(), "kind-"+clusterName))

	}

	logging.LogFatalWithTracef(
		"Unknown implmentation name '%s'",
		implementationName,
	)

	return nil
}

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

				cluster := getKubernetesByImplementationName(getCtx(), tt.implementationName)

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
