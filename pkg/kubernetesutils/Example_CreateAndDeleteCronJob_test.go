package kubernetesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/nativekubernetesoo"
)

func Test_Example_CreateAndDeleteCronJob(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// -----
	// Prepare test environment start ...

	// ... prepare test environment finished.
	// -----

	// Get Kubernetes cluster:
	cluster, err := nativekubernetesoo.GetClusterByName(ctx, "kind-"+testClusterName)
	require.NoError(t, err)

	// Define the namespace and cronjob name we use for testing
	const namespaceName = "testnamespace"
	const cronJobName = "example-cronjob"

	// Get the namespace (will be created implicitly when creating the cronjob)
	namespace, err := cluster.GetNamespaceByName(namespaceName)
	require.NoError(t, err)

	// Get the cronjob object
	cronJob, err := namespace.GetCronJobByName(cronJobName)
	require.NoError(t, err)

	// Ensure the cronjob is absent
	err = namespace.DeleteCronJobByName(ctx, cronJobName)
	require.NoError(t, err)
	exists, err := cronJob.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)

	// Create an example cronjob using YAML
	cronJobYaml := ""
	cronJobYaml += "apiVersion: batch/v1\n"
	cronJobYaml += "kind: CronJob\n"
	cronJobYaml += "metadata:\n"
	cronJobYaml += "  name: " + cronJobName + "\n"
	cronJobYaml += "  namespace: " + namespaceName + "\n"
	cronJobYaml += "spec:\n"
	cronJobYaml += "  schedule: \"*/5 * * * *\"\n"
	cronJobYaml += "  jobTemplate:\n"
	cronJobYaml += "    spec:\n"
	cronJobYaml += "      template:\n"
	cronJobYaml += "        spec:\n"
	cronJobYaml += "          containers:\n"
	cronJobYaml += "          - name: hello\n"
	cronJobYaml += "            image: busybox\n"
	cronJobYaml += "            command: [\"echo\", \"Hello World!\"]\n"
	cronJobYaml += "          restartPolicy: OnFailure\n"

	_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{
		YamlString: cronJobYaml,
	})
	require.NoError(t, err)
	exists, err = cronJob.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the cronjob exists via the cluster
	names, err := namespace.ListCronJobNames(ctx)
	require.NoError(t, err)
	require.Contains(t, names, cronJobName)

	// Delete the cronjob
	err = namespace.DeleteCronJobByName(ctx, cronJobName)
	require.NoError(t, err)
	exists, err = cronJob.Exists(ctx)
	require.NoError(t, err)
	require.False(t, exists)
}
