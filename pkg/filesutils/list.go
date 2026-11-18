package filesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func ListSubDirectories(ctx context.Context, path string, options *parameteroptions.ListDirectoryOptions) ([]string, error) {
	return nativefiles.ListSubDirectories(ctx, path, options)
}
