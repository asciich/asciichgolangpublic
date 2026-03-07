package commandexecutorinstall

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpcommandexecutorclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func installFromSourceUrl(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *installoptions.InstallOptions) error {
	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	srcUrl, err := options.GetSrcUrl()
	if err != nil {
		return err
	}

	installPath, err := options.GetInstallPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' on '%s' started.", srcUrl, installPath, hostDescription)

	httpClient, err := httpcommandexecutorclientoo.NewClient(commandExecutor)
	if err != nil {
		return err
	}

	downloaded, err := httpClient.DownloadAsFile(ctx, &httpoptions.DownloadAsFileOptions{
		RequestOptions: &httpoptions.RequestOptions{
			Url: srcUrl,
		},
		OutputPath:        installPath,
		OverwriteExisting: true,
		Sha256Sum:         options.Sha256Sum,
	})
	if err != nil {
		return err
	}

	if options.Mode != "" {
		err := downloaded.Chmod(ctx, &filesoptions.ChmodOptions{
			PermissionsString: options.Mode,
		})
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' on '%s' finished.", srcUrl, installPath, hostDescription)

	return nil
}

func Install(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *installoptions.InstallOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.IsSourceUrlSet() {
		return installFromSourceUrl(ctx, commandExecutor, options)
	}

	return tracederrors.TracedErrorf("Not implemented for '%v'", options)
}
