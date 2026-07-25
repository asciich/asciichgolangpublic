package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeleteDeployment(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	clusterName := testutils.GetKindClusterNameForTest(t)

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.CreateCluster(ctx, clusterName)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// Define the namespace and deployment name we use for testing
	const namespaceName = "testnamespace"
	const deploymentName = "example-deployment"

	// Get the namespace (will be created implicitly when creating the deployment)
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the deployment object
	deployment, err := namespace.GetDeploymentByName(deploymentName)
	require.NoError(t, err)

	// Ensure the deployment is absent
	err = deployment.Delete(ctx)
	require.NoError(t, err)
	exists, err := deployment.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example deployment using YAML
	deploymentYaml := ""
	deploymentYaml += "apiVersion: apps/v1\n"
	deploymentYaml += "kind: Deployment\n"
	deploymentYaml += "metadata:\n"
	deploymentYaml += "  name: " + deploymentName + "\n"
	deploymentYaml += "  namespace: " + namespaceName + "\n"
	deploymentYaml += "spec:\n"
	deploymentYaml += "  replicas: 2\n"
	deploymentYaml += "  selector:\n"
	deploymentYaml += "    matchLabels:\n"
	deploymentYaml += "      app: " + deploymentName + "\n"
	deploymentYaml += "  template:\n"
	deploymentYaml += "    metadata:\n"
	deploymentYaml += "      labels:\n"
	deploymentYaml += "        app: " + deploymentName + "\n"
	deploymentYaml += "    spec:\n"
	deploymentYaml += "      containers:\n"
	deploymentYaml += "      - name: ubuntu\n"
	deploymentYaml += "        image: ubuntu\n"
	deploymentYaml += "        command: [\"bash\", \"-c\", \"sleep 1m\"]\n"

	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: deploymentYaml,
	})
	require.NoError(t, err)
	exists, err = deployment.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the deployment exists via the cluster
	names, err := namespace.ListObjectNames(&kubernetesparameteroptions.ListKubernetesObjectsOptions{
		Namespace:  namespaceName,
		ObjectType: "deployment",
	})
	require.NoError(t, err)
	require.Contains(t, names, deploymentName)

	// Delete the deployment
	err = deployment.Delete(ctx)
	require.NoError(t, err)
	exists, err = deployment.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)
}
