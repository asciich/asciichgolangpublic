package testsshserver

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"golang.org/x/crypto/ssh"
)

type TestSshServer struct {
	Username string
	Password string
	Port     int

	cancelMux sync.Mutex
	cancel    func()
}

func (t *TestSshServer) StartSshServerInBackground(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Start TestSshServer in background started.")

	err := t.WaitUntilPortUnused(ctx)
	if err != nil {
		return err
	}

	var finished bool
	go func() {
		err := t.StartSshServer(ctx)
		if err == nil {
			logging.LogInfoByCtxf(ctx, "TestSshSever exited successfully.")
		} else {
			logging.LogErrorf("TestSshServer exited with error: %v", err)
		}
		finished = true
	}()

	var isOpen bool
	for i := 0; i < 10; i++ {
		isOpen, err = netutils.IsTcpPortOpen(contextutils.WithSilent(ctx), "localhost", t.Port)
		if err != nil {
			return err
		}

		if isOpen {
			break
		}

		if finished {
			break
		}

		logging.LogInfoByCtx(ctx, "Wait until TestSshServer port is open")
		time.Sleep(time.Millisecond * 300)
	}

	if !isOpen {
		return tracederrors.TracedError("Failed to start TestSshServer")
	}

	logging.LogInfoByCtxf(ctx, "Start TestSshServer in background started.")

	return nil
}

func (t *TestSshServer) WaitUntilStopped(ctx context.Context) error {
	logging.LogInfoByCtx(ctx, "Wait for TestSshServer to be stopped started")

	var isOpen bool
	var err error

	for i := 0; i < 10; i++ {
		isOpen, err = netutils.IsTcpPortOpen(ctx, "localhost", t.Port)
		if err != nil {
			return err
		}

		if !isOpen {
			break
		}
	}

	if isOpen {
		return tracederrors.TracedError("Wait until stopped failed: TestSshServer still running")
	}

	logging.LogInfoByCtx(ctx, "TestSshServer stopped.")

	return err
}

func (t *TestSshServer) Stop(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Stopping TestSshServer started.")

	t.cancelMux.Lock()
	defer t.cancelMux.Unlock()

	if t.cancel == nil {
		logging.LogInfoByCtx(ctx, "TestSshServer already stopped.")
		return nil
	}

	t.cancel()
	logging.LogInfoByCtx(ctx, "cancel to stop TestSshServer called.")

	err := t.WaitUntilStopped(ctx)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Stopping TestSshServer finished.")

	return nil
}

func (t *TestSshServer) WaitUntilPortUnused(ctx context.Context) error {
	const hostname = "localhost"

	nRetry := 10
	var isOpen bool
	var err error
	for i := 0; i < nRetry; i++ {
		isOpen, err = netutils.IsTcpPortOpen(contextutils.WithSilent(ctx), hostname, t.Port)
		if err != nil {
			return err
		}

		if !isOpen {
			break
		}

		if i+1 == nRetry {
			break
		}

		waitDuration := time.Millisecond * 100
		logging.LogInfoByCtxf(ctx, "Port '%d' is not free to listen on '%s', wait another '%s' (%d/%d).", t.Port, hostname, waitDuration, i+1, nRetry)
		time.Sleep(waitDuration)
	}

	if !isOpen {
		logging.LogInfoByCtxf(ctx, "Port '%d' on '%s' is unused.", t.Port, hostname)
		return nil
	}

	return tracederrors.TracedErrorf("Port '%d' on '%s' is not free to use.", t.Port, hostname)
}

