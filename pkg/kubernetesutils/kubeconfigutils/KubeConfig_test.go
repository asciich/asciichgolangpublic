package kubeconfigutils_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path)
				require.NoError(t, err)

				entry, err := kubeConfig.GetUserEntryByName(tt.userName)
				require.NoError(t, err)
				require.EqualValues(t, tt.userName, entry.Name)
			},
		)
	}

	t.Run("Unknown user name", func(t *testing.T) {
		ctx := getCtx()
		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-c.yaml")
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path)
				require.NoError(t, err)

				userName, err := kubeConfig.GetUserNameByContextName(ctx, tt.contextName)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedUserName, userName)
			},
		)
	}

	t.Run("Unknown context name", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-c.yaml")
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
				ctx := getCtx()
				kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path)
				require.NoError(t, err)

				require.EqualValues(t, []string{tt.expectedClusterName}, mustutils.Must(kubeConfig.GetClusterNames()))

				require.EqualValues(t, tt.expectedServerNames, mustutils.Must(kubeConfig.GetServerNames()))
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
				ctx := getCtx()

				isLoadable, err := kubeconfigutils.IsFilePathLoadableByKubectl(ctx, tt.path)
				require.NoError(t, err)
				require.True(t, isLoadable)

				kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path)
				require.NoError(t, err)

				tempFilePath, err := kubeConfig.WriteToTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempFilePath, &filesoptions.DeleteOptions{})

				isLoadable, err = kubeconfigutils.IsFilePathLoadableByKubectl(ctx, tt.path)
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
				ctx := getCtx()

				require.EqualValues(t, tt.expectedContextNames, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(ctx, tt.path)))

				kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path)
				require.NoError(t, err)

				tempFilePath, err := kubeConfig.WriteToTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempFilePath, &filesoptions.DeleteOptions{})

				require.EqualValues(t, tt.expectedContextNames, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(ctx, tempFilePath)))

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
				ctx := getCtx()

				kubeConfig1, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path1)
				require.NoError(t, err)

				kubeConfig2, err := kubeconfigutils.LoadFromFilePath(ctx, tt.path2)
				require.NoError(t, err)

				merged, err := kubeconfigutils.MergeConfig(kubeConfig1, kubeConfig2)
				require.NoError(t, err)

				require.EqualValues(t, tt.expectedNames, mustutils.Must(merged.GetClusterNames()))

				tempFilePath, err := merged.WriteToTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)

				defer nativefiles.Delete(ctx, tempFilePath, &filesoptions.DeleteOptions{})

				require.EqualValues(t, tt.expectedNames, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(ctx, tempFilePath)))
			},
		)
	}
}

func TestKubeConfig_MergeThreeConfigs(t *testing.T) {
	ctx := getCtx()

	kubeConfig1, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
	require.NoError(t, err)

	kubeConfig2, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-b.yaml")
	require.NoError(t, err)

	kubeConfig3, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-c.yaml")
	require.NoError(t, err)

	merged1, err := kubeconfigutils.MergeConfig(kubeConfig2, kubeConfig3)
	require.NoError(t, err)

	require.EqualValues(t, []string{"kind-cluster-b", "kind-cluster-c"}, mustutils.Must(merged1.GetClusterNames()))

	merged2, err := kubeconfigutils.MergeConfig(merged1, kubeConfig2, kubeConfig3, kubeConfig1)
	require.NoError(t, err)

	require.EqualValues(t, []string{"kind-cluster-a", "kind-cluster-b", "kind-cluster-c"}, mustutils.Must(merged2.GetClusterNames()))

	tempFilePath, err := merged2.WriteToTemporaryFileAndGetPath(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tempFilePath, &filesoptions.DeleteOptions{})

	require.EqualValues(t, []string{"kind-cluster-a", "kind-cluster-b", "kind-cluster-c"}, mustutils.Must(kubeconfigutils.ListContextNamesUsingKubectl(ctx, tempFilePath)))
}

func TestKubeConfig_UpdateUserByMerge(t *testing.T) {
	t.Run("Update user token and cert by merge", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		clientKeyData, err := kubeConfig.GetClientKeyDataForUser("kind-cluster-a")
		require.NoError(t, err)
		require.NotEqualValues(t, "NewToken", clientKeyData)

		kubeConfigUpdate, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a_update_user.yaml")
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

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		username, err := kubeConfig.GetUserNameByContextName(ctx, "kind-cluster-a")
		require.NoError(t, err)
		require.NotEqualValues(t, "kind-cluster-b", username)

		kubeConfigUpdate, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a_update_context.yaml")
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
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		serverUrl, err := kubeConfig.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		require.NotEqualValues(t, "https://127.0.0.1:36436", serverUrl)

		kubeConfigUpdate, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a_update_server.yaml")
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
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		contextName, err := kubeConfig.GetContextNameByClusterName(ctx, "kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "kind-cluster-a", contextName)
	})
}

