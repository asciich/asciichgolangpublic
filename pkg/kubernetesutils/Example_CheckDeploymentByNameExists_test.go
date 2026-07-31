package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CheckDeploymentByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Create an example namespace and deployment.
	const namespaceName = "testnamespace"
	const deploymentName = "example-deployment"
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	_, err = namespace.CreateDeployment(ctx, &kubernetesparameteroptions.RunCommandOptions{
		Image:    "busybox",
		Command:  []string{"sleep", "3600"},
		DeploymentName:     deploymentName,
		Replicas: 1,
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created deployment exists - CheckDeploymentByNameExists returns nil:
	err = cluster.CheckDeploymentByNameExists(ctx, namespaceName, deploymentName)
	require.NoError(t, err)

	// The same deployment name in the default namespace does not exist - CheckDeploymentByNameExists returns error:
	err = cluster.CheckDeploymentByNameExists(ctx, "default", deploymentName)
	require.Error(t, err)

	// This deployment is expected to be in the same namespace but does not exist:
	err = cluster.CheckDeploymentByNameExists(ctx, namespaceName, "deployment-does-not-exist")
	require.Error(t, err)

	// If we delete our deployment again...
	err = namespace.DeleteDeploymentByName(ctx, deploymentName)
	require.NoError(t, err)

	// ... our deployment becomes absent and CheckDeploymentByNameExists returns error:
	err = cluster.CheckDeploymentByNameExists(ctx, namespaceName, deploymentName)
	require.Error(t, err)
}
