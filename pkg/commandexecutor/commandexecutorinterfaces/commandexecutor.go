package commandexecutorinterfaces

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetDeepCopyAsCommandExecutor() CommandExecutor

	GetHostDescription() (string, error)
	
	// Run a command, wait until it's finished and get the whole output as CommandOutput. 
	RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error)

	// Run a command in background giving you the possibility to read the stdout as stream.
	RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error)

	// Run a command in background givining you the possibility to write the stdout as stream.
	RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (bool, error)
	GetCPUArchitecture(ctx context.Context) (string, error)
	RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]byte, error)
	RunCommandAndGetStdoutAsFloat64(ctx context.Context, options *parameteroptions.RunCommandOptions) (float64, error)
	RunCommandAndGetStdoutAsInt64(ctx context.Context, options *parameteroptions.RunCommandOptions) (int64, error)
	RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]string, error)
	RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (string, error)
}
