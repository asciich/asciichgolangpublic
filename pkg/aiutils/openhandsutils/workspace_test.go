package openhandsutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateDeleteAndListWorkspace(t *testing.T) {
	ctx := getCtx()
	openHands := startTestOpenHands(ctx, t)

	const workspaceName = "workspace"

	err := openHands.DeleteAllWorkspaces(ctx)
	require.NoError(t, err)

	exists, err := openHands.WorkspaceExists(ctx, workspaceName)
	require.NoError(t, err)
	require.False(t, exists)

	workspaceNames, err := openHands.ListWorkspaceNames(ctx)
	require.NoError(t, err)
	require.Empty(t, workspaceNames)

	for range 2 { // run twice to check idempotence
		err = openHands.CreateWorkspace(ctx, workspaceName, "/workspace")
		require.NoError(t, err)

		workspaceNames, err = openHands.ListWorkspaceNames(ctx)
		require.NoError(t, err)
		require.EqualValues(t, []string{workspaceName}, workspaceNames)

		exists, err = openHands.WorkspaceExists(ctx, workspaceName)
		require.NoError(t, err)
		require.True(t, exists)
	}
}
