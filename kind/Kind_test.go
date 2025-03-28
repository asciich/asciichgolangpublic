package kind

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getClusterName() (clusterName string) {
	return "kind-ci-test"
}

func getKindByImplementationName(implementationName string) (kind Kind) {
	if implementationName == "commandExecutorKind" {
		return MustGetLocalCommandExecutorKind()
	}

	logging.LogFatalWithTracef(
		"Unknown implmentation name '%s'",
		implementationName,
	)

	return nil
}

func TestKind_CreateAndDeleteCluster(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKind"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true
				clusterName := getClusterName()

				kind := getKindByImplementationName(tt.implementationName)

				kind.MustDeleteClusterByName(clusterName, verbose)
				require.False(kind.MustClusterByNameExists(clusterName, verbose))

				for i := 0; i < 2; i++ {
					kind.MustCreateClusterByName(clusterName, verbose)
					require.True(kind.MustClusterByNameExists(clusterName, verbose))
				}

				for i := 0; i < 2; i++ {
					kind.MustDeleteClusterByName(clusterName, verbose)
					require.False(kind.MustClusterByNameExists(clusterName, verbose))
				}
			},
		)
	}
}

func TestKind_CreateNamespace(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorKind"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true
				clusterName := getClusterName()

				kind := getKindByImplementationName(tt.implementationName)

				cluster := kind.MustCreateClusterByName(clusterName, verbose)

				namespaceName := "test-namespace"

				cluster.MustDeleteNamespaceByName(namespaceName, verbose)
				require.False(cluster.MustNamespaceByNameExists(namespaceName, verbose))

				for i := 0; i < 2; i++ {
					cluster.MustCreateNamespaceByName(namespaceName, verbose)
					require.True(cluster.MustNamespaceByNameExists(namespaceName, verbose))
				}

				for i := 0; i < 2; i++ {
					cluster.MustDeleteNamespaceByName(namespaceName, verbose)
					require.False(cluster.MustNamespaceByNameExists(namespaceName, verbose))
				}
			},
		)
	}
}
