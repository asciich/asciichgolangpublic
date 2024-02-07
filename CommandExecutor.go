package asciichgolangpublic

// A CommandExecutor is able to run a command like Exec or bash.
type CommandExecutor interface {
	RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error)
	MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	MustRunCommandAndGetStdoutAsBytes(options *RunCommandOptions) (stdout []byte)
	MustRunCommandAndGetStdoutAsFloat64(options *RunCommandOptions) (stdout float64)
	MustRunCommandandGetStdoutAsLines(options *RunCommandOptions) (stdoutLines []string)
	MustRunCommandAndGetStdoutAsString(options *RunCommandOptions) (stdout string)
	RunCommandAndGetStdoutAsBytes(options *RunCommandOptions) (stdout []byte, err error)
	RunCommandAndGetStdoutAsFloat64(options *RunCommandOptions) (stdout float64, err error)
	RunCommandandGetStdoutAsLines(options *RunCommandOptions) (stdoutLines []string, err error)
	RunCommandAndGetStdoutAsString(options *RunCommandOptions) (stdout string, err error)
}
