package kubeconfigutils_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeconfigutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_GetUserEntryByUserName(t *testing.T) {
	tests := []struct {
		path     string
		userName string
	}{
		{"./testdata/cluster-a.yaml", "kind-cluster-a"},
		{"./testdata/cluster-b.yaml", "kind-cluster-b"},
		{"./testdata/cluster-c.yaml", "clusteruser"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				kubeConfig, err := kubeconfigutils.LoadFromFilePath(tt.path, true)
				require.NoError(t, err)

				entry, err := kubeConfig.GetUserEntryByName(tt.userName)
				require.NoError(t, err)
				require.EqualValues(t, tt.userName, entry.Name)
			},
		)
	}

	t.Run("Unknown user name", func(t *testing.T) {
		kubeConfig, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-c.yaml", true)
		require.NoError(t, err)

		entry, err := kubeConfig.GetUserEntryByName("this-user-does-not-exist")
		require.Error(t, err)
		require.Nil(t, entry)
	})
}

func Test_GetUserNameByContextName(t *testing.T) {
	tests := []struct {
		path             string
		contextName      string
		expectedUserName string
	}{
		{"./testdata/cluster-a.yaml", "kind-cluster-a", "kind-cluster-a"},
		{"./testdata/cluster-b.yaml", "kind-cluster-b", "kind-cluster-b"},
		{"./testdata/cluster-c.yaml", "kind-cluster-c", "clusteruser"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				kubeConfig, err := kubeconfigutils.LoadFromFilePath(tt.path, true)
				require.NoError(t, err)

				userName, err := kubeConfig.GetUserNameByContextName(getCtx(), tt.contextName)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedUserName, userName)
			},
		)
	}

	t.Run("Unknown context name", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-c.yaml", true)
		require.NoError(t, err)

		entry, err := kubeConfig.GetUserNameByContextName(ctx, "this-context-does-not-exist")
		require.Error(t, err)
		require.EqualValues(t, entry, "")
	})
}

func TestKubeConfig_LoadFromPath(t *testing.T) {

	tests := []struct {
		path                string
		expectedClusterName string
		expectedServerNames []string
	}{
		{"./testdata/cluster-a.yaml", "kind-cluster-a", []string{"https://127.0.0.1:36435"}},
		{"./testdata/cluster-b.yaml", "kind-cluster-b", []string{"https://127.0.0.1:40889"}},
		{"./testdata/cluster-c.yaml", "kind-cluster-c", []string{"https://127.0.0.1:44935"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				kubeConfig, err := kubeconfigutils.LoadFromFilePath(tt.path, true)
				require.NoError(t, err)

				require.EqualValues(
					t,
					[]string{tt.expectedClusterName},
					mustutils.Must(kubeConfig.GetClusterNames()),
				)

				require.EqualValues(
					t,
					tt.expectedServerNames,
					mustutils.Must(kubeConfig.GetServerNames()),
				)
			},
		)
	}
}

func TestKubeConfig_IsLoadableByKubectl(t *testing.T) {
	tests := []struct {
		path string
	}{
		{"./testdata/cluster-a.yaml"},
		{"./testdata/cluster-b.yaml"},
		{"./testdata/cluster-c.yaml"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				isLoadable, err := kubeconfigutils.IsFilePathLoadableByKubectl(tt.path, verbose)
				require.NoError(t, err)
				require.True(t, isLoadable)

				kubeConfig, err := kubeconfigutils.LoadFromFilePath(tt.path, verbose)
				require.NoError(t, err)

				tempFilePath, err := kubeConfig.WriteToTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)
				defer files.MustDeleteFileByPath(tempFilePath, verbose)

				isLoadable, err = kubeconfigutils.IsFilePathLoadableByKubectl(tt.path, verbose)
				require.NoError(t, err)
				require.True(t, isLoadable)

			},
		)
	}
}

