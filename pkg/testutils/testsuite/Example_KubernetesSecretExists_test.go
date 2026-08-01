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

func Test_Example_KubernetesSecretExists(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "default"
	const secretName = "example-secret-test"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the secret object
	secret, err := namespace.GetSecretByName(secretName)
	require.NoError(t, err)

	// Ensure the secret is absent before testing
	err = namespace.DeleteSecretByName(ctx, secretName)
	require.NoError(t, err)
	exists, err := secret.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example secret using YAML
	secretYaml := ""
	secretYaml += "apiVersion: v1\n"
	secretYaml += "kind: Secret\n"
	secretYaml += "metadata:\n"
	secretYaml += "  name: " + secretName + "\n"
	secretYaml += "  namespace: " + namespaceName + "\n"
	secretYaml += "type: Opaque\n"
	secretYaml += "data:\n"
	secretYaml += "  username: YWRtaW4=\n"
	secretYaml += "  password: MWYyZDFlMmU2N2Rm\n"

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: secretYaml,
	})
	require.NoError(t, err)

	// Clean up the secret after the test
	defer func() {
		err := namespace.DeleteSecretByName(ctx, secretName)
		if err != nil {
			t.Logf("Warning: failed to delete secret: %v", err)
		}
	}()

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes secret exists"
test_cases:
  - name: "Test secret exists"
    test_type: kubernetes_secret_exists
    resource_name: example-secret-test
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that an existing secret is detected"

  - name: "Test nonexistent secret"
    test_type: kubernetes_secret_exists
    resource_name: secret-does-not-exist
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check that a nonexistent secret is detected"
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
