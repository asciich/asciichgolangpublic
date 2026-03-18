package nativefilesoo

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (f *File) CopyToFile(ctx context.Context, destFile filesinterfaces.File) error {
	if destFile == nil {
		return tracederrors.TracedErrorNil("destFile")
	}

	srcPath, srcHostDescription, err := f.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	destPath, destHostDescription, err := f.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Copy '%s' on '%s' to '%s' on '%s' started.", srcPath, srcHostDescription, destPath, destHostDescription)

	src, err := f.OpenAsReadCloser(ctx)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := destFile.OpenAsWriteCloser(ctx)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to copy file '%s' on '%s' to '%s' on '%s': %w", srcPath, srcHostDescription, destPath, destHostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Copy '%s' on '%s' to '%s' on '%s' finished.", srcPath, srcHostDescription, destPath, destHostDescription)

	return nil
}
