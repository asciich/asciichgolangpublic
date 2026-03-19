package commandexecutorinstall

import (
	"context"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpcommandexecutorclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpnativeclientoo"
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

	var installedFile filesinterfaces.File
	if options.ViaLocalTempDirectory {
		tempDirPath, err := tempfiles.CreateTempDir(ctx)
		if err != nil {
			return err
		}
		defer nativefiles.Delete(ctx, tempDirPath, &filesoptions.DeleteOptions{})

		downloadedFilePath := filepath.Join(tempDirPath, "download")

		logging.LogInfoByCtxf(ctx, "Perform installation using temporary file '%s'.", downloadedFilePath)

		httpClient := httpnativeclientoo.NewNativeClient()

		downloadeFile, err := httpClient.DownloadAsFile(ctx, &httpoptions.DownloadAsFileOptions{
			RequestOptions: &httpoptions.RequestOptions{
				Url:               srcUrl,
				SkipTLSvalidation: options.SkipTLSvalidation,
			},
			OutputPath:        downloadedFilePath,
			OverwriteExisting: true,
			Sha256Sum:         options.Sha256Sum,
			UseSudo:           false,
		})

		installedFile, err = commandexecutorfileoo.New(commandExecutor, installPath)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Copy locally downloaded temporary file '%s' to '%s' as '%s'.", downloadedFilePath, hostDescription, installPath)
		err = downloadeFile.CopyToFile(
			ctx, 
			installedFile,
			&filesoptions.CopyOptions{
				UseSudo: options.UseSudo,
			},
		)
		if err != nil {
			return err
		}
	} else {
		httpClient, err := httpcommandexecutorclientoo.NewClient(commandExecutor)
		if err != nil {
			return err
		}

		installedFile, err = httpClient.DownloadAsFile(ctx, &httpoptions.DownloadAsFileOptions{
			RequestOptions: &httpoptions.RequestOptions{
				Url:               srcUrl,
				SkipTLSvalidation: options.SkipTLSvalidation,
			},
			OutputPath:        installPath,
			OverwriteExisting: true,
			Sha256Sum:         options.Sha256Sum,
			UseSudo:           options.UseSudo,
		})
		if err != nil {
			return err
		}
	}

	if options.Mode != "" {
		err := installedFile.Chmod(ctx, &filesoptions.ChmodOptions{
			PermissionsString: options.Mode,
			UseSudo:           options.UseSudo,
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
