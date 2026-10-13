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
// The whole desired state is read from the source and validated before anything is changed.
// If the source fails or returns an invalid state, the host is left untouched.
//
// The following steps are performed:
//   - Applies the desired hostname.
//   - Reconciles all desired network interfaces:
//   - Deactivates offloading if requested (see networkcardutils.DeactivateOffloading).
//   - Ensures the layer 2 bridge exists if requested (see networkcardutils.CreateLayer2Bridge).
//
// Reconciliation of network interfaces is additive: Removed bridges or offloading settings are NOT reverted.
//
// Only changes which differ from the desired state are applied, so it is safe
// to call this function repeatedly (idempotent).
//
// Requires root privileges if changes need to be applied or at least one network interface action is configured.
func ReconcileOnce(ctx context.Context, source Source) error {
	if source == nil {
		return tracederrors.TracedErrorNil("source")
	}

	logging.LogInfoByCtxf(ctx, "Reconcile host started.")

	// Read and validate the complete desired state before anything is changed:
	hostname, err := source.GetHostName(ctx)
	if err != nil {
		return err
	}

	if hostname == "" {
		return tracederrors.TracedError("source returned an empty hostname")
	}

	networkInterfaces, err := source.GetNetworkInterfaces(ctx)
	if err != nil {
		return err
	}

	err = validateNetworkInterfaces(networkInterfaces)
	if err != nil {
		return err
	}

	// Apply the desired state:
	err = reconcileHostname(ctx, hostname)
	if err != nil {
		return err
	}

	err = reconcileNetworkInterfaces(ctx, networkInterfaces)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Reconcile host finished.")

	return nil
}

// reconcileHostname applies the desired 'hostname' to the current host.
func reconcileHostname(ctx context.Context, hostname string) error {
	if hostname == "" {
		return tracederrors.TracedErrorEmptyString("hostname")
	}

	logging.LogInfoByCtxf(ctx, "Reconcile hostname '%s' started.", hostname)

	host := nativehost.NewNativeHost()
	err := host.SetHostName(ctx, hostname, &hostsutilsoptions.SetHostnameOptions{})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Reconcile hostname '%s' finished.", hostname)

	return nil
}