func TestKubeConfig_CheckContextsUsingKubectl(t *testing.T) {
	tests := []struct {
		path                 string
		expectedContextNames []string
	}{
		{"./testdata/cluster-a.yaml", []string{"kind-cluster-a"}},
		{"./testdata/cluster-b.yaml", []string{"kind-cluster-b"}},
		{"./testdata/cluster-c.yaml", []string{"kind-cluster-c"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				require.EqualValues(t, tt.expectedContextNames, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(tt.path, verbose)))

				kubeConfig, err := kubeconfigutils.LoadFromFilePath(tt.path, verbose)
				require.NoError(t, err)

				tempFilePath, err := kubeConfig.WriteToTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)
				defer files.MustDeleteFileByPath(tempFilePath, verbose)

				require.EqualValues(t, tt.expectedContextNames, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(tempFilePath, verbose)))

			},
		)
	}
}

func TestKubeConfig_MergeTwoConfigs(t *testing.T) {

	tests := []struct {
		path1         string
		path2         string
		expectedNames []string
	}{
		{"./testdata/cluster-a.yaml", "./testdata/cluster-b.yaml", []string{"kind-cluster-a", "kind-cluster-b"}},
		{"./testdata/cluster-b.yaml", "./testdata/cluster-a.yaml", []string{"kind-cluster-a", "kind-cluster-b"}},
		{"./testdata/cluster-a.yaml", "./testdata/cluster-c.yaml", []string{"kind-cluster-a", "kind-cluster-c"}},
		{"./testdata/cluster-c.yaml", "./testdata/cluster-b.yaml", []string{"kind-cluster-b", "kind-cluster-c"}},
		{"./testdata/cluster-a.yaml", "./testdata/cluster-a.yaml", []string{"kind-cluster-a"}},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				kubeConfig1, err := kubeconfigutils.LoadFromFilePath(tt.path1, true)
				require.NoError(t, err)

				kubeConfig2, err := kubeconfigutils.LoadFromFilePath(tt.path2, true)
				require.NoError(t, err)

				merged, err := kubeconfigutils.MergeConfig(kubeConfig1, kubeConfig2)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedNames, mustutils.Must(merged.GetClusterNames()))

				tempFilePath, err := merged.WriteToTemporaryFileAndGetPath(verbose)
				require.NoError(t, err)

				defer files.MustDeleteFileByPath(tempFilePath, verbose)

				require.EqualValues(t, tt.expectedNames, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(tempFilePath, verbose)))
			},
		)
	}
}

func TestKubeConfig_MergeThreeConfigs(t *testing.T) {
	const verbose = true

	kubeConfig1, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a.yaml", true)
	require.NoError(t, err)

	kubeConfig2, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-b.yaml", true)
	require.NoError(t, err)

	kubeConfig3, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-c.yaml", true)
	require.NoError(t, err)

	merged1, err := kubeconfigutils.MergeConfig(kubeConfig2, kubeConfig3)
	require.NoError(t, err)

	require.EqualValues(t, []string{"kind-cluster-b", "kind-cluster-c"}, mustutils.Must(merged1.GetClusterNames()))

	merged2, err := kubeconfigutils.MergeConfig(merged1, kubeConfig2, kubeConfig3, kubeConfig1)
	require.NoError(t, err)

	require.EqualValues(
		t,
		[]string{"kind-cluster-a", "kind-cluster-b", "kind-cluster-c"},
		mustutils.Must(merged2.GetClusterNames()),
	)

	tempFilePath, err := merged2.WriteToTemporaryFileAndGetPath(verbose)
	require.NoError(t, err)
	defer files.MustDeleteFileByPath(tempFilePath, verbose)

	require.EqualValues(
		t,
		[]string{"kind-cluster-a", "kind-cluster-b", "kind-cluster-c"},
		mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(tempFilePath, verbose)),
	)
}