func Test_ListContextNames(t *testing.T) {
	t.Run("kind-cluster-a", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		contextNames, err := kubeConfig.ListContextNames(ctx)
		require.NoError(t, err)

		require.Contains(t, contextNames, "kind-cluster-a")
		require.Len(t, contextNames, 1)
	})

	t.Run("kind-cluster-b", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-b.yaml")
		require.NoError(t, err)

		contextNames, err := kubeConfig.ListContextNames(ctx)
		require.NoError(t, err)

		require.Contains(t, contextNames, "kind-cluster-b")
		require.Len(t, contextNames, 1)
	})
}

func Test_GetAndSetCurrentContext(t *testing.T) {
	t.Run("Get kind-cluster-a", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		context, err := kubeConfig.GetCurrentContext(ctx)
		require.NoError(t, err)

		require.EqualValues(t, "kind-cluster-a", context)
	})

	t.Run("Get kind-cluster-b", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-b.yaml")
		require.NoError(t, err)

		context, err := kubeConfig.GetCurrentContext(ctx)
		require.NoError(t, err)

		require.EqualValues(t, "kind-cluster-b", context)
	})

	t.Run("Get after merge", func(t *testing.T) {
		ctx := getCtx()

		kubeConfigA, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		kubeConfigB, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-b.yaml")
		require.NoError(t, err)

		kubeConfig, err := kubeconfigutils.MergeConfig(kubeConfigA, kubeConfigB)
		require.NoError(t, err)

		context, err := kubeConfig.GetCurrentContext(ctx)
		require.NoError(t, err)

		require.EqualValues(t, "kind-cluster-a", context)
	})

	t.Run("Set to nonexisting context fails", func(t *testing.T) {
		ctx := getCtx()

		kubeConfig, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		err = kubeConfig.SetCurrentContext(ctx, "non-existing")
		require.Error(t, err)
	})

	t.Run("Set after merge", func(t *testing.T) {
		ctx := getCtx()

		kubeConfigA, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		kubeConfigB, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-b.yaml")
		require.NoError(t, err)

		kubeConfig, err := kubeconfigutils.MergeConfig(kubeConfigA, kubeConfigB)
		require.NoError(t, err)

		context, err := kubeConfig.GetCurrentContext(ctx)
		require.NoError(t, err)

		require.EqualValues(t, "kind-cluster-a", context)

		err = kubeConfig.SetCurrentContext(ctx, "kind-cluster-b")
		require.NoError(t, err)

		context, err = kubeConfig.GetCurrentContext(ctx)
		require.NoError(t, err)

		require.EqualValues(t, "kind-cluster-b", context)
	})

}

