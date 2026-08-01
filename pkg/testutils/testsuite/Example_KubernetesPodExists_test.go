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

func Test_Example_KubernetesPodExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const podName = "example-pod-test"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the pod object
	pod, err := namespace.GetPodByName(podName)
	require.NoError(t, err)

	// Ensure the pod is absent before testing
	err = namespace.DeletePodByName(ctx, podName)
	require.NoError(t, err)
	exists, err := pod.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example pod using YAML
	podYaml := ""
	podYaml += "apiVersion: v1\n"
	podYaml += "kind: Pod\n"
	podYaml += "metadata:\n"
	podYaml += "  name: " + podName + "\n"
	podYaml += "  namespace: " + namespaceName + "\n"
	podYaml += "spec:\n"
	podYaml += "  containers:\n"
	podYaml += "  - name: ubuntu\n"
	podYaml += "    image: ubuntu\n"
	podYaml += "    command: [\"bash\", \"-c\", \"sleep 1m\"]\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: podYaml,
	})
	require.NoError(t, err)

	// Clean up the pod after the test
	defer func() {
		err := namespace.DeletePodByName(ctx, podName)
		if err != nil {
			t.Logf("Warning: failed to delete pod: %v", err)
		}
	}()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes pod exists"
test_cases:
  - name: "Test pod exists"
    test_type: kubernetes_pod_exists
    resource_name: example-pod-test
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that an existing pod is detected"

  - name: "Test nonexistent pod"
    test_type: kubernetes_pod_exists
    resource_name: pod-does-not-exist
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that a nonexistent pod is detected"
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
