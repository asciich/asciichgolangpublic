package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) Delete(ctx context.Context, options *filesoptions.DeleteOptions) error {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Delete(ctx, path, options)
}
