package fluxutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils"
	"github.com/asciich/asciichgolangpublic/pkg/fluxutils/fluxparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
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

	// Get the namespace where flux is installed.
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Define a fluxcd GitRepository:
	gitRepoYaml := "---\n"
	gitRepoYaml += "apiVersion: source.toolkit.fluxcd.io/v1\n"
	gitRepoYaml += "kind: GitRepository\n"
	gitRepoYaml += "metadata:\n"
	gitRepoYaml += "  name: example-repo\n"
	gitRepoYaml += "  namespace: " + namespaceName + "\n"
	gitRepoYaml += "spec:\n"
	gitRepoYaml += "  interval: 5m0s\n"
	gitRepoYaml += "  url: https://asciich.ch/example/repo\n"
	gitRepoYaml += "  ref:\n"
	gitRepoYaml += "    branch: master\n"
	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: gitRepoYaml})
	require.NoError(t, err)
}
