package tailscaleoptions_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_GetStateDir(t *testing.T) {
	t.Run("by hostname", func(t *testing.T) {
		ctx := getCtx()

		options := &tailscaleoptions.ConnectOptions{
			HostName: "myhostname",
		}

		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)

		expected := filepath.Join(homeDir, ".config", "tsnet-myhostname")

		stateDir, err := options.GetStateDir(ctx)
		require.NoError(t, err)
		require.EqualValues(t, expected, stateDir)
	})
}
