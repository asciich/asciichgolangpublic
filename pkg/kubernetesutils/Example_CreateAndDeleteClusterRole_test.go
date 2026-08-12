package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeleteClusterRole(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the ClusterRole name
	const roleName = "example-clusterrole"

	// Get the ClusterRole object
	role, err := cluster.GetClusterRoleByName(roleName)
	require.NoError(t, err)

	// Ensure the ClusterRole is absent
	err = role.Delete(ctx)
	require.NoError(t, err)
	exists, err := role.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example ClusterRole
	createOptions := &kubernetesparameteroptions.CreateClusterRoleOptions{}
	err = createOptions.SetName(roleName)
	require.NoError(t, err)

	err = createOptions.SetVerbs([]string{"get", "list", "watch"})

	err = createOptions.SetAPIGroups([]string{""})
	require.NoError(t, err)

	err = createOptions.SetResorces([]string{"pods", "services"})
	require.NoError(t, err)

	_, err = cluster.CreateClusterRole(ctx, createOptions)
	require.NoError(t, err)

	// Verify the ClusterRole exists
	exists, err = role.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the ClusterRole can be retrieved
	retrievedRole, err := cluster.GetClusterRoleByName(roleName)
	require.NoError(t, err)
	require.NotNil(t, retrievedRole)

	// List ClusterRoles to verify creation
	roleList, err := cluster.ListClusterRoleNames(ctx)
	require.NoError(t, err)
	require.Contains(t, roleList, roleName)

	// Delete the ClusterRole
	err = role.Delete(ctx)
	require.NoError(t, err)
	exists, err = role.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Verify deletion via cluster
	exists, err = cluster.ClusterRoleByNameExists(ctx, roleName)
	require.NoError(t, err)
	require.False(t, exists)
}

func Test_Example_CreateAndDeleteRoleBinding(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the namespace and RoleBinding name
	const namespaceName = "testnamespace"
	const bindingName = "example-rolebinding"

	// Get the namespace
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the RoleBinding object
	binding, err := namespace.GetRoleBindingByName(bindingName)
	require.NoError(t, err)

	// Ensure the RoleBinding is absent
	err = binding.Delete(ctx)
	require.NoError(t, err)
	exists, err := binding.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create a Role first (prerequisite for RoleBinding)
	roleName := "example-role"
	role, err := namespace.GetRoleByName(roleName)
	require.NoError(t, err)

	createRoleOptions := &kubernetesparameteroptions.CreateRoleOptions{}
	err = createRoleOptions.SetName(roleName)
	require.NoError(t, err)
	err = createRoleOptions.SetVerbs([]string{"get", "list"})

	err = createRoleOptions.SetAPIGroups([]string{""})

	err = createRoleOptions.SetAPIGroups([]string{""})
	require.NoError(t, err)
	err = createRoleOptions.SetResorces([]string{"pods"})
	require.NoError(t, err)

	_, err = namespace.CreateRole(ctx, createRoleOptions)
	require.NoError(t, err)

	// Create an example RoleBinding
	createBindingOptions := &kubernetesparameteroptions.CreateRoleBindingOptions{}
	err = createBindingOptions.SetName(bindingName)
	require.NoError(t, err)

	err = createBindingOptions.SetRoleRef(roleName)
	require.NoError(t, err)

	err = createBindingOptions.SetSubjects([]string{"default"})
	require.NoError(t, err)

	err = createBindingOptions.SetSubjectKind("ServiceAccount")
	require.NoError(t, err)

	_, err = namespace.CreateRoleBinding(ctx, createBindingOptions)
	require.NoError(t, err)

	// Verify the RoleBinding exists
	exists, err = binding.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the RoleBinding can be retrieved
	retrievedBinding, err := namespace.GetRoleBindingByName(bindingName)
	require.NoError(t, err)
	require.NotNil(t, retrievedBinding)

	// List RoleBindings to verify creation
	bindingList, err := namespace.ListRoleBindingNames(ctx)
	require.NoError(t, err)
	require.Contains(t, bindingList, bindingName)

	// Delete the RoleBinding
	err = binding.Delete(ctx)
	require.NoError(t, err)
	exists, err = binding.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Cleanup: delete the Role
	err = role.Delete(ctx)
	require.NoError(t, err)
}

func Test_Example_CreateAndDeleteClusterRoleBinding(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the ClusterRoleBinding name
	const bindingName = "example-clusterrolebinding"

	// Get the ClusterRoleBinding object
	binding, err := cluster.GetClusterRoleBindingByName(bindingName)
	require.NoError(t, err)

	// Ensure the ClusterRoleBinding is absent
	err = binding.Delete(ctx)
	require.NoError(t, err)
	exists, err := binding.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create a ClusterRole first (prerequisite for ClusterRoleBinding)
	roleName := "example-clusterrole-for-binding"
	role, err := cluster.GetClusterRoleByName(roleName)
	require.NoError(t, err)

	createRoleOptions := &kubernetesparameteroptions.CreateClusterRoleOptions{}
	err = createRoleOptions.SetName(roleName)
	require.NoError(t, err)
	err = createRoleOptions.SetVerbs([]string{"get", "list", "watch"})

	err = createRoleOptions.SetAPIGroups([]string{""})

	err = createRoleOptions.SetAPIGroups([]string{""})
	require.NoError(t, err)
	err = createRoleOptions.SetResorces([]string{"pods", "services", "deployments"})
	require.NoError(t, err)

	_, err = cluster.CreateClusterRole(ctx, createRoleOptions)
	require.NoError(t, err)

	// Create an example ClusterRoleBinding
	createBindingOptions := &kubernetesparameteroptions.CreateClusterRoleBindingOptions{}
	err = createBindingOptions.SetName(bindingName)
	require.NoError(t, err)

	err = createBindingOptions.SetRoleRef(roleName)
	require.NoError(t, err)

	err = createBindingOptions.SetSubjects([]string{"default"})
	require.NoError(t, err)

	err = createBindingOptions.SetSubjectKind("ServiceAccount")
	require.NoError(t, err)

	err = createBindingOptions.SetSubjectNamespace("default")
	require.NoError(t, err)

	_, err = cluster.CreateClusterRoleBinding(ctx, createBindingOptions)
	require.NoError(t, err)

	// Verify the ClusterRoleBinding exists
	exists, err = binding.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the ClusterRoleBinding can be retrieved
	retrievedBinding, err := cluster.GetClusterRoleBindingByName(bindingName)
	require.NoError(t, err)
	require.NotNil(t, retrievedBinding)

	// List ClusterRoleBindings to verify creation
	bindingList, err := cluster.ListClusterRoleBindingNames(ctx)
	require.NoError(t, err)
	require.Contains(t, bindingList, bindingName)

	// Delete the ClusterRoleBinding
	err = binding.Delete(ctx)
	require.NoError(t, err)
	exists, err = binding.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Cleanup: delete the ClusterRole
	err = role.Delete(ctx)
	require.NoError(t, err)
}
