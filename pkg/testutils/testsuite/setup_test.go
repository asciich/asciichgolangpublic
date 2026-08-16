package testsuite_test

// There is meanwhile the kubernetestestsshserver package available.
// TODO: All functionality of this file must go into the kubernetestestsshserver pacakge to make reduce code duplication and make it reuseable.

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubectlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
)

type SSHTestSetupResult struct {
	Namespace     kubernetesinterfaces.Namespace
	PodName       string
	SSHPort       int
	LocalPort     int
	KeyPair       *sshutils.SSHKeyPair
	CancelFunc    context.CancelFunc
	NamespaceName string
}

func SetupSSHServerInKind(ctx context.Context, t testing.TB, cluster kubernetesinterfaces.KubernetesCluster, namespaceName string, podName string) (*SSHTestSetupResult, func(), error) {
	t.Helper()

	namespace, err := cluster.GetNamespaceByName(namespaceName)
	if err != nil {
		namespace, err = cluster.CreateNamespaceByName(ctx, namespaceName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create namespace: %w", err)
		}
	}

	// Create ServiceAccount
	serviceAccountName := "ssh-readonly-sa"
	serviceAccountYaml := fmt.Sprintf(`
apiVersion: v1
kind: ServiceAccount
metadata:
  name: %s
  namespace: %s
`, serviceAccountName, namespaceName)

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: serviceAccountYaml,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ServiceAccount: %w", err)
	}
	t.Logf("Created ServiceAccount '%s' in namespace '%s'", serviceAccountName, namespaceName)

	// Create ClusterRole with permissions needed for test execution
	// Includes read-only for most resources, plus create/delete/patch for pods (needed by ValidateSSHKeyInSecret)
	roleName := "ssh-readonly-clusterrole"
	clusterRoleYaml := fmt.Sprintf(`
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: %s
rules:
- apiGroups: [""]
  resources: ["pods", "pods/log"]
  verbs: ["get", "list", "watch", "create", "delete", "patch"]
- apiGroups: [""]
  resources: ["services", "configmaps", "secrets", "namespaces", "nodes"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["batch"]
  resources: ["cronjobs"]
  verbs: ["get", "list", "watch"]
`, roleName)

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: clusterRoleYaml,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ClusterRole: %w", err)
	}
	t.Logf("Created ClusterRole '%s'", roleName)

	// Create ClusterRoleBinding to bind ServiceAccount to ClusterRole
	roleBindingName := "ssh-readonly-clusterrolebinding"
	roleBindingYaml := fmt.Sprintf(`
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: %s
subjects:
- kind: ServiceAccount
  name: %s
  namespace: %s
roleRef:
  kind: ClusterRole
  name: %s
  apiGroup: rbac.authorization.k8s.io
`, roleBindingName, serviceAccountName, namespaceName, roleName)

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: roleBindingYaml,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create ClusterRoleBinding: %w", err)
	}
	t.Logf("Created ClusterRoleBinding '%s' binding ServiceAccount '%s' to ClusterRole '%s'", roleBindingName, serviceAccountName, roleName)

	keyPair, err := sshutils.GenerateKeyPair(sshutils.SSH_KEY_TYPE_ED25519, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate SSH key pair: %w", err)
	}
	if keyPair == nil {
		return nil, nil, fmt.Errorf("generated SSH key pair is nil")
	}

	publicKey, err := keyPair.GetPublicKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public key: %w", err)
	}
	publicKeyString, err := publicKey.GetKeyMaterialAsString()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public key string: %w", err)
	}
	t.Logf("Generated SSH public key: %s", publicKeyString)

	err = namespace.DeletePodByName(ctx, podName)
	if err != nil {
		t.Logf("Pod %s did not exist, continuing: %v", podName, err)
	}

	// linuxserver/openssh-server listens on port 2222 by default
	const sshContainerPort = 2222
	podYaml := fmt.Sprintf(`
apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  serviceAccountName: %s
  containers:
  - name: ssh-server
    image: linuxserver/openssh-server:latest
    ports:
    - containerPort: %d
    env:
    - name: PUBLIC_KEY
      value: "%s"
    - name: USER_NAME
      value: "testuser"
    - name: PASSWORD_ACCESS
      value: "false"
    - name: KUBERNETES_SERVICE_HOST
      value: "kubernetes.default.svc"
    - name: KUBERNETES_SERVICE_PORT
      value: "443"
    volumeMounts:
    - name: kubeconfig
      mountPath: /etc/kubernetes/config
      readOnly: true
  volumes:
  - name: kubeconfig
    configMap:
      name: %s-kubeconfig
`, podName, namespaceName, serviceAccountName, sshContainerPort, publicKeyString, podName)

	// Create ConfigMap with kubeconfig that uses the auto-mounted ServiceAccount token
	// Use the cluster name from the actual cluster object to match the context name expected by tests
	clusterName, err := cluster.GetName()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get cluster name: %w", err)
	}
	kubeconfigConfigMapYaml := fmt.Sprintf(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: %s-kubeconfig
  namespace: %s
data:
  config: |
    apiVersion: v1
    kind: Config
    clusters:
    - name: %s
      cluster:
        server: https://kubernetes.default.svc
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    users:
    - name: serviceaccount
      user:
        tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
    contexts:
    - name: %s
      context:
        cluster: %s
        user: serviceaccount
    current-context: %s
`, podName, namespaceName, clusterName, clusterName, clusterName, clusterName)

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: kubeconfigConfigMapYaml,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create kubeconfig ConfigMap: %w", err)
	}
	t.Logf("Created kubeconfig ConfigMap '%s-kubeconfig' in namespace '%s'", podName, namespaceName)

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: podYaml,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create SSH server pod: %w", err)
	}

	err = namespace.WaitUntilPodReady(ctx, podName, 60*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("SSH server pod did not become ready in time: %w", err)
	}
	t.Log("SSH server pod is ready")

	t.Logf("Waiting 5 seconds for SSH server to fully initialize inside the pod...")
	time.Sleep(5 * time.Second)

	// Set up kubeconfig for testuser by copying from read-only mount to user's home directory
	t.Log("Setting up .kube directory for testuser...")
	pod, err := cluster.GetPodByNames(namespaceName, podName)

	testUserHome, err := pod.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"su", "testuser", "-c", "printenv HOME"},
	})
	testUserHome = strings.TrimRight(strings.TrimSpace(testUserHome), "/")

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get SSH server pod: %w", err)
	}
	_, err = pod.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sh", "-c", "mkdir -p " + testUserHome + "/.kube && cp /etc/kubernetes/config/config " + testUserHome + "/.kube/config && chown -R testuser:testuser " + testUserHome + "/.kube"},
	})
	if err != nil {
		return nil, nil, err
	}

	configInContainer, err := commandexecutorfile.ReadAsString(pod, testUserHome+"/.kube/config")
	if err != nil {
		return nil, nil, err
	}
	require.Contains(t, configInContainer, clusterName)

	// Install kubectl inside the SSH server pod for running kubernetes tests
	t.Log("Installing kubectl in SSH server pod...")

	installOptions := &kubectlutils.InstallKubectlOptions{
		InstallPath: "/usr/local/bin/kubectl",
		UseSudo:     false,
		Version:     "v1.36.2",
	}

	err = kubectlutils.InstallKubectlOnCommandExecutor(ctx, pod, installOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to install kubectl in SSH server pod: %w", err)
	}
	t.Log("kubectl installed successfully in SSH server pod")

	// Verify kubectl is available and executable
	t.Log("Verifying kubectl installation...")
	verifyOutput, err := pod.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"/usr/local/bin/kubectl", "version", "--client"},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify kubectl installation: %w", err)
	}
	verifyStdout, err := verifyOutput.GetStdoutAsString()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get kubectl version output: %w", err)
	}
	t.Logf("kubectl version: %s", verifyStdout)

	localPort, err := netutils.GetNextFreePort(ctx, 22222)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find free port: %w", err)
	}

	t.Logf("Setting up port forwarding: %d -> %s/%s:%d", localPort, namespaceName, podName, sshContainerPort)

	cancel, err := namespace.StartPortForwarding(ctx, podName, localPort, sshContainerPort)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start port forwarding: %w", err)
	}

	t.Logf("Waiting for SSH port %d to be accessible...", localPort)
	err = netutils.WaitTcpPortOpen(ctx, "localhost", localPort, 30*time.Second)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("SSH port did not become accessible: %w", err)
	}

	t.Logf("Verifying SSH port %d is accepting connections...", localPort)
	conn, err := net.DialTimeout("tcp", "localhost:"+strconv.Itoa(localPort), 5*time.Second)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("SSH port is not accepting TCP connections: %w", err)
	}
	conn.Close()
	t.Logf("SSH port %d verified to be accepting connections", localPort)

	result := &SSHTestSetupResult{
		Namespace:     namespace,
		PodName:       podName,
		SSHPort:       sshContainerPort,
		LocalPort:     localPort,
		KeyPair:       keyPair,
		CancelFunc:    cancel,
		NamespaceName: namespaceName,
	}

	cleanupFunc := func() {
		t.Log("Cleaning up SSH test environment...")
		if cancel != nil {
			cancel()
		}
		err := namespace.DeletePodByName(ctx, podName)
		if err != nil {
			t.Logf("Warning: failed to delete pod: %v", err)
		}

		// Delete kubeconfig ConfigMap
		err = namespace.DeleteConfigMapByName(ctx, podName+"-kubeconfig")
		if err != nil {
			t.Logf("Warning: failed to delete ConfigMap '%s-kubeconfig': %v", podName, err)
		} else {
			t.Logf("Deleted ConfigMap '%s-kubeconfig'", podName)
		}

		// Delete ClusterRoleBinding and ClusterRole using kubectl directly (cluster-scoped resources)
		kubectlContext, err := cluster.GetKubectlContext(ctx)
		if err != nil {
			t.Logf("Warning: failed to get kubectl context: %v", err)
		} else {
			// Delete ClusterRoleBinding
			cmd := exec.CommandContext(ctx, "kubectl", "delete", "clusterrolebinding", roleBindingName, "--context", kubectlContext, "--ignore-not-found=true")
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Warning: failed to delete ClusterRoleBinding '%s': %v, output: %s", roleBindingName, err, string(output))
			} else {
				t.Logf("Deleted ClusterRoleBinding '%s'", roleBindingName)
			}

			// Delete ClusterRole
			cmd = exec.CommandContext(ctx, "kubectl", "delete", "clusterrole", roleName, "--context", kubectlContext, "--ignore-not-found=true")
			output, err = cmd.CombinedOutput()
			if err != nil {
				t.Logf("Warning: failed to delete ClusterRole '%s': %v, output: %s", roleName, err, string(output))
			} else {
				t.Logf("Deleted ClusterRole '%s'", roleName)
			}
		}

		t.Log("SSH test environment cleanup completed")
	}

	return result, cleanupFunc, nil
}

func TestSetupSSHServerInKind(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	const namespaceName = "ssh-test-setup"
	const podName = "ssh-server-setup"

	_, cleanup, err := SetupSSHServerInKind(ctx, t, cluster, namespaceName, podName)
	require.NoError(t, err)
	defer cleanup()

	pod, err := cluster.GetPodByNames(namespaceName, podName)
	require.NoError(t, err)

	output, err := pod.RunCommand(ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"kubectl", "get", "ns"},
		},
	)
	require.NoError(t, err, "kubectl command should execute successfully")

	stdout, err := output.GetStdoutAsString()
	require.NoError(t, err)

	stderr, err := output.GetStderrAsString()
	require.NoError(t, err)

	exitCode, err := output.GetReturnCode()
	require.NoError(t, err)

	// Verify the command succeeded
	require.Equal(t, 0, exitCode, "kubectl get ns should exit with code 0")
	require.Empty(t, stderr, "kubectl get ns should not produce stderr output")
	require.Contains(t, stdout, "NAME", "kubectl get ns output should contain header with NAME column")
	require.Contains(t, stdout, namespaceName, "kubectl get ns output should contain the test namespace")

	t.Logf("Successfully validated SSH server setup:")
	t.Logf("  - kubectl is available in the default PATH")
	t.Logf("  - Service account is properly configured")
	t.Logf("  - kubectl get ns output:\n%s", stdout)
}
