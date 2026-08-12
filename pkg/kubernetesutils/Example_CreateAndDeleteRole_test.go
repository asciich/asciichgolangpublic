package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeleteRole(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the namespace and role name we use for testing
	const namespaceName = "testnamespace"
	const roleName = "example-role"

	// Get the namespace (will be created implicitly when creating the role)
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the role object
	role, err := namespace.GetRoleByName(roleName)
	require.NoError(t, err)

	// Ensure the role is absent
	err = role.Delete(ctx)
	require.NoError(t, err)
	exists, err := role.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example role
	_, err = namespace.CreateRole(ctx, &kubernetesparameteroptions.CreateRoleOptions{
		Name:      roleName,
		Verbs:     []string{"get", "list", "watch"},
		Resorces:  []string{"pods"},
		APIGroups: []string{""},
	})
	require.NoError(t, err)
	exists, err = role.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the role exists via the cluster
	names, err := namespace.ListRoleNames(ctx)
	require.NoError(t, err)
	require.Contains(t, names, roleName)

	// Delete the role
	err = role.Delete(ctx)
	require.NoError(t, err)
	exists, err = role.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)
}
