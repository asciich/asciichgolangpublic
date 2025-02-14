package kind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandExecutorKindCluster_MustGetLocalCommandExecutorKind(t *testing.T) {
	k := MustGetLocalCommandExecutorKind()

	kind, ok := k.(*CommandExecutorKind)
	require.True(t, ok)

	c, err := kind.GetCommandExecutor()
	require.NoError(t, err)
	require.NotNil(t, c)
}
