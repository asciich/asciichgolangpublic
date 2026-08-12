package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestRole_CreateAndDeleteRole(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const namespaceName = "testnamespace"
				const roleName = "testrole"

				kubernetes := getKubernetesByImplementationName(ctx, t, tt.implementationName)
				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				mustutils.Must0(namespace.DeleteRoleByName(ctx, roleName))
				require.False(t, mustutils.Must(namespace.RoleByNameExists(ctx, roleName)))

				for i := 0; i < 2; i++ {
					_, err := namespace.CreateRole(
						ctx,
						&kubernetesparameteroptions.CreateRoleOptions{
							Name:      roleName,
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resorces:  []string{"pod"},
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(namespace.RoleByNameExists(ctx, roleName)))
				}

				for i := 0; i < 2; i++ {
					err := namespace.DeleteRoleByName(ctx, roleName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(namespace.RoleByNameExists(ctx, roleName)))
				}
			},
		)
	}
}

func TestRole_CreateAndDeleteClusterRole(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const roleName = "testclusterrole"

				kubernetes := getKubernetesByImplementationName(ctx, t, tt.implementationName)

				mustutils.Must0(kubernetes.DeleteClusterRoleByName(ctx, roleName))
				require.False(t, mustutils.Must(kubernetes.ClusterRoleByNameExists(ctx, roleName)))

				for i := 0; i < 2; i++ {
					_, err := kubernetes.CreateClusterRole(
						ctx,
						&kubernetesparameteroptions.CreateClusterRoleOptions{
							Name:      roleName,
							Verbs:     []string{"get", "list", "watch"},
							APIGroups: []string{""},
							Resorces:  []string{"pods", "services"},
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(kubernetes.ClusterRoleByNameExists(ctx, roleName)))
				}

				for i := 0; i < 2; i++ {
					err := kubernetes.DeleteClusterRoleByName(ctx, roleName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(kubernetes.ClusterRoleByNameExists(ctx, roleName)))
				}
			},
		)
	}
}

func TestRole_CreateAndDeleteRoleBinding(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const namespaceName = "testnamespace"
				const roleName = "testrole-for-binding"
				const bindingName = "testrolebinding"

				kubernetes := getKubernetesByImplementationName(ctx, t, tt.implementationName)
				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Create a Role first
				mustutils.Must0(namespace.DeleteRoleByName(ctx, roleName))
				_, err = namespace.CreateRole(
					ctx,
					&kubernetesparameteroptions.CreateRoleOptions{
						Name:      roleName,
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resorces:  []string{"pods"},
					},
				)
				require.NoError(t, err)

				mustutils.Must0(namespace.DeleteRoleBindingByName(ctx, bindingName))
				require.False(t, mustutils.Must(namespace.RoleBindingByNameExists(ctx, bindingName)))

				for i := 0; i < 2; i++ {
					_, err := namespace.CreateRoleBinding(
						ctx,
						&kubernetesparameteroptions.CreateRoleBindingOptions{
							Name:        bindingName,
							RoleRef:     roleName,
							Subjects:    []string{"default"},
							SubjectKind: "ServiceAccount",
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(namespace.RoleBindingByNameExists(ctx, bindingName)))
				}

				for i := 0; i < 2; i++ {
					err := namespace.DeleteRoleBindingByName(ctx, bindingName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(namespace.RoleBindingByNameExists(ctx, bindingName)))
				}

				// Cleanup
				mustutils.Must0(namespace.DeleteRoleByName(ctx, roleName))
			},
		)
	}
}

func TestRole_CreateAndDeleteClusterRoleBinding(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const roleName = "testclusterrole-for-binding"
				const bindingName = "testclusterrolebinding"

				kubernetes := getKubernetesByImplementationName(ctx, t, tt.implementationName)

				// Create a ClusterRole first
				mustutils.Must0(kubernetes.DeleteClusterRoleByName(ctx, roleName))
				_, err := kubernetes.CreateClusterRole(
					ctx,
					&kubernetesparameteroptions.CreateClusterRoleOptions{
						Name:      roleName,
						Verbs:     []string{"get", "list", "watch"},
						APIGroups: []string{""},
						Resorces:  []string{"pods", "services", "deployments"},
					},
				)
				require.NoError(t, err)

				mustutils.Must0(kubernetes.DeleteClusterRoleBindingByName(ctx, bindingName))
				require.False(t, mustutils.Must(kubernetes.ClusterRoleBindingByNameExists(ctx, bindingName)))

				for i := 0; i < 2; i++ {
					_, err := kubernetes.CreateClusterRoleBinding(
						ctx,
						&kubernetesparameteroptions.CreateClusterRoleBindingOptions{
							Name:             bindingName,
							RoleRef:          roleName,
							Subjects:         []string{"default"},
							SubjectKind:      "ServiceAccount",
							SubjectNamespace: "default",
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(kubernetes.ClusterRoleBindingByNameExists(ctx, bindingName)))
				}

				for i := 0; i < 2; i++ {
					err := kubernetes.DeleteClusterRoleBindingByName(ctx, bindingName)
					require.NoError(t, err)
					require.False(t, mustutils.Must(kubernetes.ClusterRoleBindingByNameExists(ctx, bindingName)))
				}

				// Cleanup
				mustutils.Must0(kubernetes.DeleteClusterRoleByName(ctx, roleName))
			},
		)
	}
}
