package testsuite_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
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
	// Use LogRecorder to verify no SSH commands are used for localhost tests
	ctx, logRecorder := logging.WithLogRecorder(ctx)

	result, err := testsuite.RunFromFilePath(ctx, testSuitePath, &testutilsoptions.RunTestSuiteOptions{})
	require.NoError(t, err)

	// Verify no SSH commands were used (localhost test)
	logOutput := logRecorder.String()
	require.False(t, strings.Contains(logOutput, "Exec command 'ssh"), "No SSH commands should be used for localhost tests")

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

// Test_Example_KubernetesDeploymentExists_SSH tests running Deployment existence checks over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Creating a namespace and Deployment
// 3. Setting up an SSH server pod with key-based authentication
// 4. Using port forwarding to access the SSH server
// 5. Running kubernetes_deployment_exists tests over SSH
func Test_Example_KubernetesDeploymentExists_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "deployment-ssh-test"
	const podName = "ssh-server-deployment"

	setupResult, cleanup, err := SetupSSHServerInKind(ctx, t, cluster, namespaceName, podName)
	require.NoError(t, err)
	defer cleanup()

	// Write private key to temporary file (user manages lifecycle)
	tmpFile, err := os.CreateTemp("", "ssh_test_key_*")
	require.NoError(t, err)
	_, err = tmpFile.WriteString(setupResult.KeyPair.PrivateKey.KeyMaterial)
	require.NoError(t, err)
	err = tmpFile.Chmod(0600)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up temp file

	// Define test constants
	const deploymentName = "example-deployment-ssh"

	// Create an example Deployment using YAML
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

	_, err = setupResult.Namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: deploymentYaml,
	})
	require.NoError(t, err)

	// Clean up the Deployment after the test
	defer func() {
		err := setupResult.Namespace.DeleteDeploymentByName(ctx, deploymentName)
		if err != nil {
			t.Logf("Warning: failed to delete deployment: %v", err)
		}
	}()

	// Step 3: Run the test suite with SSH configuration
	testSuite := &testsuite.TestSuite{
		Name:                  "SSH Deployment exists test",
		Description:           "Test SSH kubernetes_deployment_exists execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:         "Test Deployment exists via SSH",
				TestType:     "kubernetes_deployment_exists",
				ResourceName: deploymentName,
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				Description:  "Check that an existing Deployment is detected via SSH",
			},
			{
				Name:         "Test nonexistent Deployment via SSH",
				TestType:     "kubernetes_deployment_exists",
				ResourceName: "deployment-does-not-exist-ssh",
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				Description:  "Check that a nonexistent Deployment is detected via SSH",
			},
		},
	}

	// Run the test suite
	// Use LogRecorder to verify SSH commands are used for SSH tests
	ctx, logRecorder := logging.WithLogRecorder(ctx)

	result, err := testSuite.Run(ctx)
	require.NoError(t, err, "Test suite execution failed")

	// Verify SSH commands were used (SSH test)
	logOutput := logRecorder.String()
	require.True(t, strings.Contains(logOutput, "Exec command 'ssh"), "SSH commands should be used for SSH tests")

	// Verify the test results
	passed, err := result.GetNPassed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, passed, "One test should pass (existing Deployment)")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed, "One test should fail (nonexistent Deployment)")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed, "Overall test suite should be failed because one test failed")

	t.Logf("SSH Deployment test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
