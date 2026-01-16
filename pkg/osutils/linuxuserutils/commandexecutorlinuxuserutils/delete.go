package commandexecutorlinuxuserutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/linuxuseroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Delete(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *linuxuseroptions.DeleteOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	userName, err := options.GetUSerName()
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil
	}

	logging.LogInfoByCtxf(ctx, "Delete linux user '%s' on '%s' started.", userName, hostDescription)

	exists, err := Exists(ctx, commandExecutor, userName)
	if err != nil {
		return nil
	}

	if exists {
		runOptions := &parameteroptions.RunCommandOptions{
			Command:   []string{"userdel", userName},
			RunAsRoot: options.UseSudo,
		}

		cmdJoined, err := runOptions.GetJoinedFullCommand()
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Going to delete linux user '%s' on '%s' using command '%s'.", userName, hostDescription, cmdJoined)

		_, err = commandExecutor.RunCommand(ctx, runOptions)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Linux user '%s' created on '%s'.", userName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Linux user '%s' already absent on '%s', skip delete.", userName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Delete linux user '%s' on '%s' finished.", userName, hostDescription)

	return nil
}
