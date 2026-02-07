package commandexecutorheadscale

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func UserExists(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, userName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if userName == "" {
		return false, tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	users, err := ListUserNames(ctx, commandExecutor)
	if err != nil {
		return false, err
	}

	exits := slices.Contains(users, userName)

	if exits {
		logging.LogInfoByCtxf(ctx, "Headscale user '%s' on '%s' exists.", userName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Headscale user '%s' on '%s' does not exist.", userName, hostDescription)
	}

	return exits, nil
}
