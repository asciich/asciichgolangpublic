package commandexecutor

import (
	"context"

	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A CommandExecutor is able to run a command like Exec or bash does.
type CommandExecutor interface {
	GetHostDescription() (hostDescription string, err error)
	RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *CommandOutput, err error)

	// These Commands can be implemented by embedding the `CommandExecutorBase` struct:
	IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error)
	RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout []byte, err error)
	RunCommandAndGetStdoutAsFloat64(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout float64, err error)
	RunCommandAndGetStdoutAsInt64(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout int64, err error)
	RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdoutLines []string, err error)
	RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout string, err error)
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

type ContextKeyLiveOutputOnStdout struct{}

func WithLiveOutputOnStdoutIfVerbose(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = contextutils.ContextSilent()
	}

	return WithLiveOutputOnStdoutEnabled(ctx, contextutils.GetVerboseFromContext(ctx))
}

func WithLiveOutputOnStdoutEnabled(ctx context.Context, enabled bool) context.Context {
	if ctx == nil {
		ctx = contextutils.ContextVerbose()
	}

	return context.WithValue(ctx, ContextKeyLiveOutputOnStdout{}, enabled)
}

func WithLiveOutputOnStdout(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = contextutils.ContextVerbose()
	}

	return WithLiveOutputOnStdoutEnabled(ctx, true)
}

func IsLiveOutputOnStdoutEnabled(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	val := ctx.Value(ContextKeyLiveOutputOnStdout{})
	if val == nil {
		return false
	}

	valBool, ok := val.(bool)
	if !ok {
		return false
	}

	return valBool
}
