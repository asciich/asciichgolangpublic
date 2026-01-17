package commandexecutorinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetDeepCopyAsCommandExecutor() CommandExecutor

	GetHostDescription() (string, error)
	RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (bool, error)
	RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]byte, error)
	RunCommandAndGetStdoutAsFloat64(ctx context.Context, options *parameteroptions.RunCommandOptions) (float64, error)
	RunCommandAndGetStdoutAsInt64(ctx context.Context, options *parameteroptions.RunCommandOptions) (int64, error)
	RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]string, error)
	RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (string, error)
}
