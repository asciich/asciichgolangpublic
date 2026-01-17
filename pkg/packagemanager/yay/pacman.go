package yay

import (
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Yay struct {
	commandExectuor commandexecutorinterfaces.CommandExecutor
}

func NewYay(commandExecutor commandexecutorinterfaces.CommandExecutor) (*Yay, error) {
	ret := new(Yay)

	err := ret.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *Yay) SetCommandExecutor(commandExectuor commandexecutorinterfaces.CommandExecutor) error {
	if commandExectuor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	p.commandExectuor = commandExectuor

	return nil
}

func (p *Yay) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if p.commandExectuor == nil {
		return nil, tracederrors.TracedError("commandExectutor not set")
	}

	return p.commandExectuor, nil
}
