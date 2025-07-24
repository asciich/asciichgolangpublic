package nativesshclient

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	Hostname string
	Port     int
	Username string
	Password string
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
