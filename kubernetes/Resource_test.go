package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert := assert.New(t)

				const verbose bool = true
				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(tt.implementationName)
				kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
				assert.False(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))

				k8sResource := kubernetes.MustGetResourceByNames(tt.resourceName, tt.resourceType, namespaceName)
				k8sResource.MustDelete(verbose)

				assert.False(k8sResource.MustExists(verbose))

				roleYaml := ""
				roleYaml += "apiVersion: v1\n"
				roleYaml += "kind: Secret\n"
				roleYaml += "metadata:\n"
				roleYaml += "  name: " + tt.resourceName + "\n"
				roleYaml += "  namespace: " + namespaceName + "\n"

				for i := 0; i < 2; i++ {
					k8sResource.MustCreateByYamlString(roleYaml, verbose)
					assert.True(k8sResource.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					k8sResource.MustDelete(verbose)
					assert.False(k8sResource.MustExists(verbose))
				}
			},
		)
	}
}
