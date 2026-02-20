package tailscaleutils

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalelocaldevserver"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscalehttpclient"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscalehttpserver"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
)

// This example shows how to connect two tailscale nodes.
// One is serving as HTTP server, the other one acts as a HTTP client.
//
// For the local setup a headscale running in a local docker container is uesd so this testcase runs independendly of the public control plane servers of tailscale.
func Test_Example_ServeAndRequestHttpOverTailscale_test(t *testing.T) {
	// Currently not available in Github CI:
	continuousintegration.SkipInGithubCi(t, "Expose docker port does not work yet.")

	// Use verbose output:
	ctx := contextutils.ContextVerbose()

	// Start the headscale test server:
	const headScalePort = 11236

	headscale, cancel, err := headscalelocaldevserver.RunLocalDevServer(ctx, &headscalelocaldevserver.RunOptions{Port: headScalePort, RestartAlreadyRunningDevServer: true})
	require.NoError(t, err)
	defer cancel()

	// Add a user and generate two preauth keys for this user:
	err = headscale.CreateUser(ctx, "testuser")
	require.NoError(t, err)

	preauthKeyServer, err := headscale.GeneratePreauthKeyForUser(ctx, "testuser")
	require.NoError(t, err)

	preauthKeyClient, err := headscale.GeneratePreauthKeyForUser(ctx, "testuser")
	require.NoError(t, err)

	// Define a http.Mux to build our HTTP server logic.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		logging.LogInfoByCtxf(ctx, "/hello requested")
		fmt.Fprintf(w, "Hello from Tailscale! You are connected via tsnet.\n")
	})
	mux.HandleFunc("GET /world", func(w http.ResponseWriter, r *http.Request) {
		logging.LogInfoByCtxf(ctx, "/world requested")
		fmt.Fprintf(w, "Another endpoint")
	})

	// Start a HTTP server serving on tailscale port 80
	server, err := nativetailscalehttpserver.StartHttpServer(
		ctx,
		mux,
		80,
		&tailscaleoptions.ConnectOptions{
			HostName:   "webserver",
			PreAuthKey: preauthKeyServer,
			ControlURL: "http://localhost:" + strconv.Itoa(headScalePort),
		},
	)
	require.NoError(t, err)
	defer server.Close(ctx)

	// Start the client and perform a get request.
	response, httpClient, cancel, err := nativetailscalehttpclient.SendRequest(
		ctx,
		&tailscaleoptions.ConnectOptions{
			HostName:   "httpclient",
			PreAuthKey: preauthKeyClient,
			ControlURL: "http://localhost:" + strconv.Itoa(headScalePort),
		},
		&httpoptions.RequestOptions{
			Url: "http://webserver/hello",
		},
	)
	require.NoError(t, err)
	defer cancel()
	bodyString, err := response.GetBodyAsString()
	require.NoError(t, err)

	require.EqualValues(t, "Hello from Tailscale! You are connected via tsnet.\n", bodyString)

	// The returned httpClient can be reused to perform the next request:
	response2, err := httpClient.SendRequest(ctx, &httpoptions.RequestOptions{
		Url: "http://webserver/world",
	})
	require.NoError(t, err)
	bodyString2, err := response2.GetBodyAsString()
	require.NoError(t, err)

	require.EqualValues(t, "Another endpoint", bodyString2)
}
