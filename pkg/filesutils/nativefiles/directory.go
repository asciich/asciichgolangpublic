package nativefiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsEmptyDirectory(ctx context.Context, path string) (bool, error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, tracederrors.TracedErrorf("failed to read directory '%s': %w", path, err)
	}

	isEmpty := len(entries) == 0

	if isEmpty {
		logging.LogInfoByCtxf(ctx, "The directory '%s' is empty.", path)
	} else {
		logging.LogInfoByCtxf(ctx, "The directory '%s' is not empty.", path)
	}

	return isEmpty, nil
}
