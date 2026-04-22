package nativefiles

import (
	"context"
	"os"
)

func GetSizeBytes(ctx context.Context, path string) (int64, error) {
	err := ctx.Err()
	if err != nil {
		return 0, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}