func (t *TestSshServer) StartSshServer(ctx context.Context) error {
	if t.cancel != nil {
		return tracederrors.TracedError("TestSshServer already running")
	}

	err := t.WaitUntilPortUnused(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	t.cancel = cancel

	hostKey, err := t.generateHostKey()
	if err != nil {
		tracederrors.TracedErrorf("Failed to generate host key: %w", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: false,

		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == t.Username && string(pass) == t.Password {
				return nil, nil // Authentication successful
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	config.AddHostKey(hostKey)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(t.Port))
	if err != nil {
		log.Fatalf("Failed to listen on %d: %v", t.Port, err)
	}
	defer listener.Close()

	go func() {
		select {
		case <-ctx.Done():
			logging.LogInfoByCtxf(ctx, "context for TestSshServer is done. Going to close TestSshServer.")
			listener.Close()
		}
	}()

	log.Printf("SSH server listening on %d...", t.Port)
	logging.LogInfoByCtxf(ctx, "Use 'ssh -p %d %s@localhost' to connect (password: %s)", t.Port, t.Username, t.Password)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) || err == io.EOF {
				logging.LogInfoByCtxf(ctx, "Webserver connection closed. Going to exit")
				break
			}
			logging.LogWarnf("Failed to accept incoming connection: %v", err)
			continue
		}

		go t.handleConnection(conn, config)
	}

	logging.LogInfoByCtx(ctx, "TestSshServer finished.")

	return nil
}

func (t *TestSshServer) generateHostKey() (ssh.Signer, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA private key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH signer from private key: %v", err)
	}

	return signer, nil
}

func (t *TestSshServer) handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	// Handle global out-of-band requests (e.g., keep-alives)
	go ssh.DiscardRequests(reqs)

	// Handle channels (e.g., "session" channels for shell commands)
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		go t.handleSession(newChannel)
	}
}

func (t *TestSshServer) handleSession(newChannel ssh.NewChannel) {
	channel, requests, err := newChannel.Accept()
	if err != nil {
		log.Printf("Could not accept channel: %v", err)
		return
	}
	defer channel.Close()

	var wg sync.WaitGroup
	wg.Add(1) // Wait for the session to complete or exit

	go func() {
		defer wg.Done()
		for req := range requests {
			switch req.Type {
			case "shell":
				if len(req.Payload) == 0 {
					req.Reply(true, nil)
					t.serveShell(channel)
				} else {
					req.Reply(false, nil)
				}
			case "exec":
				command := string(req.Payload[4:])
				log.Printf("Received exec command: %s", command)
				t.handleExecCommand(channel, command)
				req.Reply(true, nil)
				channel.Close()
			case "pty-req":
				// We don't really care about PTY details for this example, just acknowledge
				req.Reply(true, nil)
			case "window-change":
				// Ignore window-change requests for this simple example
				req.Reply(false, nil)
			default:
				log.Printf("Unknown channel request type: %s", req.Type)
				req.Reply(false, nil)
			}
		}
	}()

	wg.Wait()
}

func (t *TestSshServer) serveShell(channel ssh.Channel) {
	defer channel.Close()

	welcomeMsg := "Welcome to the Go SSH server!\n"
	io.WriteString(channel, welcomeMsg)
	io.WriteString(channel, "Type 'ping' or 'exit'.\n")
	io.WriteString(channel, "> ")

	buf := make([]byte, 256)
	for {
		n, err := channel.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected.")
			} else {
				log.Printf("Error reading from channel: %v", err)
			}
			break
		}

		command := string(buf[:n])
		command = strings.TrimRight(command, "\n")

		switch command {
		case "ping":
			io.WriteString(channel, "pong\n")
		case "exit":
			io.WriteString(channel, "Goodbye!\n")
			return
		default:
			io.WriteString(channel, fmt.Sprintf("Unknown command: %s\n", command))
		}
		io.WriteString(channel, "> ")
	}
}

func (t *TestSshServer) handleExecCommand(channel ssh.Channel, command string) {
	switch strings.TrimRight(command, "\n") {
	case "ping":
		io.WriteString(channel, "pong\n")
		channel.SendRequest("exit-status", false, ssh.Marshal(struct {
			Status uint32
		}{Status: 0}))
	case "exit":
		io.WriteString(channel, "Goodbye!\n")
		channel.SendRequest("exit-status", false, ssh.Marshal(struct {
			Status uint32
		}{Status: 0}))
	default:
		io.WriteString(channel, fmt.Sprintf("Unknown command: %s\n", command))
		channel.SendRequest("exit-status", false, ssh.Marshal(struct {
			Status uint32
		}{Status: 1}))
	}
}
