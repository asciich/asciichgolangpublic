package storageutils

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"golang.org/x/sys/unix"
)

// Flush OS write cache to the storage.
func Sync(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Sync OS write cache to the storage started.")

	tStart := time.Now()

	unix.Sync()

	duration := time.Since(tStart)

	logging.LogInfoByCtxf(ctx, "Sync OS write cache to the storage finished. Took '%s'.", duration)

	return nil
}
