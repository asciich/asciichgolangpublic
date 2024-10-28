package asciichgolangpublic

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetDeepCopy() (deepCopy CommandExecutor)
	GetHostDescription() (hostDescription string, err error)
	RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error)
	MustGetHostDescription() (hostDescription string)
	MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
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
