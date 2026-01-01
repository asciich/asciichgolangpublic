package osutils

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

// Does an os.Exit():
// - os.Exit(changedExistCode) when the ctx indicates a change.
// - os.Exit(0) otherwise
func ExitWithChangedExitCode(ctx context.Context, changedExistCode int) {
	if contextutils.IsChanged(ctx) {
		logging.LogInfoByCtxf(ctx, "There was a change performed. Going to exit with exit code = %d", changedExistCode)
		os.Exit(changedExistCode)
	}

	logging.LogInfoByCtxf(ctx, "There was no change performed. Going to exit with exit code = 0")
	os.Exit(0)
}
