package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) Truncate(ctx context.Context, newSizeBytes int64) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Truncate(ctx, path, newSizeBytes)
}

func (f *File) GetSizeBytes(ctx context.Context) (fileSize int64, err error) {
	path, err := f.GetPath()
	if err != nil {
		return 0, err
	}

	return nativefiles.GetSizeBytes(ctx, path)
}

