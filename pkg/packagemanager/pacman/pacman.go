package pacman

import (
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Pacman struct {
	commandExectuor commandexecutorinterfaces.CommandExecutor
}

func NewPacman(commandExecutor commandexecutorinterfaces.CommandExecutor) (*Pacman, error) {
	ret := new(Pacman)

	err := ret.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *Pacman) SetCommandExecutor(commandExectuor commandexecutorinterfaces.CommandExecutor) error {
	if commandExectuor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	p.commandExectuor = commandExectuor

	return nil
}

func (p *Pacman) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if p.commandExectuor == nil {
		return nil, tracederrors.TracedError("commandExectutor not set")
	}

	return p.commandExectuor, nil
}
