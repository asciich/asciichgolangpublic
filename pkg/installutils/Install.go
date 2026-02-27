package installutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func installFromSourcePath(ctx context.Context, options *InstallOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	sourcePath, err := options.GetSrcPath()
	if err != nil {
		return err
	}

	installPath, err := options.GetInstallPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' started.", sourcePath, installPath)

	err = nativefiles.Copy(ctx, sourcePath, installPath, &filesoptions.CopyOptions{
		UseSudo:         options.UseSudo,
		ReplaceExisting: options.ReplaceExisting,
	})
	if err != nil {
		return err
	}

	if options.IsModeSet() {
		permissions, err := options.GetMode()
		if err != nil {
			return err
		}
		err = nativefiles.Chmod(ctx, installPath, &filesoptions.ChmodOptions{
			PermissionsString: permissions,
		})
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' finished.", sourcePath, installPath)

	return nil
}

func Install(ctx context.Context, options *InstallOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.IsSourcePathSet() {
		return installFromSourcePath(ctx, options)
	}

	return tracederrors.TracedError("No source to install set.")
}
