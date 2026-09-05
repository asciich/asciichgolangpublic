package aptget

import (
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type AptGet struct {
	commandExectuor commandexecutorinterfaces.CommandExecutor
}

func NewAptGet(commandExecutor commandexecutorinterfaces.CommandExecutor) (*AptGet, error) {
	ret := new(AptGet)

	err := ret.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (a *AptGet) SetCommandExecutor(commandExectuor commandexecutorinterfaces.CommandExecutor) error {
	if commandExectuor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	a.commandExectuor = commandExectuor

	return nil
}

func (a *AptGet) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if a.commandExectuor == nil {
		return nil, tracederrors.TracedError("commandExectutor not set")
	}

	return a.commandExectuor, nil
}
