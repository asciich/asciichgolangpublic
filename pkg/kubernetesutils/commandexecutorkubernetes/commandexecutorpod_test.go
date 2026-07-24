package commandexecutorkubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils/kindparameteroptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_NonExistingPod(t *testing.T) {
	ctx := getCtx()

	kind, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	cluster, err := kind.CreateCluster(ctx, &kindparameteroptions.CreateClusterOptions{
		Name: "kubernetesutils",
	})
	require.NoError(t, err)

	exists, err := cluster.PodByNameExists(ctx, "default", "this-pod-does-not-exist")
	require.NoError(t, err)
	require.False(t, exists)
}
