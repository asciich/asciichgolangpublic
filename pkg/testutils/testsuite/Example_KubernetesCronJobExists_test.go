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

func Test_Example_KubernetesCronJobExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const cronJobName = "example-cronjob-test"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the cronJob object
	cronJob, err := namespace.GetCronJobByName(cronJobName)
	require.NoError(t, err)

	// Ensure the cronJob is absent before testing
	err = namespace.DeleteCronJobByName(ctx, cronJobName)
	require.NoError(t, err)
	exists, err := cronJob.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example cronJob using YAML
	cronJobYaml := ""
	cronJobYaml += "apiVersion: batch/v1\n"
	cronJobYaml += "kind: CronJob\n"
	cronJobYaml += "metadata:\n"
	cronJobYaml += "  name: " + cronJobName + "\n"
	cronJobYaml += "  namespace: " + namespaceName + "\n"
	cronJobYaml += "spec:\n"
	cronJobYaml += "  schedule: \"*/5 * * * *\"\n"
	cronJobYaml += "  jobTemplate:\n"
	cronJobYaml += "    spec:\n"
	cronJobYaml += "      template:\n"
	cronJobYaml += "        spec:\n"
	cronJobYaml += "          containers:\n"
	cronJobYaml += "          - name: hello\n"
	cronJobYaml += "            image: busybox\n"
	cronJobYaml += "            command: [\"echo\", \"Hello World!\"]\n"
	cronJobYaml += "          restartPolicy: OnFailure\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: cronJobYaml,
	})
	require.NoError(t, err)

	// Clean up the cronJob after the test
	defer func() {
		err := namespace.DeleteCronJobByName(ctx, cronJobName)
		if err != nil {
			t.Logf("Warning: failed to delete cronjob: %v", err)
		}
	}()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes cronjob exists"
test_cases:
  - name: "Test cronJob exists"
    test_type: kubernetes_cronjob_exists
    resource_name: example-cronjob-test
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that an existing cronJob is detected"

  - name: "Test nonexistent cronJob"
    test_type: kubernetes_cronjob_exists
    resource_name: cronjob-does-not-exist
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that a nonexistent cronJob is detected"
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

// Test_Example_KubernetesCronJobExists_SSH tests running CronJob existence checks over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Creating a namespace and CronJob
// 3. Setting up an SSH server pod with key-based authentication
// 4. Using port forwarding to access the SSH server
// 5. Running kubernetes_cronjob_exists tests over SSH
func Test_Example_KubernetesCronJobExists_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Step 2: Setup SSH server in Kind cluster
	const namespaceName = "cronjob-ssh-test"
	const podName = "ssh-server-cronjob"

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
	const cronJobName = "example-cronjob-ssh"

	// Create an example CronJob using YAML
	cronJobYaml := ""
	cronJobYaml += "apiVersion: batch/v1\n"
	cronJobYaml += "kind: CronJob\n"
	cronJobYaml += "metadata:\n"
	cronJobYaml += "  name: " + cronJobName + "\n"
	cronJobYaml += "  namespace: " + namespaceName + "\n"
	cronJobYaml += "spec:\n"
	cronJobYaml += "  schedule: \"*/5 * * * *\"\n"
	cronJobYaml += "  jobTemplate:\n"
	cronJobYaml += "    spec:\n"
	cronJobYaml += "      template:\n"
	cronJobYaml += "        spec:\n"
	cronJobYaml += "          containers:\n"
	cronJobYaml += "          - name: hello\n"
	cronJobYaml += "            image: busybox\n"
	cronJobYaml += "            command: [\"echo\", \"Hello World!\"]\n"
	cronJobYaml += "          restartPolicy: OnFailure\n"

	_, err = setupResult.Namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: cronJobYaml,
	})
	require.NoError(t, err)

	// Clean up the CronJob after the test
	defer func() {
		err := setupResult.Namespace.DeleteCronJobByName(ctx, cronJobName)
		if err != nil {
			t.Logf("Warning: failed to delete cronjob: %v", err)
		}
	}()

	// Step 3: Run the test suite with SSH configuration
	testSuite := &testsuite.TestSuite{
		Name:                  "SSH CronJob exists test",
		Description:           "Test SSH kubernetes_cronjob_exists execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:         "Test CronJob exists via SSH",
				TestType:     "kubernetes_cronjob_exists",
				ResourceName: cronJobName,
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				Description:  "Check that an existing CronJob is detected via SSH",
			},
			{
				Name:         "Test nonexistent CronJob via SSH",
				TestType:     "kubernetes_cronjob_exists",
				ResourceName: "cronjob-does-not-exist-ssh",
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				Description:  "Check that a nonexistent CronJob is detected via SSH",
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
	require.EqualValues(t, 1, passed, "One test should pass (existing CronJob)")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed, "One test should fail (nonexistent CronJob)")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed, "Overall test suite should be failed because one test failed")

	t.Logf("SSH CronJob test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
