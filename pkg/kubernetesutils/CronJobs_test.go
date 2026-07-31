package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_CronJobByNameExists(t *testing.T) {
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
				const cronJobName = "cronjobname"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = namespace.DeleteCronJobByName(ctx, cronJobName)
				require.NoError(t, err)

				exists, err := namespace.CronJobByNameExists(ctx, cronJobName)
				require.NoError(t, err)
				require.False(t, exists)

				cronJob, err := namespace.CreateCronJob(ctx, cronJobName, "*/5 * * * *", "busybox", []string{"echo", "hello"}, map[string]string{"app": "test"})
				require.NoError(t, err)

				exists, err = cronJob.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				exists, err = namespace.CronJobByNameExists(ctx, cronJobName)
				require.NoError(t, err)
				require.True(t, exists)

				for i := 0; i < 2; i++ {
					err = namespace.DeleteCronJobByName(ctx, cronJobName)
					require.NoError(t, err)

					exists, err := namespace.CronJobByNameExists(ctx, cronJobName)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}

func Test_ListCronJobNames(t *testing.T) {
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
				const cronJobName1 = "cronjob1"
				const cronJobName2 = "cronjob2"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				// Clean up any existing cronjobs
				err = namespace.DeleteCronJobByName(ctx, cronJobName1)
				require.NoError(t, err)
				err = namespace.DeleteCronJobByName(ctx, cronJobName2)
				require.NoError(t, err)

				// Create two cronjobs
				_, err = namespace.CreateCronJob(ctx, cronJobName1, "*/5 * * * *", "busybox", []string{"echo", "hello1"}, map[string]string{})
				require.NoError(t, err)
				_, err = namespace.CreateCronJob(ctx, cronJobName2, "*/10 * * * *", "busybox", []string{"echo", "hello2"}, map[string]string{})
				require.NoError(t, err)

				// List cronjobs
				names, err := namespace.ListCronJobNames(ctx)
				require.NoError(t, err)
				require.Contains(t, names, cronJobName1)
				require.Contains(t, names, cronJobName2)

				// Clean up
				err = namespace.DeleteCronJobByName(ctx, cronJobName1)
				require.NoError(t, err)
				err = namespace.DeleteCronJobByName(ctx, cronJobName2)
				require.NoError(t, err)
			},
		)
	}
}
