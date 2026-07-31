package nativekubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetes"
)

func Test_CreateAndDeleteCronJob(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	namespaceName := "test-createanddeletecronjob"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("happy path", func(t *testing.T) {
		const cronJobName = "testcronjob"

		cronJobSchedule := "*/5 * * * *"
		cronJobImage := "busybox"
		cronJobCommand := []string{"echo", "Hello from CronJob"}

		err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, cronJobName)
		require.NoError(t, err)

		exists, err := nativekubernetes.CronJobExists(ctx, clientset, namespaceName, cronJobName)
		require.NoError(t, err)
		require.False(t, exists)

		// check if consecutive create, delete, create, delete... works
		for range 3 {
			err = nativekubernetes.CreateCronJob(ctx, clientset, namespaceName, cronJobName, cronJobSchedule, cronJobImage, cronJobCommand, nil)
			require.NoError(t, err)

			exists, err = nativekubernetes.CronJobExists(ctx, clientset, namespaceName, cronJobName)
			require.NoError(t, err)
			require.True(t, exists)

			err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, cronJobName)
			require.NoError(t, err)

			exists, err = nativekubernetes.CronJobExists(ctx, clientset, namespaceName, cronJobName)
			require.NoError(t, err)
			require.False(t, exists)
		}
	})
}

func Test_ListCronJobs(t *testing.T) {
	ctx := getCtx()

	// -----
	// Prepare test environment start ...
	const clusterName = kindutils.SharedClusterName

	// Ensure a local kind cluster is available for testing:
	_, err := kindutils.GetOrCreateSharedCluster(ctx)
	require.NoError(t, err)

	config, err := nativekubernetes.GetConfig(ctx, "kind-"+clusterName)
	require.NoError(t, err)

	clientset, err := nativekubernetes.GetClientSetFromRestConfig(ctx, config)
	require.NoError(t, err)

	namespaceName := "test-listcronjobs"
	err = nativekubernetes.CreateNamespace(ctx, clientset, namespaceName)
	require.NoError(t, err)

	// ... prepare test environment finished.
	// -----

	t.Run("create and delete CronJobs with list in between", func(t *testing.T) {

		cronJobNames := []string{"listcj-1", "listcj-2", "listcj-3"}
		cronJobSchedule := "*/5 * * * *"
		cronJobImage := "busybox"
		cronJobCommand := []string{"echo", "Hello"}

		// Ensure all test CronJobs are absent before starting
		for _, name := range cronJobNames {
			err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, name)
			require.NoError(t, err)
		}

		// Create CronJobs one by one and verify list grows
		for i, name := range cronJobNames {
			err = nativekubernetes.CreateCronJob(ctx, clientset, namespaceName, name, cronJobSchedule, cronJobImage, cronJobCommand, nil)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListCronJobs(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, created := range cronJobNames[:i+1] {
				require.Contains(t, listed, created)
			}
			for _, notYetCreated := range cronJobNames[i+1:] {
				require.NotContains(t, listed, notYetCreated)
			}
		}

		// Delete CronJobs one by one and verify list shrinks
		for i, name := range cronJobNames {
			err = nativekubernetes.DeleteCronJob(ctx, clientset, namespaceName, name)
			require.NoError(t, err)

			listed, err := nativekubernetes.ListCronJobs(ctx, clientset, namespaceName)
			require.NoError(t, err)

			for _, deleted := range cronJobNames[:i+1] {
				require.NotContains(t, listed, deleted)
			}
			for _, stillPresent := range cronJobNames[i+1:] {
				require.Contains(t, listed, stillPresent)
			}
		}
	})
}
