package commandexecutorheadscale

import (
	"context"
	"strconv"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func DeleteUser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, userName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if userName == "" {
		return tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete headscale user '%s' on '%s' started.", userName, hostDescription)

	var performDelete = true
	userId, err := GetUserId(ctx, commandExecutor, userName)
	if err != nil {
		if headscalegeneric.IsErrHeadscaleUserNotFound(err) {
			performDelete = false
		} else {
			return err
		}
	}

	if performDelete {
		_, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"headscale", "users", "destroy", "--force", "--identifier", strconv.Itoa(userId)},
		})
		if err != nil {
			return err
		}
	}

	if performDelete {
		logging.LogChangedByCtxf(ctx, "Deleted headscale user '%s' on '%s'.", userName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Headscale user '%s' on '%s' is already absent. Skip delete.", userName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Delete headscale user '%s' on '%s' finished.", userName, hostDescription)

	return nil
}
