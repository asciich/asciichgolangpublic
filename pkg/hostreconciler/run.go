package hostreconciler

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/nativehost"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ReconcileOnce runs a single reconciliation of the current host against the
// desired state provided by the given source.
//
// The following steps are always performed:
//   - Reads the desired hostname and applies it.
//
// Only changes which differ from the desired state are applied, so it is safe
// to call this function repeatedly (idempotent).
//
// Requires root privileges if changes need to be applied.
func ReconcileOnce(ctx context.Context, source Source) error {
	if source == nil {
		return tracederrors.TracedErrorNil("source")
	}

	logging.LogInfoByCtxf(ctx, "Reconcile host started.")

	err := reconcileHostname(ctx, source)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Reconcile host finished.")

	return nil
}

// reconcileHostname reads the desired hostname from the source and applies it to the current host.
//
// The hostname is read before anything is changed. If the source fails, the host is left untouched.
func reconcileHostname(ctx context.Context, source Source) error {
	if source == nil {
		return tracederrors.TracedErrorNil("source")
	}

	hostname, err := source.GetHostName(ctx)
	if err != nil {
		return err
	}

	if hostname == "" {
		return tracederrors.TracedError("source returned an empty hostname")
	}

	logging.LogInfoByCtxf(ctx, "Reconcile hostname '%s' started.", hostname)

	host := nativehost.NewNativeHost()
	err = host.SetHostName(ctx, hostname, &hostsutilsoptions.SetHostnameOptions{})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Reconcile hostname '%s' finished.", hostname)

	return nil
}
