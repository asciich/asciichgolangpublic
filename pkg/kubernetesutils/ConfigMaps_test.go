package kubernetesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils"
	"github.com/asciich/asciichgolangpublic/testutils"
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

				configmap, err := namespace.CreateConfigMap(ctx, configmapName, &kubernetesutils.CreateConfigMapOptions{ConfigMapData: map[string]string{}})
				require.NoError(t, err)

				exists, err = configmap.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

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
				_, err = kubernetes.CreateConfigMap(ctx, namespaceName, configmapName, &kubernetesutils.CreateConfigMapOptions{ConfigMapData: map[string]string{"my-configmap": "value"}})
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
