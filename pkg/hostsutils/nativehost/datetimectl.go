package nativehost

import (
	"context"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// EnableTimeDateCtlNtp enables NTP time synchronization via
// "timedatectl set-ntp true" and waits until the system clock is synchronized.
// The operation is skipped when the clock is already synchronized (idempotent).
func (n *NativeHost) EnableTimeDateCtlNtp(ctx context.Context) (err error) {
	hostDescription, err := n.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable timedatectl NTP on '%s' started.", hostDescription)

	alreadySynced, err := n.isTimeDateCtlNtpSynchronized(ctx)
	if err != nil {
		return err
	}

	if alreadySynced {
		logging.LogInfoByCtxf(ctx, "Timedatectl NTP on '%s' is already synchronized. Skip enabling NTP.", hostDescription)
		logging.LogInfoByCtxf(ctx, "Enable timedatectl NTP on '%s' finished.", hostDescription)
		return nil
	}

	_, err = n.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "timedatectl", "set-ntp", "true"},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Enabled timedatectl NTP on '%s'.", hostDescription)

	// Wait until the clock is actually synchronized before returning.
	timeout := 60 * time.Second
	delay := 2 * time.Second
	tStart := time.Now()
	for {
		synced, err := n.isTimeDateCtlNtpSynchronized(ctx)
		if err != nil {
			return err
		}

		if synced {
			logging.LogInfoByCtxf(ctx, "Timedatectl NTP on '%s' is now synchronized.", hostDescription)
			break
		}

		if time.Since(tStart) > timeout {
			return tracederrors.TracedErrorf("Timedatectl NTP on '%s' did not synchronize within '%v'.", hostDescription, timeout)
		}

		logging.LogInfoByCtxf(ctx, "Wait '%v' for timedatectl NTP on '%s' to synchronize.", delay, hostDescription)
		time.Sleep(delay)
	}

	logging.LogInfoByCtxf(ctx, "Enable timedatectl NTP on '%s' finished.", hostDescription)

	return nil
}

// isTimeDateCtlNtpSynchronized returns true if the system clock is currently
// synchronized (timedatectl "NTPSynchronized" property is "yes").
func (n *NativeHost) isTimeDateCtlNtpSynchronized(ctx context.Context) (bool, error) {
	stdout, err := n.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"timedatectl", "show", "-p", "NTPSynchronized", "--value"},
		},
	)
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(stdout) == "yes", nil
}
