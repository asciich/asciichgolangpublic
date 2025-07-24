package commandexecutorsshclient

import (
	"context"
	"strings"

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
	hostName    string
	sshUserName string
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

func (s *SSHClient) GetDeepCopy() (copy commandexecutorinterfaces.CommandExecutor) {
	toReturn := NewSSHClient()

	*toReturn = *s

	err := toReturn.SetParentCommandExecutorForBaseClass(toReturn)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return toReturn
}

func (s *SSHClient) GetHostDescription() (hostDescription string, err error) {
	return s.GetHostName()
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

func (s *SSHClient) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
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

	commandToUse := options.GetDeepCopy()
	commandToUse.Command = []string{
		"ssh",
		userAtHost,
		commandString,
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
