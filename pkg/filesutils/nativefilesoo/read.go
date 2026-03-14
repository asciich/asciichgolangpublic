package nativefilesoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func (f *File) ReadAsBytes(ctx context.Context) (content []byte, err error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	return nativefiles.ReadAsBytes(ctx, path)
}
