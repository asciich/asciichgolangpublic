package nativekubernetes_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_CreateAndDeleteDeployment(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("happy path", func(t *testing.T) {

		const deploymentName = "testdeployment"
		const namespaceName = "default"

		err = nativekubernetes.DeleteDeployment(ctx, clientset, deploymentName, namespaceName)
		require.NoError(t, err)

		exists, err := nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
		require.NoError(t, err)
		require.False(t, exists)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreateDeployment(ctx, config, &kubernetesparameteroptions.RunCommandOptions{
				Namespace:      namespaceName,
				DeploymentName: deploymentName,
				Image:          "ubuntu",
				Command:        []string{"bash", "-c", "sleep 1m"},
				Replicas:       2,
			})
			require.NoError(t, err)

			exists, err = nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
			require.NoError(t, err)
			require.True(t, exists)

			err = nativekubernetes.DeleteDeployment(ctx, clientset, deploymentName, namespaceName)
			require.NoError(t, err)

			exists, err = nativekubernetes.DeploymentExists(ctx, clientset, deploymentName, namespaceName)
			require.NoError(t, err)
			require.False(t, exists)
		}
	})
}

func Test_WaitForDeploymentDeleted(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("already deleted Deployment", func(t *testing.T) {
		deploymentName := "testdeployment"
		namespaceName := "default"

		// Ensure Deployment is absent
		err = nativekubernetes.DeleteDeployment(ctx, clientset, deploymentName, namespaceName)
		require.NoError(t, err)

		// Check there's no wait for an already deleted Deployment:
		err = nativekubernetes.WaitForDeploymentDeleted(ctx, clientset, deploymentName, namespaceName, time.Second*1)
		require.NoError(t, err)
	})
}

func Test_ListDeployments(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	clusterName := "kubernetesutils"

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("create and delete Deployments with list in between", func(t *testing.T) {
		const namespaceName = "default"

		deploymentNames := []string{"listdeploy-1", "listdeploy-2", "listdeploy-3"}

		// Ensure all test Deployments are absent before starting
		for _, name := range deploymentNames {
			err = nativekubernetes.DeleteDeployment(ctx, clientset, name, namespaceName)
			require.NoError(t, err)
		}

		// Create Deployments one by one and verify list grows
		for i, name := range deploymentNames {
			err = nativekubernetes.CreateDeployment(ctx, config, &kubernetesparameteroptions.RunCommandOptions{
				Namespace:      namespaceName,
				DeploymentName: name,
				Image:          "ubuntu",
				Command:        []string{"bash", "-c", "sleep 1m"},
				Replicas:       1,
			})
			require.NoError(t, err)

			listed, err := nativekubernetes.ListDeployments(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, created := range deploymentNames[:i+1] {
				require.Contains(t, listed, created)
			}
			for _, notYetCreated := range deploymentNames[i+1:] {
				require.NotContains(t, listed, notYetCreated)
			}
		}

		// Delete Deployments one by one and verify list shrinks
		for i, name := range deploymentNames {
			err = nativekubernetes.DeleteDeployment(ctx, clientset, name, namespaceName)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListDeployments(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, deleted := range deploymentNames[:i+1] {
				require.NotContains(t, listed, deleted)
			}
			for _, stillPresent := range deploymentNames[i+1:] {
				require.Contains(t, listed, stillPresent)
			}
		}
	})
}
