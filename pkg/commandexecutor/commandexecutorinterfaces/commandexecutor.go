package commandexecutorinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetHostDescription() (hostDescription string, err error)
	RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error)
	RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout []byte, err error)
	RunCommandAndGetStdoutAsFloat64(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout float64, err error)
	RunCommandAndGetStdoutAsInt64(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout int64, err error)
	RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdoutLines []string, err error)
	RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout string, err error)
}
