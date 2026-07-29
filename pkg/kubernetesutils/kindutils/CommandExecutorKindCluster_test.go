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
	require.EqualValues(t, "kindcluster", name)
}
