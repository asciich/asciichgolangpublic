package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
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

				kubernetes := getKubernetesByImplementationName(ctx, tt.implementationName)
				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				mustutils.Must0(namespace.DeleteRoleByName(ctx, roleName))
				require.False(t, mustutils.Must(namespace.RoleByNameExists(ctx, roleName)))

				for i := 0; i < 2; i++ {
					_, err := namespace.CreateRole(
						ctx,
						&kubernetesparameteroptions.CreateRoleOptions{
							Name:     roleName,
							Verbs:    []string{"get"},
							Resorces: []string{"pod"},
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
