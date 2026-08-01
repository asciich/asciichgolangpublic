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

func Test_Example_KubernetesConfigMapExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const configMapName = "example-configmap-test"

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

	// Create an example configMap using YAML
	configMapYaml := ""
	configMapYaml += "apiVersion: v1\n"
	configMapYaml += "kind: ConfigMap\n"
	configMapYaml += "metadata:\n"
	configMapYaml += "  name: " + configMapName + "\n"
	configMapYaml += "  namespace: " + namespaceName + "\n"
	configMapYaml += "data:\n"
	configMapYaml += "  key1: value1\n"
	configMapYaml += "  key2: value2\n"

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
name: "Kubernetes configmap exists"
test_cases:
  - name: "Test configMap exists"
    test_type: kubernetes_configmap_exists
    resource_name: example-configmap-test
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that an existing configMap is detected"

  - name: "Test nonexistent configMap"
    test_type: kubernetes_configmap_exists
    resource_name: configmap-does-not-exist
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that a nonexistent configMap is detected"
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
