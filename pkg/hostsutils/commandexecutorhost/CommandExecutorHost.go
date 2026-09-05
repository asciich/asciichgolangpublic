package commandexecutorhost

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/ftputils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
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
func GetCommandExecutorHostByCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (host hostsutilsinterfaces.Host, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	ret := NewCommandExecutorHost()

	err = ret.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func NewCommandExecutorHost() (c *CommandExecutorHost) {
	c = new(CommandExecutorHost)

	err := c.SetParentCommandExecutorForBaseClass(c)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return c
}

func (c *CommandExecutorHost) GetDeepCopy() *CommandExecutorHost {
	ret := &CommandExecutorHost{
		Comment: c.Comment,
	}

	if c.commandExecutor != nil {
		ret.commandExecutor = c.commandExecutor.GetDeepCopyAsCommandExecutor()
	}

	err := ret.SetParentCommandExecutorForBaseClass(ret)
	if err != nil {
		panic(err)
	}

	return ret
}

func (c *CommandExecutorHost) GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor {
	return c.GetDeepCopy()
}

func (c *CommandExecutorHost) GetFileInUsersHome(ctx context.Context, userName string, path string) (file filesinterfaces.File, err error) {
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

	for _, publicKeyBaseName := range []string{"id_ed25519.pub", "id_rsa.pub"} {
		sshKeyFile, err := c.GetFileInUsersHome(ctx, userName, ".ssh/"+publicKeyBaseName)
		if err != nil {
			return "", err
		}

		exists, err := sshKeyFile.Exists(ctx)
		if err != nil {
			return "", err
		}

		if exists {
			path, err := sshKeyFile.GetPath()
			if err != nil {
				return "", err
			}

			logging.LogInfoByCtxf(ctx, "SSH public key for user '%s' on host '%s' found in '%s'.", userName, hostDescription, path)

			return sshKeyFile.ReadAsString(ctx)
		}
	}

	return "", tracederrors.TracedErrorf("No SSH public key for user '%s' on host '%s' found.", userName, hostDescription)
}

func (c *CommandExecutorHost) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return c.commandExecutor, nil
}

func (c *CommandExecutorHost) GetDirectoryByPath(ctx context.Context, path string) (directory filesinterfaces.Directory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandexecutorfileoo.NewDirectory(commandExecutor, path)
}

func (c *CommandExecutorHost) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

// IsRunningOnLocalhost delegates to the wrapped command executor.
// This ensures that SSH clients (even to localhost) are correctly identified as non-local.
func (c *CommandExecutorHost) IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	return commandExecutor.IsRunningOnLocalhost()
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

func (h *CommandExecutorHost) AddSshHostKeyToKnownHosts(ctx context.Context) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	_, err = commandexecutorbashoo.Bash().RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				fmt.Sprintf("ssh-keyscan -H '%s' >> ${HOME}/.ssh/known_hosts", hostname),
			},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Added host key of '%s' from known hosts", hostname)

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

func (h *CommandExecutorHost) CheckReachable(ctx context.Context) (err error) {
	isReachable, err := h.IsReachable(ctx)
	if err != nil {
		return err
	}

	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	if isReachable {
		logging.LogInfoByCtxf(ctx, "Host '%s' is reachable by SSH.", hostname)
	} else {
		errorMessage := fmt.Sprintf("Host '%s' is reachable by SSH.", hostname)
		logging.LogErrorByCtxf(ctx, errorMessage)
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

func (h *CommandExecutorHost) GetFileByPath(path string) (file filesinterfaces.File, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	commandExecutor, err := h.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	file, err = commandexecutorfileoo.New(commandExecutor, path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (h *CommandExecutorHost) InstallBinary(ctx context.Context, installOptions *parameteroptions.InstallOptions) (installedFile filesinterfaces.File, err error) {
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

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' on host '%s' started.", sourceFilePath, binaryName, hostName)

	tempCopy, err := tempfilesoo.CreateTemporaryFileFromFile(ctx, sourceFile)
	if err != nil {
		return nil, err
	}

	destPath, err := installOptions.GetInstallationPathOrDefaultIfUnset()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "'%s' will be installed as '%s' on host '%s'.", binaryName, destPath, hostName)

	installedFile, err = tempCopy.MoveToPath(ctx, destPath, installOptions.UseSudoToInstall)
	if err != nil {
		return nil, err
	}

	err = installedFile.Chmod(ctx,
		&filesoptions.ChmodOptions{
			PermissionsString: "u=rwx,g=rx,o=rx",
			UseSudo:           installOptions.UseSudoToInstall,
		},
	)
	if err != nil {
		return nil, err
	}

	if installOptions.Group != "" || installOptions.Owner != "" {
		err = installedFile.Chown(
			ctx,
			&parameteroptions.ChownOptions{
				UserName:  installOptions.Owner,
				GroupName: installOptions.Group,
				UseSudo:   installOptions.UseSudoToInstall,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	logging.LogChangedByCtxf(ctx, "Install '%s' as '%s' on host '%s' finished.", sourceFilePath, binaryName, hostName)

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

func (h *CommandExecutorHost) IsReachable(ctx context.Context) (isReachable bool, err error) {
	_, err = h.RunCommandAndGetStdoutAsString(
		ctx,
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

func (h *CommandExecutorHost) RemoveSshHostKeyFromKnownHosts(ctx context.Context) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	_, err = commandexecutorbashoo.Bash().RunCommand(ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"ssh-keygen", "-R", hostname},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Removed host key of '%s' from known hosts", hostname)

	return nil
}

func (h *CommandExecutorHost) RenewSshHostKey(ctx context.Context) (err error) {
	err = h.RemoveSshHostKeyFromKnownHosts(ctx)
	if err != nil {
		return err
	}

	err = h.AddSshHostKeyToKnownHosts(ctx)
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

func (h *CommandExecutorHost) WaitUntilReachable(ctx context.Context, renewHostKey bool) (err error) {
	hostname, err := h.GetHostName()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		if renewHostKey {
			err = h.RenewSshHostKey(ctx)
			if err != nil {
				logging.LogWarn("Renewing host key failed, but error is ignored in WaitUntilReachableBySsh since running in a retry loop.")
			}
		}

		isReachableBySsh, err := h.IsReachable(ctx)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isReachableBySsh {
			logging.LogGoodByCtxf(ctx, "Host '%s' is reachable by SSH after '%v'", hostname, elapsedTime)
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not reachable by SSH after '%v'", hostname, elapsedTime)
			logging.LogErrorByCtx(ctx, errorMessage)
			return tracederrors.TracedError(errorMessage)
		}

		logging.LogInfoByCtxf(ctx,
			"Wait '%v' for host '%s' to get reachable by SSH. Total '%v' left, elapsed time so far: '%v'.",
			delayBetweenPings,
			hostname,
			timeout-elapsedTime,
			elapsedTime,
		)
	}
}

func (j *CommandExecutorHost) GetHostName() (hostName string, err error) {
	return j.GetHostDescription()
}

func (c *CommandExecutorHost) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	commandExecuor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecuor.RunCommandAndGetStdoutAsIoReadCloser(ctx, options)
}

func (c *CommandExecutorHost) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	commandExecuor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecuor.RunCommandAndGetStdinAsIoWriteCloser(ctx, options)
}

func (c *CommandExecutorHost) GetCPUArchitecture(ctx context.Context) (string, error) {
	commandExecuor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecuor.GetCPUArchitecture(ctx)
}
