package hosts

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/ftputils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorHost struct {
	commandexecutorgeneric.CommandExecutorBase
	commandExecutor commandexecutorinterfaces.CommandExecutor
	Comment         string
}

// Get a Host by a CommandExecutor capable of executing commands on the Host.
// E.g. for SSH a SSHCLient can be used.
func GetCommandExecutorHostByCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (host Host, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	toReturn := NewCommandExecutorHost()

	err = toReturn.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func MustGetCommandExecutorHostByCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (host Host) {
	host, err := GetCommandExecutorHostByCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return host
}

func NewCommandExecutorHost() (c *CommandExecutorHost) {
	c = new(CommandExecutorHost)

	err := c.SetParentCommandExecutorForBaseClass(c)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return c
}

func (c *CommandExecutorHost) GetFileInUsersHome(ctx context.Context, userName string, path string) (file files.File, err error) {
	if userName == "" {
		return nil, tracederrors.TracedErrorEmptyString("userName")
	}

	fullPath := filepath.Join("/home", userName, path)

	return c.GetFileByPath(fullPath)
}

func (c *CommandExecutorHost) GetSshPublicKeyOfUserAsString(ctx context.Context, userName string) (publicKey string, err error) {
	if userName == "" {
		return "", tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Getting SSH public key of user '%s' on host '%s'.", userName, hostDescription)

	verbose := contextutils.GetVerboseFromContext(ctx)

	for _, publicKeyBaseName := range []string{"id_ed25519.pub", "id_rsa.pub"} {
		sshKeyFile, err := c.GetFileInUsersHome(ctx, userName, ".ssh/"+publicKeyBaseName)
		if err != nil {
			return "", err
		}

		exists, err := sshKeyFile.Exists(verbose)
		if err != nil {
			return "", err
		}

		if exists {
			path, err := sshKeyFile.GetPath()
			if err != nil {
				return "", err
			}

			logging.LogInfoByCtxf(ctx, "SSH public key for user '%s' on host '%s' found in '%s'.", userName, hostDescription, path)

			return sshKeyFile.ReadAsString()
		}
	}

	return "", tracederrors.TracedErrorf("No SSH public key for user '%s' on host '%s' found.", userName, hostDescription)
}

func (c *CommandExecutorHost) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *CommandExecutorHost) GetDirectoryByPath(path string) (directory files.Directory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return files.GetCommandExecutorDirectoryByPath(commandExecutor, path)
}

func (c *CommandExecutorHost) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

func (c *CommandExecutorHost) MustGetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorHost) MustIsReachable(verbose bool) (isReachable bool) {
	isReachable, err := c.IsReachable(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isReachable
}

func (c *CommandExecutorHost) MustSetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorHost) MustWaitUntilReachable(renewHostKey bool, verbose bool) {
	err := c.WaitUntilReachable(renewHostKey, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorHost) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, options)
}

func (c *CommandExecutorHost) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (h *CommandExecutorHost) AddSshHostKeyToKnownHosts(verbose bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	_, err = commandexecutorbashoo.Bash().RunCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				fmt.Sprintf("ssh-keyscan -H '%s' >> ${HOME}/.ssh/known_hosts", hostname),
			},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Added host key of '%s' from known hosts", hostname)
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
			logging.LogInfof("FTP port on host '%s' is open.", hostname)
		}
	} else {
		errorMessage := fmt.Sprintf("FTP port on host '%s' is not open.", hostname)
		if verbose {
			logging.LogError(errorMessage)
		}

		return tracederrors.TracedError(errorMessage)
	}

	return nil
}

func (h *CommandExecutorHost) CheckReachable(verbose bool) (err error) {
	isReachable, err := h.IsReachable(verbose)
	if err != nil {
		return err
	}

	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	if isReachable {
		if verbose {
			logging.LogInfof("Host '%s' is reachable by SSH.", hostname)
		}
	} else {
		errorMessage := fmt.Sprintf("Host '%s' is reachable by SSH.", hostname)
		if verbose {
			logging.LogError(errorMessage)
		}

		return tracederrors.TracedError(errorMessage)
	}

	return nil
}

func (h *CommandExecutorHost) GetComment() (comment string, err error) {
	if h.Comment == "" {
		return "", tracederrors.TracedErrorf("Comment not set")
	}

	return h.Comment, nil
}

