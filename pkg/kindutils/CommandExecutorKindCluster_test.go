package kindutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kindutils"
)

func TestCommandExecutorKindCluster_MustGetLocalCommandExecutorKind(t *testing.T) {
	k := kindutils.MustGetLocalCommandExecutorKind()

	kind, ok := k.(*kindutils.CommandExecutorKind)
	require.True(t, ok)

	c, err := kind.GetCommandExecutor()
	require.NoError(t, err)
	require.NotNil(t, c)
}
