package kind

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic"
)

func getKindByImplementationName(implementationName string) (kind Kind) {
	if implementationName == "commandExecutorKind" {
		return MustGetLocalCommandExecutorKind()
	}

	asciichgolangpublic.LogFatalWithTracef(
		"Unknwon implmentation name '%s'",
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
			asciichgolangpublic.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true
				const clusterName = "kind-ci-test"

				kind := getKindByImplementationName(tt.implementationName)

				kind.MustDeleteClusterByName(clusterName, verbose)
				assert.False(kind.MustClusterByNameExists(clusterName, verbose))

				for i := 0; i < 2; i++ {
					kind.MustCreateClusterByName(clusterName, verbose)
					assert.True(kind.MustClusterByNameExists(clusterName, verbose))
				}

				for i := 0; i < 2; i++ {
					kind.MustDeleteClusterByName(clusterName, verbose)
					assert.False(kind.MustClusterByNameExists(clusterName, verbose))
				}
			},
		)
	}
}
