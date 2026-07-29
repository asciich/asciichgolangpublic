package filesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
)

func NewFileByPath(path string) (filesinterfaces.File, error) {
	return nativefilesoo.NewFileByPath(path)
}

func Delete(ctx context.Context, path string, options *filesoptions.DeleteOptions) error {
	return nativefiles.Delete(ctx, path, options)
}

func WriteString(ctx context.Context, pathToWrite string, content string) error {
	return nativefiles.WriteString(ctx, pathToWrite, content)
}
