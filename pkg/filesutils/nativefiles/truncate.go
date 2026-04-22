package nativefiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Truncate(ctx context.Context, path string, newSizeBytes int64) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString(path)
	}

	if newSizeBytes < 0 {
		return tracederrors.TracedErrorf("Invalid new file size to truncate '%s': %d", path, newSizeBytes)
	}

	err := os.Truncate(path, newSizeBytes)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to truncate '%s' to %d bytes: %w", path, newSizeBytes, err)
	}

	logging.LogChangedByCtxf(ctx, "Truncated '%s' to %d bytes.", path, newSizeBytes)

	return nil
}
