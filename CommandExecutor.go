package asciichgolangpublic

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetHostDescription() (hostDescription string, err error)
	RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error)
	MustGetHostDescription() (hostDescription string)
	MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error)
	MustIsRunningOnLocalhost() (isRunningOnLocalhost bool)
	MustRunCommandAndGetStdoutAsBytes(options *RunCommandOptions) (stdout []byte)
	MustRunCommandAndGetStdoutAsFloat64(options *RunCommandOptions) (stdout float64)
	MustRunCommandAndGetStdoutAsInt64(options *RunCommandOptions) (stdout int64)
	MustRunCommandAndGetStdoutAsLines(options *RunCommandOptions) (stdoutLines []string)
	MustRunCommandAndGetStdoutAsString(options *RunCommandOptions) (stdout string)
	RunCommandAndGetStdoutAsBytes(options *RunCommandOptions) (stdout []byte, err error)
	RunCommandAndGetStdoutAsFloat64(options *RunCommandOptions) (stdout float64, err error)
	RunCommandAndGetStdoutAsInt64(options *RunCommandOptions) (stdout int64, err error)
	RunCommandAndGetStdoutAsLines(options *RunCommandOptions) (stdoutLines []string, err error)
	RunCommandAndGetStdoutAsString(options *RunCommandOptions) (stdout string, err error)
}

func GetDeepCopyOfCommandExecutor(commandExectuor CommandExecutor) (copy CommandExecutor, err error) {
	if commandExectuor == nil {
		return nil, TracedErrorNil("commandExecutor")
	}

	withDeepCopy, ok := commandExectuor.(interface{ GetDeepCopy() CommandExecutor })
	if !ok {
		typeName, err := Types().GetTypeName(commandExectuor)
		if err != nil {
			return nil, err
		}

		return nil, TracedErrorf(
			"CommandExecutor implementation '%s' has no GetDeepCopyFunction!",
			typeName,
		)
	}

	return withDeepCopy.GetDeepCopy(), nil
}

func MustGetDeepCopyOfCommandExecutor(commandExectuor CommandExecutor) (copy CommandExecutor) {
	copy, err := GetDeepCopyOfCommandExecutor(commandExectuor)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return copy
}
