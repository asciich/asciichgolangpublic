package hostsutilsinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// Host like a VM, Laptop, Desktop, Server.
type Host interface {
	commandexecutorinterfaces.CommandExecutor

	CheckReachable(verbose bool) (err error)

	GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor

	GetDirectoryByPath(ctx context.Context, path string) (directory filesinterfaces.Directory, err error)
	GetHostDescription() (hostDescription string, err error)
	GetHostName() (hostName string, err error)
	GetSshPublicKeyOfUserAsString(ctx context.Context, username string) (publicKey string, err error)
	InstallBinary(ctx context.Context, installOptions *parameteroptions.InstallOptions) (installedFile filesinterfaces.File, err error)
}
