package nativetailscalehttpclient

import (
	"context"
	"net/http"

	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpnativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscaleclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"tailscale.com/tsnet"
)

type HttpClient struct {
	httpnativeclientoo.NativeClient

	tailscaleClient *nativetailscaleclientoo.Client
}

func NewClient() (*HttpClient, error) {
	client := &HttpClient{}

	return client, nil
}

func Connect(ctx context.Context, options *tailscaleoptions.ConnectOptions) (*HttpClient, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx, options)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Directly establish a connection and send the http request.
func SendRequest(ctx context.Context, connectOptions *tailscaleoptions.ConnectOptions, requestOptions *httpoptions.RequestOptions) (response httputilsinterfaces.Response, client *HttpClient, cancel func() error, err error) {
	if connectOptions == nil {
		return nil, nil, nil, tracederrors.TracedErrorNil("connectOptions")
	}

	if requestOptions == nil {
		return nil, nil, nil, tracederrors.TracedErrorNil("requestOptions")
	}

	client, err = Connect(ctx, connectOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	response, err = client.SendRequest(ctx, requestOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	cancel = func() error {
		return client.Close()
	}

	return response, client, cancel, nil
}

func (h *HttpClient) Close() error {
	if h.tailscaleClient != nil {
		err := h.tailscaleClient.Cancel()
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HttpClient) Connect(ctx context.Context, options *tailscaleoptions.ConnectOptions) error {
	logging.LogInfoByCtxf(ctx, "Connect tailscale HttpClient started.")

	err := h.Close()
	if err != nil {
		return err
	}

	tailscaleClient, _, err := nativetailscaleclientoo.Connect(ctx, options)
	if err != nil {
		return err
	}

	h.tailscaleClient = tailscaleClient

	logging.LogInfoByCtxf(ctx, "Connect tailscale HttpClient finished.")

	return nil
}

func (h *HttpClient) GetTailscaleClient() (*nativetailscaleclientoo.Client, error) {
	if h.tailscaleClient == nil {
		return nil, tracederrors.TracedError("tailscaleClient not set. Did you run Connect()?")
	}

	return h.tailscaleClient, nil
}

func (h *HttpClient) GetNativeTailscaleClient() (*tsnet.Server, error) {
	tailscaleClient, err := h.GetTailscaleClient()
	if err != nil {
		return nil, err
	}

	return tailscaleClient.GetNativeServer()
}

func (h *HttpClient) SendRequest(ctx context.Context, options *httpoptions.RequestOptions) (response httputilsinterfaces.Response, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	srv, err := h.GetNativeTailscaleClient()
	if err != nil {
		return nil, err
	}

	optionsToUse := options.GetDeepCopy()

	// By using the the tailscale client in the transport the request is send over tailscale:
	optionsToUse.TransportToUse = &http.Transport{
		DialContext: srv.Dial,
	}

	return h.NativeClient.SendRequest(ctx, optionsToUse)
}
