package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesparameteroptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func Test_ConfigMapByNameExists(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		// {"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const configmapName = "configmapname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				namespace, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				err = namespace.DeleteConfigMapByName(ctx, configmapName)
				require.NoError(t, err)

				exists, err := namespace.ConfigMapByNameExists(ctx, configmapName)
				require.NoError(t, err)
				require.False(t, exists)

				for i := 0; i < 2; i++ {
					ctx := contextutils.WithChangeIndicator(ctx)
					configmap, err := namespace.CreateConfigMap(ctx, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{ConfigMapData: map[string]string{}})
					require.NoError(t, err)

					if i == 0 {
						// Creating the config map is considered a change:
						require.True(t, contextutils.IsChanged(ctx))
					} else {
						// Update the same ConfigMap again and again with the same values must be idempotent and not indicate a change:
						require.False(t, contextutils.IsChanged(ctx))
					}

					exists, err = configmap.Exists(ctx)
					require.NoError(t, err)
					require.True(t, exists)
				}

				exists, err = namespace.ConfigMapByNameExists(ctx, configmapName)
				require.NoError(t, err)
				require.True(t, exists)

				for i := 0; i < 2; i++ {
					err = namespace.DeleteConfigMapByName(ctx, configmapName)
					require.NoError(t, err)

					exists, err := namespace.ConfigMapByNameExists(ctx, configmapName)
					require.NoError(t, err)
					require.False(t, exists)
				}
			},
		)
	}
}

func Test_CreateConfigMapInNonExistentNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		// {"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const configmapName = "configmapname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				// ensure namespace is absent:
				err := kubernetes.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)

				exists, err := kubernetes.NamespaceByNameExists(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, exists)

				// create the configmap in the absent namespace:
				_, err = kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{ConfigMapData: map[string]string{"my-configmap": "value"}})
				require.NoError(t, err)

				// Namespace is implicitly generated:
				exists, err = kubernetes.NamespaceByNameExists(ctx, namespaceName)
				require.NoError(t, err)
				require.True(t, exists)

				// and configmap is generated as well:
				exists, err = kubernetes.ConfigMapByNameExists(ctx, namespaceName, configmapName)
				require.NoError(t, err)
				require.True(t, exists)
			},
		)
	}
}

func Test_ReadAndWriteConfigMap(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeKubernetes"},
		// {"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"
				const configmapName = "configmapname"

				kubernetes := getKubernetesByImplementationName(getCtx(), tt.implementationName)

				labels := map[string]string{"label1": "value1"}
				content := map[string]string{"file.txt": "hello_world"}

				for i := 0; i < 2; i++ {
					configMap, err := kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{
						ConfigMapData: content,
						Labels:        labels,
					})
					require.NoError(t, err)

					currentContent, err := configMap.GetAllData(ctx)
					require.NoError(t, err)
					require.EqualValues(t, content, currentContent)

					currentLabels, err := configMap.GetAllLabels(ctx)
					require.NoError(t, err)
					require.EqualValues(t, labels, currentLabels)
				}

				labels2 := map[string]string{
					"label2": "value2",
					"label3": "value3",
				}

				for i := 0; i < 2; i++ {
					configMap, err := kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{
						ConfigMapData: content,
						Labels:        labels2,
					})
					require.NoError(t, err)

					currentContent, err := configMap.GetAllData(ctx)
					require.NoError(t, err)
					require.EqualValues(t, content, currentContent)

					currentLabels, err := configMap.GetAllLabels(ctx)
					require.NoError(t, err)
					require.EqualValues(t, labels2, currentLabels)
				}

				content2 := map[string]string{
					"file.txt":  "hello_world",
					"file2.txt": "hello_world2",
				}

				for i := 0; i < 2; i++ {
					configMap, err := kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{
						ConfigMapData: content2,
						Labels:        labels2,
					})
					require.NoError(t, err)

					currentContent, err := configMap.GetAllData(ctx)
					require.NoError(t, err)
					require.EqualValues(t, content2, currentContent)

					currentLabels, err := configMap.GetAllLabels(ctx)
					require.NoError(t, err)
					require.EqualValues(t, labels2, currentLabels)
				}

				for i := 0; i < 2; i++ {
					configMap, err := kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{
						ConfigMapData: content,
						Labels:        labels,
					})
					require.NoError(t, err)

					currentContent, err := configMap.GetAllData(ctx)
					require.NoError(t, err)
					require.EqualValues(t, content, currentContent)

					currentLabels, err := configMap.GetAllLabels(ctx)
					require.NoError(t, err)
					require.EqualValues(t, labels, currentLabels)
				}
			},
		)
	}
}
