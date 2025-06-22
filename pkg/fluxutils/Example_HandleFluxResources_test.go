package fluxutils_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// This example shows how to handle flux resources.
//
// To start this test use:
//
//	bash -c "cd pkg/fluxutils && go test -v -run Test_HandleFluxResources"
func Test_HandleFluxResources(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := continuousintegration.GetDefaultKindClusterName()
	const namespaceName = "flux-system"
	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)
	defer kindutils.DeleteClusterByNameIfInContinuousIntegration(ctx, clusterName)
	// Get Kubernetes cluster:
	cluster, err := nativekubernetes.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)
	// Install flux using flux-operator
	_, err = fluxutils.InstallFlux(ctx, &fluxparameteroptions.InstalFluxOptions{
		KubernetesCluster: cluster,
		Namespace:         namespaceName,
	})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Get the deployed flux
	fluxDeployment, err := fluxutils.GetFluxDeployment(cluster, namespaceName)
	require.NoError(t, err)

	// Define example resource names:
	const gitRepoName = "example-repo"
	const kustomizationName = "example-kustomization"

	// Ensure example resources absent:
	err = fluxDeployment.DeleteGitRepository(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	err = fluxDeployment.DeleteKustomization(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)

	// Check the example GitRepository is absent before registering callback functions:
	exists, err := fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check the example GitRepository is absent before registering callback functions:
	exists, err = fluxDeployment.KustomizationExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Define counters and context to watch the flux objects:
	var grCreateCounter, grUpdateCounter, grDeleteCounter int // gr for GitRepository
	var kuCreateCounter, kuUpdateCounter, kuDeleteCounter int // ku for Kustomization
	ctxWatch, cancel := context.WithCancel(ctx)               // ensure we can cancel the watching

	// Start watching the GitRepository
	err = fluxDeployment.WatchGitRepository(
		ctxWatch,
		gitRepoName,
		namespaceName,
		func(*unstructured.Unstructured) { grCreateCounter++ },
		func(*unstructured.Unstructured) { grUpdateCounter++ },
		func(*unstructured.Unstructured) { grDeleteCounter++ },
	)
	require.NoError(t, err)
	defer cancel()

	// Start watching the Kustomization
	err = fluxDeployment.WatchKustomization(
		ctxWatch,
		kustomizationName,
		namespaceName,
		func(*unstructured.Unstructured) { kuCreateCounter++ },
		func(*unstructured.Unstructured) { kuUpdateCounter++ },
		func(*unstructured.Unstructured) { kuDeleteCounter++ },
	)
	require.NoError(t, err)
	defer cancel()

	// check no callback called for GitRepository
	time.Sleep(100 * time.Millisecond)
	require.EqualValues(t, 0, grCreateCounter)
	require.EqualValues(t, 0, grUpdateCounter)
	require.EqualValues(t, 0, grDeleteCounter)

	// check no callback called for Kustomization
	require.EqualValues(t, 0, kuCreateCounter)
	require.EqualValues(t, 0, kuUpdateCounter)
	require.EqualValues(t, 0, kuDeleteCounter)

	// Define a fluxcd GitRepository:
	gitRepoYaml := "---\n"
	gitRepoYaml += "apiVersion: source.toolkit.fluxcd.io/v1\n"
	gitRepoYaml += "kind: GitRepository\n"
	gitRepoYaml += "metadata:\n"
	gitRepoYaml += "  name: " + gitRepoName + "\n"
	gitRepoYaml += "  namespace: " + namespaceName + "\n"
	gitRepoYaml += "spec:\n"
	gitRepoYaml += "  interval: 5m0s\n"
	gitRepoYaml += "  url: https://asciich.ch/example/repo\n"
	gitRepoYaml += "  ref:\n"
	gitRepoYaml += "    branch: master\n"
	// Create the GitRepository:
	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: gitRepoYaml})
	require.NoError(t, err)

	// Check the example GitRepository exists:
	exists, err = fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Check Create counter increased for GitRepository:
	require.EqualValues(t, 1, grCreateCounter)
	require.EqualValues(t, 0, grUpdateCounter)
	require.EqualValues(t, 0, grDeleteCounter)

	// Define a fluxcd Kustomization:
	kustomizationYaml := "apiVersion: kustomize.toolkit.fluxcd.io/v1\n"
	kustomizationYaml += "kind: Kustomization\n"
	kustomizationYaml += "metadata:\n"
	kustomizationYaml += "  name: " + kustomizationName + "\n"
	kustomizationYaml += "  namespace: " + namespaceName + "\n"
	kustomizationYaml += "spec:\n"
	kustomizationYaml += "  interval: 10m\n"
	kustomizationYaml += "  targetNamespace: " + namespaceName + "\n"
	kustomizationYaml += "  sourceRef:\n"
	kustomizationYaml += "    kind: GitRepository\n"
	kustomizationYaml += "    name: " + gitRepoName + "\n"
	kustomizationYaml += "  path: \"./kustomize\"\n"
	kustomizationYaml += "  prune: true\n"
	kustomizationYaml += "  timeout: 1m\n"

	// Create the Kustomization:
	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: kustomizationYaml})
	require.NoError(t, err)

	// Check the example Kustomization exists:
	exists, err = fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.True(t, exists)

	// Check Create counter increased for GitRepository:
	require.EqualValues(t, 1, kuCreateCounter)
	require.EqualValues(t, 0, kuUpdateCounter)
	require.EqualValues(t, 0, kuDeleteCounter)

	// Give the resources some time to settle:
	time.Sleep(time.Second * 5)

	// Get the status of the GitRepository:
	status, err := fluxDeployment.GetGitRepositoryStatusMessage(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.Contains(t, status, "failed to checkout and determine revision: unable to clone ") // The repo of this example does not exist.

	// Get the status of the Kustomization:
	status, err = fluxDeployment.GetKustomizationStatusMessage(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)
	require.Contains(t, status, "Source artifact not found, retrying in 5s") // The repo of this example does not exist.

	// Check update counter increased for GitRepository:
	require.EqualValues(t, 1, grCreateCounter)
	require.GreaterOrEqual(t, grUpdateCounter, 3) // Every status update generates an updated version
	require.EqualValues(t, 0, grDeleteCounter)

	// Check update counter increased for Kustomization:
	require.EqualValues(t, 1, kuCreateCounter)
	require.GreaterOrEqual(t, kuUpdateCounter, 2) // Every status update generates an updated version
	require.EqualValues(t, 0, kuDeleteCounter)

	// Delete the example GitRepository:
	err = fluxDeployment.DeleteGitRepository(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)

	// Check if example GitRepository is absent:
	exists, err = fluxDeployment.GitRepositoryExists(ctx, gitRepoName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check delete counter increased for GitRepository:
	require.EqualValues(t, 1, grCreateCounter)
	require.GreaterOrEqual(t, grUpdateCounter, 3) // Every status update generates an updated version
	require.EqualValues(t, 1, grDeleteCounter)

	// Delete the example Kustomization:
	err = fluxDeployment.DeleteKustomization(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)

	// Check if example GitRepository is absent:
	exists, err = fluxDeployment.KustomizationExists(ctx, kustomizationName, namespaceName)
	require.NoError(t, err)
	require.False(t, exists)

	// Check delete counter increased for GitRepository:
	require.EqualValues(t, 1, kuCreateCounter)
	require.GreaterOrEqual(t, kuUpdateCounter, 2) // Every status update generates an updated version
	require.EqualValues(t, 1, kuDeleteCounter)
}
