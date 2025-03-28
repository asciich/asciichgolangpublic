package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getKubernetesByImplementationName(implementationName string) (kubernetes KubernetesCluster) {
	if implementationName == "commandExecutorKubernetes" {
		// Directly call kind binary to avoid cyclic import...
		commandexecutor.Bash().RunOneLiner(
			"kind create cluster -n kind || true",
			true,
		)

		return MustGetLocalCommandExecutorKubernetesByName("kind-kind")
	}

	logging.LogFatalWithTracef(
		"Unknown implmentation name '%s'",
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true
				const namespaceName = "testnamespace"

				kubernetes := getKubernetesByImplementationName(tt.implementationName)

				kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
				require.False(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))

				for i := 0; i < 2; i++ {
					kubernetes.MustCreateNamespaceByName(namespaceName, verbose)
					require.True(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))
				}

				for i := 0; i < 2; i++ {
					kubernetes.MustDeleteNamespaceByName(namespaceName, verbose)
					require.False(kubernetes.MustNamespaceByNameExists(namespaceName, verbose))
				}
			},
		)
	}
}
