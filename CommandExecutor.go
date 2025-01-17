package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetHostDescription() (hostDescription string, err error)
	RunCommand(options *parameteroptions.RunCommandOptions) (commandOutput *CommandOutput, err error)
	MustGetHostDescription() (hostDescription string)
	MustRunCommand(options *parameteroptions.RunCommandOptions) (commandOutput *CommandOutput)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error)
	MustIsRunningOnLocalhost() (isRunningOnLocalhost bool)
	MustRunCommandAndGetStdoutAsBytes(options *parameteroptions.RunCommandOptions) (stdout []byte)
	MustRunCommandAndGetStdoutAsFloat64(options *parameteroptions.RunCommandOptions) (stdout float64)
	MustRunCommandAndGetStdoutAsInt64(options *parameteroptions.RunCommandOptions) (stdout int64)
	MustRunCommandAndGetStdoutAsLines(options *parameteroptions.RunCommandOptions) (stdoutLines []string)
	MustRunCommandAndGetStdoutAsString(options *parameteroptions.RunCommandOptions) (stdout string)
	RunCommandAndGetStdoutAsBytes(options *parameteroptions.RunCommandOptions) (stdout []byte, err error)
	RunCommandAndGetStdoutAsFloat64(options *parameteroptions.RunCommandOptions) (stdout float64, err error)
	RunCommandAndGetStdoutAsInt64(options *parameteroptions.RunCommandOptions) (stdout int64, err error)
	RunCommandAndGetStdoutAsLines(options *parameteroptions.RunCommandOptions) (stdoutLines []string, err error)
	RunCommandAndGetStdoutAsString(options *parameteroptions.RunCommandOptions) (stdout string, err error)
}

func GetDeepCopyOfCommandExecutor(commandExectuor CommandExecutor) (copy CommandExecutor, err error) {
	if commandExectuor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	withDeepCopy, ok := commandExectuor.(interface{ GetDeepCopy() CommandExecutor })
	if !ok {
		typeName, err := datatypes.GetTypeName(commandExectuor)
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"CommandExecutor implementation '%s' has no GetDeepCopyFunction!",
			typeName,
		)
	}

	return withDeepCopy.GetDeepCopy(), nil
}

func MustGetDeepCopyOfCommandExecutor(commandExectuor CommandExecutor) (copy CommandExecutor) {
	copy, err := GetDeepCopyOfCommandExecutor(commandExectuor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return copy
}
