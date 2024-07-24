package asciichgolangpublic

import (
	"fmt"
	"strings"
)

type SSHClient struct {
	host        *Host
	sshUserName string
}

func NewSSHClient() (s *SSHClient) {
	return new(SSHClient)
}

func (s *SSHClient) CheckReachable(verbose bool) (isReachable bool, err error) {
	hostname, err := s.GetHostName()
	if err != nil {
		return false, err
	}

	isReachable, err = s.IsReachable(verbose)
	if err != nil {
		return false, err
	}

	if isReachable {
		return true, nil
	}

	errorMessage := fmt.Sprintf("host '%v' is not reachable", hostname)
	if verbose {
		LogError(errorMessage)
	}
	return false, TracedError(errorMessage)
}

func (s *SSHClient) GetHost() (host *Host, err error) {
	if s.host == nil {
		return nil, TracedError("host not set")
	}

	return s.host, nil
}

func (s *SSHClient) GetHostName() (hostname string, err error) {
	host, err := s.GetHost()
	if err != nil {
		return "", err
	}

	hostname, err = host.GetHostname()
	if err != nil {
		return "", err
	}

	return hostname, nil
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

func (s *SSHClient) MustCheckReachable(verbose bool) (isReachable bool) {
	isReachable, err := s.CheckReachable(verbose)
	if err != nil {
		LogFatalf("SshClient.CheckReachableBySsh failed: '%v'", err)
	}

	return isReachable
}

func (s *SSHClient) MustGetHost() (host *Host) {
	host, err := s.GetHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return host
}

func (s *SSHClient) MustGetHostName() (hostname string) {
	hostname, err := s.GetHostName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostname
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

func (s *SSHClient) MustSetHost(host *Host) {
	err := s.SetHost(host)
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

func (s *SSHClient) SetHost(host *Host) (err error) {
	if host == nil {
		return TracedError("host is nil")
	}

	s.host = host

	return nil
}

func (s *SSHClient) SetSshUserName(sshUserName string) (err error) {
	if len(sshUserName) <= 0 {
		return TracedError("sshUserName is nil")
	}

	s.sshUserName = sshUserName

	return nil
}
