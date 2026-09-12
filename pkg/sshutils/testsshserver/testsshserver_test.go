package testsshserver_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
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

// generateTestSshKey generates a test SSH key pair and returns the private key signer and public key
func generateTestSshKey(t *testing.T) (ssh.Signer, ssh.PublicKey) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	signer, err := ssh.NewSignerFromKey(privateKey)
	require.NoError(t, err)

	return signer, signer.PublicKey()
}

// dialSshClientWithKey connects to the running TestSshServer using SSH key authentication
func dialSshClientWithKey(t *testing.T, port int, username string, signer ssh.Signer) *ssh.Client {
	t.Helper()

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
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

// dialSshClientWithPasswordAndKey connects to the running TestSshServer using both password and SSH key authentication
// The client will try password first, then key
func dialSshClientWithPasswordAndKey(t *testing.T, port int, username, password string, signer ssh.Signer) *ssh.Client {
	t.Helper()

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.PublicKeys(signer),
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

func Test_TestSshServer_PasswordAuthenticationOnly(t *testing.T) {
	ctx := getCtx()

	const (
		port     = 2223
		username = "user"
		password = "pass"
	)

	testSshServer := &testsshserver.TestSshServer{
		Username: username,
		Password: password,
		Port:     port,
		// No AuthorizedKeys - password only
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := testSshServer.StartSshServerInBackground(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testSshServer.Stop(ctx))
	}()

	// Test successful password authentication
	t.Run("password auth succeeds", func(t *testing.T) {
		client := dialSshClient(t, port, username, password)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("ping")
		require.NoError(t, err)
		require.Equal(t, "pong", strings.TrimRight(string(output), "\n"))
	})

	// Test that SSH key authentication fails when only password is configured
	t.Run("key auth fails when only password configured", func(t *testing.T) {
		signer, _ := generateTestSshKey(t)

		clientConfig := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		address := net.JoinHostPort("localhost", strconv.Itoa(port))
		_, err := ssh.Dial("tcp", address, clientConfig)
		require.Error(t, err, "SSH key authentication should fail when only password is configured")
	})
}

func Test_TestSshServer_SshKeyAuthenticationOnly(t *testing.T) {
	ctx := getCtx()

	const (
		port     = 2224
		username = "user"
	)

	// Generate a test SSH key pair
	signer, publicKey := generateTestSshKey(t)

	testSshServer := &testsshserver.TestSshServer{
		Username:       username,
		Password:       "", // No password - key only
		Port:           port,
		AuthorizedKeys: []ssh.PublicKey{publicKey},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := testSshServer.StartSshServerInBackground(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testSshServer.Stop(ctx))
	}()

	// Test successful SSH key authentication
	t.Run("key auth succeeds", func(t *testing.T) {
		client := dialSshClientWithKey(t, port, username, signer)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("ping")
		require.NoError(t, err)
		require.Equal(t, "pong", strings.TrimRight(string(output), "\n"))
	})

	// Test that password authentication fails when only SSH key is configured
	t.Run("password auth fails when only key configured", func(t *testing.T) {
		clientConfig := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password("wrongpassword"),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		address := net.JoinHostPort("localhost", strconv.Itoa(port))
		_, err := ssh.Dial("tcp", address, clientConfig)
		require.Error(t, err, "password authentication should fail when only SSH key is configured")
	})

	// Test that wrong SSH key fails
	t.Run("wrong key fails", func(t *testing.T) {
		wrongSigner, _ := generateTestSshKey(t)

		clientConfig := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(wrongSigner),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		address := net.JoinHostPort("localhost", strconv.Itoa(port))
		_, err := ssh.Dial("tcp", address, clientConfig)
		require.Error(t, err, "wrong SSH key should be rejected")
	})
}

func Test_TestSshServer_BothPasswordAndSshKeyAuthentication(t *testing.T) {
	ctx := getCtx()

	const (
		port     = 2225
		username = "user"
		password = "pass"
	)

	// Generate a test SSH key pair
	signer, publicKey := generateTestSshKey(t)

	testSshServer := &testsshserver.TestSshServer{
		Username:       username,
		Password:       password,
		Port:           port,
		AuthorizedKeys: []ssh.PublicKey{publicKey},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := testSshServer.StartSshServerInBackground(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testSshServer.Stop(ctx))
	}()

	// Test successful password authentication
	t.Run("password auth succeeds", func(t *testing.T) {
		client := dialSshClient(t, port, username, password)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("ping")
		require.NoError(t, err)
		require.Equal(t, "pong", strings.TrimRight(string(output), "\n"))
	})

	// Test successful SSH key authentication
	t.Run("key auth succeeds", func(t *testing.T) {
		client := dialSshClientWithKey(t, port, username, signer)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("ping")
		require.NoError(t, err)
		require.Equal(t, "pong", strings.TrimRight(string(output), "\n"))
	})

	// Test that client can connect with both methods available (tries password first)
	t.Run("client with both methods succeeds", func(t *testing.T) {
		client := dialSshClientWithPasswordAndKey(t, port, username, password, signer)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("echo hello")
		require.NoError(t, err)
		require.Equal(t, "hello", strings.TrimRight(string(output), "\n"))
	})

	// Test that wrong password fails even with correct key available
	t.Run("wrong password fails", func(t *testing.T) {
		wrongSigner, _ := generateTestSshKey(t)

		clientConfig := &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password("wrongpassword"),
				ssh.PublicKeys(wrongSigner),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         5 * time.Second,
		}

		address := net.JoinHostPort("localhost", strconv.Itoa(port))
		_, err := ssh.Dial("tcp", address, clientConfig)
		require.Error(t, err, "wrong password and wrong key should both fail")
	})
}

func Test_TestSshServer_MultipleAuthorizedKeys(t *testing.T) {
	ctx := getCtx()

	const (
		port     = 2226
		username = "user"
	)

	// Generate multiple test SSH key pairs
	signer1, publicKey1 := generateTestSshKey(t)
	signer2, publicKey2 := generateTestSshKey(t)

	testSshServer := &testsshserver.TestSshServer{
		Username:       username,
		Password:       "", // Key only
		Port:           port,
		AuthorizedKeys: []ssh.PublicKey{publicKey1, publicKey2},
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := testSshServer.StartSshServerInBackground(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testSshServer.Stop(ctx))
	}()

	// Test that first key works
	t.Run("first key succeeds", func(t *testing.T) {
		client := dialSshClientWithKey(t, port, username, signer1)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("ping")
		require.NoError(t, err)
		require.Equal(t, "pong", strings.TrimRight(string(output), "\n"))
	})

	// Test that second key works
	t.Run("second key succeeds", func(t *testing.T) {
		client := dialSshClientWithKey(t, port, username, signer2)
		defer client.Close()

		session, err := client.NewSession()
		require.NoError(t, err)
		defer session.Close()

		output, err := session.Output("ping")
		require.NoError(t, err)
		require.Equal(t, "pong", strings.TrimRight(string(output), "\n"))
	})
}
