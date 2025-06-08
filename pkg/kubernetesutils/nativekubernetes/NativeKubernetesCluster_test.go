package nativekubernetes_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestListKindNames(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctx := getCtx()

		// Ensure a local kind cluster is available for testing:
		_, err := commandexecutor.Bash().RunOneLiner(ctx, "kind create cluster -n 'kind' || true")
		require.NoError(t, err)
		time.Sleep(1 * time.Second)

		cluster, err := nativekubernetes.GetDefaultCluster(ctx)
		require.NoError(t, err)

		apiVersions, err := cluster.ListKindNames(ctx)
		require.NoError(t, err)
		require.Contains(t, apiVersions, "Pod")
	})
}
