package commandexecutorlinuxuserutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Exists(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, userName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if userName == "" {
		return false, tracederrors.TracedErrorEmptyString("userName")
	}

	output, err := commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command:           []string{"id", userName},
			AllowAllExitCodes: true,
		},
	)
	if err != nil {
		return false, err
	}

	exists := output.IsExitSuccess()

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Linux user '%s' does exist on '%s'.", userName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Linux user '%s' does not exist on '%s'.", userName, hostDescription)
	}

	return exists, err
}
