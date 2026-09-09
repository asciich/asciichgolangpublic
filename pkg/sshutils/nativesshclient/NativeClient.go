package nativesshclient

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/sshutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	commandexecutorgeneric.CommandExecutorBase
	Hostname string
	Port     int
	Username string
	Password string
}

func NewSshClientByHostName(hostname string) (*SshClient, error) {
	if hostname == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostname")
	}

	client := &SshClient{
		Hostname: hostname,
	}

	err := client.SetParentCommandExecutorForBaseClass(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
func (s *SshClient) SetSshUserName(userName string) error {
	if userName == "" {
		return tracederrors.TracedErrorEmptyString("userName")
	}

	s.Username = userName

	return nil
}

func (s *SshClient) GetSshUserName() (userName string, err error) {
	if s.Username == "" {
		return "", tracederrors.TracedError("Username not set")
	}

	return s.Username, nil
}

func (s *SshClient) SetSshPort(port int) error {
	if port <= 0 {
		return tracederrors.TracedErrorf("Invalid ssh port: '%d'", port)
	}

	if port > 65535 {
		return tracederrors.TracedErrorf("Invalid ssh port: '%d'", port)
	}

	s.Port = port

	return nil
}

func (s *SshClient) GetSshPort() (port int, err error) {
	if s.Port <= 0 {
		return 0, tracederrors.TracedError("Port not set")
	}

	return s.Port, nil
}

func (s *SshClient) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	serverAddress := fmt.Sprintf("%s:%d", s.Hostname, s.Port)
	cmd, err := options.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}

	logging.LogInfoByCtxf(ctx, "Connecting to SSH server at %s...", serverAddress)
	client, err := ssh.Dial("tcp", serverAddress, config)
	if err != nil {
		if sshutilsgeneric.IsSshConnectionRefused(err, nil) {
			return nil, tracederrors.TracedErrorf(
				"%w: Failed to dial SSH server '%s': %w",
				netutilserrors.ErrConnectionRefused,
				serverAddress,
				err,
			)
		}

		if sshutilsgeneric.IsNoRouteToHost(err, nil) {
			return nil, tracederrors.TracedErrorf(
				"%w: %w",
				netutilserrors.ErrNoRouteToHost,
				err,
			)
		}

		return nil, tracederrors.TracedErrorf("Failed to dial SSH server: %w", err)
	}
	defer client.Close()
	logging.LogInfoByCtx(ctx, "Successfully connected to SSH server.")

	session, err := client.NewSession()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create SSH session: %v", err)
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(cmd)
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			return nil, tracederrors.TracedErrorf("Command exited with non-zero status: %d", exitErr.ExitStatus())
		} else {
			return nil, tracederrors.TracedErrorf("Failed to run command: %w", err)
		}
	}

	output := &commandoutput.CommandOutput{}
	err = output.SetReturnCode(0)
	if err != nil {
		return nil, err
	}

	err = output.SetStdout(stdoutBuf.Bytes())
	if err != nil {
		return nil, err
	}

	err = output.SetStderr(stderrBuf.Bytes())
	if err != nil {
		return nil, err
	}

	return output, nil
}
