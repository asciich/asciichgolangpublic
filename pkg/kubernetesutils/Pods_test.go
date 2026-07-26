package kubernetesutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_PodsRunSingleCommand_echoHelloWorld(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const podName = "podname"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Wait until default sa in created namespace exists.
				time.Sleep(10 * time.Second)

				output, err := kubernetes.RunCommandInTemporaryPod(
					ctx,
					namespaceName,
					&kubernetesparameteroptions.RunCommandOptions{
						Image:                    "ubuntu",
						PodName:                  podName,
						DeleteAlreadyExistingPod: true,
						Command:                  []string{"bash", "-c", "echo hello_world"},
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, "hello_world\n", stdout)

				stderr, err := output.GetStderrAsString()
				require.NoError(t, err)
				require.EqualValues(t, "", stderr)

				retVal, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, 0, retVal)
			},
		)
	}
}

func Test_ListPodNames(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				namespace, err := kubernetes.GetNamespaceByName(namespaceName)
				require.NoError(t, err)

				podNames := []string{"listpod-1", "listpod-2", "listpod-3"}

				// Ensure all test pods are absent before starting
				for _, name := range podNames {
					err := kubernetes.DeletePodByNames(ctx, namespaceName, name)
					require.NoError(t, err)
				}

				// List pods in empty namespace
				names, err := namespace.ListPodNames(ctx)
				require.NoError(t, err)
				for _, name := range podNames {
					require.NotContains(t, names, name)
				}

				for i, podName := range podNames {
					_, err := namespace.CreatePod(
						ctx,
						&kubernetesparameteroptions.RunCommandOptions{
							Image:                    "ubuntu",
							PodName:                  podName,
							DeleteAlreadyExistingPod: true,
							Command:                  []string{"bash", "-c", "echo hello_world"},
						},
					)
					require.NoError(t, err)

					names, err = namespace.ListPodNames(ctx)
					require.NoError(t, err)

					for _, created := range podNames[:i+1] {
						require.Contains(t, names, created)
					}
					for _, notYetCreated := range podNames[i+1:] {
						require.NotContains(t, names, notYetCreated)
					}
				}

				// Delete pods one by one and verify list shrinks
				for i, name := range podNames {
					pod, err := kubernetes.GetPodByNames(namespaceName, name)
					require.NoError(t, err)

					err = pod.Delete(ctx)
					require.NoError(t, err)

					names, err = namespace.ListPodNames(ctx)
					require.NoError(t, err)

					for _, deleted := range podNames[:i+1] {
						require.NotContains(t, names, deleted)
					}
					for _, stillPresent := range podNames[i+1:] {
						require.Contains(t, names, stillPresent)
					}
				}
			},
		)
	}
}

func Test_CreateAndDeletePod(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const podName = "testpod"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = kubernetes.DeletePodByNames(ctx, namespaceName, podName)
				require.NoError(t, err)

				pod, err := kubernetes.GetPodByNames(namespaceName, podName)
				require.NoError(t, err)

				exists, err := pod.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				for range 3 {
					_, err := kubernetes.CreatePod(
						ctx,
						namespaceName,
						&kubernetesparameteroptions.RunCommandOptions{
							Image:                    "ubuntu",
							PodName:                  podName,
							DeleteAlreadyExistingPod: true,
							Command:                  []string{"bash", "-c", "echo hello_world"},
						},
					)
					require.NoError(t, err)

					exists, err = pod.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					err = pod.Delete(ctx)
					require.NoError(t, err)

					exists, err = pod.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}

func Test_RunCommandInTemporaryPod_WithSecretAsEnvVar(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const podName = "podname-secret-env"
				const secretName = "test-secret"
				const secretKey = "mykey"
				const secretValue = "my-secret-value"
				const envVarName = "MY_SECRET_VAR"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				// Create namespace
				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Wait until default sa in created namespace exists.
				time.Sleep(10 * time.Second)

				// Create secret with test data
				_, err = kubernetes.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{
						secretKey: []byte(secretValue),
					},
				})
				require.NoError(t, err)

				// Ensure secret is cleaned up after test
				defer func() {
					_ = kubernetes.DeleteSecretByName(ctx, namespaceName, secretName)
				}()

				// Run command in temporary pod with secret as environment variable
				output, err := kubernetes.RunCommandInTemporaryPod(
					ctx,
					namespaceName,
					&kubernetesparameteroptions.RunCommandOptions{
						Image:                    "ubuntu",
						PodName:                  podName,
						DeleteAlreadyExistingPod: true,
						Command:                  []string{"bash", "-c", "printenv " + envVarName},
						SecretEnvVars: map[string]kubernetesparameteroptions.SecretEnvVarSource{
							envVarName: {
								SecretName: secretName,
								SecretKey:  secretKey,
							},
						},
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, secretValue+"\n", stdout)

				stderr, err := output.GetStderrAsString()
				require.NoError(t, err)
				require.EqualValues(t, "", stderr)

				retVal, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, 0, retVal)
			},
		)
	}
}

func Test_RunCommandInTemporaryPod_WithSecretAsFile(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const podName = "podname-secret-file"
				const secretName = "test-secret-file"
				const secretKey = "mykey"
				const secretValue = "my-secret-file-value"
				const mountPath = "/etc/secret"
				const secretFilePath = mountPath + "/" + secretKey

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				// Create namespace
				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Wait until default sa in created namespace exists.
				time.Sleep(10 * time.Second)

				// Create secret with test data
				_, err = kubernetes.CreateSecret(ctx, namespaceName, secretName, &kubernetesparameteroptions.CreateSecretOptions{
					SecretData: map[string][]byte{
						secretKey: []byte(secretValue),
					},
				})
				require.NoError(t, err)

				// Ensure secret is cleaned up after test
				defer func() {
					_ = kubernetes.DeleteSecretByName(ctx, namespaceName, secretName)
				}()

				// Run command in temporary pod with secret mounted as file
				output, err := kubernetes.RunCommandInTemporaryPod(
					ctx,
					namespaceName,
					&kubernetesparameteroptions.RunCommandOptions{
						Image:                    "ubuntu",
						PodName:                  podName,
						DeleteAlreadyExistingPod: true,
						Command:                  []string{"bash", "-c", "cat " + secretFilePath},
						SecretMounts: map[string]kubernetesparameteroptions.SecretMountSource{
							mountPath: {
								SecretName: secretName,
							},
						},
					},
				)
				require.NoError(t, err)

				stdout, err := output.GetStdoutAsString()
				require.NoError(t, err)
				require.EqualValues(t, secretValue, stdout)

				stderr, err := output.GetStderrAsString()
				require.NoError(t, err)
				require.EqualValues(t, "", stderr)

				retVal, err := output.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, 0, retVal)
			},
		)
	}
}
