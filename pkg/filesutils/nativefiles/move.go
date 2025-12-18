package nativefiles

import (
	"context"
	"errors"
	"os"
	"syscall"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Copy the file and delete the source.
// Usefull when not on the same device an an os.Rename fails.
func moveUsingCopyAndDelete(ctx context.Context, src string, dst string) error {
	err := Copy(contextutils.WithSilent(ctx), src, dst, &filesoptions.CopyOptions{})
	if err != nil {
		return err
	}

	// Delete the original file
	err = os.Remove(src)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to remove file '%s' after copying it for move: %w", src, err)
	}

	return nil
}

// Move the file 'src' to 'dst'.
//
// If a simple os.Rename fails the file is moved by copy it first and then delete it.
// So this function works even when src and dest are not on the same filesystem.
func Move(ctx context.Context, src string, dst string, options *filesoptions.MoveOptions) error {
	if src == "" {
		return tracederrors.TracedErrorEmptyString("src")
	}

	if dst == "" {
		return tracederrors.TracedErrorEmptyString("dst")
	}

	if options == nil {
		options = new(filesoptions.MoveOptions)
	}

	if options.UseSudo {
		logging.LogInfoByCtxf(ctx, "Move '%s' to '%s' using sudo started.", src, dst)
		_, err := commandexecutorexec.RunCommand(ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"sudo", "mv", src, dst},
			},
		)
		if err != nil {
			return err
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Move '%s' to '%s' started.", src, dst)

		err := os.Rename(src, dst)
		if err != nil {
			// Check if the error is specifically "cross-device link".
			// This happens when when the src and dst are not on the same file system.
			//
			// os.Rename returns a *os.LinkError which contains the underlying syscall error
			var linkErr *os.LinkError
			if errors.As(err, &linkErr) {
				if lerr, ok := linkErr.Err.(syscall.Errno); ok && lerr == syscall.EXDEV {
					logging.LogInfoByCtxf(ctx, "Move '%s' to '%s' using copy and delete as the source and the destination are not on the same file system.", src, dst)
					err = moveUsingCopyAndDelete(ctx, src, dst)
					if err != nil {
						return err
					}
				}
			} else {
				// If it's a different error (e.g. permission denied), return it
				return tracederrors.TracedErrorf("Failed to move '%s' to '%s': %w", src, dst, err)
			}
		}
	}

	logging.LogChangedByCtxf(ctx, "Move '%s' to '%s' finished.", src, dst)

	return nil
}
