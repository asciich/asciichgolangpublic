package nativefilesoo

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) ReadAsBytes(ctx context.Context) (content []byte, err error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	return nativefiles.ReadAsBytes(ctx, path)
}

func (f *File) OpenAsReadCloser(ctx context.Context) (io.ReadCloser, error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	return nativefiles.OpenAsReadCloser(ctx, path)
}