// Test_GetDeepCopy_ShallowCopyBug reproduces the bug where GetDeepCopy() creates a shallow copy
// instead of a deep copy, causing the original config to be modified when the copy is modified.
// This bug causes kubectl to fail when loading merged kubeconfigs because the original config
// data gets corrupted during the merge process.
//
// These tests will FAIL until the GetDeepCopy() function is fixed to create a true deep copy.
func Test_GetDeepCopy_ShallowCopyBug(t *testing.T) {
	t.Run("GetDeepCopy should not share cluster slice with original", func(t *testing.T) {
		ctx := getCtx()

		// Load original config
		original, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		// Get the original cluster server URL before copy
		originalServerBefore, err := original.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "https://127.0.0.1:36435", originalServerBefore)

		// Create a copy
		copy := original.GetDeepCopy()

		// Modify the copy by updating the cluster entry
		updatedCluster := &kubeconfigutils.KubeConfigCluster{
			Name: "kind-cluster-a",
			Cluster: kubeconfigutils.KubeConfigClusterCluster{
				Server:                   "https://modified-server:9999",
				CertificateAuthorityData: "modified-cert-data",
			},
		}
		err = copy.AddClusterEntry(updatedCluster)
		require.NoError(t, err)

		// The original should NOT be modified
		originalServerAfter, err := original.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)

		// This assertion will FAIL with current implementation, proving the bug exists
		require.EqualValues(t, originalServerBefore, originalServerAfter,
			"BUG: Original config was modified when modifying the copy! GetDeepCopy() creates a shallow copy, not a deep copy.")
	})

	t.Run("GetDeepCopy should not share context slice with original", func(t *testing.T) {
		ctx := getCtx()

		// Load original config
		original, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		// Get the original context user before copy
		originalUserBefore, err := original.GetUserNameByContextName(ctx, "kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "kind-cluster-a", originalUserBefore)

		// Create a copy
		copy := original.GetDeepCopy()

		// Modify the copy by updating the context entry
		updatedContext := &kubeconfigutils.KubeConfigContext{
			Name: "kind-cluster-a",
			Context: struct {
				Cluster   string `yaml:"cluster"`
				Namespace string `yaml:"namespace"`
				User      string `yaml:"user"`
			}{
				Cluster: "kind-cluster-a",
				User:    "modified-user",
			},
		}
		err = copy.AddContextEntry(updatedContext)
		require.NoError(t, err)

		// The original should NOT be modified
		originalUserAfter, err := original.GetUserNameByContextName(ctx, "kind-cluster-a")
		require.NoError(t, err)

		// This assertion will FAIL with current implementation, proving the bug exists
		require.EqualValues(t, originalUserBefore, originalUserAfter,
			"BUG: Original config context was modified when modifying the copy!")
	})

	t.Run("GetDeepCopy should not share user slice with original", func(t *testing.T) {
		ctx := getCtx()

		// Load original config
		original, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		// Get the original user key data before copy
		originalKeyDataBefore, err := original.GetClientKeyDataForUser("kind-cluster-a")
		require.NoError(t, err)
		require.NotEmpty(t, originalKeyDataBefore)

		// Create a copy
		copy := original.GetDeepCopy()

		// Modify the copy by updating the user entry
		updatedUser := &kubeconfigutils.KubeConfigUser{
			Name: "kind-cluster-a",
			User: struct {
				ClientCertificateData string `yaml:"client-certificate-data"`
				ClientKeyData         string `yaml:"client-key-data"`
				Username              string `yaml:"username"`
				Password              string `yaml:"password"`
			}{
				ClientKeyData: "modified-key-data",
			},
		}
		err = copy.AddUserEntry(updatedUser)
		require.NoError(t, err)

		// The original should NOT be modified
		originalKeyDataAfter, err := original.GetClientKeyDataForUser("kind-cluster-a")
		require.NoError(t, err)

		// This assertion will FAIL with current implementation, proving the bug exists
		require.EqualValues(t, originalKeyDataBefore, originalKeyDataAfter,
			"BUG: Original config user was modified when modifying the copy!")
	})

	t.Run("MergeConfig should not corrupt original config", func(t *testing.T) {
		ctx := getCtx()

		// Load two configs
		configA, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		configB, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-b.yaml")
		require.NoError(t, err)

		// Capture original state of configA before merge
		originalAServerBefore, err := configA.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		originalAContextsBefore, err := configA.ListContextNames(ctx)
		require.NoError(t, err)

		// Merge configs (this should NOT modify configA or configB)
		_, err = kubeconfigutils.MergeConfig(configA, configB)
		require.NoError(t, err)

		// Check if configA was corrupted by the merge
		originalAServerAfter, err := configA.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)

		originalAContextsAfter, err := configA.ListContextNames(ctx)
		require.NoError(t, err)

		// These assertions will FAIL with current implementation, proving the bug exists
		require.EqualValues(t, originalAServerBefore, originalAServerAfter,
			"BUG: MergeConfig corrupted configA's cluster data!")
		require.EqualValues(t, originalAContextsBefore, originalAContextsAfter,
			"BUG: MergeConfig corrupted configA's context data!")
	})

	t.Run("MergeConfig with overlapping cluster names corrupts original", func(t *testing.T) {
		ctx := getCtx()

		// Load two configs where both have the same cluster name
		configA, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a.yaml")
		require.NoError(t, err)

		// Load update config that has the same cluster name but different server
		configUpdate, err := kubeconfigutils.LoadFromFilePath(ctx, "./testdata/cluster-a_update_server.yaml")
		require.NoError(t, err)

		// Capture original state before merge
		originalServerBefore, err := configA.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "https://127.0.0.1:36435", originalServerBefore)

		// The update config should have a different server
		updateServer, err := configUpdate.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)
		require.EqualValues(t, "https://127.0.0.1:36436", updateServer)
		require.NotEqualValues(t, originalServerBefore, updateServer)

		// Merge configs - this SHOULD NOT modify configA
		_, err = kubeconfigutils.MergeConfig(configA, configUpdate)
		require.NoError(t, err)

		// Check if configA was corrupted
		originalServerAfter, err := configA.GetClusterServerUrlAsString("kind-cluster-a")
		require.NoError(t, err)

		// BUG: configA's server will be changed to the update server value
		require.EqualValues(t, originalServerBefore, originalServerAfter,
			"BUG: MergeConfig corrupted configA's cluster data when merging overlapping cluster names!")
	})
}

