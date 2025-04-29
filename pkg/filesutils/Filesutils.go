package filesutils

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/logging"
)

func IsFile(ctx context.Context, pathToCheck string) bool {
	stat, err := os.Stat(pathToCheck)
	if err != nil {
		logging.LogInfoByCtxf(ctx, "'%s' is not a file", pathToCheck)
		return false
	}

	if stat.IsDir() {
		logging.LogInfoByCtxf(ctx, "'%s' is a directory, not a file.", pathToCheck)
		return false
	}

	logging.LogInfoByCtxf(ctx, "'%s' is a file", pathToCheck)
	return true
}
