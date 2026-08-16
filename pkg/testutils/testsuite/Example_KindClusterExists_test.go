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
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testsuite"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsoptions"
)

func Test_Example_KindClusterExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Get the cluster name for testing
	clusterName, err := cluster.GetName()
	require.NoError(t, err)

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kind cluster exists"
test_cases:
  - name: "Test existing kind cluster"
    test_type: kind_cluster_exists
    cluster: `+clusterName+`
    description: "Check that an existing kind cluster is detected"

  - name: "Test nonexistent kind cluster"
    test_type: kind_cluster_exists
    cluster: kind-cluster-does-not-exist
    description: "Check that a nonexistent kind cluster is detected"
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

// Test_Example_KindClusterExists_SSH tests running kind cluster existence checks over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Setting up an SSH server pod with key-based authentication
// 3. Using port forwarding to access the SSH server
// 4. Running kind_cluster_exists tests over SSH
func Test_Example_KindClusterExists_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Get the cluster name for testing
	clusterName, err := cluster.GetName()
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "kind-cluster-ssh-test"
	const podName = "ssh-server-kind-cluster"

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

	// Step 3: Run the test suite with SSH configuration
	testSuite := &testsuite.TestSuite{
		Name:                  "SSH Kind cluster exists test",
		Description:           "Test SSH kind_cluster_exists execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:        "Test kind cluster exists via SSH",
				TestType:    "kind_cluster_exists",
				Cluster:     clusterName,
				Description: "Check that an existing kind cluster is detected via SSH",
			},
			{
				Name:        "Test nonexistent kind cluster via SSH",
				TestType:    "kind_cluster_exists",
				Cluster:     "kind-cluster-does-not-exist-ssh",
				Description: "Check that a nonexistent kind cluster is detected via SSH",
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
	require.EqualValues(t, 1, passed, "One test should pass (existing kind cluster)")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed, "One test should fail (nonexistent kind cluster)")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed, "Overall test suite should be failed because one test failed")

	t.Logf("SSH Kind cluster test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
