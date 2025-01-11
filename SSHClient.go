package asciichgolangpublic

import (
	"strings"
)

type SSHClient struct {
	CommandExecutorBase
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

func MustGetSshClientByHostName(hostName string) (sshClient *SSHClient) {
	sshClient, err := GetSshClientByHostName(hostName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshClient
}

func NewSSHClient() (s *SSHClient) {
	s = new(SSHClient)

	s.SetParentCommandExecutorForBaseClass(s)

	return s
}

func (s *SSHClient) CheckReachable(verbose bool) (err error) {
	hostname, err := s.GetHostName()
	if err != nil {
		return err
	}

	isReachable, err := s.IsReachable(verbose)
	if err != nil {
		return err
	}

	if isReachable {
		return nil
	}

	return TracedErrorf("host '%v' is not reachable", hostname)
}

func (s *SSHClient) GetDeepCopy() (copy CommandExecutor) {
	toReturn := NewSSHClient()

	*toReturn = *s

	toReturn.MustSetParentCommandExecutorForBaseClass(toReturn)

	return toReturn
}

func (s *SSHClient) GetHostDescription() (hostDescription string, err error) {
	return s.GetHostName()
}

func (s *SSHClient) GetHostName() (hostName string, err error) {
	if s.hostName == "" {
		return "", TracedErrorf("hostName not set")
	}

	return s.hostName, nil
}

func (s *SSHClient) GetSshUserName() (sshUserName string, err error) {
	if !s.IsSshUserNameSet() {
		return "", TracedError("sshUserName not set")
	}

	return s.sshUserName, nil
}

func (s *SSHClient) IsReachable(verbose bool) (isReachable bool, err error) {
	hostname, err := s.GetHostName()
	if err != nil {
		return false, err
	}

	commandOutput, err := s.RunCommand(
		&RunCommandOptions{
			Command:           []string{"echo", "hello"},
			TimeoutString:     "5 seconds",
			AllowAllExitCodes: true,
			Verbose:           verbose,
		},
	)
	if err != nil {
		if commandOutput == nil {
			return false, TracedErrorf("commandOutput is nil and '%v'", err)
		}

		isTimedOut, err := commandOutput.IsTimedOut()
		if err != nil {
			return false, err
		}

		if isTimedOut {
			if verbose {
				LogInfof("'%v' is NOT reachable by SSH.", hostname)
			}
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
		return false, TracedErrorf(
			"Unexpected stdout: '%s', stderr is '%s', return value is '%d'",
			stdout,
			commandOutput.GetStderrAsStringOrEmptyIfUnset(),
			returnValue,
		)
	}

	if verbose {
		LogInfof("'%v' is reachable by SSH.", hostname)
	}
	return true, nil
}

func (s *SSHClient) IsSshUserNameSet() (isSet bool) {
	return len(s.sshUserName) > 0
}

func (s *SSHClient) MustCheckReachable(verbose bool) {
	err := s.CheckReachable(verbose)
	if err != nil {
		LogFatalf("SshClient.CheckReachableBySsh failed: '%v'", err)
	}
}

func (s *SSHClient) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := s.GetHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostDescription
}

func (s *SSHClient) MustGetHostName() (hostName string) {
	hostName, err := s.GetHostName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostName
}

func (s *SSHClient) MustGetSshUserName() (sshUserName string) {
	sshUserName, err := s.GetSshUserName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshUserName
}

func (s *SSHClient) MustIsReachable(verbose bool) (isReachavble bool) {
	isReachavble, err := s.IsReachable(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isReachavble
}

func (s *SSHClient) MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput) {
	commandOutput, err := s.RunCommand(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (s *SSHClient) MustSetHostName(hostName string) {
	err := s.SetHostName(hostName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHClient) MustSetSshUserName(sshUserName string) {
	err := s.SetSshUserName(sshUserName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHClient) RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error) {
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

	commandString, err := ShellLineHandler().Join(options.Command)
	if err != nil {
		return nil, err
	}

	commandToUse := options.GetDeepCopy()
	commandToUse.Command = []string{
		"ssh",
		userAtHost,
		commandString,
	}

	commandOutput, err = Exec().RunCommand(commandToUse)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (s *SSHClient) SetHostName(hostName string) (err error) {
	if hostName == "" {
		return TracedErrorf("hostName is empty string")
	}

	s.hostName = hostName

	return nil
}

func (s *SSHClient) SetSshUserName(sshUserName string) (err error) {
	if len(sshUserName) <= 0 {
		return TracedError("sshUserName is nil")
	}

	s.sshUserName = sshUserName

	return nil
}
