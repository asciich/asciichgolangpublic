package hosts

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/commandexecutorsshclient"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// Host like a VM, Laptop, Desktop, Server.
type Host interface {
	CheckReachable(verbose bool) (err error)

	GetDirectoryByPath(path string) (directory files.Directory, err error)
	GetHostDescription() (hostDescription string, err error)
	GetHostName() (hostName string, err error)
	GetSshPublicKeyOfUserAsString(ctx context.Context, username string) (publicKey string, err error)
	InstallBinary(installOptions *parameteroptions.InstallOptions) (installedFile files.File, err error)
	RunCommand(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutorgeneric.CommandOutput, err error)

	// All methods below this line can be implemented by embedding the `CommandExecutorBase` struct:
	RunCommandAndGetStdoutAsString(ctx context.Context, runCommandOptions *parameteroptions.RunCommandOptions) (stdout string, err error)
}

func GetHostByHostname(hostname string) (host Host, err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return nil, tracederrors.TracedError("hostname is empty string")
	}

	var commandExecutor commandexecutorinterfaces.CommandExecutor
	if hostname == "localhost" {
		commandExecutor = commandexecutor.Bash()
	} else {
		commandExecutor, err = commandexecutorsshclient.GetSshClientByHostName(hostname)
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
