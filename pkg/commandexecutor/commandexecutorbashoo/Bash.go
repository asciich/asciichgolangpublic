package commandexecutorbashoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbash"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

type BashService struct {
	commandexecutorgeneric.CommandExecutorBase
}

// Can be used to run commands in bash on localhost.
func Bash() (b *BashService) {
	return NewBashService()
}

func NewBashService() (b *BashService) {
	b = new(BashService)
	b.SetParentCommandExecutorForBaseClass(b)
	return b
}

func (b *BashService) GetDeepCopyAsCommandExecutor() (deepCopy commandexecutorinterfaces.CommandExecutor) {
	d := NewBashService()

	*d = *b

	deepCopy = d

	return deepCopy
}

func (b *BashService) GetHostDescription() (hostDescription string, err error) {
	return "localhost", err
}

func (b *BashService) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	return commandexecutorbash.RunCommand(ctx, options)
}

func (b *BashService) RunOneLiner(ctx context.Context, oneLiner string) (output *commandoutput.CommandOutput, err error) {
	return commandexecutorbash.RunOneLiner(ctx, oneLiner)
}

func (b *BashService) RunOneLinerAndGetStdoutAsLines(ctx context.Context, oneLiner string) (stdoutLines []string, err error) {
	return commandexecutorbash.RunOneLinerAndGetStdoutAsLines(ctx, oneLiner)
}

func (b *BashService) RunOneLinerAndGetStdoutAsString(ctx context.Context, oneLiner string) (stdout string, err error) {
	return commandexecutorbash.RunOneLinerAndGetStdoutAsString(ctx, oneLiner)
}
