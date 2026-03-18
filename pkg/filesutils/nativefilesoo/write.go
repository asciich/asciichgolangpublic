package nativefilesoo

import (
	"context"
	"io"

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

func (f *File) OpenAsWriteCloser(ctx context.Context) (io.WriteCloser, error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	return nativefiles.OpenAsWriteCloser(ctx, path)
}
