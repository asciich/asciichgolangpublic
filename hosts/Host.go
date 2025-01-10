package hosts

import "github.com/asciich/asciichgolangpublic"

// Host like a VM, Laptop, Desktop, Server.
type Host interface {
	GetHostDescription() (hostDescription string, err error)
	GetHostName() (hostName string, err error)
	MustGetHostDescription() (hostDescription string)
	MustGetHostName() (hostName string)
	MustRunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput)
	MustSetHostName(hostName string)
	RunCommand(runCommandOptions *asciichgolangpublic.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error)
	SetHostName(hostName string) (err error)

	// All methods below this line can be implemented by embedding the `CommandExecutorBase` struct:
	RunCommandAndGetStdoutAsString(runCommandOptions *asciichgolangpublic.RunCommandOptions) (stdout string, err error)
	MustRunCommandAndGetStdoutAsString(runCommandOptions *asciichgolangpublic.RunCommandOptions) (stdout string)
}
