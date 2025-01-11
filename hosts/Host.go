package hosts

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
)

// Host like a VM, Laptop, Desktop, Server.
type Host interface {
	CheckReachable(verbose bool) (err error)

	GetDirectoryByPath(path string) (directory asciichgolangpublic.Directory, err error)
	GetHostDescription() (hostDescription string, err error)
	GetHostName() (hostName string, err error)
	InstallBinary(installOptions *asciichgolangpublic.InstallOptions) (installedFile asciichgolangpublic.File, err error)
	MustCheckReachable(verbose bool)
	MustGetDirectoryByPath(path string) (directory asciichgolangpublic.Directory)
	MustGetHostDescription() (hostDescription string)
	MustGetHostName() (hostName string)
	MustInstallBinary(installOptions *asciichgolangpublic.InstallOptions) (installedFile asciichgolangpublic.File)
	MustRunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput)
	RunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error)

	// All methods below this line can be implemented by embedding the `CommandExecutorBase` struct:
	RunCommandAndGetStdoutAsString(runCommandOptions *asciichgolangpublic.RunCommandOptions) (stdout string, err error)
	MustRunCommandAndGetStdoutAsString(runCommandOptions *asciichgolangpublic.RunCommandOptions) (stdout string)
}

func GetHostByHostname(hostname string) (host Host, err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return nil, asciichgolangpublic.TracedError("hostname is empty string")
	}

	var commandExecutor asciichgolangpublic.CommandExecutor
	if hostname == "localhost" {
		commandExecutor = asciichgolangpublic.Bash()
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

func MustGetLocalHost() (host Host) {
	host, err := GetLocalHost()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return host
}
