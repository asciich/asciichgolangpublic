package kindutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
)

func TestCommandExecutorKindCluster_GetLocalCommandExecutorKind(t *testing.T) {
	k, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	kind, ok := k.(*kindutils.CommandExecutorKind)
	require.True(t, ok)

	c, err := kind.GetCommandExecutor()
	require.NoError(t, err)
	require.NotNil(t, c)
}

func Test_LocalKindCluster(t *testing.T) {
	kindCluster, err := kindutils.GetLocalKindCluster("kindcluster")
	require.NoError(t, err)

	var k8sCluster kubernetesinterfaces.KubernetesCluster = kindCluster

	name, err := k8sCluster.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "kind-kindcluster", name)
}

func TestGetCommandExecutorKindCluster_NilCommandExecutor(t *testing.T) {
	cluster, err := kindutils.GetCommandExecutorKindCluster(nil, "test-cluster")
	require.Error(t, err)
	require.Nil(t, cluster)
	require.Contains(t, err.Error(), "commandExectuor")
}

func TestGetCommandExecutorKindCluster_EmptyClusterName(t *testing.T) {
	kind, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	commandExecutorKind, ok := kind.(*kindutils.CommandExecutorKind)
	require.True(t, ok)

	commandExecutor, err := commandExecutorKind.GetCommandExecutor()
	require.NoError(t, err)

	cluster, err := kindutils.GetCommandExecutorKindCluster(commandExecutor, "")
	require.Error(t, err)
	require.Nil(t, cluster)
	require.Contains(t, err.Error(), "clusterName")
}

func TestGetCommandExecutorKindCluster_ValidInputs(t *testing.T) {
	kind, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	commandExecutorKind, ok := kind.(*kindutils.CommandExecutorKind)
	require.True(t, ok)

	commandExecutor, err := commandExecutorKind.GetCommandExecutor()
	require.NoError(t, err)

	cluster, err := kindutils.GetCommandExecutorKindCluster(commandExecutor, "test-cluster")
	require.NoError(t, err)
	require.NotNil(t, cluster)

	var k8sCluster kubernetesinterfaces.KubernetesCluster = cluster
	name, err := k8sCluster.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "kind-test-cluster", name)

	// Verify it's the correct type
	kindCluster, ok := cluster.(*kindutils.CommandExecutorKindCluster)
	require.True(t, ok)

	// Verify cluster name is set
	clusterName, err := kindCluster.GetClusterName()
	require.NoError(t, err)
	require.EqualValues(t, "test-cluster", clusterName)

	// Verify kind is set
	retrievedKind, err := kindCluster.GetKind()
	require.NoError(t, err)
	require.NotNil(t, retrievedKind)
}

func TestGetCommandExecutorKindCluster_WithDifferentClusterNames(t *testing.T) {
	kind, err := kindutils.GetLocalCommandExecutorKind()
	require.NoError(t, err)

	commandExecutorKind, ok := kind.(*kindutils.CommandExecutorKind)
	require.True(t, ok)

	commandExecutor, err := commandExecutorKind.GetCommandExecutor()
	require.NoError(t, err)

	testCases := []string{
		"cluster1",
		"my-cluster",
		"test_cluster_123",
	}

	for _, clusterName := range testCases {
		t.Run(clusterName, func(t *testing.T) {
			cluster, err := kindutils.GetCommandExecutorKindCluster(commandExecutor, clusterName)
			require.NoError(t, err)
			require.NotNil(t, cluster)

			kindCluster, ok := cluster.(*kindutils.CommandExecutorKindCluster)
			require.True(t, ok)

			retrievedName, err := kindCluster.GetClusterName()
			require.NoError(t, err)
			require.EqualValues(t, clusterName, retrievedName)
		})
	}
}
