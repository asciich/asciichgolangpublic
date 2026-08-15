package kubernetestestsshserver

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// StartTestSshServerInCluster deploys a SSH server as a pod inside the Kubernetes cluster for testing purposes.
// It returns a Pod interface that represents the SSH server pod.
// The function waits until the SSH server pod is running and ready before returning.
func StartTestSshServerInCluster(ctx context.Context, kubernetesCluster kubernetesinterfaces.KubernetesCluster, options *StartTestSshServerOptions) (kubernetesinterfaces.Pod, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if options.KubernetesNamespace == "" {
		return nil, tracederrors.TracedErrorEmptyString("KubernetesNamespace")
	}

	if options.PodName == "" {
		return nil, tracederrors.TracedErrorEmptyString("PodName")
	}

	if options.SSHUsername == "" {
		return nil, tracederrors.TracedErrorEmptyString("SSHUsername")
	}

	// Set default image if not provided
	image := options.Image
	if image == "" {
		image = "linuxserver/openssh-server:latest"
	}

	// Generate random password if not provided
	sshPassword := options.SSHPassword
	if sshPassword == "" {
		sshPassword = generateRandomPassword(16)
		logging.LogInfoByCtxf(ctx, "Generated random SSH password for test server")
	}

	// Get or create the namespace
	namespace, err := kubernetesCluster.GetNamespaceByName(options.KubernetesNamespace)
	if err != nil {
		return nil, err
	}

	// Check if namespace exists, create if not
	exists, err := namespace.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = namespace.Create(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Create a secret with SSH credentials if using password auth or public key
	if options.SSHPassword != "" || options.SSHPublicKey != "" {
		secretName := options.PodName + "-ssh-secret"
		secretData := make(map[string][]byte)

		if options.SSHPassword != "" {
			secretData["USER_PASSWORD"] = []byte(options.SSHPassword)
			secretData["PASSWORD_ACCESS"] = []byte("true")
		}

		if options.SSHPublicKey != "" {
			secretData["PUBLIC_KEY"] = []byte(options.SSHPublicKey)
		}

		// Set required environment variables for the SSH server container
		secretData["USER_NAME"] = []byte(options.SSHUsername)
		secretData["PUID"] = []byte("1000")
		secretData["PGID"] = []byte("1000")

		// Create the secret
		secretOptions := &kubernetesparameteroptions.CreateSecretOptions{
			SecretData: secretData,
		}

		_, err = namespace.CreateSecret(ctx, secretName, secretOptions)
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Created SSH secret '%s' in namespace '%s'", secretName, options.KubernetesNamespace)
	}

	// Create the SSH server pod
	podOptions := &kubernetesparameteroptions.KubernetesRunCommandOptions{
		PodName:           options.PodName,
		Image:             image,
		ContainerName:     options.PodName,
		WaitForPodRunning: true,
		RunCommandOptions: &parameteroptions.RunCommandOptions{
			Command: []string{"/init"},
		},
		DeleteAlreadyExistingPod: true,
	}

	// Add environment variables from secret
	secretName := options.PodName + "-ssh-secret"
	podOptions.SecretEnvVars = map[string]kubernetesparameteroptions.SecretEnvVarSource{
		"USER_NAME": {
			SecretName: secretName,
			SecretKey:  "USER_NAME",
		},
		"PUID": {
			SecretName: secretName,
			SecretKey:  "PUID",
		},
		"PGID": {
			SecretName: secretName,
			SecretKey:  "PGID",
		},
	}

	if options.SSHPassword != "" {
		podOptions.SecretEnvVars["USER_PASSWORD"] = kubernetesparameteroptions.SecretEnvVarSource{
			SecretName: secretName,
			SecretKey:  "USER_PASSWORD",
		}
		podOptions.SecretEnvVars["PASSWORD_ACCESS"] = kubernetesparameteroptions.SecretEnvVarSource{
			SecretName: secretName,
			SecretKey:  "PASSWORD_ACCESS",
		}
	}

	if options.SSHPublicKey != "" {
		podOptions.SecretEnvVars["PUBLIC_KEY"] = kubernetesparameteroptions.SecretEnvVarSource{
			SecretName: secretName,
			SecretKey:  "PUBLIC_KEY",
		}
	}

	// Create the pod
	pod, err := namespace.CreatePod(ctx, podOptions)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "SSH test server pod '%s' created and running in namespace '%s'", options.PodName, options.KubernetesNamespace)

	// If public key auth is used, manually configure authorized_keys to ensure it works.
	// The linuxserver/openssh-server init scripts may not always handle the PUBLIC_KEY env var reliably.
	if options.SSHPublicKey != "" {
		pubKeyBase64 := base64.StdEncoding.EncodeToString([]byte(options.SSHPublicKey))
		setupCommand := []string{"sh", "-c",
			"mkdir -p /config/.ssh && " +
				"echo '" + pubKeyBase64 + "' | base64 -d > /config/.ssh/authorized_keys && " +
				"chmod 600 /config/.ssh/authorized_keys && " +
				"chmod 700 /config/.ssh && " +
				"chown -R 1000:1000 /config/.ssh"}

		_, err = pod.RunCommandInContainer(ctx, &kubernetesparameteroptions.KubernetesRunCommandOptions{
			ContainerName: options.PodName,
			RunCommandOptions: &parameteroptions.RunCommandOptions{
				Command: setupCommand,
			},
		})
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to configure authorized_keys on SSH server: %w", err)
		}

		logging.LogInfoByCtxf(ctx, "Configured authorized_keys for user '%s' on SSH server pod '%s'", options.SSHUsername, options.PodName)
	}

	// Create a Service to make the SSH server pod reachable via DNS
	// The selector "run: <podName>" matches the label automatically added by "kubectl run"
	serviceYaml := fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: %s
  namespace: %s
spec:
  selector:
    run: %s
  ports:
  - protocol: TCP
    port: 22
    targetPort: 2222
`, options.PodName, options.KubernetesNamespace, options.PodName)

	_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: serviceYaml,
	})
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Created Service '%s' in namespace '%s' for SSH server pod", options.PodName, options.KubernetesNamespace)

	return pod, nil
}

// generateRandomPassword generates a random password of specified length
func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to a simple password if random generation fails
		return "testpassword123"
	}
	return base64.StdEncoding.EncodeToString(bytes)[:length]
}
