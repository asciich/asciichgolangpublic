package filesutils

import (
	"context"
	"os"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func CreateTempDir(ctx context.Context) (string, error) {
	dirPath, err := os.MkdirTemp("", "empty")
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to create temporary directory: %w", err)
	}

	logging.LogChangedByCtxf(ctx, "Created temporary directory '%s'", dirPath)

	return dirPath, nil
}
