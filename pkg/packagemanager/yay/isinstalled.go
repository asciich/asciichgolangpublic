package yay

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsInstalled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	isInstalled, err := osutils.IsCommandAvailable(contextutils.WithSilent(ctx), commandExecutor, "yay")
	if err != nil {
		return false, err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "yay is installed on '%s'.", hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "yay is not installed on '%s'.", hostDescription)
	}

	return isInstalled, nil
}
