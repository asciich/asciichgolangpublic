package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic"
)

func getKubernetesByImplementationName(implementationName string) (kubernetes KubernetesCluster) {
	if implementationName == "commandExecutorKubernetes" {
		// Directly call kind binary to avoid cyclic import...
		asciichgolangpublic.Bash().RunOneLiner(
			"kind create cluster -n kind || true",
			true,
		)

		return MustGetLocalCommandExecutorKubernetesByName("kind-kind")
	}

	asciichgolangpublic.LogFatalWithTracef(
		"Unknwon implmentation name '%s'",
		implementationName,
	)

	return nil
}

func TestNamespace_CreateAndDeleteNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKubernetes"},
	}

	for _, tt := range tests {
		t.Run(
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true
				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(tt.implementationName)

				kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
				assert.False(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))

				for i := 0; i < 2; i++ {
					kubernetes.MustCreateNamespaceByName(namespaceName, verbose)
					assert.True(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))
				}

				for i := 0; i < 2; i++ {
					kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
					assert.False(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))
				}
			},
		)
	}
}
