package kuberneteshost

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// defaultApiServerPort is the default kube-apiserver port.
	defaultApiServerPort = 6443

	// defaultWaitUntilApiReachableTimeout is the maximum time to wait until the
	// API is reachable (InitialDelay not included).
	defaultWaitUntilApiReachableTimeout = 5 * time.Minute

	// defaultWaitUntilApiReachableInitialDelay gives the kubelet time to detect
	// a changed static-pod manifest (file check frequency is 20s by default)
	// and restart the pod before the first check.
	defaultWaitUntilApiReachableInitialDelay = 30 * time.Second

	// defaultWaitUntilApiReachablePollInterval is the time between two checks.
	defaultWaitUntilApiReachablePollInterval = 5 * time.Second

	// defaultWaitUntilApiReachableRequiredSuccesses is the number of
	// consecutive successful checks needed. A single success directly after a
	// VIP failover is not considered stable.
	defaultWaitUntilApiReachableRequiredSuccesses = 3

	// apiReadyzRequestTimeoutSeconds is the timeout of a single readyz request.
	apiReadyzRequestTimeoutSeconds = 5
)

// WaitUntilApiReachableOptions contains the options used to wait until the
// kube-apiserver is reachable through the given address (e.g. the KubeVip VIP).
type WaitUntilApiReachableOptions struct {
	// Port is the kube-apiserver port.
	// Default: 6443
	Port int

	// Timeout is the maximum time to wait until the API is reachable.
	// The InitialDelay is not included.
	// Default: 5 minutes
	Timeout time.Duration

	// InitialDelay is waited before the first check. This gives the kubelet
	// time to restart a static pod (e.g. KubeVip) after its manifest changed.
	// Otherwise the still running old instance answers and the check succeeds
	// before the restart even happened.
	// Default: 30 seconds
	InitialDelay time.Duration

	// SkipInitialDelay disables the InitialDelay (e.g. when nothing was
	// changed before the wait).
	// Default: false
	SkipInitialDelay bool

	// PollInterval is the time between two checks.
	// Default: 5 seconds
	PollInterval time.Duration

	// RequiredConsecutiveSuccesses is the number of consecutive successful
	// checks needed until the API is considered reachable.
	// Default: 3
	RequiredConsecutiveSuccesses int
}

// DefaultWaitUntilApiReachableOptions returns the default options.
func DefaultWaitUntilApiReachableOptions() *WaitUntilApiReachableOptions {
	return &WaitUntilApiReachableOptions{
		Port:                         defaultApiServerPort,
		Timeout:                      defaultWaitUntilApiReachableTimeout,
		InitialDelay:                 defaultWaitUntilApiReachableInitialDelay,
		SkipInitialDelay:             false,
		PollInterval:                 defaultWaitUntilApiReachablePollInterval,
		RequiredConsecutiveSuccesses: defaultWaitUntilApiReachableRequiredSuccesses,
	}
}

// prepareWaitUntilApiReachableOptions applies the defaults to a copy of the
// given options, so the options of the caller are never modified.
func prepareWaitUntilApiReachableOptions(options *WaitUntilApiReachableOptions) *WaitUntilApiReachableOptions {
	if options == nil {
		return DefaultWaitUntilApiReachableOptions()
	}

	ret := *options

	if ret.Port == 0 {
		ret.Port = defaultApiServerPort
	}

	if ret.Timeout == 0 {
		ret.Timeout = defaultWaitUntilApiReachableTimeout
	}

	if ret.InitialDelay == 0 {
		ret.InitialDelay = defaultWaitUntilApiReachableInitialDelay
	}

	if ret.PollInterval == 0 {
		ret.PollInterval = defaultWaitUntilApiReachablePollInterval
	}

	if ret.RequiredConsecutiveSuccesses == 0 {
		ret.RequiredConsecutiveSuccesses = defaultWaitUntilApiReachableRequiredSuccesses
	}

	return &ret
}

// validateWaitUntilApiReachableOptions validates the prepared options.
func validateWaitUntilApiReachableOptions(options *WaitUntilApiReachableOptions) error {
	if options.Port < 1 || options.Port > 65535 {
		return tracederrors.TracedErrorf("Invalid port '%d': must be between 1 and 65535.", options.Port)
	}

	if options.Timeout < 0 {
		return tracederrors.TracedErrorf("Invalid timeout '%s': must not be negative.", options.Timeout)
	}

	if options.InitialDelay < 0 {
		return tracederrors.TracedErrorf("Invalid initial delay '%s': must not be negative.", options.InitialDelay)
	}

	if options.PollInterval < 0 {
		return tracederrors.TracedErrorf("Invalid poll interval '%s': must not be negative.", options.PollInterval)
	}

	if options.RequiredConsecutiveSuccesses < 1 {
		return tracederrors.TracedErrorf("Invalid number of required consecutive successes '%d': must be at least 1.", options.RequiredConsecutiveSuccesses)
	}

	return nil
}

