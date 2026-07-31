package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CronJobByNameExists(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Create an example cronjob. This implicitly generates the namespace if it does not exist.
	const namespaceName = "testnamespace"
	const cronJobName = "example-cronjob"
	_, err = cluster.CreateCronJob(ctx, namespaceName, cronJobName, "*/5 * * * *", "busybox", []string{"echo", "Hello"}, nil)
	require.NoError(t, err)
	// ... prepare test environment finished.
	// -----

	// Our created cronjob exists:
	exists, err := cluster.CronJobByNameExists(ctx, namespaceName, cronJobName)
	require.NoError(t, err)
	require.True(t, exists)

	// The same cronjob name in the default namespace does not exist:
	exists, err = cluster.CronJobByNameExists(ctx, "default", cronJobName)
	require.NoError(t, err)
	require.False(t, exists)

	// This cronjob is expected to be in the same namespace but does not exist:
	exists, err = cluster.CronJobByNameExists(ctx, namespaceName, "cronjob-does-not-exist")
	require.NoError(t, err)
	require.False(t, exists)

	// If we delete our cronjob again...
	err = cluster.DeleteCronJobByName(ctx, namespaceName, cronJobName)
	require.NoError(t, err)

	// ... our cronjob becomes absent:
	exists, err = cluster.CronJobByNameExists(ctx, namespaceName, cronJobName)
	require.NoError(t, err)
	require.False(t, exists)
}
