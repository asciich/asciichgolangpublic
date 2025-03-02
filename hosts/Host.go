package hosts

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Host like a VM, Laptop, Desktop, Server.
type Host interface {
	CheckReachable(verbose bool) (err error)

	GetDirectoryByPath(path string) (directory files.Directory, err error)
	GetHostDescription() (hostDescription string, err error)
	GetHostName() (hostName string, err error)
	GetSshPublicKeyOfUser(ctx context.Context, username string) (publicKey string, err error)
	InstallBinary(installOptions *parameteroptions.InstallOptions) (installedFile files.File, err error)
	MustCheckReachable(verbose bool)
	MustGetDirectoryByPath(path string) (directory files.Directory)
	MustGetHostDescription() (hostDescription string)
	MustGetHostName() (hostName string)
	MustInstallBinary(installOptions *parameteroptions.InstallOptions) (installedFile files.File)
	MustRunCommand(runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput)
	RunCommand(runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput, err error)

	// All methods below this line can be implemented by embedding the `CommandExecutorBase` struct:
	RunCommandAndGetStdoutAsString(runCommandOptions *parameteroptions.RunCommandOptions) (stdout string, err error)
	MustRunCommandAndGetStdoutAsString(runCommandOptions *parameteroptions.RunCommandOptions) (stdout string)
}

func GetHostByHostname(hostname string) (host Host, err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return nil, tracederrors.TracedError("hostname is empty string")
	}

	var commandExecutor commandexecutor.CommandExecutor
	if hostname == "localhost" {
		commandExecutor = commandexecutor.Bash()
	} else {
		commandExecutor, err = asciichgolangpublic.GetSshClientByHostName(hostname)
		if err != nil {
			return nil, err
		}
	}

	return GetCommandExecutorHostByCommandExecutor(commandExecutor)
}

func GetLocalCommandExecutorHost() (host Host, err error) {
	return GetHostByHostname("localhost")
}

func GetLocalHost() (host Host, err error) {
	return GetHostByHostname("localhost")
}

func MustGetHostByHostname(hostname string) (host Host) {
	host, err := GetHostByHostname(hostname)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return host
}

func MustGetLocalCommandExecutorHost() (host Host) {
	host, err := GetLocalCommandExecutorHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return host
}

func MustGetLocalHost() (host Host) {
	host, err := GetLocalHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return host
}
