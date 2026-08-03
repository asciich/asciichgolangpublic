package testsuite_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testsuite"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsoptions"
)

func Test_Example_KubernetesDeploymentExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const deploymentName = "example-deployment-test"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the deployment object
	deployment, err := namespace.GetDeploymentByName(deploymentName)
	require.NoError(t, err)

	// Ensure the deployment is absent before testing
	err = namespace.DeleteDeploymentByName(ctx, deploymentName)
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
	deploymentYaml += "  replicas: 1\n"
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

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: deploymentYaml,
	})
	require.NoError(t, err)

	// Clean up the deployment after the test
	defer func() {
		err := namespace.DeleteDeploymentByName(ctx, deploymentName)
		if err != nil {
			t.Logf("Warning: failed to delete deployment: %v", err)
		}
	}()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes deployment exists"
test_cases:
  - name: "Test deployment exists"
    test_type: kubernetes_deployment_exists
    resource_name: example-deployment-test
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that an existing deployment is detected"

  - name: "Test nonexistent deployment"
    test_type: kubernetes_deployment_exists
    resource_name: deployment-does-not-exist
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that a nonexistent deployment is detected"
`)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, testSuitePath, &filesoptions.DeleteOptions{})

	// Run the test suite
	result, err := testsuite.RunFromFilePath(ctx, testSuitePath, &testutilsoptions.RunTestSuiteOptions{})
	require.NoError(t, err)

	// We can get the number of passed and failed test cases from the result:
	passed, err := result.GetNPassed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, passed)

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed)

	// We can log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed)
}