// sleepWithContext waits for the given duration or until the context is done.
func sleepWithContext(ctx context.Context, duration time.Duration) error {
	if duration <= 0 {
		return nil
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return tracederrors.TracedErrorf("Waiting was aborted: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

// isApiReadyz checks once if "https://<address>:<port>/readyz" answers "ok"
// from the target host. The output "yes"/"no" is evaluated instead of the
// exit code, so a failed execution is never mistaken for "not ready".
//
// The TLS certificate is not verified ("-k"): this is a pure reachability
// check, no data is sent or trusted. "/readyz" is readable without
// authentication by default (system:public-info-viewer).
func isApiReadyz(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, address string, port int) (bool, error) {
	// address is validated by iputils, so it is safe to embed into the script.
	url := "https://" + net.JoinHostPort(address, strconv.Itoa(port)) + "/readyz"

	script := `command -v curl >/dev/null 2>&1 || { echo missing-curl; exit 0; }
out="$(curl -sk --max-time ` + strconv.Itoa(apiReadyzRequestTimeoutSeconds) + ` '` + url + `' 2>/dev/null || true)"
if [ "$out" = "ok" ]; then echo yes; else echo no; fi`

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", script},
		},
	)
	if err != nil {
		return false, err
	}

	result := strings.TrimSpace(stdout)

	switch result {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	case "missing-curl":
		return false, tracederrors.TracedError("'curl' is not installed on the target: unable to check if the API is reachable.")
	default:
		return false, tracederrors.TracedErrorf("Unexpected output while checking '%s': '%s'", url, result)
	}
}

// WaitUntilApiReachable waits until the kube-apiserver answers "ok" on
// "https://<address>:<port>/readyz" when requested from the given host.
// Typically address is the KubeVip control-plane VIP.
//
// Use it after changing a KubeVip manifest before continuing with the next
// control plane, so the API VIP is never unavailable on all nodes at the
// same time.
//
// Read-only: does not require root privileges.
func WaitUntilApiReachable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, address string, options *WaitUntilApiReachableOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	// Handles the empty string as well (TracedErrorEmptyString).
	err := iputils.CheckValidIP(ctx, address)
	if err != nil {
		return err
	}

	options = prepareWaitUntilApiReachableOptions(options)

	err = validateWaitUntilApiReachableOptions(options)
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	endpoint := net.JoinHostPort(address, strconv.Itoa(options.Port))

	logging.LogInfoByCtxf(ctx, "Wait until API '%s' is reachable from '%s' started.", endpoint, hostDescription)

	if !options.SkipInitialDelay {
		logging.LogInfoByCtxf(ctx, "Wait initial delay of '%s' to let the kubelet restart changed static pods.", options.InitialDelay)

		err = sleepWithContext(ctx, options.InitialDelay)
		if err != nil {
			return err
		}
	}

	startTime := time.Now()
	deadline := startTime.Add(options.Timeout)
	consecutiveSuccesses := 0
	attempts := 0

	for {
		attempts++

		ready, err := isApiReadyz(ctx, commandExecutor, address, options.Port)
		if err != nil {
			return err
		}

		if ready {
			consecutiveSuccesses++
			if consecutiveSuccesses >= options.RequiredConsecutiveSuccesses {
				break
			}

			logging.LogInfoByCtxf(ctx, "API '%s' is ready ('%d/%d' consecutive successes).", endpoint, consecutiveSuccesses, options.RequiredConsecutiveSuccesses)
		} else {
			if consecutiveSuccesses > 0 {
				logging.LogInfoByCtxf(ctx, "API '%s' became unreachable again after '%d' successes. Restart counting.", endpoint, consecutiveSuccesses)
			} else {
				logging.LogInfoByCtxf(ctx, "API '%s' not reachable yet from '%s' (attempt '%d').", endpoint, hostDescription, attempts)
			}

			consecutiveSuccesses = 0
		}

		if time.Now().After(deadline) {
			return tracederrors.TracedErrorf(
				"API '%s' was not reachable from '%s' within '%s' ('%d' attempts, '%d/%d' consecutive successes at the end).",
				endpoint, hostDescription, options.Timeout, attempts, consecutiveSuccesses, options.RequiredConsecutiveSuccesses,
			)
		}

		err = sleepWithContext(ctx, options.PollInterval)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(
		ctx,
		"Wait until API '%s' is reachable from '%s' finished. Reachable after '%s' and '%d' attempts.",
		endpoint, hostDescription, time.Since(startTime).Round(time.Second), attempts,
	)

	return nil
}
