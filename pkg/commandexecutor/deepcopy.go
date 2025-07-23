package commandexecutor

import (
	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetDeepCopyOfCommandExecutor(commandExectuor commandexecutorinterfaces.CommandExecutor) (copy commandexecutorinterfaces.CommandExecutor, err error) {
	if commandExectuor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	withDeepCopy, ok := commandExectuor.(interface {
		GetDeepCopy() commandexecutorinterfaces.CommandExecutor
	})
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
