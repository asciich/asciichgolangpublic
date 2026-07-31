package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CheckCronJobByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Create an example namespace and cronJob.
	const namespaceName = "testnamespace"
	const cronJobName = "example-cronjob"
	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	_, err = namespace.CreateCronJob(ctx, cronJobName, "*/5 * * * *", "busybox", []string{"echo", "hello"}, map[string]string{})
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created cronJob exists - CheckCronJobByNameExists returns nil:
	err = cluster.CheckCronJobByNameExists(ctx, namespaceName, cronJobName)
	require.NoError(t, err)

	// The same cronJob name in the default namespace does not exist - CheckCronJobByNameExists returns error:
	err = cluster.CheckCronJobByNameExists(ctx, "default", cronJobName)
	require.Error(t, err)

	// This cronJob is expected to be in the same namespace but does not exist:
	err = cluster.CheckCronJobByNameExists(ctx, namespaceName, "cronjob-does-not-exist")
	require.Error(t, err)

	// If we delete our cronJob again...
	err = namespace.DeleteCronJobByName(ctx, cronJobName)
	require.NoError(t, err)

	// ... our cronJob becomes absent and CheckCronJobByNameExists returns error:
	err = cluster.CheckCronJobByNameExists(ctx, namespaceName, cronJobName)
	require.Error(t, err)
}
