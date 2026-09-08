package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_ListDaemonSetNames(t *testing.T) {
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

				daemonSetNames := []string{"listdaemonset-1", "listdaemonset-2", "listdaemonset-3"}

				// Ensure all test daemonsets are absent before starting
				for _, name := range daemonSetNames {
					err := kubernetes.DeleteDaemonSetByNames(ctx, namespaceName, name)
					require.NoError(t, err)
				}

				// List daemonsets in empty namespace
				names, err := namespace.ListDaemonSetNames(ctx)
				require.NoError(t, err)
				for _, name := range daemonSetNames {
					require.NotContains(t, names, name)
				}

				for i, daemonSetName := range daemonSetNames {
					_, err := namespace.CreateDaemonSet(
						ctx,
						&kubernetesparameteroptions.KubernetesRunCommandOptions{
							Image:                          "ubuntu",
							DaemonSetName:                  daemonSetName,
							DeleteAlreadyExistingDaemonSet: true,
							RunCommandOptions: &parameteroptions.RunCommandOptions{
								Command: []string{"bash", "-c", "echo hello_world"},
							},
						},
					)
					require.NoError(t, err)

					names, err = namespace.ListDaemonSetNames(ctx)
					require.NoError(t, err)

					for _, created := range daemonSetNames[:i+1] {
						require.Contains(t, names, created)
					}
					for _, notYetCreated := range daemonSetNames[i+1:] {
						require.NotContains(t, names, notYetCreated)
					}
				}

				// Delete daemonsets one by one and verify list shrinks
				for i, name := range daemonSetNames {
					daemonSet, err := kubernetes.GetDaemonSetByNames(namespaceName, name)
					require.NoError(t, err)

					err = daemonSet.Delete(ctx)
					require.NoError(t, err)

					names, err = namespace.ListDaemonSetNames(ctx)
					require.NoError(t, err)

					for _, deleted := range daemonSetNames[:i+1] {
						require.NotContains(t, names, deleted)
					}
					for _, stillPresent := range daemonSetNames[i+1:] {
						require.Contains(t, names, stillPresent)
					}
				}
			},
		)
	}
}

func Test_CreateAndDeleteDaemonSet(t *testing.T) {
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
				const daemonSetName = "testdaemonset"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = kubernetes.DeleteDaemonSetByNames(ctx, namespaceName, daemonSetName)
				require.NoError(t, err)

				daemonSet, err := kubernetes.GetDaemonSetByNames(namespaceName, daemonSetName)
				require.NoError(t, err)

				exists, err := daemonSet.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				for range 3 {
					_, err := kubernetes.CreateDaemonSet(
						ctx,
						namespaceName,
						&kubernetesparameteroptions.KubernetesRunCommandOptions{
							Image:                          "ubuntu",
							DaemonSetName:                  daemonSetName,
							DeleteAlreadyExistingDaemonSet: true,
							RunCommandOptions: &parameteroptions.RunCommandOptions{
								Command: []string{"bash", "-c", "echo hello_world"},
							},
						},
					)
					require.NoError(t, err)

					exists, err = daemonSet.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					err = daemonSet.Delete(ctx)
					require.NoError(t, err)

					exists, err = daemonSet.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
