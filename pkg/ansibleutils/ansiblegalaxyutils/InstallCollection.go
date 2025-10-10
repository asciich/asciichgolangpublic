package ansiblegalaxyutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func InstallGalaxyCollection(ctx context.Context, options *InstallCollectionOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.LocalCollectionPath == "" {
		return tracederrors.TracedError("Only implemented for LocalCollectionPath at the moment")
	}

	ansiblePath, err := options.GetAnsibleGalaxyPath()
	if err != nil {
		return err
	}

	path := options.LocalCollectionPath

	logging.LogInfoByCtxf(ctx, "Install ansible collection '%s' using ansible-galaxy '%s' started.", path, ansiblePath)

	_, err = commandexecutorexec.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{ansiblePath, "collection", "install", options.LocalCollectionPath},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Install ansible collection '%s' using ansible-galaxy '%s' finished.", path, ansiblePath)

	return nil
}
