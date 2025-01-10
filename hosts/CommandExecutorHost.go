package hosts

import (
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic"
)

type CommandExecutorHost struct {
	asciichgolangpublic.CommandExecutorBase
	hostname    string
	sshUsername string
	Comment     string
}

func GetHostByHostname(hostname string) (host *CommandExecutorHost, err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return nil, asciichgolangpublic.TracedError("hostname is empty string")
	}

	host = NewHost()
	err = host.SetHostName(hostname)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func GetLocalCommandExecutorHost() (host Host, err error) {
	return GetHostByHostname("localhost")
}

func MustGetHostByHostname(hostname string) (host *CommandExecutorHost) {
	host, err := GetHostByHostname(hostname)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return host
}

func MustGetLocalCommandExecutorHost() (host Host) {
	host, err := GetLocalCommandExecutorHost()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return host
}

func NewCommandExecutorHost() (c *CommandExecutorHost) {
	return new(CommandExecutorHost)
}

func NewHost() (host *CommandExecutorHost) {
	host = new(CommandExecutorHost)

	host.MustSetParentCommandExecutorForBaseClass(host)

	return host
}

func (c *CommandExecutorHost) GetHostname() (hostname string, err error) {
	if c.hostname == "" {
		return "", asciichgolangpublic.TracedErrorf("hostname not set")
	}

	return c.hostname, nil
}

func (c *CommandExecutorHost) MustGetHostname() (hostname string) {
	hostname, err := c.GetHostname()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return hostname
}

func (c *CommandExecutorHost) MustSetHostname(hostname string) {
	err := c.SetHostname(hostname)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorHost) SetHostname(hostname string) (err error) {
	if hostname == "" {
		return asciichgolangpublic.TracedErrorf("hostname is empty string")
	}

	c.hostname = hostname

	return nil
}

func (h *CommandExecutorHost) AddSshHostKeyToKnownHosts(verbose bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	_, err = asciichgolangpublic.Bash().RunCommand(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{
				fmt.Sprintf("ssh-keyscan -H '%s' >> ${HOME}/.ssh/known_hosts", hostname),
			},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		asciichgolangpublic.LogInfof("Added host key of '%s' from known hosts", hostname)
	}

	return nil
}

func (h *CommandExecutorHost) CheckFtpPortOpen(verbose bool) (err error) {
	isOpen, err := h.IsFtpPortOpen(verbose)
	if err != nil {
		return err
	}

	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	if isOpen {
		if verbose {
			asciichgolangpublic.LogInfof("FTP port on host '%s' is open.", hostname)
		}
	} else {
		errorMessage := fmt.Sprintf("FTP port on host '%s' is not open.", hostname)
		if verbose {
			asciichgolangpublic.LogError(errorMessage)
		}

		return asciichgolangpublic.TracedError(errorMessage)
	}

	return nil
}

func (h *CommandExecutorHost) CheckReachableBySsh(verbose bool) (err error) {
	isReachable, err := h.IsReachableBySsh(verbose)
	if err != nil {
		return err
	}

	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	if isReachable {
		if verbose {
			asciichgolangpublic.LogInfof("Host '%s' is reachable by SSH.", hostname)
		}
	} else {
		errorMessage := fmt.Sprintf("Host '%s' is reachable by SSH.", hostname)
		if verbose {
			asciichgolangpublic.LogError(errorMessage)
		}

		return asciichgolangpublic.TracedError(errorMessage)
	}

	return nil
}

func (h *CommandExecutorHost) GetComment() (comment string, err error) {
	if h.Comment == "" {
		return "", asciichgolangpublic.TracedErrorf("Comment not set")
	}

	return h.Comment, nil
}

func (h *CommandExecutorHost) GetDeepCopy() (deepCopy asciichgolangpublic.CommandExecutor) {
	d := NewHost()

	*d = *h

	deepCopy = d

	return deepCopy
}

