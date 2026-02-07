package commandexecutorheadscale

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateUser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, username string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if username == "" {
		return tracederrors.TracedErrorEmptyString("username")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create headscale user '%s' on '%s' started.", username, hostDescription)

	var created = true
	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"headscale", "users", "create", username},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "creating user: constraint failed: UNIQUE constraint failed: users.name") {
			created = false
		} else {
			return err
		}
	}

	if created {
		logging.LogChangedByCtxf(ctx, "Created headscale user '%s'.", username)
	} else {
		logging.LogInfoByCtxf(ctx, "Headscale user '%s' already exists. Skip creation.", username)
	}

	logging.LogInfoByCtxf(ctx, "Create headscale user '%s' on '%s' finished.", username, hostDescription)

	return nil
}
