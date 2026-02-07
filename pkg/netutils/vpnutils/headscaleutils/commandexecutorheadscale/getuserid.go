package commandexecutorheadscale

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetUserId(ctx context.Context, commandExectuor commandexecutorinterfaces.CommandExecutor, userName string) (int, error) {
	if commandExectuor == nil {
		return 0, tracederrors.TracedErrorNil("commandExecutor")
	}

	if userName == "" {
		return 0, tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := commandExectuor.GetHostDescription()
	if err != nil {
		return 0, err
	}

	logging.LogInfoByCtxf(ctx, "Get headscale user Id of user '%s' on '%s' started.", userName, hostDescription)

	rawList, err := ListUsersRaw(ctx, commandExectuor)
	if err != nil {
		return 0, err
	}

	var userId int
	for _, u := range rawList {
		if u.Name == userName {
			userId = u.ID
			break
		}
	}

	if userId == 0 {
		return 0, tracederrors.TracedErrorf("%w: username='%s' on '%s'", headscalegeneric.ErrHeadscaleUserNotFound, userName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Get headscale user Id of user '%s' on '%s' finished. ID is '%d'", userName, hostDescription, userId)

	return userId, nil
}
