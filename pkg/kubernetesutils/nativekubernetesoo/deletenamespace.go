package nativekubernetesoo

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteFluxInstance deletes the given FluxInstance custom resource and waits
// until it is fully removed.
//
// This must be done while the flux-operator is still running: the FluxInstance
// carries the finalizer 'fluxcd.controlplane.io/finalizer', which is only removed
// by the operator itself. If the namespace (and with it the operator) were deleted
// first, the FluxInstance would keep its finalizer forever and the namespace would
// stay stuck in 'Terminating' (the chicken-and-egg deadlock).
func (n *NativeKubernetesCluster) deleteFluxInstance(ctx context.Context, namespaceName string, fluxInstanceName string) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if fluxInstanceName == "" {
		return tracederrors.TracedErrorEmptyString("fluxInstanceName")
	}

	dynamicClient, err := n.GetDynamicClient()
	if err != nil {
		return err
	}

	resourceClient := dynamicClient.Resource(fluxInstanceGVR).Namespace(namespaceName)

	_, err = resourceClient.Get(ctx, fluxInstanceName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "FluxInstance '%s' in namespace '%s' already absent. Skip delete.", fluxInstanceName, namespaceName)
			return nil
		}

		// The FluxInstance CRD is not installed on this cluster (e.g. flux-operator
		// was never deployed). There is nothing to delete in that case.
		if meta.IsNoMatchError(err) {
			logging.LogInfoByCtxf(ctx, "FluxInstance CRD not registered on cluster. Skip FluxInstance delete in namespace '%s'.", namespaceName)
			return nil
		}

		return tracederrors.TracedErrorf("Failed to get FluxInstance '%s' in namespace '%s': %w", fluxInstanceName, namespaceName, err)
	}

	// Foreground propagation ensures dependents are removed before the FluxInstance
	// object itself disappears, giving the operator time to run its teardown logic.
	deletePolicy := metav1.DeletePropagationForeground
	err = resourceClient.Delete(ctx, fluxInstanceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		if errors.IsNotFound(err) {
			logging.LogInfoByCtxf(ctx, "FluxInstance '%s' in namespace '%s' already absent. Skip delete.", fluxInstanceName, namespaceName)
			return nil
		}

		return tracederrors.TracedErrorf("Failed to delete FluxInstance '%s' in namespace '%s': %w", fluxInstanceName, namespaceName, err)
	}

	logging.LogChangedByCtxf(ctx, "FluxInstance '%s' in namespace '%s' deletion triggered.", fluxInstanceName, namespaceName)

	err = n.WaitUntilFluxInstanceDeleted(ctx, namespaceName, fluxInstanceName)
	if err != nil {
		return err
	}

	return nil
}

// WaitUntilFluxInstanceDeleted blocks until the given FluxInstance custom resource
// no longer exists or the timeout is reached. This guarantees the operator has
// finished processing the deletion (and removed its finalizer) before the
// namespace itself is deleted.
func (n *NativeKubernetesCluster) WaitUntilFluxInstanceDeleted(ctx context.Context, namespaceName string, fluxInstanceName string) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if fluxInstanceName == "" {
		return tracederrors.TracedErrorEmptyString("fluxInstanceName")
	}

	dynamicClient, err := n.GetDynamicClient()
	if err != nil {
		return err
	}

	resourceClient := dynamicClient.Resource(fluxInstanceGVR).Namespace(namespaceName)

	timeout := time.Second * 120

	logging.LogInfoByCtxf(ctx, "Wait for FluxInstance '%s' in namespace '%s' to be deleted started (timeout = %s).", fluxInstanceName, namespaceName, timeout)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tStart := time.Now()
	for {
		if ctx.Err() != nil {
			return tracederrors.TracedErrorf("Wait until FluxInstance '%s' in namespace '%s' deleted failed: %w", fluxInstanceName, namespaceName, ctx.Err())
		}

		_, err := resourceClient.Get(ctx, fluxInstanceName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				break
			}

			return tracederrors.TracedErrorf("Failed to get FluxInstance '%s' in namespace '%s' while waiting for deletion: %w", fluxInstanceName, namespaceName, err)
		}

		waitTime := time.Second * 1
		elapsed := time.Since(tStart)
		logging.LogInfoByCtxf(ctx, "Wait another %s until the FluxInstance '%s' in namespace '%s' is deleted (%s/%s).", waitTime, fluxInstanceName, namespaceName, elapsed, timeout)
		time.Sleep(waitTime)
	}

	logging.LogInfoByCtxf(ctx, "Wait for FluxInstance '%s' in namespace '%s' to be deleted finished.", fluxInstanceName, namespaceName)

	return nil
}

func (n *NativeKubernetesCluster) deleteResourcesBeforeNamespaceDeletion(ctx context.Context, namespaceName string) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if namespaceName == "flux-system" {
		// Delete the FluxInstance 'flux' before the namespace itself. This lets the
		// still-running flux-operator remove its finalizer, avoiding the deadlock
		// where the namespace gets stuck in 'Terminating' because the controller
		// that owns the finalizer was already torn down.
		err := n.deleteFluxInstance(ctx, namespaceName, "flux")
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *NativeKubernetesCluster) DeleteNamespaceByName(ctx context.Context, namespaceName string) (err error) {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	exists, err := n.NamespaceByNameExists(ctx, namespaceName)
	if err != nil {
		return err
	}

	if exists {
		clientset, err := n.GetClientSet()
		if err != nil {
			return err
		}

		err = n.deleteResourcesBeforeNamespaceDeletion(ctx, namespaceName)
		if err != nil {
			return err
		}

		deletePolicy := metav1.DeletePropagationForeground // This ensures child objects are deleted before the namespace
		deleteOptions := metav1.DeleteOptions{
			PropagationPolicy:  &deletePolicy,
			GracePeriodSeconds: nil, // Use default graceful termination period
		}

		err = clientset.CoreV1().Namespaces().Delete(ctx, namespaceName, deleteOptions)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete kubernetes namespace '%s': %w", namespaceName, err)
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' deleted.", namespaceName)

		err = n.WaitUntilNamespaceDeleted(ctx, namespaceName)
		if err != nil {
			return err
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already absent. Skip delete.", namespaceName)
	}

	return nil
}
