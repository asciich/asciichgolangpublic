package nativetailscaleclientoo_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalelocaldevserver"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscaleclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_ConnectAsNode(t *testing.T) {
	continuousintegration.SkipInGithubCi(t, "Expose docker port does not work yet.")

	ctx := getCtx()

	const headScalePort = 11235

	headscale, cancel, err := headscalelocaldevserver.RunLocalDevServer(ctx, &headscalelocaldevserver.RunOptions{Port: headScalePort, RestartAlreadyRunningDevServer: true})
	require.NoError(t, err)
	defer cancel()

	err = headscale.CreateUser(ctx, "testuser")
	require.NoError(t, err)
	preAuthKey, err := headscale.GeneratePreauthKeyForUser(ctx, "testuser")
	require.NoError(t, err)

	const hostname = "nativeclient1"
	client, cancel, err := nativetailscaleclientoo.Connect(ctx, &tailscaleoptions.ConnectOptions{
		HostName:   hostname,
		PreAuthKey: preAuthKey,
		ControlURL: "http://localhost:" + strconv.Itoa(headScalePort),
	})
	require.NoError(t, err)
	defer cancel()

	isConnected, err := client.IsConnected(ctx)
	require.NoError(t, err)
	require.True(t, isConnected)

	nodesNames, err := headscale.ListNodeNames(ctx)
	require.NoError(t, err)
	require.Len(t, nodesNames, 1)
	require.Contains(t, nodesNames, "nativeclient1")
}
