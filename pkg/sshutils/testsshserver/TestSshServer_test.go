package testsshserver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/netutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/sshutils/testsshserver"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_TestSshServer(t *testing.T) {
	ctx := getCtx()

	t.Run("start and cancel", func(t *testing.T) {
		const port = 2222

		testSshServer := &testsshserver.TestSshServer{
			Username: "user",
			Password: "pass",
			Port:     port,
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		err := testSshServer.StartSshServerInBackground(ctx)
		require.NoError(t, err)

		isOpen, err := netutils.IsTcpPortOpen(ctx, "localhost", port)
		require.NoError(t, err)
		require.True(t, isOpen)

		cancel()
		err = testSshServer.WaitUntilStopped(ctx)
		require.NoError(t, err)

		isOpen, err = netutils.IsTcpPortOpen(ctx, "localhost", port)
		require.NoError(t, err)
		require.False(t, isOpen)
	})

	t.Run("start and stop", func(t *testing.T) {
		const port = 2222

		testSshServer := &testsshserver.TestSshServer{
			Username: "user",
			Password: "pass",
			Port:     port,
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		err := testSshServer.StartSshServerInBackground(ctx)
		require.NoError(t, err)

		isOpen, err := netutils.IsTcpPortOpen(ctx, "localhost", port)
		require.NoError(t, err)
		require.True(t, isOpen)

		testSshServer.Stop(ctx)

		isOpen, err = netutils.IsTcpPortOpen(ctx, "localhost", port)
		require.NoError(t, err)
		require.False(t, isOpen)
	})

}
