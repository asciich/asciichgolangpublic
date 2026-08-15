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
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetestestsshserver"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testcase"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testsuite"
	"github.com/asciich/asciichgolangpublic/pkg/testutils/testutilsoptions"
)

func Test_Example_KubernetesValidateSshKeyInSecret(t *testing.T) {
	// Use a context with verbose output:
	ctx := contextutils.ContextVerbose()

	// Ensure a local kind cluster is available for testing:
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Generate an SSH key pair for testing
	keyPair, err := sshutils.GenerateKeyPair(sshutils.SSH_KEY_TYPE_ED25519, nil)
	require.NoError(t, err)

	privateKey, err := keyPair.GetPrivateKey()
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "test-validate-ssh-key"
	const secretName = "ssh-private-key-secret"
	const secretKey = "id_ed25519"
	const sshServerPodName = "ssh-server-validate-key"
	const sshUsername = "testuser"

	// Create namespace
	namespace, err := cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)
	defer namespace.Delete(ctx)

	// Deploy a test SSH server with public key authentication
	publicKeyString := keyPair.PublicKey.KeyMaterial

	sshServerPod, err := kubernetestestsshserver.StartTestSshServerInCluster(ctx, cluster, &kubernetestestsshserver.StartTestSshServerOptions{
		KubernetesNamespace: namespaceName,
		PodName:             sshServerPodName,
		SSHUsername:         sshUsername,
		SSHPublicKey:        publicKeyString,
	})
	require.NoError(t, err)
	require.NotNil(t, sshServerPod)
	defer sshServerPod.Delete(ctx)

	// Create a secret containing the SSH private key
	err = namespace.DeleteSecretByName(ctx, secretName)
	require.NoError(t, err)

	_, err = namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{
			secretKey: []byte(privateKey.KeyMaterial),
		},
	})
	require.NoError(t, err)

	defer func() {
		err := namespace.DeleteSecretByName(ctx, secretName)
		if err != nil {
			t.Logf("Warning: failed to delete secret: %v", err)
		}
	}()

	// The SSH server is reachable via its Service DNS
	sshServerHost := sshServerPodName + "." + namespaceName + ".svc.cluster.local"

	// Define the testsuite as temporary file:
	testSuitePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, `---
name: "Kubernetes validate SSH key in secret"
test_cases:
  - name: "Test SSH key validation succeeds"
    test_type: kubernetes_validate_ssh_key_in_secret
    resource_name: `+secretName+`
    secret_key: `+secretKey+`
    namespace: `+namespaceName+`
    cluster: kind-asciichgolangpublic
    target_host: `+sshServerHost+`
    target_user: `+sshUsername+`
    target_port: 22
    description: "Check that a valid SSH key in a secret can authenticate"

  - name: "Test SSH key validation with nonexistent secret"
    test_type: kubernetes_validate_ssh_key_in_secret
    resource_name: nonexistent-ssh-secret
    secret_key: `+secretKey+`
    namespace: `+namespaceName+`
    cluster: kind-asciichgolangpublic
    target_host: `+sshServerHost+`
    target_user: `+sshUsername+`
    target_port: 22
    description: "Check that a nonexistent secret is handled correctly"
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

	// The overall status is failed (because one test failed as expected - nonexistent secret):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed)
}

// Test_Example_KubernetesValidateSshKeyInSecret_SSH tests running SSH key validation over SSH to a pod in a Kind cluster.
// It demonstrates:
// 1. Starting a Kind cluster
// 2. Creating a namespace, SSH server pod, and secret with private key
// 3. Setting up an SSH server pod for test suite execution over SSH
// 4. Using port forwarding to access the SSH server
// 5. Running kubernetes_validate_ssh_key_in_secret tests over SSH
func Test_Example_KubernetesValidateSshKeyInSecret_SSH(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Step 1: Get or create Kind cluster
	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	// Generate an SSH key pair for the target SSH server (the one we validate against)
	targetKeyPair, err := sshutils.GenerateKeyPair(sshutils.SSH_KEY_TYPE_ED25519, nil)
	require.NoError(t, err)

	targetPrivateKey, err := targetKeyPair.GetPrivateKey()
	require.NoError(t, err)

	// Define test constants
	const namespaceName = "validate-ssh-key-ssh-test"
	const secretName = "ssh-private-key-secret"
	const secretKey = "id_ed25519"
	const sshTargetPodName = "ssh-target-server"
	const sshTargetUsername = "targetuser"

	// Step 2: Setup SSH server for test suite execution (the "jump host")
	const jumpPodName = "ssh-server-jump"

	setupResult, cleanup, err := SetupSSHServerInKind(ctx, t, cluster, namespaceName, jumpPodName)
	require.NoError(t, err)
	defer cleanup()

	// Deploy the target SSH server with public key authentication
	targetPublicKeyString := targetKeyPair.PublicKey.KeyMaterial

	sshTargetPod, err := kubernetestestsshserver.StartTestSshServerInCluster(ctx, cluster, &kubernetestestsshserver.StartTestSshServerOptions{
		KubernetesNamespace: namespaceName,
		PodName:             sshTargetPodName,
		SSHUsername:         sshTargetUsername,
		SSHPublicKey:        targetPublicKeyString,
	})
	require.NoError(t, err)
	require.NotNil(t, sshTargetPod)
	defer sshTargetPod.Delete(ctx)

	// Create a secret containing the SSH private key for validation
	err = setupResult.Namespace.DeleteSecretByName(ctx, secretName)
	require.NoError(t, err)

	_, err = setupResult.Namespace.CreateSecret(ctx, secretName, &kubernetesparameteroptions.CreateSecretOptions{
		SecretData: map[string][]byte{
			secretKey: []byte(targetPrivateKey.KeyMaterial),
		},
	})
	require.NoError(t, err)

	defer func() {
		err := setupResult.Namespace.DeleteSecretByName(ctx, secretName)
		if err != nil {
			t.Logf("Warning: failed to delete secret: %v", err)
		}
	}()

	// Write jump host private key to temporary file
	tmpFile, err := os.CreateTemp("", "ssh_test_key_*")
	require.NoError(t, err)
	_, err = tmpFile.WriteString(setupResult.KeyPair.PrivateKey.KeyMaterial)
	require.NoError(t, err)
	err = tmpFile.Chmod(0600)
	require.NoError(t, err)
	err = tmpFile.Close()
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// The target SSH server is reachable via its Service DNS
	sshTargetHost := sshTargetPodName + "." + namespaceName + ".svc.cluster.local"

	// Step 3: Run the test suite with SSH configuration
	testSuite := &testsuite.TestSuite{
		Name:                  "SSH validate SSH key in secret test",
		Description:           "Test SSH kubernetes_validate_ssh_key_in_secret execution on Kubernetes pod",
		SSHHost:               "localhost",
		SSHUser:               "testuser",
		SSHPort:               setupResult.LocalPort,
		SSHSkipHostValidation: true,
		SSHPrivateKeyFile:     tmpFile.Name(),
		TestCases: []*testcase.TestCase{
			{
				Name:         "Test SSH key validation succeeds via SSH",
				TestType:     "kubernetes_validate_ssh_key_in_secret",
				ResourceName: secretName,
				SecretKey:    secretKey,
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				TargetHost:   sshTargetHost,
				TargetUser:   sshTargetUsername,
				TargetPort:   22,
				Description:  "Check that a valid SSH key in a secret can authenticate via SSH",
			},
			{
				Name:         "Test SSH key validation with nonexistent secret via SSH",
				TestType:     "kubernetes_validate_ssh_key_in_secret",
				ResourceName: "nonexistent-ssh-secret",
				SecretKey:    secretKey,
				Namespace:    namespaceName,
				Cluster:      "kind-asciichgolangpublic",
				TargetHost:   sshTargetHost,
				TargetUser:   sshTargetUsername,
				TargetPort:   22,
				Description:  "Check that a nonexistent secret is handled correctly via SSH",
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
	require.EqualValues(t, 1, passed, "One test should pass (valid SSH key)")

	failed, err := result.GetNFailed(ctx)
	require.NoError(t, err)
	require.EqualValues(t, 1, failed, "One test should fail (nonexistent secret)")

	// Log the result
	err = result.LogResult(ctx)
	require.NoError(t, err)

	// The overall status is failed (because one test failed as expected):
	isPassed, err := result.IsPassed(ctx)
	require.NoError(t, err)
	require.False(t, isPassed, "Overall test suite should be failed because one test failed")

	t.Logf("SSH validate SSH key in secret test completed successfully on %s:%d!", setupResult.NamespaceName, setupResult.LocalPort)
}
