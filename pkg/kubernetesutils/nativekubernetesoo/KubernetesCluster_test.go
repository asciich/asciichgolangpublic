package nativekubernetesoo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestListKindNames(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctx := getCtx()

		// Ensure a local kind cluster is available for testing:
		_, err := kindutils.CreateCluster(ctx, testClusterName)

		require.NoError(t, err)

		cluster, err := nativekubernetesoo.GetDefaultCluster(ctx)
		require.NoError(t, err)

		apiVersions, err := cluster.ListKindNames(ctx)
		require.NoError(t, err)
		require.Contains(t, apiVersions, "Pod")
	})
}
