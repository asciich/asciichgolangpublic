package headscaleutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalelocaldevserver"
)

// This test sets up a headscale development server with one user.
// Then two containers were started as tailscale clients/ nodes.
// After connecting the two nodes using the same user a `tailscale ping` is performed to check connectivity.
func Test_ConnectTwoNodesWithOneUser(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Restart the headscale dev server:
	headscale, cancel, err := headscalelocaldevserver.RunLocalDevServer(ctx, &headscalelocaldevserver.RunOptions{
		RestartAlreadyRunningDevServer: true,
	})
	require.NoError(t, err)
	defer cancel()

	// create the test user:
	const userName = "testuser"
	err = headscale.CreateUser(ctx, userName)
	require.NoError(t, err)

	// Get two keys to connect two nodes:
	_, err = headscale.GeneratePreauthKeyForUser(ctx, userName)
	require.NoError(t, err)

	_, err = headscale.GeneratePreauthKeyForUser(ctx, userName)
	require.NoError(t, err)
}
