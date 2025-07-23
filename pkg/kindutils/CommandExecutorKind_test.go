package kindutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kindutils"
)

func TestCommandExeuctorKind_GetClusterByName(t *testing.T) {
	kind := kindutils.MustGetLocalCommandExecutorKind()

	cluster, err := kind.GetClusterByName("abc")
	require.NoError(t, err)

	nativeCluster, ok := cluster.(*kindutils.CommandExecutorKindCluster)
	require.True(t, ok)

	commandExecutor, err := nativeCluster.GetCommandExecutor()
	require.NoError(t, err)
	require.NotNil(t, commandExecutor)
}
