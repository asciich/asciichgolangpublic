package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) Chmod(ctx context.Context, options *filesoptions.ChmodOptions) error {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Chmod(ctx, path, options)
}

func (f *File) GetAccessPermissionsString() (string, error) {
	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	return nativefiles.GetAccessPermissionsString(path)
}

func (f *File) GetAccessPermissions() (int, error) {
	path, err := f.GetPath()
	if err != nil {
		return 0, err
	}

	return nativefiles.GetAccessPermissions(path)
}
