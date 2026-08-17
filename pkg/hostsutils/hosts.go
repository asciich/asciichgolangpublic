package hostsutils

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhost"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/nativehost"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils/commandexecutorsshclient"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetHostByHostname(hostname string) (host hostsutilsinterfaces.Host, err error) {
	hostname = strings.TrimSpace(hostname)
	if len(hostname) <= 0 {
		return nil, tracederrors.TracedError("hostname is empty string")
	}

	if hostname == "localhost" {
		return nativehost.NewNativeHost(), nil
	}

	var commandExecutor commandexecutorinterfaces.CommandExecutor
	commandExecutor, err = commandexecutorsshclient.GetSshClientByHostName(hostname)
	if err != nil {
		return nil, err
	}

	return commandexecutorhost.GetCommandExecutorHostByCommandExecutor(commandExecutor)
}

func GetLocalCommandExecutorHost() (host hostsutilsinterfaces.Host, err error) {
	commandExecutor := commandexecutorbashoo.Bash()
	return commandexecutorhost.GetCommandExecutorHostByCommandExecutor(commandExecutor)
}

func GetLocalHost() (host hostsutilsinterfaces.Host, err error) {
	return nativehost.NewNativeHost(), nil
}
