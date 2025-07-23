package kubernetesutils_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestKubernetesObject_CreateAndDelete(t *testing.T) {
	tests := []struct {
		implementationName string
		objectKind         string
		objectName         string
	}{
		{"commandExecutorKubernetes", "secret", "object-test-secret"},
		{"nativeKubernetes", "secret", "object-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"

				cluster := getKubernetesByImplementationName(ctx, tt.implementationName)
				err := cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				k8sObject, err := cluster.GetObjectByNames(tt.objectName, tt.objectKind, namespaceName)
				require.NoError(t, err)
				err = k8sObject.Delete(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(k8sObject.Exists(ctx)))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.objectName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				for i := 0; i < 2; i++ {
					err = k8sObject.CreateByYamlString(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: roleYaml})
					require.NoError(t, err)

					require.True(t, mustutils.Must(k8sObject.Exists(ctx)))
				}

				for i := 0; i < 2; i++ {
					err = k8sObject.Delete(ctx)
					require.NoError(t, err)
					require.False(t, mustutils.Must(k8sObject.Exists(ctx)))
				}
			},
		)
	}
}

func TestKubernetesObject_ListObjects(t *testing.T) {
	tests := []struct {
		implementationName string
		objectType         string
		objectName         string
	}{
		{"commandExecutorKubernetes", "secret", "object-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const namespaceName = "testnamespace"

				cluster := getKubernetesByImplementationName(ctx, tt.implementationName)
				err := cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				require.Len(
					t,
					mustutils.Must(cluster.ListObjects(
						&kubernetesparameteroptions.ListKubernetesObjectsOptions{
							Namespace:  namespaceName,
							ObjectType: tt.objectType,
						},
					)),
					0,
				)

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.objectName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				for i := 0; i < 3; i++ {
					k8sObject, err := cluster.GetObjectByNames(tt.objectName+strconv.Itoa(i), tt.objectType, namespaceName)
					require.NoError(t, err)
					err = k8sObject.CreateByYamlString(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: roleYaml})
					require.NoError(t, err)
					require.True(t, mustutils.Must(k8sObject.Exists(ctx)))
				}

				require.Len(
					t,
					mustutils.Must(cluster.ListObjects(
						&kubernetesparameteroptions.ListKubernetesObjectsOptions{
							Namespace:  namespaceName,
							ObjectType: tt.objectType,
						},
					)),
					3,
				)

				expectedNames := []string{}
				for i := 0; i < 3; i++ {
					expectedNames = append(expectedNames, tt.objectName+strconv.Itoa(i))
				}

				require.EqualValues(
					t,
					expectedNames,
					mustutils.Must(cluster.ListObjectNames(
						&kubernetesparameteroptions.ListKubernetesObjectsOptions{
							Namespace:  namespaceName,
							ObjectType: tt.objectType,
						},
					)),
				)
			},
		)
	}
}

func TestKubernetesObject_GetAsYamlString(t *testing.T) {
	tests := []struct {
		implementationName string
		objectType         string
		objectName         string
	}{
		{"commandExecutorKubernetes", "secret", "object-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				const namespaceName = "testnamespace"

				cluster := getKubernetesByImplementationName(ctx, tt.implementationName)
				err := cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				k8sObject, err := cluster.GetObjectByNames(tt.objectName, tt.objectType, namespaceName)
				require.NoError(t, err)
				err = k8sObject.Delete(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(k8sObject.Exists(ctx)))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.objectName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				err = k8sObject.CreateByYamlString(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: roleYaml})
				require.NoError(t, err)
				defer k8sObject.Delete(ctx)

				yamlString, err := k8sObject.GetAsYamlString()
				require.NoError(t, err)
				require.EqualValues(
					t,
					tt.objectName,
					yamlutils.MustRunYqQueryAginstYamlStringAsString(yamlString, ".metadata.name"),
				)
			},
		)
	}
}

func TestKubernetesObject_NamespaceCreateObject(t *testing.T) {
	tests := []struct {
		implementationName string
		objectKind         string
		objectName         string
	}{
		// {"commandExecutorKubernetes", "secret", "object-test-secret"},
		{"nativeKubernetes", "secret", "object-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"

				cluster := getKubernetesByImplementationName(ctx, tt.implementationName)
				err := cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				k8sObject, err := cluster.GetObjectByNames(tt.objectName, tt.objectKind, namespaceName)
				require.NoError(t, err)
				err = k8sObject.Delete(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(k8sObject.Exists(ctx)))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.objectName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				namespace, err := cluster.GetNamespaceByName(namespaceName)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					_, err = namespace.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: roleYaml})
					require.NoError(t, err)

					require.True(t, mustutils.Must(k8sObject.Exists(ctx)))
				}

				for i := 0; i < 2; i++ {
					err = k8sObject.Delete(ctx)
					require.NoError(t, err)
					require.False(t, mustutils.Must(k8sObject.Exists(ctx)))
				}
			},
		)
	}
}

func TestKubernetesObject_ClusterCreateObject(t *testing.T) {
	tests := []struct {
		implementationName string
		objectKind         string
		objectName         string
	}{
		// {"commandExecutorKubernetes", "secret", "object-test-secret"},
		{"nativeKubernetes", "secret", "object-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const namespaceName = "testnamespace"

				cluster := getKubernetesByImplementationName(ctx, tt.implementationName)
				err := cluster.DeleteNamespaceByName(ctx, namespaceName)
				require.NoError(t, err)
				require.False(t, mustutils.Must(cluster.NamespaceByNameExists(ctx, namespaceName)))

				k8sObject, err := cluster.GetObjectByNames(tt.objectName, tt.objectKind, namespaceName)
				require.NoError(t, err)
				err = k8sObject.Delete(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(k8sObject.Exists(ctx)))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.objectName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				for i := 0; i < 2; i++ {
					_, err = cluster.CreateObject(ctx, &kubernetesparameteroptions.CreateObjectOptions{YamlString: roleYaml})
					require.NoError(t, err)

					require.True(t, mustutils.Must(k8sObject.Exists(ctx)))
				}

				for i := 0; i < 2; i++ {
					err = k8sObject.Delete(ctx)
					require.NoError(t, err)
					require.False(t, mustutils.Must(k8sObject.Exists(ctx)))
				}
			},
		)
	}
}
