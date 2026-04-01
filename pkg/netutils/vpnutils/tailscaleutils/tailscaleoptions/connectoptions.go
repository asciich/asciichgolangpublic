package tailscaleoptions

import (
	"context"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ConnectOptions struct {
	// The hostname used on tailscale
	HostName string

	// PreAuthKey to use to connect
	PreAuthKey string

	// URL of the control server
	ControlURL string
}

func (c *ConnectOptions) GetHostName() (string, error) {
	if c.HostName == "" {
		return "", tracederrors.TracedError("HostName not set")
	}

	return c.HostName, nil
}

func (c *ConnectOptions) GetAuthKey() (string, error) {
	if c.PreAuthKey == "" {
		return "", tracederrors.TracedError("PreAuthKey not set")
	}

	return c.PreAuthKey, nil
}

func (c *ConnectOptions) GetControlUrl() (string, error) {
	if c.ControlURL == "" {
		return "", tracederrors.TracedError("ControlURL not set")
	}

	return c.ControlURL, nil
}

func (c *ConnectOptions) GetStateDir(ctx context.Context) (string, error) {
	var stateDir string
	const envName = "TS_STATE_DIR"
	if os.Getenv(envName) == "" {
		hostname, err := c.GetHostName()
		if err != nil {
			return "", err
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", tracederrors.TracedErrorf("Failed to get users home: %w", err)
		}

		stateDir = filepath.Join(homeDir, ".config", "tsnet-"+hostname)
	} else {
		stateDir = os.Getenv(envName)
		logging.LogInfoByCtxf(ctx, "Taking tailscale client state dir from env var '%s'='%s'", envName, stateDir)
	}

	logging.LogInfoByCtxf(ctx, "Tailscale client state dir by ConnectOptions is '%s'.", stateDir)

	return stateDir, nil
}
