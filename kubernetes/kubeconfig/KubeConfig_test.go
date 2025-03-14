package kubeconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/testutils"
)

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
				kubeConfig := MustLoadFromFilePath(tt.path, true)

				entry, err := kubeConfig.GetUserEntryByName(tt.userName)
				require.NoError(t, err)
				require.EqualValues(t, tt.userName, entry.Name)
			},
		)
	}

	t.Run("Unknown user name", func(t *testing.T) {
		kubeConfig := MustLoadFromFilePath("./testdata/cluster-c.yaml", true)
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
				kubeConfig := MustLoadFromFilePath(tt.path, true)

				userName, err := kubeConfig.GetUserNameByContextName(tt.contextName)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedUserName, userName)
			},
		)
	}

	t.Run("Unknown context name", func(t *testing.T) {
		kubeConfig := MustLoadFromFilePath("./testdata/cluster-c.yaml", true)
		entry, err := kubeConfig.GetUserNameByContextName("this-context-does-not-exist")
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
				kubeConfig := MustLoadFromFilePath(tt.path, true)

				require.EqualValues(
					t,
					[]string{tt.expectedClusterName},
					kubeConfig.MustGetClusterNames(),
				)

				require.EqualValues(
					t,
					tt.expectedServerNames,
					kubeConfig.MustGetServerNames(),
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

				require.True(t, MustIsFilePathLoadableByKubectl(tt.path, verbose))

				kubeConfig := MustLoadFromFilePath(tt.path, verbose)

				tempFilePath := kubeConfig.MustWriteToTemporaryFileAndGetPath(verbose)
				defer files.MustDeleteFileByPath(tempFilePath, verbose)

				require.True(t, MustIsFilePathLoadableByKubectl(tempFilePath, verbose))
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

				require.EqualValues(
					t,
					tt.expectedContextNames,
					MustListContextNamesUsingKubectl(tt.path, verbose),
				)

				kubeConfig := MustLoadFromFilePath(tt.path, verbose)

				tempFilePath := kubeConfig.MustWriteToTemporaryFileAndGetPath(verbose)
				defer files.MustDeleteFileByPath(tempFilePath, verbose)

				require.EqualValues(
					t,
					tt.expectedContextNames,
					MustListContextNamesUsingKubectl(tempFilePath, verbose),
				)

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

				kubeConfig1 := MustLoadFromFilePath(tt.path1, true)
				kubeConfig2 := MustLoadFromFilePath(tt.path2, true)

				merged := MustMergeConfig(kubeConfig1, kubeConfig2)

				require.EqualValues(
					t,
					tt.expectedNames,
					merged.MustGetClusterNames(),
				)

				tempFilePath := merged.MustWriteToTemporaryFileAndGetPath(verbose)
				defer files.MustDeleteFileByPath(tempFilePath, verbose)

				require.EqualValues(
					t,
					tt.expectedNames,
					MustListContextNamesUsingKubectl(tempFilePath, verbose),
				)
			},
		)
	}
}

func TestKubeConfig_MergeThreeConfigs(t *testing.T) {
	const verbose = true

	kubeConfig1 := MustLoadFromFilePath("./testdata/cluster-a.yaml", true)
	kubeConfig2 := MustLoadFromFilePath("./testdata/cluster-b.yaml", true)
	kubeConfig3 := MustLoadFromFilePath("./testdata/cluster-c.yaml", true)

	merged1 := MustMergeConfig(kubeConfig2, kubeConfig3)

	require.EqualValues(
		t,
		[]string{"kind-cluster-b", "kind-cluster-c"},
		merged1.MustGetClusterNames(),
	)

	merged2 := MustMergeConfig(merged1, kubeConfig2, kubeConfig3, kubeConfig1)
	require.EqualValues(
		t,
		[]string{"kind-cluster-a", "kind-cluster-b", "kind-cluster-c"},
		merged2.MustGetClusterNames(),
	)

	tempFilePath := merged2.MustWriteToTemporaryFileAndGetPath(verbose)
	defer files.MustDeleteFileByPath(tempFilePath, verbose)

	require.EqualValues(
		t,
		[]string{"kind-cluster-a", "kind-cluster-b", "kind-cluster-c"},
		MustListContextNamesUsingKubectl(tempFilePath, verbose),
	)
}
