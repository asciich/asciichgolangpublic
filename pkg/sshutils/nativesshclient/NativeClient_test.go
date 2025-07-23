package nativesshclient_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/parameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/sshutils/nativesshclient"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/sshutils/testsshserver"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_NativeClient(t *testing.T) {
	ctx := getCtx()

	t.Run("ping", func(t *testing.T) {
		const user = "user"
		const password = "pass"
		const port = 2222

		testSshServer := &testsshserver.TestSshServer{
			Username: "user",
			Password: "pass",
			Port:     port,
		}

		err := testSshServer.StartSshServerInBackground(ctx)
		require.NoError(t, err)

		sshClient := &nativesshclient.SshClient{
			Hostname: "localhost",
			Port:     port,
			Username: user,
			Password: password,
		}

		output, err := sshClient.RunCommand(ctx, &parameteroptions.RunCommandOptions{Command: []string{"ping"}})
		require.NoError(t, err)
		require.NotNil(t, output)

		stdout, err := output.GetStdoutAsString()
		require.NoError(t, err)
		require.EqualValues(t, "pong\n", stdout)

		err = testSshServer.Stop(ctx)
		require.NoError(t, err)
	})

}
