package nativehost

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type NativeHost struct {
	commandexecutorgeneric.CommandExecutorBase
}

func NewNativeHost() (host hostsutilsinterfaces.Host) {
	n := &NativeHost{}
	err := n.SetParentCommandExecutorForBaseClass(n)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return n
}

func (n *NativeHost) GetDeepCopy() *NativeHost {
	ret := &NativeHost{}
	err := ret.SetParentCommandExecutorForBaseClass(ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func (n *NativeHost) GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor {
	return n.GetDeepCopy()
}

func (n *NativeHost) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
}

func (n *NativeHost) GetHostName() (hostName string, err error) {
	return n.GetHostDescription()
}

func (n *NativeHost) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if len(options.Command) == 0 {
		return nil, tracederrors.TracedError("options.Command is empty")
	}

	cmd := exec.CommandContext(ctx, options.Command[0], options.Command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Create CommandOutput with error
		commandOutput = commandoutput.NewCommandOutput()
		err2 := commandOutput.SetStdoutByString(string(output))
		if err2 != nil {
			return nil, err2
		}
		commandOutput.SetCmdRunError(err)
		// Try to get return code from exit error
		if exitErr, ok := err.(*exec.ExitError); ok {
			commandOutput.SetReturnCode(exitErr.ExitCode())
		}
		return commandOutput, err
	}

	// Create CommandOutput with success
	commandOutput = commandoutput.NewCommandOutput()
	err = commandOutput.SetStdoutByString(string(output))
	if err != nil {
		return nil, err
	}
	err = commandOutput.SetReturnCode(0)
	if err != nil {
		return nil, err
	}
	return commandOutput, nil
}

func (n *NativeHost) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if len(options.Command) == 0 {
		return nil, tracederrors.TracedError("options.Command is empty")
	}

	cmd := exec.CommandContext(ctx, options.Command[0], options.Command[1:]...)
	return cmd.StdoutPipe()
}

func (n *NativeHost) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if len(options.Command) == 0 {
		return nil, tracederrors.TracedError("options.Command is empty")
	}

	cmd := exec.CommandContext(ctx, options.Command[0], options.Command[1:]...)
	return cmd.StdinPipe()
}

func (n *NativeHost) GetDirectoryByPath(ctx context.Context, path string) (directory filesinterfaces.Directory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return nativefilesoo.NewDirectoryByPath(path)
}

func (n *NativeHost) GetFileByPath(path string) (file filesinterfaces.File, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return nativefilesoo.NewFileByPath(path)
}

func (n *NativeHost) GetFileInUsersHome(ctx context.Context, userName string, path string) (file filesinterfaces.File, err error) {
	if userName == "" {
		return nil, tracederrors.TracedErrorEmptyString("userName")
	}

	fullPath := filepath.Join("/home", userName, path)
	return n.GetFileByPath(fullPath)
}

