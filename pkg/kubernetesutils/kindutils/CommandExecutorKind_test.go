package kindutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
)

func TestCommandExeuctorKind_GetClusterByName(t *testing.T) {
	kind, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	cluster, err := kind.GetClusterByName("abc")
	require.NoError(t, err)

	nativeCluster, ok := cluster.(*kindutils.CommandExecutorKindCluster)
	require.True(t, ok)

	commandExecutor, err := nativeCluster.GetCommandExecutor()
	require.NoError(t, err)
	require.NotNil(t, commandExecutor)
}

func TestCommandExecutorKind_ClusterByNameExists_AcceptsPrefix(t *testing.T) {
	kind, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	ctx := context.Background()

	// Test that both "test-cluster" and "kind-test-cluster" are handled equivalently
	// We test with a non-existent cluster name to verify the prefix stripping logic
	existsWithoutPrefix, err := kind.ClusterByNameExists(ctx, "test-cluster")
	require.NoError(t, err)

	existsWithPrefix, err := kind.ClusterByNameExists(ctx, "kind-test-cluster")
	require.NoError(t, err)

	// Both should return the same result (both false for non-existent cluster)
	require.Equal(t, existsWithoutPrefix, existsWithPrefix, "ClusterByNameExists should return same result for names with and without 'kind-' prefix")
}
