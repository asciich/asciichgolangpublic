package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				require := require.New(t)

				const verbose bool = true
				const namespaceName = "testnamespace"
				const roleName = "testrole"

				kubernetes := getKubernetesByImplementationName(tt.implementationName)
				namespace := kubernetes.MustCreateNamespaceByName(namespaceName, verbose)

				namespace.MustDeleteRoleByName(roleName, verbose)
				require.False(namespace.MustRoleByNameExists(roleName, verbose))

				for i := 0; i < 2; i++ {
					namespace.MustCreateRole(
						&CreateRoleOptions{
							Name:     roleName,
							Verbs:    []string{"get"},
							Resorces: []string{"pod"},
						},
					)
					require.True(namespace.MustRoleByNameExists(roleName, verbose))
				}

				for i := 0; i < 2; i++ {
					namespace.MustDeleteRoleByName(roleName, verbose)
					require.False(namespace.MustRoleByNameExists(roleName, verbose))
				}
			},
		)
	}
}
