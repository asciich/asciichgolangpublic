package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/continuousintegration"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kindutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/nativekubernetes"
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
		defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)

		cluster, err := nativekubernetes.GetDefaultCluster(ctx)
		require.NoError(t, err)

		apiVersions, err := cluster.ListKindNames(ctx)
		require.NoError(t, err)
		require.Contains(t, apiVersions, "Pod")
	})
}
