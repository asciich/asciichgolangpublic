package hostreconciler

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

// reconcileNetworkInterfaces brings every network interface in 'interfaces' on localhost to its desired state:
//   - 'DeactivateOffloading': see networkcardutils.DeactivateOffloading.
//   - 'Layer2Bridge': see networkcardutils.CreateLayer2Bridge.
//
// All changes are instantly applied and persistent across reboots.
//
// Reconciliation is additive: Bridges or offloading settings no longer configured are NOT reverted.
//
// All network interfaces are validated before any change is performed.
//
// Requires root privileges if at least one action is configured.
// This function is idempotent.
func reconcileNetworkInterfaces(ctx context.Context, interfaces []NetworkInterface) error {
	err := validateNetworkInterfaces(interfaces)
	if err != nil {
		return err
	}

	toReconcile := []NetworkInterface{}
	for _, iface := range interfaces {
		if iface.hasActions() {
			toReconcile = append(toReconcile, iface)
		} else {
			logging.LogInfoByCtxf(ctx, "Network interface '%s' has no actions configured. Nothing to reconcile.", iface.Name)
		}
	}

	if len(toReconcile) == 0 {
		logging.LogInfoByCtxf(ctx, "No network interface actions configured. Nothing to reconcile for network interfaces.")
		return nil
	}

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	if err != nil {
		return err
	}

	if !isRoot {
		return tracederrors.TracedErrorf("Reconcile '%d' network interfaces requires root privileges.", len(toReconcile))
	}

	logging.LogInfoByCtxf(ctx, "Reconcile '%d' network interfaces started.", len(toReconcile))

	for _, iface := range toReconcile {
		err = reconcileNetworkInterface(ctx, iface)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Reconcile '%d' network interfaces finished.", len(toReconcile))

	return nil
}

// reconcileNetworkInterface applies the configured actions of a single, already validated network interface.
func reconcileNetworkInterface(ctx context.Context, iface NetworkInterface) error {
	if iface.Name == "" {
		return tracederrors.TracedErrorEmptyString("iface.Name")
	}

	logging.LogInfoByCtxf(ctx, "Reconcile network interface '%s' started.", iface.Name)

	// Offloading is deactivated first so the network card is stable before it is moved into the bridge.
	if iface.DeactivateOffloading {
		err := networkcardutils.DeactivateOffloading(ctx, iface.Name)
		if err != nil {
			return tracederrors.TracedErrorf("Deactivate offloading for network interface '%s' failed: %w", iface.Name, err)
		}
	}

	if iface.Layer2Bridge != "" {
		err := networkcardutils.CreateLayer2Bridge(ctx, iface.Layer2Bridge, iface.Name)
		if err != nil {
			return tracederrors.TracedErrorf("Reconcile layer 2 bridge '%s' with network interface '%s' failed: %w", iface.Layer2Bridge, iface.Name, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Reconcile network interface '%s' finished.", iface.Name)

	return nil
}