func (h *CommandExecutorHost) GetDirectoryByPath(path string) (directory asciichgolangpublic.Directory, err error) {
	if path == "" {
		return nil, asciichgolangpublic.TracedErrorEmptyString("path")
	}

	commandExecutorDir, err := asciichgolangpublic.NewCommandExecutorDirectory(h)
	if err != nil {
		return nil, err
	}

	err = commandExecutorDir.SetDirPath(path)
	if err != nil {
		return nil, err
	}

	return commandExecutorDir, nil
}

func (h *CommandExecutorHost) GetHostDescription() (hostDescription string, err error) {
	hostDescription, err = h.GetHostName()
	if err != nil {
		return "", err
	}

	return hostDescription, nil
}

func (h *CommandExecutorHost) GetHostName() (hostname string, err error) {
	if len(h.hostname) <= 0 {
		return "", asciichgolangpublic.TracedError("hostname not set")
	}

	return h.hostname, nil
}

func (h *CommandExecutorHost) GetSSHClient() (sshClient *asciichgolangpublic.SSHClient, err error) {
	sshClient = asciichgolangpublic.NewSSHClient()

	hostname, err := h.GetHostName()
	if err != nil {
		return nil, err
	}

	err = sshClient.SetHostName(hostname)
	if err != nil {
		return nil, err
	}

	if h.IsSshUserNameSet() {
		sshUserName, err := h.GetSshUserName()
		if err != nil {
			return nil, err
		}

		err = sshClient.SetSshUserName(sshUserName)
		if err != nil {
			return nil, err
		}
	}

	return sshClient, nil
}

func (h *CommandExecutorHost) GetSshUserName() (sshUserName string, err error) {
	if !h.IsSshUserNameSet() {
		return "", asciichgolangpublic.TracedError("sshUserName not set")
	}

	return h.sshUsername, nil
}

func (h *CommandExecutorHost) GetSshUsername() (sshUsername string, err error) {
	if h.sshUsername == "" {
		return "", asciichgolangpublic.TracedErrorf("sshUsername not set")
	}

	return h.sshUsername, nil
}

