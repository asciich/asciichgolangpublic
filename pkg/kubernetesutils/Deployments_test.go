package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_ListDeploymentNames(t *testing.T) {
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

				deploymentNames := []string{"listdeployment-1", "listdeployment-2", "listdeployment-3"}

				// Ensure all test deployments are absent before starting
				for _, name := range deploymentNames {
					err := kubernetes.DeleteDeploymentByNames(ctx, namespaceName, name)
					require.NoError(t, err)
				}

				// List deployments in empty namespace
				names, err := namespace.ListDeploymentNames(ctx)
				require.NoError(t, err)
				for _, name := range deploymentNames {
					require.NotContains(t, names, name)
				}

				for i, deploymentName := range deploymentNames {
					_, err := namespace.CreateDeployment(
						ctx,
						&kubernetesparameteroptions.RunCommandOptions{
							Image:                           "ubuntu",
							DeploymentName:                  deploymentName,
							DeleteAlreadyExistingDeployment: true,
							Command:                         []string{"bash", "-c", "echo hello_world"},
						},
					)
					require.NoError(t, err)

					names, err = namespace.ListDeploymentNames(ctx)
					require.NoError(t, err)

					for _, created := range deploymentNames[:i+1] {
						require.Contains(t, names, created)
					}
					for _, notYetCreated := range deploymentNames[i+1:] {
						require.NotContains(t, names, notYetCreated)
					}
				}

				// Delete deployments one by one and verify list shrinks
				for i, name := range deploymentNames {
					deployment, err := kubernetes.GetDeploymentByNames(namespaceName, name)
					require.NoError(t, err)

					err = deployment.Delete(ctx)
					require.NoError(t, err)

					names, err = namespace.ListDeploymentNames(ctx)
					require.NoError(t, err)

					for _, deleted := range deploymentNames[:i+1] {
						require.NotContains(t, names, deleted)
					}
					for _, stillPresent := range deploymentNames[i+1:] {
						require.Contains(t, names, stillPresent)
					}
				}
			},
		)
	}
}

func Test_CreateAndDeleteDeployment(t *testing.T) {
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
				const deploymentName = "testdeployment"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				_, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = kubernetes.DeleteDeploymentByNames(ctx, namespaceName, deploymentName)
				require.NoError(t, err)

				deployment, err := kubernetes.GetDeploymentByNames(namespaceName, deploymentName)
				require.NoError(t, err)

				exists, err := deployment.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)

				for range 3 {
					_, err := kubernetes.CreateDeployment(
						ctx,
						namespaceName,
						&kubernetesparameteroptions.RunCommandOptions{
							Image:                           "ubuntu",
							DeploymentName:                  deploymentName,
							DeleteAlreadyExistingDeployment: true,
							Command:                         []string{"bash", "-c", "echo hello_world"},
						},
					)
					require.NoError(t, err)

					exists, err = deployment.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)

					err = deployment.Delete(ctx)
					require.NoError(t, err)

					exists, err = deployment.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}
