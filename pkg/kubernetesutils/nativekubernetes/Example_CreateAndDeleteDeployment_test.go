package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

// This example shows how to create and delete a Deployment.
func Test_Example_CreateAndDeleteDeployment(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// Create a unique namespace for this test to avoid conflicts
	namespaceName := "test-example-createanddeletedeployment"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)
	defer func() {
		// Clean up namespace after test
		_ = nativekubernetes.DeleteNamespace(ctx, clientset, namespaceName)
	}()

	// ... prepare test environment finished.
	// -----

	// Define the name of the Deployment we use for testing
	const deploymentName = "example-deployment"

	// Ensure the Deployment is absent:
	err = nativekubernetes.DeleteDeployment(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)

	exists, err := nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Create the Deployment with 2 replicas:
	err = nativekubernetes.CreateDeployment(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			DeploymentName: deploymentName,
			Image:          "ubuntu",
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "sleep 1m"},
			},
			Replicas: 2,
		},
	)
	require.NoError(t, err)

	exists, err = nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// The create function is idempotent if DeleteAlreadyExistingDeployment is set:
	err = nativekubernetes.CreateDeployment(
		ctx,
		clientset,
		namespaceName,
		&kubernetesparameteroptions.KubernetesRunCommandOptions{
			DeploymentName: deploymentName,
			Image:          "ubuntu",
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: []string{"bash", "-c", "sleep 1m"},
			},
			Replicas:                        2,
			DeleteAlreadyExistingDeployment: true,
		},
	)
	require.NoError(t, err)

	exists, err = nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the Deployment again so we have no leftovers:
	err = nativekubernetes.DeleteDeployment(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// The delete function is idempotent as well.
	// If the Deployment is already absent no error will be raised.
	err = nativekubernetes.DeleteDeployment(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)

	exists, err = nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)
}