func (h *CommandExecutorHost) InstallBinary(installOptions *asciichgolangpublic.InstallOptions) (installedFile asciichgolangpublic.File, err error) {
	if installOptions == nil {
		return nil, asciichgolangpublic.TracedErrorNil("installOptions")
	}

	hostName, err := h.GetHostName()
	if err != nil {
		return nil, err
	}

	sourceFile, err := installOptions.GetSourceFile()
	if err != nil {
		return nil, err
	}

	sourceFilePath, err := sourceFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	binaryName, err := installOptions.GetBinaryName()
	if err != nil {
		return nil, err
	}

	if installOptions.Verbose {
		asciichgolangpublic.LogInfof(
			"Install '%s' as '%s' on host '%s' started.",
			sourceFilePath,
			binaryName,
			hostName,
		)
	}

	tempCopy, err := asciichgolangpublic.TemporaryFiles().CreateTemporaryFileFromFile(sourceFile, installOptions.Verbose)
	if err != nil {
		return nil, err
	}

	destPath, err := installOptions.GetInstallationPathOrDefaultIfUnset()
	if err != nil {
		return nil, err
	}

	installedFile, err = tempCopy.MoveToPath(destPath, installOptions.UseSudoToInstall, installOptions.Verbose)
	if err != nil {
		return nil, err
	}

	err = installedFile.Chmod(
		&asciichgolangpublic.ChmodOptions{
			PermissionsString: "u=rwx,g=rx,o=rx",
			UseSudo:           installOptions.UseSudoToInstall,
			Verbose:           installOptions.Verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	err = installedFile.Chown(
		&asciichgolangpublic.ChownOptions{
			UserName:  "root",
			GroupName: "root",
			UseSudo:   installOptions.UseSudoToInstall,
			Verbose:   installOptions.Verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	if installOptions.Verbose {
		asciichgolangpublic.LogChangedf(
			"Install '%s' as '%s' on host '%s' finished.",
			sourceFilePath,
			binaryName,
			hostName,
		)
	}

	return installedFile, nil
}

func (h *CommandExecutorHost) IsFtpPortOpen(verbose bool) (isOpen bool, err error) {
	isOpen, err = h.IsTcpPortOpen(asciichgolangpublic.FTP().GetDefaultPort(), verbose)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (h *CommandExecutorHost) IsPingable(verbose bool) (isPingable bool, err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return false, err
	}

	stdout, err := asciichgolangpublic.Bash().RunCommandAndGetStdoutAsString(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{"bash", "-c", fmt.Sprintf("ping -c 1 '%s' &>/dev/null && echo yes || echo no", hostname)},
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)
	if stdout == "yes" {
		return true, nil
	}
	if stdout == "no" {
		return false, nil
	}

	return false, asciichgolangpublic.TracedErrorf("Unexpected stdout: '%v'", stdout)
}

func (h *CommandExecutorHost) IsReachableBySsh(verbose bool) (isReachable bool, err error) {
	sshClient, err := h.GetSSHClient()
	if err != nil {
		return false, err
	}

	isReachable, err = sshClient.IsReachable(verbose)
	if err != nil {
		return false, err
	}

	return isReachable, nil
}

func (h *CommandExecutorHost) IsSshUserNameSet() (isSet bool) {
	return len(h.sshUsername) > 0
}

func (h *CommandExecutorHost) IsTcpPortOpen(portNumber int, verbose bool) (isOpen bool, err error) {
	if portNumber <= 0 {
		return false, asciichgolangpublic.TracedErrorf("Invalid portNumber: '%d'", portNumber)
	}

	hostname, err := h.GetHostName()
	if err != nil {
		return false, err
	}

	isOpen, err = asciichgolangpublic.TcpPorts().IsPortOpen(hostname, portNumber, verbose)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (h *CommandExecutorHost) MustAddSshHostKeyToKnownHosts(verbose bool) {
	err := h.AddSshHostKeyToKnownHosts(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustCheckFtpPortOpen(verbose bool) {
	err := h.CheckFtpPortOpen(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustCheckReachableBySsh(verbose bool) {
	err := h.CheckReachableBySsh(verbose)
	if err != nil {
		asciichgolangpublic.LogFatalf("host.CheckReachableBySsh failed: '%v'", err)
	}
}

func (h *CommandExecutorHost) MustGetComment() (comment string) {
	comment, err := h.GetComment()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return comment
}

func (h *CommandExecutorHost) MustGetDirectoryByPath(path string) (directory asciichgolangpublic.Directory) {
	directory, err := h.GetDirectoryByPath(path)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return directory
}

func (h *CommandExecutorHost) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := h.GetHostDescription()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (h *CommandExecutorHost) MustGetHostName() (hostname string) {
	hostname, err := h.GetHostName()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return hostname
}

func (h *CommandExecutorHost) MustGetSSHClient() (sshClient *asciichgolangpublic.SSHClient) {
	sshClient, err := h.GetSSHClient()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return sshClient
}

func (h *CommandExecutorHost) MustGetSshUserName() (sshUserName string) {
	sshUserName, err := h.GetSshUserName()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return sshUserName
}

func (h *CommandExecutorHost) MustGetSshUsername() (sshUsername string) {
	sshUsername, err := h.GetSshUsername()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return sshUsername
}

func (h *CommandExecutorHost) MustInstallBinary(installOptions *asciichgolangpublic.InstallOptions) (installedFile asciichgolangpublic.File) {
	installedFile, err := h.InstallBinary(installOptions)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return installedFile
}

func (h *CommandExecutorHost) MustIsFtpPortOpen(verbose bool) (isOpen bool) {
	isOpen, err := h.IsFtpPortOpen(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isOpen
}

func (h *CommandExecutorHost) MustIsPingable(verbose bool) (isPingable bool) {
	isPingable, err := h.IsPingable(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isPingable
}

func (h *CommandExecutorHost) MustIsReachableBySsh(verbose bool) (isReachable bool) {
	isReachable, err := h.IsReachableBySsh(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isReachable
}

func (h *CommandExecutorHost) MustIsTcpPortOpen(portNumber int, verbose bool) (isOpen bool) {
	isOpen, err := h.IsTcpPortOpen(portNumber, verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isOpen
}

func (h *CommandExecutorHost) MustRemoveSshHostKeyFromKnownHosts(verbose bool) {
	err := h.RemoveSshHostKeyFromKnownHosts(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustRenewSshHostKey(verbose bool) {
	err := h.RenewSshHostKey(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustRunCommand(options *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput) {
	commandOutput, err := h.RunCommand(options)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (h *CommandExecutorHost) MustSetComment(comment string) {
	err := h.SetComment(comment)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustSetHostName(hostname string) {
	err := h.SetHostName(hostname)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustSetSshUserName(username string) {
	err := h.SetSshUserName(username)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustSetSshUsername(sshUsername string) {
	err := h.SetSshUsername(sshUsername)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustWaitUntilPingable(verbose bool) {
	err := h.WaitUntilPingable(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustWaitUntilReachableBySsh(renewHostKey bool, verbose bool) {
	err := h.WaitUntilReachableBySsh(renewHostKey, verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) RemoveSshHostKeyFromKnownHosts(verbose bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	_, err = asciichgolangpublic.Bash().RunCommand(
		&asciichgolangpublic.RunCommandOptions{
			Command: []string{"ssh-keygen", "-R", hostname},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		asciichgolangpublic.LogInfof("Removed host key of '%s' from known hosts", hostname)
	}

	return nil
}

func (h *CommandExecutorHost) RenewSshHostKey(verbose bool) (err error) {
	err = h.RemoveSshHostKeyFromKnownHosts(verbose)
	if err != nil {
		return err
	}

	err = h.AddSshHostKeyToKnownHosts(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (h *CommandExecutorHost) RunCommand(options *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error) {
	if options == nil {
		return nil, asciichgolangpublic.TracedErrorNil("options")
	}

	sshClient, err := h.GetSSHClient()
	if err != nil {
		return nil, err
	}

	commandOutput, err = sshClient.RunCommand(options)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (h *CommandExecutorHost) SetComment(comment string) (err error) {
	if comment == "" {
		return asciichgolangpublic.TracedErrorf("comment is empty string")
	}

	h.Comment = comment

	return nil
}

func (h *CommandExecutorHost) SetHostName(hostname string) (err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return asciichgolangpublic.TracedError("hostname is empty string")
	}

	h.hostname = hostname

	return nil
}

func (h *CommandExecutorHost) SetSshUserName(username string) (err error) {
	if len(username) <= 0 {
		return asciichgolangpublic.TracedError("username is empty string")
	}

	h.sshUsername = username

	return nil
}

func (h *CommandExecutorHost) SetSshUsername(sshUsername string) (err error) {
	if sshUsername == "" {
		return asciichgolangpublic.TracedErrorf("sshUsername is empty string")
	}

	h.sshUsername = sshUsername

	return nil
}

func (h *CommandExecutorHost) WaitUntilPingable(verbose bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		isPingable, err := h.IsPingable(verbose)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isPingable {
			if verbose {
				asciichgolangpublic.LogGoodf("Host '%s' is pingable after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not pingable after '%v'", hostname, elapsedTime)
			if verbose {
				asciichgolangpublic.LogError(errorMessage)
			}
			return asciichgolangpublic.TracedError(errorMessage)
		}

		if verbose {
			asciichgolangpublic.LogInfof(
				"Wait '%v' for host '%s' to get reachable by ping. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}
	}
}

func (h *CommandExecutorHost) WaitUntilReachableBySsh(renewHostKey bool, verbose bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		if renewHostKey {
			err = h.RenewSshHostKey(verbose)
			if err != nil {
				asciichgolangpublic.LogWarn("Renewing host key failed, but error is ignored in WaitUntilReachableBySsh since running in a retry loop.")
			}
		}

		isReachableBySsh, err := h.IsReachableBySsh(verbose)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isReachableBySsh {
			if verbose {
				asciichgolangpublic.LogGoodf("Host '%s' is reachable by SSH after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not reachable by SSH after '%v'", hostname, elapsedTime)
			if verbose {
				asciichgolangpublic.LogError(errorMessage)
			}
			return asciichgolangpublic.TracedError(errorMessage)
		}

		if verbose {
			asciichgolangpublic.LogInfof(
				"Wait '%v' for host '%s' to get reachable by SSH. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}
	}
}
