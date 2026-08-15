package commandexecutorsshclient

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type SSHClient struct {
	commandexecutorgeneric.CommandExecutorBase
	hostName            string
	sshUserName         string
	sshPort             int  // SSH port (default: 0 means use SSH default)
	skipHostKeyChecking bool // Only for test environments - must be explicitly enabled
	sshPrivateKeyFile   string
}

func GetSshClientByHostName(hostName string) (sshClient *SSHClient, err error) {
	sshClient = NewSSHClient()

	err = sshClient.SetHostName(hostName)
	if err != nil {
		return nil, err
	}

	return sshClient, err
}

func NewSSHClient() (s *SSHClient) {
	s = new(SSHClient)

	s.SetParentCommandExecutorForBaseClass(s)

	return s
}

func (s *SSHClient) CheckReachable(ctx context.Context) (err error) {
	hostname, err := s.GetHostName()
	if err != nil {
		return err
	}

	isReachable, err := s.IsReachable(ctx)
	if err != nil {
		return err
	}

	if isReachable {
		return nil
	}

	return tracederrors.TracedErrorf("host '%v' is not reachable", hostname)
}

func (s *SSHClient) GetDeepCopyAsCommandExecutor() (copy commandexecutorinterfaces.CommandExecutor) {
	ret := NewSSHClient()

	*ret = *s

	err := ret.SetParentCommandExecutorForBaseClass(ret)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return ret
}

func (s *SSHClient) GetHostDescription() (hostDescription string, err error) {
	return s.GetHostName()
}

// IsRunningOnLocalhost always returns false for SSH clients.
// Even if the SSH host is "localhost", commands are executed through an SSH session,
// not natively on the local machine. This distinction is important for test cases
// that need to verify command execution paths.
func (s *SSHClient) IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error) {
	return false, nil
}

func (s *SSHClient) GetHostName() (hostName string, err error) {
	if s.hostName == "" {
		return "", tracederrors.TracedErrorf("hostName not set")
	}

	return s.hostName, nil
}

func (s *SSHClient) GetSshUserName() (sshUserName string, err error) {
	if !s.IsSshUserNameSet() {
		return "", tracederrors.TracedError("sshUserName not set")
	}

	return s.sshUserName, nil
}

func (s *SSHClient) IsReachable(ctx context.Context) (isReachable bool, err error) {
	hostname, err := s.GetHostName()
	if err != nil {
		return false, err
	}

	commandOutput, err := s.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command:           []string{"echo", "hello"},
			TimeoutString:     "5 seconds",
			AllowAllExitCodes: true,
		},
	)
	if err != nil {
		if commandOutput == nil {
			return false, tracederrors.TracedErrorf("commandOutput is nil and '%v'", err)
		}

		isTimedOut, err := commandOutput.IsTimedOut()
		if err != nil {
			return false, err
		}

		if isTimedOut {
			logging.LogInfoByCtxf(ctx, "'%v' is NOT reachable by SSH.", hostname)
			return false, nil
		}

		return false, err
	}

	returnValue, err := commandOutput.GetReturnCode()
	if err != nil {
		return false, err
	}

	if returnValue != 0 {
		return false, nil
	}

	stdout, err := commandOutput.GetStdoutAsString()
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout != "hello" {
		return false, tracederrors.TracedErrorf(
			"Unexpected stdout: '%s', stderr is '%s', return value is '%d'",
			stdout,
			commandOutput.GetStderrAsStringOrEmptyIfUnset(),
			returnValue,
		)
	}

	logging.LogInfoByCtxf(ctx, "'%v' is reachable by SSH.", hostname)
	return true, nil
}

func (s *SSHClient) IsSshUserNameSet() (isSet bool) {
	return len(s.sshUserName) > 0
}

// Get the full CLI command including ssh ... to be executed.
func (s *SSHClient) getCommandToUse(options *parameteroptions.RunCommandOptions) (*parameteroptions.RunCommandOptions, error) {
	userAtHost, err := s.GetHostName()
	if err != nil {
		return nil, err
	}

	if s.IsSshUserNameSet() {
		username, err := s.GetSshUserName()
		if err != nil {
			return nil, err
		}

		userAtHost = username + "@" + userAtHost
	}

	commandString, err := shelllinehandler.Join(options.Command)
	if err != nil {
		return nil, err
	}

	// Build SSH command with optional port and private key
	sshArgs := []string{"ssh"}

	// Add private key file if configured
	if s.sshPrivateKeyFile != "" {
		sshArgs = append(sshArgs, "-i", s.sshPrivateKeyFile)
	}

	// Only add non-interactive options when explicitly enabled (for test environments)
	if s.skipHostKeyChecking {
		sshArgs = append(sshArgs,
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "BatchMode=yes",
		)
	}

	if s.sshPort != 0 {
		sshArgs = append(sshArgs, "-p", fmt.Sprintf("%d", s.sshPort))
	}
	sshArgs = append(sshArgs, userAtHost, commandString)

	commandToUse := options.GetDeepCopy()
	commandToUse.Command = sshArgs

	return commandToUse, nil
}

func (s *SSHClient) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	commandToUse, err := s.getCommandToUse(options)
	if err != nil {
		return nil, err
	}

	commandOutput, err = commandexecutorexecoo.Exec().RunCommand(ctx, commandToUse)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (s *SSHClient) SetHostName(hostName string) (err error) {
	if hostName == "" {
		return tracederrors.TracedErrorf("hostName is empty string")
	}

	s.hostName = hostName

	return nil
}

func (s *SSHClient) SetSshUserName(sshUserName string) (err error) {
	if len(sshUserName) <= 0 {
		return tracederrors.TracedError("sshUserName is nil")
	}

	s.sshUserName = sshUserName

	return nil
}

func (s *SSHClient) SetSshPort(sshPort int) (err error) {
	if sshPort <= 0 {
		return tracederrors.TracedErrorf("Invalid SSH port: %d", sshPort)
	}
	s.sshPort = sshPort
	return nil
}

func (s *SSHClient) SetSkipHostKeyChecking(skip bool) {
	s.skipHostKeyChecking = skip
}

func (s *SSHClient) SetSshPrivateKeyFile(privateKeyFilePath string) (err error) {
	if privateKeyFilePath == "" {
		return tracederrors.TracedErrorEmptyString(privateKeyFilePath)
	}
	s.sshPrivateKeyFile = privateKeyFilePath
	return nil
}

func (s *SSHClient) GetSshPrivateKeyFile() string {
	return s.sshPrivateKeyFile
}

func (s *SSHClient) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	commandToUse, err := s.getCommandToUse(options)
	if err != nil {
		return nil, err
	}

	return commandexecutorexec.RunCommandAndGetStdoutAsIoReadCloser(ctx, commandToUse)
}

func (s *SSHClient) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	commandToUse, err := s.getCommandToUse(options)
	if err != nil {
		return nil, err
	}

	return commandexecutorexec.RunCommandAndGetStdinAsIoWriteCloser(ctx, commandToUse)
}
