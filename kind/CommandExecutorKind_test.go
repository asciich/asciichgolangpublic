package kind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandExeuctorKind_GetClusterByName(t *testing.T) {
	kind := MustGetLocalCommandExecutorKind()

	cluster := kind.MustGetClusterByName("abc")

	nativeCluster, ok := cluster.(*CommandExecutorKindCluster)
	require.True(t, ok)

	commandExecutor, err := nativeCluster.GetCommandExecutor()
	require.NoError(t, err)
	require.NotNil(t, commandExecutor)
}
