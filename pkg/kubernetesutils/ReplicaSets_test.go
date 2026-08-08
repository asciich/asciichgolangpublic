package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_ListReplicaSetNames(t *testing.T) {
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

				replicaSetNames := []string{"listreplicaset-1", "listreplicaset-2", "listreplicaset-3"}

				// Ensure all test replicasets are absent before starting
				for _, name := range replicaSetNames {
					err := kubernetes.DeleteReplicaSetByNames(ctx, namespaceName, name)
					require.NoError(t, err)
				}

				// List replicasets in empty namespace
				names, err := namespace.ListReplicaSetNames(ctx)
				require.NoError(t, err)
				for _, name := range replicaSetNames {
					require.NotContains(t, names, name)
				}

				for i, replicaSetName := range replicaSetNames {
					_, err := namespace.CreateReplicaSet(
						ctx,
						&kubernetesparameteroptions.KubernetesRunCommandOptions{
							Image:                           "ubuntu",
							ReplicaSetName:                  replicaSetName,
							DeleteAlreadyExistingReplicaSet: true,
							RunCommandOptions: &parameteroptions.RunCommandOptions{
								Command: []string{"bash", "-c", "echo hello_world"},
							},
						},
					)
					require.NoError(t, err)

					names, err = namespace.ListReplicaSetNames(ctx)
					require.NoError(t, err)

					for _, created := range replicaSetNames[:i+1] {
						require.Contains(t, names, created)
					}
					for _, notYetCreated := range replicaSetNames[i+1:] {
						require.NotContains(t, names, notYetCreated)
					}
				}

				// Delete replicasets one by one and verify list shrinks
				for i, name := range replicaSetNames {
					replicaSet, err := kubernetes.GetReplicaSetByNames(namespaceName, name)
					require.NoError(t, err)

					err = replicaSet.Delete(ctx)
					require.NoError(t, err)

					names, err = namespace.ListReplicaSetNames(ctx)
					require.NoError(t, err)

					for _, deleted := range replicaSetNames[:i+1] {
						require.NotContains(t, names, deleted)
					}
					for _, stillPresent := range replicaSetNames[i+1:] {
						require.Contains(t, names, stillPresent)
					}
				}
			},
		)
	}
}

func Test_CreateAndDeleteReplicaSet(t *testing.T) {
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
				const replicaSetName = "testreplicaset"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = kubernetes.DeleteReplicaSetByNames(ctx, namespaceName, replicaSetName)
				require.NoError(t, err)

				replicaSet, err := kubernetes.GetReplicaSetByNames(namespaceName, replicaSetName)
				require.NoError(t, err)

				exists, err := replicaSet.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				for range 3 {
					_, err := kubernetes.CreateReplicaSet(
						ctx,
						namespaceName,
						&kubernetesparameteroptions.KubernetesRunCommandOptions{
							Image:                           "ubuntu",
							ReplicaSetName:                  replicaSetName,
							DeleteAlreadyExistingReplicaSet: true,
							RunCommandOptions: &parameteroptions.RunCommandOptions{
								Command: []string{"bash", "-c", "echo hello_world"},
							},
						},
					)
					require.NoError(t, err)

					exists, err = replicaSet.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					err = replicaSet.Delete(ctx)
					require.NoError(t, err)

					exists, err = replicaSet.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
