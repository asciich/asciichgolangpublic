package nativetailscaleclient

import (
	"context"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tsnet"
)

func GetStatus(ctx context.Context, srv *tsnet.Server) (*ipnstate.Status, error) {
	if srv == nil {
		return nil, tracederrors.TracedErrorNil("srv")
	}

	localClient, err := srv.LocalClient()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get tailscale local client: %w", err)
	}

	status, err := localClient.Status(ctx)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get tailscale status: %w", err)
	}

	return status, nil
}

func GetBackendStateString(ctx context.Context, srv *tsnet.Server) (string, error) {
	if srv == nil {
		return "", tracederrors.TracedError("srv")
	}

	status, err := GetStatus(ctx, srv)
	if err != nil {
		return "", err
	}

	backendState := status.BackendState
	if backendState == "" {
		return "", tracederrors.TracedError("backendState string is empty string after evaluation")
	}

	logging.LogInfoByCtxf(ctx, "Tailscale backend state is '%s'.", backendState)

	return backendState, nil
}

func IsConnected(ctx context.Context, srv *tsnet.Server) (bool, error) {
	if srv == nil {
		return false, tracederrors.TracedError("srv")
	}

	state, err := GetBackendStateString(ctx, srv)
	if err != nil {
		return false, err
	}

	var isConnected = false
	if strings.EqualFold(state, "Running") {
		isConnected = true
	}

	if isConnected {
		logging.LogInfoByCtxf(ctx, "Tailscale client is connected. State is '%s'.", state)
	} else {
		logging.LogInfoByCtxf(ctx, "Tailscale client is not connected. State is '%s'.", state)
	}

	return isConnected, nil
}

func WaitUntilRunning(ctx context.Context, srv *tsnet.Server, timeout time.Duration) error {
	if srv == nil {
		return tracederrors.TracedErrorNil("srv")
	}

	logging.LogInfoByCtxf(ctx, "Wait until tailscale client is connected started.")

	tStart := time.Now()
	var isConnected bool
	var err error
	for {
		if time.Since(tStart) > timeout {
			logging.LogInfoByCtxf(ctx, "Timeout exceeded while waiting for tailscale client to be connected.")
			break
		}

		isConnected, err = IsConnected(contextutils.WithSilent(ctx), srv)
		if err != nil {
			return err
		}

		if isConnected {
			break
		}

		delay := time.Millisecond * 500
		logging.LogInfoByCtxf(ctx, "Tailscale client is not yet connected. Going to retry in '%s'.", delay)
		time.Sleep(delay)
	}

	if isConnected {
		logging.LogInfoByCtxf(ctx, "Tailscale client is connected.")
	} else {
		return tracederrors.TracedErrorf("Tailscale client is still not connected after '%v'.", timeout)
	}

	logging.LogInfoByCtxf(ctx, "Wait until tailscale client is connected finished.")

	return nil
}
