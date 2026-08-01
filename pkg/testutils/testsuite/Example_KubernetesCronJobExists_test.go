package testsuite_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
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
