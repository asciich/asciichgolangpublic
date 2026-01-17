package commandexecutorlinuxuserutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/linuxuseroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Create(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *linuxuseroptions.CreateOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	userName, err := options.GetUserName()
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create linux user '%s' on '%s' started.", userName, hostDescription)

	exists, err := Exists(ctx, commandExecutor, userName)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Linux user '%s' already exists on '%s'. Skip create.", userName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Going to create user '%s' on '%s'.", userName, hostDescription)

		cmd := []string{"useradd"}

		if options.CreateHomeDirectory {
			cmd = append(cmd, "-m")
		}

		cmd = append(cmd, userName)

		if options.UseSudo {
			cmd = append([]string{"sudo"}, cmd...)
		}

		_, err := commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: cmd,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Created linux user '%s' on '%s'.", userName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Create linux user '%s' on '%s' finished.", userName, hostDescription)

	return nil
}
