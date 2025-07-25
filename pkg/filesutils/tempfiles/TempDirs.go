package tempfiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateTempDir(ctx context.Context) (string, error) {
	dirPath, err := os.MkdirTemp("", "empty")
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to create temporary directory: %w", err)
	}

	logging.LogChangedByCtxf(ctx, "Created temporary directory '%s'", dirPath)

	return dirPath, nil
}
