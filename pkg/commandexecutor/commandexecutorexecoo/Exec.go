package commandexecutorexecoo

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

type ExecService struct {
	commandexecutorgeneric.CommandExecutorBase
}

func Exec() (e *ExecService) {
	return NewExec()
}

func NewExec() (e *ExecService) {
	e = new(ExecService)
	err := e.SetParentCommandExecutorForBaseClass(e)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
	return e
}

func NewExecService() (e *ExecService) {
	return new(ExecService)
}

func (e *ExecService) GetDeepCopyAsCommandExecutor() (deepCopy commandexecutorinterfaces.CommandExecutor) {
	d := NewExec()
	*d = *e
	deepCopy = d
	return deepCopy
}

func (e *ExecService) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
}

func (e *ExecService) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	return commandexecutorexec.RunCommand(ctx, options)
}

func (e *ExecService) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	return commandexecutorexec.RunCommandAndGetStdoutAsIoReadCloser(ctx, options)
}

func (e *ExecService) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	return commandexecutorexec.RunCommandAndGetStdinAsIoWriteCloser(ctx, options)
}