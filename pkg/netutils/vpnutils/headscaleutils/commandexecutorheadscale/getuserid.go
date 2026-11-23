package commandexecutorheadscale

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetUserId returns the ID of the headscale user with the given username.
//
// Returns headscalegeneric.ErrHeadscaleUserNotFound (wrapped) if there is no such user.
func GetUserId(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, userName string) (int, error) {
	if commandExecutor == nil {
		return 0, tracederrors.TracedErrorNil("commandExecutor")
	}

	if userName == "" {
		return 0, tracederrors.TracedErrorEmptyString("userName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return 0, err
	}

	logging.LogInfoByCtxf(ctx, "Get headscale user Id of user '%s' on '%s' started.", userName, hostDescription)

	users, err := ListUsers(ctx, commandExecutor)
	if err != nil {
		return 0, err
	}

	userId := 0
	for _, user := range users {
		username, err := user.GetUsername()
		if err != nil {
			return 0, err
		}

		if username != userName {
			continue
		}

		id, err := user.GetID()
		if err != nil {
			return 0, err
		}

		userId = int(id)
		break
	}

	if userId == 0 {
		return 0, tracederrors.TracedErrorf("%w: username='%s' on '%s'", headscalegeneric.ErrHeadscaleUserNotFound, userName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Get headscale user Id of user '%s' on '%s' finished. ID is '%d'", userName, hostDescription, userId)

	return userId, nil
}