func (n *NativeHost) GetSshPublicKeyOfUserAsString(ctx context.Context, userName string) (publicKey string, err error) {
	if userName == "" {
		return "", tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := n.GetHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Getting SSH public key of user '%s' on host '%s'.", userName, hostDescription)

	for _, publicKeyBaseName := range []string{"id_ed25519.pub", "id_rsa.pub"} {
		sshKeyFile, err := n.GetFileInUsersHome(ctx, userName, ".ssh/"+publicKeyBaseName)
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

func (n *NativeHost) CheckReachable(verbose bool) (err error) {
	hostDescription, err := n.GetHostDescription()
	if err != nil {
		return err
	}

	isReachable, err := n.IsReachable(verbose)
	if err != nil {
		return err
	}

	if !isReachable {
		return tracederrors.TracedErrorf("Host '%s' is not reachable", hostDescription)
	}

	return nil
}

func (n *NativeHost) IsReachable(verbose bool) (isReachable bool, err error) {
	_, err = n.RunCommandAndGetStdoutAsString(
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

func (n *NativeHost) InstallBinary(ctx context.Context, installOptions *parameteroptions.InstallOptions) (installedFile filesinterfaces.File, err error) {
	if installOptions == nil {
		return nil, tracederrors.TracedErrorNil("installOptions")
	}

	hostName, err := n.GetHostName()
	if err != nil {
		return nil, err
	}

	sourceFilePath, err := installOptions.GetSourcePath()
	if err != nil {
		return nil, err
	}

	sourceFile, err := n.GetFileByPath(sourceFilePath)
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

func (n *NativeHost) IsPingable(verbose bool) (isPingable bool, err error) {
	hostname, err := n.GetHostName()
	if err != nil {
		return false, err
	}

	stdout, err := n.RunCommandAndGetStdoutAsString(
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

func (n *NativeHost) IsTcpPortOpen(portNumber int, verbose bool) (isOpen bool, err error) {
	if portNumber <= 0 {
		return false, tracederrors.TracedErrorf("Invalid portNumber: '%d'", portNumber)
	}

	hostname, err := n.GetHostName()
	if err != nil {
		return false, err
	}

	isOpen, err = netutils.IsTcpPortOpen(contextutils.GetVerbosityContextByBool(verbose), hostname, portNumber)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (n *NativeHost) IsFtpPortOpen(verbose bool) (isOpen bool, err error) {
	isOpen, err = n.IsTcpPortOpen(21, verbose)
	if err != nil {
		return false, err
	}

	return isOpen, nil
}

func (n *NativeHost) WaitUntilPingable(verbose bool) (err error) {
	hostname, err := n.GetHostName()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		isPingable, err := n.IsPingable(verbose)
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

		time.Sleep(delayBetweenPings)
	}
}

func (n *NativeHost) WaitUntilReachable(renewHostKey bool, verbose bool) (err error) {
	hostname, err := n.GetHostName()
	if err != nil {
		return err
	}

	t_start := time.Now()
	timeout := 60 * time.Second
	delayBetweenPings := 2 * time.Second

	for {
		isReachable, err := n.IsReachable(verbose)
		if err != nil {
			return nil
		}

		elapsedTime := time.Since(t_start)

		if isReachable {
			if verbose {
				logging.LogGoodf("Host '%s' is reachable after '%v'", hostname, elapsedTime)
			}
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not reachable after '%v'", hostname, elapsedTime)
			if verbose {
				logging.LogError(errorMessage)
			}
			return tracederrors.TracedError(errorMessage)
		}

		if verbose {
			logging.LogInfof(
				"Wait '%v' for host '%s' to get reachable. Total '%v' left, elapsed time so far: '%v'.",
				delayBetweenPings,
				hostname,
				timeout-elapsedTime,
				elapsedTime,
			)
		}

		time.Sleep(delayBetweenPings)
	}
}

func (n *NativeHost) CheckFtpPortOpen(verbose bool) (err error) {
	isOpen, err := n.IsFtpPortOpen(verbose)
	if err != nil {
		return err
	}

	hostname, err := n.GetHostName()
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

func (n *NativeHost) GetComment() (comment string, err error) {
	return "", tracederrors.TracedError("comment not implemented for NativeHost")
}

func (n *NativeHost) SetComment(comment string) (err error) {
	return tracederrors.TracedError("SetComment not implemented for NativeHost")
}

// AddSshHostKeyToKnownHosts is not applicable for NativeHost since it represents localhost.
// SSH connections to localhost use the bash command executor, not SSH protocol.
func (n *NativeHost) AddSshHostKeyToKnownHosts(verbose bool) (err error) {
	return tracederrors.TracedError("AddSshHostKeyToKnownHosts not applicable for NativeHost (localhost)")
}

// RemoveSshHostKeyFromKnownHosts is not applicable for NativeHost since it represents localhost.
func (n *NativeHost) RemoveSshHostKeyFromKnownHosts(verbose bool) (err error) {
	return tracederrors.TracedError("RemoveSshHostKeyFromKnownHosts not applicable for NativeHost (localhost)")
}

// RenewSshHostKey is not applicable for NativeHost since it represents localhost.
func (n *NativeHost) RenewSshHostKey(verbose bool) (err error) {
	return tracederrors.TracedError("RenewSshHostKey not applicable for NativeHost (localhost)")
}
