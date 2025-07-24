package commandexecutorpowershelloo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorpowershell"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

type PowerShellService struct {
	commandexecutorgeneric.CommandExecutorBase
}

func NewPowerShell() (p *PowerShellService) {
	return new(PowerShellService)
}

func NewPowerShellService() (p *PowerShellService) {
	return new(PowerShellService)
}

func PowerShell() (p *PowerShellService) {
	return NewPowerShell()
}

func (b *PowerShellService) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	return commandexecutorpowershell.RunCommand(ctx, options)
}

func (p *PowerShellService) RunOneLiner(ctx context.Context, oneLiner string) (output *commandoutput.CommandOutput, err error) {
	return commandexecutorpowershell.RunOneLiner(ctx, oneLiner)
}

func (p *PowerShellService) RunOneLinerAndGetStdoutAsString(ctx context.Context, oneLiner string) (stdout string, err error) {
	return commandexecutorpowershell.RunOneLinerAndGetStdoutAsString(ctx, oneLiner)
}
