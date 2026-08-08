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

func Test_Example_KubernetesReplicaSetExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const replicaSetName = "example-replicaset-test"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the replicaSet object
	replicaSet, err := namespace.GetReplicaSetByName(replicaSetName)
	require.NoError(t, err)

	// Ensure the replicaSet is absent before testing
	err = namespace.DeleteReplicaSetByName(ctx, replicaSetName)
	require.NoError(t, err)
	exists, err := replicaSet.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example replicaSet using YAML
	replicaSetYaml := ""
	replicaSetYaml += "apiVersion: apps/v1\n"
	replicaSetYaml += "kind: ReplicaSet\n"
	replicaSetYaml += "metadata:\n"
	replicaSetYaml += "  name: " + replicaSetName + "\n"
	replicaSetYaml += "  namespace: " + namespaceName + "\n"
	replicaSetYaml += "spec:\n"
	replicaSetYaml += "  replicas: 1\n"
	replicaSetYaml += "  selector:\n"
	replicaSetYaml += "    matchLabels:\n"
	replicaSetYaml += "      app: " + replicaSetName + "\n"
	replicaSetYaml += "  template:\n"
	replicaSetYaml += "    metadata:\n"
	replicaSetYaml += "      labels:\n"
	replicaSetYaml += "        app: " + replicaSetName + "\n"
	replicaSetYaml += "    spec:\n"
	replicaSetYaml += "      containers:\n"
	replicaSetYaml += "      - name: ubuntu\n"
	replicaSetYaml += "        image: ubuntu\n"
	replicaSetYaml += "        command: [\"bash\", \"-c\", \"sleep 1m\"]\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: replicaSetYaml,
	})
	require.NoError(t, err)

	// Clean up the replicaSet after the test
	defer func() {
		err := namespace.DeleteReplicaSetByName(ctx, replicaSetName)
		if err != nil {
			t.Logf("Warning: failed to delete replicaset: %v", err)
		}
	}()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes replicaset exists"
test_cases:
  - name: "Test replicaSet exists"
    test_type: kubernetes_replicaset_exists
    resource_name: example-replicaset-test
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that an existing replicaSet is detected"

  - name: "Test nonexistent replicaSet"
    test_type: kubernetes_replicaset_exists
    resource_name: replicaset-does-not-exist
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that a nonexistent replicaSet is detected"
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

// Test_Example_KubernetesReplicaSetExists_SSH tests running ReplicaSet existence checks over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Creating a namespace and ReplicaSet
// 3. Setting up an SSH server pod with key-based authentication
// 4. Using port forwarding to access the SSH server
// 5. Running kubernetes_replicaset_exists tests over SSH
func Test_Example_KubernetesReplicaSetExists_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "replicaset-ssh-test"
	const podName = "ssh-server-replicaset"

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
	const replicaSetName = "example-replicaset-ssh"

	// Create an example ReplicaSet using YAML
	replicaSetYaml := ""
	replicaSetYaml += "apiVersion: apps/v1\n"
	replicaSetYaml += "kind: ReplicaSet\n"
	replicaSetYaml += "metadata:\n"
	replicaSetYaml += "  name: " + replicaSetName + "\n"
	replicaSetYaml += "  namespace: " + namespaceName + "\n"
	replicaSetYaml += "spec:\n"
	replicaSetYaml += "  replicas: 1\n"
	replicaSetYaml += "  selector:\n"
	replicaSetYaml += "    matchLabels:\n"
	replicaSetYaml += "      app: " + replicaSetName + "\n"
	replicaSetYaml += "  template:\n"
	replicaSetYaml += "    metadata:\n"
	replicaSetYaml += "      labels:\n"
	replicaSetYaml += "        app: " + replicaSetName + "\n"
	replicaSetYaml += "    spec:\n"
	replicaSetYaml += "      containers:\n"
	replicaSetYaml += "      - name: ubuntu\n"
	replicaSetYaml += "        image: ubuntu\n"
	replicaSetYaml += "        command: [\"bash\", \"-c\", \"sleep 1m\"]\n"

	_, err = setupResult.Namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: replicaSetYaml,
	})
	require.NoError(t, err)

	// Clean up the ReplicaSet after the test
	defer func() {
		err := setupResult.Namespace.DeleteReplicaSetByName(ctx, replicaSetName)
		if err != nil {
			t.Logf("Warning: failed to delete replicaset: %v", err)
		}
	}()

	// Step 3: Run the test suite with SSH configuration
	testSuite := &testsuite.TestSuite{
		Name:                  "SSH ReplicaSet exists test",
		Description:           "Test SSH kubernetes_replicaset_exists execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:         "Test ReplicaSet exists via SSH",
				TestType:     "kubernetes_replicaset_exists",
				ResourceName: replicaSetName,
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				Description:  "Check that an existing ReplicaSet is detected via SSH",
			},
			{
				Name:         "Test nonexistent ReplicaSet via SSH",
				TestType:     "kubernetes_replicaset_exists",
				ResourceName: "replicaset-does-not-exist-ssh",
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				Description:  "Check that a nonexistent ReplicaSet is detected via SSH",
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
	require.EqualValues(t, 1, passed, "One test should pass (existing ReplicaSet)")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed, "One test should fail (nonexistent ReplicaSet)")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed, "Overall test suite should be failed because one test failed")

	t.Logf("SSH ReplicaSet test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
