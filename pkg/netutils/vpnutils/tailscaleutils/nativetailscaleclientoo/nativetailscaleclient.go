package nativetailscaleclientoo

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscaleclient"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"

	"tailscale.com/tsnet"
)

type Client struct {
	server *tsnet.Server
}

func NewClient() *Client {
	return &Client{}
}

func Connect(ctx context.Context, options *tailscaleoptions.ConnectOptions) (client *Client, cancel func() error, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	hostname, err := options.GetHostName()
	if err != nil {
		return nil, nil, err
	}

	controlUrl, err := options.GetControlUrl()
	if err != nil {
		return nil, nil, err
	}

	logging.LogInfoByCtxf(ctx, "Connect tailscale client with hostname '%s' on control server '%s' started.", hostname, controlUrl)

	authKey, err := options.GetAuthKey()
	if err != nil {
		return nil, nil, err
	}

	stateDir, err := options.GetStateDir(ctx)
	if err != nil {
		return nil, nil, err
	}

	srv := &tsnet.Server{
		Hostname:   hostname,
		AuthKey:    authKey,
		ControlURL: controlUrl,
		Dir:        stateDir,
	}

	err = srv.Start()
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to start tailscale client: %w", err)
	}

	err = nativetailscaleclient.WaitUntilRunning(ctx, srv, time.Second*30)

	isConnected, err := nativetailscaleclient.IsConnected(ctx, srv)
	if err != nil {
		return nil, nil, err
	}

	if isConnected {
		logging.LogChangedByCtxf(ctx, "Connected tailscale client.")
	} else {
		return nil, nil, tracederrors.TracedErrorf("Failed to connnect tailscale client")
	}

	client = &Client{
		server: srv,
	}

	cancel = func() error {
		return client.Cancel()
	}

	logging.LogInfoByCtxf(ctx, "Connect tailscale client with hostname '%s' on control server '%s' finished.", hostname, controlUrl)

	return client, cancel, nil
}

func (c *Client) Cancel() error {
	if c.server != nil {
		c.server.Close()
	}

	return nil
}

func (c *Client) GetNativeServer() (*tsnet.Server, error) {
	if c.server == nil {
		return nil, tracederrors.TracedError("server not set")
	}

	return c.server, nil
}

func (c *Client) IsConnected(ctx context.Context) (bool, error) {
	srv, err := c.GetNativeServer()
	if err != nil {
		return false, err
	}

	return nativetailscaleclient.IsConnected(ctx, srv)
}

func (c *Client) Listen(port int) (net.Listener, error) {
	srv, err := c.GetNativeServer()
	if err != nil {
		return nil, err
	}

	listener, err := srv.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create listener: %w", err)
	}

	return listener, nil
}
