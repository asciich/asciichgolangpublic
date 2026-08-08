package testsuite_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testsuite"
)

// Test_HostnameCommand_SSH tests running a command over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Creating a namespace
// 3. Setting up an SSH server pod with key-based authentication
// 4. Using port forwarding to access the SSH server
// 5. Running a command test over SSH
func Test_HostnameCommand_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "ssh-test"
	const podName = "ssh-server-test"

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
		Name:                  "SSH hostname test",
		Description:           "Test SSH command execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort, // Use the dynamically assigned local port
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(), // Path to private key file (user manages lifecycle)
		TestCases: []*testcase.TestCase{
			{
				Name:     "hostname",
				TestType: "command",
				Command:  "hostname",
			},
		},
	}

	// Run the test suite
	result, err := testSuite.Run(ctx)
	require.NoError(t, err, "Test suite execution failed")

	// Verify the test passed
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.True(t, isPassed, "SSH command test should pass")

	t.Logf("SSH test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
