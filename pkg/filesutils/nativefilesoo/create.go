package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) Create(ctx context.Context, options *filesoptions.CreateOptions) error {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Create(ctx, path)
}

func (f *File) Exists(ctx context.Context) (bool, error) {
	path, err := f.GetPath()
	if err != nil {
		return false, err
	}

	return nativefiles.Exists(ctx, path), nil
}

