package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestListKindNames(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctx := getCtx()

		// Ensure a local kind cluster is available for testing:
		clusterName := continuousintegration.GetDefaultKindClusterName()
		_, err := kindutils.CreateCluster(ctx, clusterName)
		require.NoError(t, err)
		defer kindutils.DeleteClusterByName(ctx, clusterName)

		cluster, err := nativekubernetes.GetDefaultCluster(ctx)
		require.NoError(t, err)

		apiVersions, err := cluster.ListKindNames(ctx)
		require.NoError(t, err)
		require.Contains(t, apiVersions, "Pod")
	})
}
