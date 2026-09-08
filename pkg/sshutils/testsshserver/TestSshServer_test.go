package testsshserver_test

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/testsshserver"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// runExecCommandOverSsh connects to the running TestSshServer, executes the
// given command via an "exec" request and returns the combined stdout output.
func runExecCommandOverSsh(t *testing.T, port int, username, password, command string) string {
	t.Helper()

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// This is a test server with an ephemeral host key, so we don't verify it.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := net.JoinHostPort("localhost", strconv.Itoa(port))

	client, err := ssh.Dial("tcp", address, clientConfig)
	require.NoError(t, err)
	defer client.Close()

	session, err := client.NewSession()
	require.NoError(t, err)
	defer session.Close()

	output, err := session.Output(command)
	require.NoError(t, err)

	return strings.TrimRight(string(output), "\n")
}

func Test_TestSshServer_ExecCommands(t *testing.T) {
	ctx := getCtx()

	const (
		port     = 2222
		username = "user"
		password = "pass"
	)

	tests := []struct {
		command        string
		expectedOutput string
	}{
		{"echo hello", "hello"},
		{"echo hallo", "hallo"},
		{"ping", "pong"},
	}

	testSshServer := &testsshserver.TestSshServer{
		Username: username,
		Password: password,
		Port:     port,
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := testSshServer.StartSshServerInBackground(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testSshServer.Stop(ctx))
	}()

	isOpen, err := netutils.IsTcpPortOpen(ctx, "localhost", port)
	require.NoError(t, err)
	require.True(t, isOpen)

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			output := runExecCommandOverSsh(t, port, username, password, tt.command)
			require.Equal(t, tt.expectedOutput, output)
		})
	}
}

// dialSshClient connects to the running TestSshServer and returns a ready-to-use
// SSH client. The caller is responsible for closing it.
func dialSshClient(t *testing.T, port int, username, password string) *ssh.Client {
	t.Helper()

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// This is a test server with an ephemeral host key, so we don't verify it.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := net.JoinHostPort("localhost", strconv.Itoa(port))

	client, err := ssh.Dial("tcp", address, clientConfig)
	require.NoError(t, err)

	return client
}

func Test_TestSshServer_UnknownCommandReturnsError(t *testing.T) {
	ctx := getCtx()

	const (
		port     = 2222
		username = "user"
		password = "pass"
	)

	testSshServer := &testsshserver.TestSshServer{
		Username: username,
		Password: password,
		Port:     port,
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := testSshServer.StartSshServerInBackground(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testSshServer.Stop(ctx))
	}()

	client := dialSshClient(t, port, username, password)
	defer client.Close()

	session, err := client.NewSession()
	require.NoError(t, err)
	defer session.Close()

	err = session.Run("this-command-does-not-exist")

	var exitErr *ssh.ExitError
	require.ErrorAs(t, err, &exitErr)
	require.Equal(t, 1, exitErr.ExitStatus())
}
