package filesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func Exists(ctx context.Context, path string) (bool, error) {
	return nativefiles.Exists(ctx, path)
}

func CheckExists(ctx context.Context, path string) error {
	return nativefiles.CheckExists(ctx, path)
}
