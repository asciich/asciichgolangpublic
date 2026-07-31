package nativekubernetes_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

// This example shows how to create and delete a CronJob.
func Test_Example_CreateAndDeleteCronJob(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	namespaceName := "test-example-createanddeletecronjob"

	cluster, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	_, err = cluster.CreateNamespaceByName(ctx, namespaceName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSet(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	// Define the name of the CronJob we use for testing
	const cronJobName = "example-cronjob"

	cronJobSchedule := "*/5 * * * *"
	cronJobImage := "busybox"
	cronJobCommand := []string{"echo", "Hello from CronJob"}

	// Ensure the CronJob is absent:
	err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, cronJobName)
	require.NoError(t, err)

	cronJobNames, err := nativekubernetes.ListCronJobs(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, cronJobNames, cronJobName)

	// Create the CronJob:
	err = nativekubernetes.CreateCronJob(ctx, clientset, namespaceName, cronJobName, cronJobSchedule, cronJobImage, cronJobCommand, nil)
	require.NoError(t, err)

	cronJobNames, err = nativekubernetes.ListCronJobs(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.Contains(t, cronJobNames, cronJobName)

	// The create function is idempotent if we handle AlreadyExists:
	exists, err := nativekubernetes.CronJobExists(ctx, clientset, namespaceName, cronJobName)
	require.NoError(t, err)
	require.True(t, exists)

	// Delete the CronJob again so we have no leftovers:
	err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, cronJobName)
	require.NoError(t, err)

	cronJobNames, err = nativekubernetes.ListCronJobs(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, cronJobNames, cronJobName)

	// The delete function is idempotent as well.
	// If the CronJob is already absent no error will be raised.
	err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, cronJobName)
	require.NoError(t, err)

	cronJobNames, err = nativekubernetes.ListCronJobs(ctx, clientset, namespaceName)
	require.NoError(t, err)
	require.NotContains(t, cronJobNames, cronJobName)
}
