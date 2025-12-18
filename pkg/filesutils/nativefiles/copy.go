package nativefiles

import (
	"context"
	"io"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Copy(ctx context.Context, src string, dst string, options *filesoptions.CopyOptions) error {
	if src == "" {
		return tracederrors.TracedErrorEmptyString("src")
	}

	if dst == "" {
		return tracederrors.TracedErrorEmptyString("dst")
	}

	if options == nil {
		options = new(filesoptions.CopyOptions)
	}

	if options.UseSudo {
		logging.LogInfoByCtxf(ctx, "Copy '%s' to '%s' started using sudo started.", src, dst)
		_, err := commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "cp", src, dst},
		})
		if err != nil {
			return err
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Copy '%s' to '%s' started.", src, dst)

		sourceFile, err := os.Open(src)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to open source file '%s' for move.", src)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(dst)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to open dest file '%s' for move.", dst)
		}
		defer destFile.Close()

		// Stream data from source to destination
		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return err
		}

		// Ensure data is physically written to disk
		err = destFile.Sync()
		if err != nil {
			return err
		}
	}

	logging.LogChangedByCtxf(ctx, "Copied '%s' to '%s'.", src, dst)

	return nil
}
