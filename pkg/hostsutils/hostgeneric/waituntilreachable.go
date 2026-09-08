package hostgeneric

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// WaitUntilReachable blocks until the given host reports as reachable or the
// timeout is exceeded. If renewHostKey is true, the SSH host key is renewed
// on every attempt (errors while renewing are ignored, since we're in a retry loop).
func WaitUntilReachable(ctx context.Context, host hostsutilsinterfaces.Host, renewHostKey bool) (err error) {
	if host == nil {
		return tracederrors.TracedErrorNil("host")
	}

	hostname, err := host.GetHostName()
	if err != nil {
		return err
	}

	const (
		timeout           = 60 * time.Second
		delayBetweenPings = 2 * time.Second
	)

	tStart := time.Now()

	for {
		if renewHostKey {
			renewErr := host.RenewSshHostKey(ctx)
			if renewErr != nil {
				logging.LogWarnByCtxf(ctx,
					"Renewing host key for '%s' failed, but error is ignored in WaitUntilReachable since running in a retry loop.",
					hostname,
				)
			}
		}

		isReachable, err := host.IsReachable(ctx)
		if err != nil {
			if !netutilserrors.IsConnectionRefusedError(err) {
				return err
			}
		}

		elapsedTime := time.Since(tStart)

		if isReachable {
			logging.LogGoodByCtxf(ctx, "Host '%s' is reachable after '%v'", hostname, elapsedTime)
			return nil
		}

		if elapsedTime > timeout {
			errorMessage := fmt.Sprintf("Host '%s' is not reachable after '%v'", hostname, elapsedTime)
			logging.LogErrorByCtx(ctx, errorMessage)
			return tracederrors.TracedError(errorMessage)
		}

		logging.LogInfoByCtxf(ctx,
			"Wait '%v' for host '%s' to get reachable. Total '%v' left, elapsed time so far: '%v'.",
			delayBetweenPings,
			hostname,
			timeout-elapsedTime,
			elapsedTime,
		)

		select {
		case <-ctx.Done():
			return tracederrors.TracedErrorf(
				"context cancelled while waiting for host '%s' to become reachable: %w",
				hostname, ctx.Err(),
			)
		case <-time.After(delayBetweenPings):
		}
	}
}
