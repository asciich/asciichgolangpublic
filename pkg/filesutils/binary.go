package filesutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func IsStaticallyLinkedBinary(ctx context.Context, path string) (bool, error) {
	return nativefiles.IsStaticallyLinkedBinary(ctx, path)
}