func Test_ReadCurrentKubeConfigAsString(t *testing.T) {
	t.Run("Read default kubeconfig path", func(t *testing.T) {
		ctx := getCtx()

		// Test reading from default path (should fail gracefully if no config exists)
		content, err := kubeconfigutils.ReadCurrentKubeConfigAsString(ctx)
		// May succeed if default config exists, or fail if not
		if err == nil {
			require.NotEmpty(t, content)
			require.Contains(t, content, "apiVersion:")
		}
	})

	t.Run("Read with KUBECONFIG env var set to path", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Set KUBECONFIG to test file
		os.Setenv("KUBECONFIG", "./testdata/cluster-a.yaml")

		content, err := kubeconfigutils.ReadCurrentKubeConfigAsString(ctx)
		require.NoError(t, err)
		require.Contains(t, content, "kind-cluster-a")
	})

	t.Run("Read with inline KUBECONFIG", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Set KUBECONFIG to inline config
		inlineConfig := `apiVersion: v1
kind: Config
clusters:
- name: test-cluster
  cluster:
    server: https://test-server:6443
contexts:
- name: test-context
  context:
    cluster: test-cluster
`
		os.Setenv("KUBECONFIG", inlineConfig)

		content, err := kubeconfigutils.ReadCurrentKubeConfigAsString(ctx)
		require.NoError(t, err)
		require.Contains(t, content, "test-cluster")
	})
}

func Test_ReadCurrentKubeConfig(t *testing.T) {
	t.Run("Read current kubeconfig as struct", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Set KUBECONFIG to test file
		os.Setenv("KUBECONFIG", "./testdata/cluster-a.yaml")

		config, err := kubeconfigutils.ReadCurrentKubeConfig(ctx)
		require.NoError(t, err)
		require.NotNil(t, config)

		clusterNames, err := config.GetClusterNames()
		require.NoError(t, err)
		require.Contains(t, clusterNames, "kind-cluster-a")
	})
}

