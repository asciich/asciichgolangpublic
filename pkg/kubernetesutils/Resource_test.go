package kubernetesutils_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestKubernetesResource_CreateAndDelete(t *testing.T) {
	tests := []struct {
		implementationName string
		resourceType       string
		resourceName       string
	}{
		{"commandExecutorKubernetes", "secret", "resource-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()
				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(ctx, tt.implementationName)
				kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
				require.False(t, kubernetes.MustNamespaceByNameExists(namespaceName, verbose))

				k8sResource := kubernetes.MustGetResourceByNames(tt.resourceName, tt.resourceType, namespaceName)
				err := k8sResource.Delete(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(k8sResource.Exists(ctx)))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.resourceName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				for i := 0; i < 2; i++ {
					err = k8sResource.CreateByYamlString(ctx, roleYaml)
					require.NoError(t, err)

					require.True(t, mustutils.Must(k8sResource.Exists(ctx)))
				}

				for i := 0; i < 2; i++ {
					err = k8sResource.Delete(ctx)
					require.NoError(t, err)
					require.False(t, mustutils.Must(k8sResource.Exists(ctx)))
				}
			},
		)
	}
}

func TestKubernetesResource_ListResources(t *testing.T) {
	tests := []struct {
		implementationName string
		resourceType       string
		resourceName       string
	}{
		{"commandExecutorKubernetes", "secret", "resource-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(ctx, tt.implementationName)
				kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
				require.False(t, kubernetes.MustNamespaceByNameExists(namespaceName, verbose))

				require.Len(
					t,
					kubernetes.MustListResources(
						&parameteroptions.ListKubernetesResourcesOptions{
							Namespace:    namespaceName,
							ResourceType: tt.resourceType,
						},
					),
					0,
				)

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.resourceName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				for i := 0; i < 3; i++ {
					k8sResource := kubernetes.MustGetResourceByNames(tt.resourceName+strconv.Itoa(i), tt.resourceType, namespaceName)
					err := k8sResource.CreateByYamlString(ctx, roleYaml)
					require.NoError(t, err)
					require.True(t, mustutils.Must(k8sResource.Exists(ctx)))
				}

				require.Len(
					t,
					kubernetes.MustListResources(
						&parameteroptions.ListKubernetesResourcesOptions{
							Namespace:    namespaceName,
							ResourceType: tt.resourceType,
						},
					),
					3,
				)

				expectedNames := []string{}
				for i := 0; i < 3; i++ {
					expectedNames = append(expectedNames, tt.resourceName+strconv.Itoa(i))
				}

				require.EqualValues(
					t,
					expectedNames,
					kubernetes.MustListResourceNames(
						&parameteroptions.ListKubernetesResourcesOptions{
							Namespace:    namespaceName,
							ResourceType: tt.resourceType,
						},
					),
				)
			},
		)
	}
}

func TestKubernetesResource_GetAsYamlString(t *testing.T) {
	tests := []struct {
		implementationName string
		resourceType       string
		resourceName       string
	}{
		{"commandExecutorKubernetes", "secret", "resource-test-secret"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				verbose := true
				ctx := getCtx()

				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(ctx, tt.implementationName)
				kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
				require.False(t, kubernetes.MustNamespaceByNameExists(namespaceName, verbose))

				k8sResource := kubernetes.MustGetResourceByNames(tt.resourceName, tt.resourceType, namespaceName)
				err := k8sResource.Delete(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(k8sResource.Exists(ctx)))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.resourceName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				err = k8sResource.CreateByYamlString(ctx, roleYaml)
				require.NoError(t, err)
				defer k8sResource.Delete(ctx)

				yamlString, err := k8sResource.GetAsYamlString()
				require.NoError(t, err)
				require.EqualValues(
					t,
					tt.resourceName,
					yamlutils.MustRunYqQueryAginstYamlStringAsString(yamlString, ".metadata.name"),
				)
			},
		)
	}
}
