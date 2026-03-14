package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) WriteBytes(ctx context.Context, toWrite []byte, options *filesoptions.WriteOptions) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.WriteBytes(ctx, path, toWrite, options)
}