func (h *CommandExecutorHost) GetFileByPath(path string) (file files.File, err error) {
	if path == "" {
		return nil, err
	}

	commandExecutor, err := h.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	file, err = files.GetCommandExecutorFileByPath(commandExecutor, path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (h *CommandExecutorHost) InstallBinary(installOptions *parameteroptions.InstallOptions) (installedFile files.File, err error) {
	if installOptions == nil {
		return nil, tracederrors.TracedErrorNil("installOptions")
	}

	hostName, err := h.GetHostName()
	if err != nil {
		return nil, err
	}

	sourceFilePath, err := installOptions.GetSourcePath()
	if err != nil {
		return nil, err
	}

	sourceFile, err := h.GetFileByPath(sourceFilePath)
	if err != nil {
		return nil, err
	}

	binaryName, err := installOptions.GetBinaryName()
	if err != nil {
		return nil, err
	}

	if installOptions.Verbose {
		logging.LogInfof(
			"Install '%s' as '%s' on host '%s' started.",
			sourceFilePath,
			binaryName,
			hostName,
		)
	}

	tempCopy, err := tempfiles.CreateTemporaryFileFromFile(sourceFile, installOptions.Verbose)
	if err != nil {
		return nil, err
	}

	destPath, err := installOptions.GetInstallationPathOrDefaultIfUnset()
	if err != nil {
		return nil, err
	}

	if installOptions.Verbose {
		logging.LogInfof(
			"'%s' will be installed as '%s' on host '%s'.",
			binaryName,
			destPath,
			hostName,
		)
	}

	installedFile, err = tempCopy.MoveToPath(destPath, installOptions.UseSudoToInstall, installOptions.Verbose)
	if err != nil {
		return nil, err
	}

	err = installedFile.Chmod(
		&parameteroptions.ChmodOptions{
			PermissionsString: "u=rwx,g=rx,o=rx",
			UseSudo:           installOptions.UseSudoToInstall,
			Verbose:           installOptions.Verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	err = installedFile.Chown(
		&parameteroptions.ChownOptions{
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
		logging.LogChangedf(
			"Install '%s' as '%s' on host '%s' finished.",
			sourceFilePath,
			binaryName,
			hostName,
		)
	}

	return installedFile, nil
}

func (h *CommandExecutorHost) IsFtpPortOpen(verbose bool) (isOpen bool, err error) {
	isOpen, err = h.IsTcpPortOpen(ftputils.DEFAUT_PORT, verbose)
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

	stdout, err := commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsString(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
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

	return false, tracederrors.TracedErrorf("Unexpected stdout: '%v'", stdout)
}

func (h *CommandExecutorHost) IsReachable(verbose bool) (isReachable bool, err error) {
	_, err = h.RunCommandAndGetStdoutAsString(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{"echo", "hello"},
		},
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *CommandExecutorHost) IsTcpPortOpen(portNumber int, verbose bool) (isOpen bool, err error) {
	if portNumber <= 0 {
		return false, tracederrors.TracedErrorf("Invalid portNumber: '%d'", portNumber)
	}

	hostname, err := h.GetHostName()
	if err != nil {
		return false, err
	}

	isOpen, err = netutils.IsTcpPortOpen(contextutils.GetVerbosityContextByBool(verbose), hostname, portNumber)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (h *CommandExecutorHost) MustAddSshHostKeyToKnownHosts(verbose bool) {
	err := h.AddSshHostKeyToKnownHosts(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustCheckFtpPortOpen(verbose bool) {
	err := h.CheckFtpPortOpen(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustCheckReachable(verbose bool) {
	err := h.CheckReachable(verbose)
	if err != nil {
		logging.LogFatalf("host.CheckReachableBySsh failed: '%v'", err)
	}
}

func (h *CommandExecutorHost) MustGetComment() (comment string) {
	comment, err := h.GetComment()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return comment
}

func (h *CommandExecutorHost) MustGetDirectoryByPath(path string) (directory files.Directory) {
	directory, err := h.GetDirectoryByPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return directory
}

func (h *CommandExecutorHost) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := h.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (h *CommandExecutorHost) MustGetHostName() (hostname string) {
	hostname, err := h.GetHostName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostname
}

func (h *CommandExecutorHost) MustInstallBinary(installOptions *parameteroptions.InstallOptions) (installedFile files.File) {
	installedFile, err := h.InstallBinary(installOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return installedFile
}

func (h *CommandExecutorHost) MustIsFtpPortOpen(verbose bool) (isOpen bool) {
	isOpen, err := h.IsFtpPortOpen(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isOpen
}

func (h *CommandExecutorHost) MustIsPingable(verbose bool) (isPingable bool) {
	isPingable, err := h.IsPingable(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isPingable
}

func (h *CommandExecutorHost) MustIsTcpPortOpen(portNumber int, verbose bool) (isOpen bool) {
	isOpen, err := h.IsTcpPortOpen(portNumber, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isOpen
}

func (h *CommandExecutorHost) MustRemoveSshHostKeyFromKnownHosts(verbose bool) {
	err := h.RemoveSshHostKeyFromKnownHosts(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustRenewSshHostKey(verbose bool) {
	err := h.RenewSshHostKey(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustSetComment(comment string) {
	err := h.SetComment(comment)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) MustWaitUntilPingable(verbose bool) {
	err := h.WaitUntilPingable(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (h *CommandExecutorHost) RemoveSshHostKeyFromKnownHosts(verbose bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	_, err = commandexecutorbashoo.Bash().RunCommand(
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{"ssh-keygen", "-R", hostname},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Removed host key of '%s' from known hosts", hostname)
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

func (h *CommandExecutorHost) SetComment(comment string) (err error) {
	if comment == "" {
		return tracederrors.TracedErrorf("comment is empty string")
	}

	h.Comment = comment

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
				logging.LogGoodf("Host '%s' is pingable after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not pingable after '%v'", hostname, elapsedTime)
			if verbose {
				logging.LogError(errorMessage)
			}
			return tracederrors.TracedError(errorMessage)
		}

		if verbose {
			logging.LogInfof(
				"Wait '%v' for host '%s' to get reachable by ping. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}
	}
}

func (h *CommandExecutorHost) WaitUntilReachable(renewHostKey bool, verbose bool) (err error) {
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
				logging.LogWarn("Renewing host key failed, but error is ignored in WaitUntilReachableBySsh since running in a retry loop.")
			}
		}

		isReachableBySsh, err := h.IsReachable(verbose)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isReachableBySsh {
			if verbose {
				logging.LogGoodf("Host '%s' is reachable by SSH after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not reachable by SSH after '%v'", hostname, elapsedTime)
			if verbose {
				logging.LogError(errorMessage)
			}
			return tracederrors.TracedError(errorMessage)
		}

		if verbose {
			logging.LogInfof(
				"Wait '%v' for host '%s' to get reachable by SSH. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}
	}
}

func (j *CommandExecutorHost) GetHostName() (hostName string, err error) {
	return j.GetHostDescription()
}
