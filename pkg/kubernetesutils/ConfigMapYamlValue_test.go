package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_ReadYamlValueAsStringFromConfigMap(t *testing.T) {
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
				const configmapName = "configmapname"

				kubernetes := getKubernetesByImplementationName(getCtx(), t, tt.implementationName)

				configMapData := map[string]string{
					"file.txt":    "hello_world",
					"values.yaml": "image:\n  tag: v1.2.3\n  repository: my-repo\n",
				}

				// Create the configmap under test. CreateConfigMap is idempotent:
				configMap, err := kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesparameteroptions.CreateConfigMapOptions{
					ConfigMapData: configMapData,
				})
				require.NoError(t, err)

				t.Run("reads plain value", func(t *testing.T) {
					value, err := configMap.ReadYamlValueAsString(ctx, ".\"file.txt\"")
					require.NoError(t, err)
					require.EqualValues(t, "hello_world", value)
				})

				t.Run("reads value from embedded yaml", func(t *testing.T) {
					value, err := configMap.ReadYamlValueAsString(ctx, ".\"values.yaml\".image.tag")
					require.NoError(t, err)
					require.EqualValues(t, "v1.2.3", value)
				})

				t.Run("errors on empty yamlPath", func(t *testing.T) {
					_, err := configMap.ReadYamlValueAsString(ctx, "")
					require.Error(t, err)
				})

				t.Run("errors on unknown yamlPath", func(t *testing.T) {
					_, err := configMap.ReadYamlValueAsString(ctx, ".not.existing.path")
					require.Error(t, err)
				})
			},
		)
	}
}