func TestKubeConfig_UpdateUserByMerge(t *testing.T) {
	t.Run("Update user token and cert by merge", func(t *testing.T) {
		kubeConfig, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a.yaml", true)
		require.NoError(t, err)

		clientKeyData, err := kubeConfig.GetClientKeyDataForUser("kind-cluster-a")
		require.NoError(t, err)
		require.NotEqualValues(t, "NewToken", clientKeyData)

		kubeConfigUpdate, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a_update_user.yaml", true)
		require.NoError(t, err)

		clientKeyData, err = kubeConfigUpdate.GetClientKeyDataForUser("kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "NewToken", clientKeyData)

		merged, err := kubeconfigutils.MergeConfig(kubeConfig, kubeConfigUpdate)
		require.NoError(t, err)
		require.NotNil(t, merged)

		user, err := merged.GetUserEntryByName("kind-cluster-a")
		require.NoError(t, err)
		require.NotNil(t, user)

		clientKeyData, err = user.GetClientKeyData()
		require.NoError(t, err)
		require.EqualValues(t, "NewToken", clientKeyData)
	})
}

func TestKubeConfig_UpdateContextByMerge(t *testing.T) {
	t.Run("Update context by merge", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a.yaml", true)
		require.NoError(t, err)

		username, err := kubeConfig.GetUserNameByContextName(ctx, "kind-cluster-a")
		require.NoError(t, err)
		require.NotEqualValues(t, "kind-cluster-b", username)

		kubeConfigUpdate, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a_update_context.yaml", true)
		require.NoError(t, err)

		username, err = kubeConfigUpdate.GetUserNameByContextName(ctx, "kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "kind-cluster-new-name", username)

		merged, err := kubeconfigutils.MergeConfig(kubeConfig, kubeConfigUpdate)
		require.NoError(t, err)
		require.NotNil(t, merged)

		kubeConfigContext, err := merged.GetContextEntryByName("kind-cluster-a")
		require.NoError(t, err)
		require.NotNil(t, kubeConfigContext)

		username, err = kubeConfigContext.GetUserName()
		require.NoError(t, err)
		require.EqualValues(t, "kind-cluster-new-name", username)
	})
}

func TestKubeConfig_UpdateClusterByMerge(t *testing.T) {
	t.Run("Update cluster by merge", func(t *testing.T) {
		kubeConfig, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a.yaml", true)
		require.NoError(t, err)

		serverUrl, err := kubeConfig.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		require.NotEqualValues(t, "https://127.0.0.1:36436", serverUrl)

		kubeConfigUpdate, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a_update_server.yaml", true)
		require.NoError(t, err)

		serverUrl, err = kubeConfigUpdate.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "https://127.0.0.1:36436", serverUrl)

		merged, err := kubeconfigutils.MergeConfig(kubeConfig, kubeConfigUpdate)
		require.NoError(t, err)
		require.NotNil(t, merged)

		kubeCluster, err := merged.GetClusterEntryByName("kind-cluster-a")
		require.NoError(t, err)
		require.NotNil(t, kubeCluster)

		serverUrl, err = kubeCluster.GetServerUrlAsString()
		require.NoError(t, err)
		require.EqualValues(t, "https://127.0.0.1:36436", serverUrl)
	})
}

func Test_GetDefaultKubeConfigPath(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		path, err := kubeconfigutils.GetDefaultKubeConfigPath(getCtx())
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(path, "/.kube/config"))
		require.True(t, pathsutils.IsAbsolutePath(path))
	})
}

func Test_GetContextNameByClusterName(t *testing.T) {
	t.Run("cluster-a", func(t *testing.T) {
		kubeConfig, err := kubeconfigutils.LoadFromFilePath("./testdata/cluster-a.yaml", true)
		require.NoError(t, err)

		contextName, err := kubeConfig.GetContextNameByClusterName(getCtx(), "kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "kind-cluster-a", contextName)
	})
}
