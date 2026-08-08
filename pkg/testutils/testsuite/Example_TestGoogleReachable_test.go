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

func Test_Example_GoogleReachable(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Google reachable"
test_cases:
  - name: "Test HTTPS port open"
    test_type: tcp_port_open
    port: 443
    host: google.com
    description: "Check if the HTTPS port 443 is open on google.com"

  - name: "Test google reachable using another command"
    command: curl --fail google.com
    test_type: command
    description: "Check if we can get the main page of google.com using curl. You can use any tool installed on your machine."
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

	// We can get the number of passend and failed test cases from the result:
	passed, err := result.GetNPassed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 2, passed)

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, failed)

	// We can log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is passed:
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.True(t, isPassed)
}

// Test_Example_GoogleReachable_SSH tests running network reachability tests over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Setting up an SSH server pod with key-based authentication
// 3. Using port forwarding to access the SSH server
// 4. Running command and tcp_port_open tests over SSH
func Test_Example_GoogleReachable_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "google-ssh-test"
	const podName = "ssh-server-google"

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
		Name:                  "SSH Google reachable test",
		Description:           "Test SSH command and tcp_port_open execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:        "Test curl google via SSH",
				TestType:    "command",
				Command:     "curl --fail google.com",
				Description: "Check if we can get the main page of google.com using curl via SSH",
			},
			{
				Name:        "Test HTTPS port open via SSH",
				TestType:    "tcp_port_open",
				Port:        "443",
				Host:        "google.com",
				Description: "Check if the HTTPS port 443 is open on google.com via SSH",
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
	require.EqualValues(t, 2, passed, "Both tests should pass")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 0, failed, "No tests should fail")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is passed:
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.True(t, isPassed, "Overall test suite should be passed")

	t.Logf("SSH Google reachable test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
