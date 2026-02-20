package nativetailscalehttpserver

import (
	"context"
	"net"
	"net/http"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscaleclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type TailScaleHttpServer struct {
	tailscaleClient *nativetailscaleclientoo.Client
	listener        net.Listener
	httpserver      *http.Server
}

func StartHttpServer(ctx context.Context, mux *http.ServeMux, port int, options *tailscaleoptions.ConnectOptions) (*TailScaleHttpServer, error) {
	if mux == nil {
		return nil, tracederrors.TracedErrorNil("mux")
	}

	if port <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid port '%d'.", port)
	}

	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	hostname, err := options.GetHostName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Start tailscale http server '%s' started.", hostname)

	client, _, err := nativetailscaleclientoo.Connect(ctx, options)
	if err != nil {
		return nil, err
	}

	listener, err := client.Listen(port)
	if err != nil {
		return nil, err
	}

	httpServer := http.Server{
		Handler: mux,
	}

	go func() {
		err := httpServer.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			logging.LogErrorByCtxf(ctx, "Failed to server over tailscale: %s", err.Error())
		}
	}()

	tailscaleServer := &TailScaleHttpServer{
		tailscaleClient: client,
		listener:        listener,
		httpserver:      &httpServer,
	}

	logging.LogInfoByCtxf(ctx, "Start tailscale http server '%s' finished.", hostname)

	return tailscaleServer, nil
}

func (t *TailScaleHttpServer) Close(ctx context.Context) error {
	if t.httpserver != nil {
		err := t.httpserver.Close()
		if err != nil {
			return tracederrors.TracedErrorf("Failed to close TailScaleHttpServer httpserver: %w", err)
		}

		t.httpserver = nil
	}

	if t.listener != nil {
		err := t.listener.Close()
		if err != nil {
			return tracederrors.TracedErrorf("Failed to close TailScaleHttpServer listener : %w", err)
		}

		t.listener = nil
	}

	if t.tailscaleClient != nil {
		err := t.tailscaleClient.Cancel()
		if err != nil {
			return tracederrors.TracedErrorf("Failed to cancel tailscaleClient : %w", err)
		}

		t.tailscaleClient = nil
	}

	return nil
}
