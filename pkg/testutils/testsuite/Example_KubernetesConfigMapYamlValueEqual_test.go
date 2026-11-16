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

func Test_Example_KubernetesConfigMapYamlValueEqual(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const configMapName = "example-configmap-yamlvalue-test"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the configMap object
	configMap, err := namespace.GetConfigMapByName(configMapName)
	require.NoError(t, err)

	// Ensure the configMap is absent before testing
	err = namespace.DeleteConfigMapByName(ctx, configMapName)
	require.NoError(t, err)
	exists, err := configMap.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example configMap holding a yaml document:
	configMapYaml := ""
	configMapYaml += "apiVersion: v1\n"
	configMapYaml += "kind: ConfigMap\n"
	configMapYaml += "metadata:\n"
	configMapYaml += "  name: " + configMapName + "\n"
	configMapYaml += "  namespace: " + namespaceName + "\n"
	configMapYaml += "data:\n"
	configMapYaml += "  values.yaml: |\n"
	configMapYaml += "    server:\n"
	configMapYaml += "      name: my-server\n"
	configMapYaml += "      listeners:\n"
	configMapYaml += "        - port: 8080\n"
	configMapYaml += "          host: localhost\n"
	configMapYaml += "    replicas: 3\n"
	configMapYaml += "    enabled: true\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: configMapYaml,
	})
	require.NoError(t, err)

	// Clean up the configMap after the test
	defer func() {
		err := namespace.DeleteConfigMapByName(ctx, configMapName)
		if err != nil {
			t.Logf("Warning: failed to delete configmap: %v", err)
		}
	}()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes configmap yaml value equal"
test_cases:
  - name: "Test string value equals expected value"
    test_type: kubernetes_configmap_yamlvalue_equal
    resource_name: example-configmap-yamlvalue-test
    namespace: default
    cluster: kind-asciichgolangpublic
    configmap_key: values.yaml
    yaml_path: server.name
    expected_value: my-server
    description: "Check that the server name equals the expected value"

  - name: "Test array element value equals expected value"
    test_type: kubernetes_configmap_yamlvalue_equal
    resource_name: example-configmap-yamlvalue-test
    namespace: default
    cluster: kind-asciichgolangpublic
    configmap_key: values.yaml
    yaml_path: server.listeners[0].port
    expected_value: 8080
    description: "Check that the first listener port equals the expected value"

  - name: "Test value differs from expected value"
    test_type: kubernetes_configmap_yamlvalue_equal
    resource_name: example-configmap-yamlvalue-test
    namespace: default
    cluster: kind-asciichgolangpublic
    configmap_key: values.yaml
    yaml_path: server.name
    expected_value: other-server
    description: "Check that a differing value is detected"

  - name: "Test yaml path resolves to no value"
    test_type: kubernetes_configmap_yamlvalue_equal
    resource_name: example-configmap-yamlvalue-test
    namespace: default
    cluster: kind-asciichgolangpublic
    configmap_key: values.yaml
    yaml_path: server.does.not.exist
    expected_value: some-value
    description: "Check that a non resolving yaml path is detected"
`)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, testSuitePath, &filesoptions.DeleteOptions{})

	// Use LogRecorder to verify no SSH commands are used for localhost tests
	ctx, logRecorder := logging.WithLogRecorder(ctx)

	// Run the test suite
	result, err := testsuite.RunFromFilePath(ctx, testSuitePath, &testutilsoptions.RunTestSuiteOptions{})
	require.NoError(t, err)

	// Verify no SSH commands were used (localhost test)
	logOutput := logRecorder.String()
	require.False(t, strings.Contains(logOutput, "Exec command 'ssh"), "No SSH commands should be used for localhost tests")

	// We can get the number of passed and failed test cases from the result:
	passed, err := result.GetNPassed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 2, passed)

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 2, failed)

	// We can log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because two tests failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed)
}

// Test_Example_KubernetesConfigMapYamlValueEqual_SSH tests running yaml value equality checks over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Creating a namespace and ConfigMap holding a yaml document
// 3. Setting up an SSH server pod with key-based authentication
// 4. Using port forwarding to access the SSH server
// 5. Running kubernetes_configmap_yamlvalue_equal tests over SSH
func Test_Example_KubernetesConfigMapYamlValueEqual_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "yamlvalue-ssh-test"
	const podName = "ssh-server-yamlvalue"

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
	const configMapName = "example-configmap-yamlvalue-ssh"

	// Create an example ConfigMap holding a yaml document:
	configMapYaml := ""
	configMapYaml += "apiVersion: v1\n"
	configMapYaml += "kind: ConfigMap\n"
	configMapYaml += "metadata:\n"
	configMapYaml += "  name: " + configMapName + "\n"
	configMapYaml += "  namespace: " + namespaceName + "\n"
	configMapYaml += "data:\n"
	configMapYaml += "  values.yaml: |\n"
	configMapYaml += "    server:\n"
	configMapYaml += "      name: my-server\n"
	configMapYaml += "      listeners:\n"
	configMapYaml += "        - port: 8080\n"
	configMapYaml += "          host: localhost\n"
	configMapYaml += "    replicas: 3\n"
	configMapYaml += "    enabled: true\n"

	_, err = setupResult.Namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: configMapYaml,
	})
	require.NoError(t, err)

	// Clean up the ConfigMap after the test
	defer func() {
		err := setupResult.Namespace.DeleteConfigMapByName(ctx, configMapName)
		if err != nil {
			t.Logf("Warning: failed to delete configmap: %v", err)
		}
	}()

	// Step 3: Run the test suite with SSH configuration
	testSuite := &testsuite.TestSuite{
		Name:                  "SSH configmap yaml value equal test",
		Description:           "Test SSH kubernetes_configmap_yamlvalue_equal execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:          "Test value equals expected value via SSH",
				TestType:      "kubernetes_configmap_yamlvalue_equal",
				ResourceName:  configMapName,
				Namespace:     namespaceName,
				Cluster:       "kind-asciichgolangpublic",
				ConfigMapKey:  "values.yaml",
				YamlPath:      "server.name",
				ExpectedValue: "my-server",
				Description:   "Check that the server name equals the expected value via SSH",
			},
			{
				Name:          "Test value differs from expected value via SSH",
				TestType:      "kubernetes_configmap_yamlvalue_equal",
				ResourceName:  configMapName,
				Namespace:     namespaceName,
				Cluster:       "kind-asciichgolangpublic",
				ConfigMapKey:  "values.yaml",
				YamlPath:      "server.listeners[0].port",
				ExpectedValue: 9999,
				Description:   "Check that a differing array element value is detected via SSH",
			},
		},
	}

	// Use LogRecorder to verify SSH commands are used for SSH tests
	ctx, logRecorder := logging.WithLogRecorder(ctx)

	// Run the test suite
	result, err := testSuite.Run(ctx)
	require.NoError(t, err, "Test suite execution failed")

	// Verify SSH commands were used (SSH test)
	logOutput := logRecorder.String()
	require.True(t, strings.Contains(logOutput, "Exec command 'ssh"), "SSH commands should be used for SSH tests")

	// Verify the test results
	passed, err := result.GetNPassed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, passed, "One test should pass (matching value)")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed, "One test should fail (differing value)")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed, "Overall test suite should be failed because one test failed")

	t.Logf("SSH configmap yaml value equal test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
