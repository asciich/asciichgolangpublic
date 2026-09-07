package commandexecutorinstall

import (
	"context"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils/commandexecutortarutils"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils/commandexecutorchecksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfile"
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

// isAlreadyInstalled returns true if the file at installPath already exists on
// the command executor's host and its sha256 sum matches the expected one.
// Returns false when no checksum is given so the caller always performs the
// installation in that case.
func isAlreadyInstalled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, installPath string, sha256Sum string) (bool, error) {
	if sha256Sum == "" {
		return false, nil
	}

	exists, err := commandexecutorfile.Exists(ctx, commandExecutor, installPath)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	actualSha256Sum, err := commandexecutorchecksumutils.GetSha256SumFromFile(ctx, commandExecutor, installPath)
	if err != nil {
		return false, err
	}

	return actualSha256Sum == sha256Sum, nil
}

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
		// The download is performed locally (see InstallOptions.ViaLocalTempDirectory),
		// so the temporary directory must be local as well.
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

func installFromSourceUrlArchive(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *installoptions.InstallOptions) error {
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

	archiveMember, err := options.GetSourceArchivePath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' (archive member '%s') as '%s' on '%s' started.", srcUrl, archiveMember, installPath, hostDescription)

	alreadyInstalled, err := isAlreadyInstalled(ctx, commandExecutor, installPath, options.Sha256Sum)
	if err != nil {
		return err
	}
	if alreadyInstalled {
		logging.LogInfoByCtxf(ctx, "Install '%s' (archive member '%s') as '%s' on '%s' finished. Already in place with matching checksum.", srcUrl, archiveMember, installPath, hostDescription)
		return nil
	}

	// The archive is downloaded, extracted and installed on the command
	// executor's host. This keeps the whole flow on the target host and works
	// for remote executors too.
	tempDirPath, err := commandexecutortempfile.CreateEmptyTemporaryDirectory(ctx, commandExecutor)
	if err != nil {
		return err
	}
	defer commandexecutorfile.Delete(ctx, commandExecutor, tempDirPath, &filesoptions.DeleteOptions{})

	downloadedArchivePath := filepath.Join(tempDirPath, "downloaded_archive")

	httpClient, err := httpcommandexecutorclientoo.NewClient(commandExecutor)
	if err != nil {
		return err
	}

	_, err = httpClient.DownloadAsFile(ctx, &httpoptions.DownloadAsFileOptions{
		RequestOptions: &httpoptions.RequestOptions{
			Url:               srcUrl,
			SkipTLSvalidation: options.SkipTLSvalidation,
		},
		OutputPath:        downloadedArchivePath,
		OverwriteExisting: true,
	})
	if err != nil {
		return err
	}

	// ExtractFileFromTarArchive auto-detects gzip by magic bytes, so both
	// .tar and .tar.gz archives are supported.
	extractedPath := filepath.Join(tempDirPath, "extracted")
	err = commandexecutortarutils.ExtractFileFromTarArchive(ctx, commandExecutor, downloadedArchivePath, archiveMember, extractedPath)
	if err != nil {
		return err
	}

	// The checksum refers to the extracted file, not the archive itself.
	if options.Sha256Sum != "" {
		actualSha256Sum, err := commandexecutorchecksumutils.GetSha256SumFromFile(ctx, commandExecutor, extractedPath)
		if err != nil {
			return err
		}

		if actualSha256Sum != options.Sha256Sum {
			return tracederrors.TracedErrorf(
				"Sha256Sum mismatch for extracted file '%s'. Expected '%s' but got '%s'.",
				extractedPath,
				options.Sha256Sum,
				actualSha256Sum,
			)
		}
	}

	logging.LogInfoByCtxf(ctx, "Copy extracted file '%s' to '%s' as '%s'.", extractedPath, hostDescription, installPath)
	err = commandexecutorfile.Copy(ctx, commandExecutor, extractedPath, installPath, &filesoptions.CopyOptions{
		UseSudo:         options.UseSudo,
		ReplaceExisting: options.ReplaceExisting,
	})
	if err != nil {
		return err
	}

	if options.Mode != "" {
		err = commandexecutorfile.Chmod(ctx, commandExecutor, installPath, &filesoptions.ChmodOptions{
			PermissionsString: options.Mode,
			UseSudo:           options.UseSudo,
		})
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' (archive member '%s') as '%s' on '%s' finished.", srcUrl, archiveMember, installPath, hostDescription)

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
		if options.SrcArchivePath != "" {
			return installFromSourceUrlArchive(ctx, commandExecutor, options)
		}

		return installFromSourceUrl(ctx, commandExecutor, options)
	}

	return tracederrors.TracedErrorf("Not implemented for '%v'", options)
}
