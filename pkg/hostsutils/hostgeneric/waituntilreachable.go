package hostgeneric

import (
	"context"
	"fmt"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/netutilserrors"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// WaitUntilReachable blocks until the given host reports as reachable or the
// timeout is exceeded. If renewHostKey is true, the SSH host key is renewed
// on every attempt (errors while renewing are ignored, since we're in a retry loop).
func WaitUntilReachable(ctx context.Context, host hostsutilsinterfaces.Host, options *hostsutilsoptions.WaitUntilReachableOptions) (err error) {
	if host == nil {
		return tracederrors.TracedErrorNil("host")
	}

	if options == nil {
		options = &hostsutilsoptions.WaitUntilReachableOptions{}
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
		if options.RenewHostKey {
			renewErr := host.RenewSshHostKey(ctx)
			if renewErr != nil {
				logging.LogWarnByCtxf(ctx,
					"Renewing host key for '%s' failed, but error is ignored in WaitUntilReachable since running in a retry loop.",
					hostname,
				)
			}
		}

		if options.AddHostKeyToKnownHosts {
			renewErr := host.AddSshHostKeyToKnownHosts(ctx)
			if renewErr != nil {
				logging.LogWarnByCtxf(ctx,
					"Adding host key for '%s' failed, but error is ignored in WaitUntilReachable since running in a retry loop.",
					hostname,
				)
			}
		}

		isReachable, err := host.IsReachable(ctx)
		if err != nil {
			if !netutilserrors.IsConnectionRefusedError(err) {
				logging.LogInfoByCtxf(ctx, "Host '%s' not reachable: connection refused.", hostname)
			} else if !netutilserrors.IsNoRouteToHostError(err) {
				logging.LogInfoByCtxf(ctx, "Host '%s' not reachable: No route to host.", hostname)
			} else {
				logging.LogInfoByCtxf(ctx, "Host not reachable: %v", err)
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
