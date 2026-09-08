package commandexecutorkubernetes

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (c *CommandExecutorKubernetes) DeleteNamespaceByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	exists, err := c.NamespaceByNameExists(ctx, name)
	if err != nil {
		return err
	}

	clusterName, err := c.GetName()
	if err != nil {
		return err
	}

	if exists {

		context, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return err
		}

		// Delete resources that own finalizers processed by controllers living
		// inside this namespace, before the namespace (and its controllers) are gone.
		err = c.deleteResourcesBeforeNamespaceDeletion(ctx, name)
		if err != nil {
			return err
		}

		_, err = c.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"kubectl",
					"--context",
					context,
					"delete",
					"namespace",
					name,
				},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Namespace '%s' in cluster '%s' deleted.", name, clusterName)
	} else {
		logging.LogInfoByCtxf(ctx, "Namespace '%s' already absent in cluster '%s'.", name, clusterName)
	}

	return nil
}

// deleteResourcesBeforeNamespaceDeletion deletes resources that must be removed
// while the namespace (and its controllers) still exist. This avoids the
// chicken-and-egg deadlock where a namespace gets stuck in 'Terminating' because
// a custom resource inside it still carries a finalizer that only its controller
// can remove — but that controller is torn down together with the namespace.
func (c *CommandExecutorKubernetes) deleteResourcesBeforeNamespaceDeletion(ctx context.Context, namespaceName string) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if namespaceName == "flux-system" {
		// Delete the FluxInstance 'flux' first so the still-running flux-operator
		// can remove its finalizer 'fluxcd.controlplane.io/finalizer'.
		err := c.deleteFluxInstance(ctx, namespaceName, "flux")
		if err != nil {
			return err
		}
	}

	return nil
}

// deleteFluxInstance deletes the given FluxInstance custom resource
// (fluxinstances.fluxcd.controlplane.io) and waits until it is fully gone.
//
// This must happen while the flux-operator is still running, because the
// FluxInstance carries a finalizer that only the operator can remove. 'kubectl
// delete' blocks until the object is actually removed (i.e. until the operator
// has processed the finalizer), so the wait is implicit but we add an explicit
// timeout to avoid hanging forever.
//
// The operation is a no-op if the resource, its namespace, or the FluxInstance
// CRD do not exist (e.g. on clusters where flux-operator was never installed).
func (c *CommandExecutorKubernetes) deleteFluxInstance(ctx context.Context, namespaceName string, fluxInstanceName string) error {
	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if fluxInstanceName == "" {
		return tracederrors.TracedErrorEmptyString("fluxInstanceName")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	cmd := []string{"kubectl"}

	if kubernetesutils.IsInClusterAuthenticationAvailable(ctx) {
		logging.LogInfoByCtxf(ctx, "Kubernetes in cluster authentication is used. cluster context is not used.")
	} else {
		kubeContext, err := c.GetCachedKubectlContext(ctx)
		if err != nil {
			return err
		}

		cmd = append(cmd, "--context", kubeContext)
	}

	cmd = append(cmd,
		"delete", "fluxinstance", fluxInstanceName,
		"--namespace", namespaceName,
		"--ignore-not-found=true",
		"--timeout=120s",
	)

	logging.LogInfoByCtxf(ctx, "Delete FluxInstance '%s' in namespace '%s' started.", fluxInstanceName, namespaceName)

	output, err := commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command:           cmd,
			AllowAllExitCodes: true,
			TimeoutString:     "130 seconds",
		},
	)
	if err != nil {
		return err
	}

	if output.IsExitSuccess() {
		logging.LogChangedByCtxf(ctx, "FluxInstance '%s' in namespace '%s' deleted.", fluxInstanceName, namespaceName)
		return nil
	}

	stderr, err := output.GetStderrAsString()
	if err != nil {
		return err
	}

	// Gracefully ignore the cases where there is simply nothing to delete:
	//   - the FluxInstance CRD is not installed on this cluster
	//   - the namespace does not exist
	// '--ignore-not-found=true' already covers the missing FluxInstance object itself.
	if isFluxInstanceHarmlessDeleteError(stderr) {
		logging.LogInfoByCtxf(ctx, "FluxInstance '%s' in namespace '%s' not present (or CRD not installed). Skip delete.", fluxInstanceName, namespaceName)
		return nil
	}

	return tracederrors.TracedErrorf("Failed to delete FluxInstance '%s' in namespace '%s': %s", fluxInstanceName, namespaceName, stderr)
}

// isFluxInstanceHarmlessDeleteError returns true if the given kubectl stderr
// indicates there was nothing meaningful to delete (missing CRD or namespace),
// which should not be treated as a failure.
func isFluxInstanceHarmlessDeleteError(stderr string) bool {
	harmlessMarkers := []string{
		// CRD not installed:
		"the server doesn't have a resource type",
		"doesn't have a resource type \"fluxinstance",
		"could not find the requested resource",
		"the server could not find the requested resource",
		"no matches for kind",
		"unable to recognize",
		// Namespace missing:
		"namespaces \"",
	}

	lowerStderr := strings.ToLower(stderr)
	for _, marker := range harmlessMarkers {
		if strings.Contains(lowerStderr, strings.ToLower(marker)) {
			return true
		}
	}

	return false
}
