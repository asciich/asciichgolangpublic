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