func Test_AddAdditionalConfigFromString(t *testing.T) {
	t.Run("Add config to non-existing default config", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Set KUBECONFIG to a non-existing temp file
		tempDir := t.TempDir()
		tempConfigPath := filepath.Join(tempDir, "config")
		os.Setenv("KUBECONFIG", tempConfigPath)

		additionalConfig := `apiVersion: v1
kind: Config
clusters:
- name: new-cluster
  cluster:
    server: https://new-server:6443
contexts:
- name: new-context
  context:
    cluster: new-cluster
`

		err := kubeconfigutils.AddAdditionalConfigFromString(ctx, additionalConfig)
		require.NoError(t, err)

		// Verify the config was written
		config, err := kubeconfigutils.LoadFromFilePath(ctx, tempConfigPath)
		require.NoError(t, err)

		clusterNames, err := config.GetClusterNames()
		require.NoError(t, err)
		require.Contains(t, clusterNames, "new-cluster")
	})

	t.Run("Merge config with existing config", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Copy existing config to temp location
		tempDir := t.TempDir()
		tempConfigPath := filepath.Join(tempDir, "config")

		sourceContent, err := os.ReadFile("./testdata/cluster-a.yaml")
		require.NoError(t, err)

		err = os.WriteFile(tempConfigPath, sourceContent, 0644)
		require.NoError(t, err)

		os.Setenv("KUBECONFIG", tempConfigPath)

		// Note: AddAdditionalConfigFromString expects cluster, context, and user to all have the same name
		additionalConfig := `apiVersion: v1
kind: Config
clusters:
- name: additional-cluster
  cluster:
    server: https://additional-server:6443
contexts:
- name: additional-cluster
  context:
    cluster: additional-cluster
    user: additional-cluster
users:
- name: additional-cluster
  user:
    client-certificate-data: cert-data
    client-key-data: key-data
`

		err = kubeconfigutils.AddAdditionalConfigFromString(ctx, additionalConfig)
		require.NoError(t, err)

		// Verify both configs are present
		config, err := kubeconfigutils.LoadFromFilePath(ctx, tempConfigPath)
		require.NoError(t, err)

		clusterNames, err := config.GetClusterNames()
		require.NoError(t, err)
		require.Contains(t, clusterNames, "kind-cluster-a")
		require.Contains(t, clusterNames, "additional-cluster")
	})

	t.Run("Merge config with distinct context user and cluster names", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Copy existing config to temp location
		tempDir := t.TempDir()
		tempConfigPath := filepath.Join(tempDir, "config")

		sourceContent, err := os.ReadFile("./testdata/cluster-a.yaml")
		require.NoError(t, err)

		err = os.WriteFile(tempConfigPath, sourceContent, 0644)
		require.NoError(t, err)

		os.Setenv("KUBECONFIG", tempConfigPath)

		// Config with distinct names for context, user, and cluster
		additionalConfig := `apiVersion: v1
kind: Config
clusters:
- name: distinct-cluster
  cluster:
    server: https://distinct-cluster-server:6443
contexts:
- name: distinct-context
  context:
    cluster: distinct-cluster
    user: distinct-user
users:
- name: distinct-user
  user:
    client-certificate-data: distinct-cert-data
    client-key-data: distinct-key-data
`

		err = kubeconfigutils.AddAdditionalConfigFromString(ctx, additionalConfig)
		require.NoError(t, err)

		// Verify both configs are present
		config, err := kubeconfigutils.LoadFromFilePath(ctx, tempConfigPath)
		require.NoError(t, err)

		clusterNames, err := config.GetClusterNames()
		require.NoError(t, err)
		require.Contains(t, clusterNames, "kind-cluster-a")
		require.Contains(t, clusterNames, "distinct-cluster")

		contextNames, err := config.ListContextNames(ctx)
		require.NoError(t, err)
		require.Contains(t, contextNames, "kind-cluster-a")
		require.Contains(t, contextNames, "distinct-context")
	})

	t.Run("Merge config where context references different named cluster and user", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Copy existing config to temp location
		tempDir := t.TempDir()
		tempConfigPath := filepath.Join(tempDir, "config")

		sourceContent, err := os.ReadFile("./testdata/cluster-a.yaml")
		require.NoError(t, err)

		err = os.WriteFile(tempConfigPath, sourceContent, 0644)
		require.NoError(t, err)

		os.Setenv("KUBECONFIG", tempConfigPath)

		// Config where context name differs from cluster and user names
		additionalConfig := `apiVersion: v1
kind: Config
clusters:
- name: my-cluster
  cluster:
    server: https://my-cluster-server:6443
- name: another-cluster
  cluster:
    server: https://another-cluster-server:6443
contexts:
- name: my-context
  context:
    cluster: my-cluster
    user: my-user
- name: another-context
  context:
    cluster: another-cluster
    user: my-user
users:
- name: my-user
  user:
    client-certificate-data: my-cert-data
    client-key-data: my-key-data
`

		err = kubeconfigutils.AddAdditionalConfigFromString(ctx, additionalConfig)
		require.NoError(t, err)

		// Verify all entries are present
		config, err := kubeconfigutils.LoadFromFilePath(ctx, tempConfigPath)
		require.NoError(t, err)

		clusterNames, err := config.GetClusterNames()
		require.NoError(t, err)
		require.Contains(t, clusterNames, "kind-cluster-a")
		require.Contains(t, clusterNames, "my-cluster")
		require.Contains(t, clusterNames, "another-cluster")

		contextNames, err := config.ListContextNames(ctx)
		require.NoError(t, err)
		require.Contains(t, contextNames, "kind-cluster-a")
		require.Contains(t, contextNames, "my-context")
		require.Contains(t, contextNames, "another-context")
	})

	t.Run("Error when KUBECONFIG contains inline config", func(t *testing.T) {
		ctx := getCtx()

		// Save original KUBECONFIG
		originalKubeconfig := os.Getenv("KUBECONFIG")
		defer os.Setenv("KUBECONFIG", originalKubeconfig)

		// Set KUBECONFIG to inline config
		inlineConfig := `apiVersion: v1
kind: Config
clusters:
- name: inline-cluster
  cluster:
    server: https://inline-server:6443
`
		os.Setenv("KUBECONFIG", inlineConfig)

		additionalConfig := `apiVersion: v1
kind: Config
clusters:
- name: new-cluster
  cluster:
    server: https://new-server:6443
`

		err := kubeconfigutils.AddAdditionalConfigFromString(ctx, additionalConfig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Cannot merge config when KUBECONFIG contains inline config")
	})

	t.Run("Error on invalid YAML", func(t *testing.T) {
		ctx := getCtx()

		invalidConfig := `this is not valid yaml: [`

		err := kubeconfigutils.AddAdditionalConfigFromString(ctx, invalidConfig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid additional config YAML")
	})
}
