package kindparameteroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils/kindparameteroptions"
)

func TestCreateClusterOptions_GetNodeImageOrDefaultIfUnset(t *testing.T) {
	t.Run("returns default node image if unset", func(t *testing.T) {
		options := &kindparameteroptions.CreateClusterOptions{Name: "test-cluster"}

		require.Equal(t, kindparameteroptions.DefaultNodeImage, options.GetNodeImageOrDefaultIfUnset())
	})

	t.Run("returns configured node image if set", func(t *testing.T) {
		options := &kindparameteroptions.CreateClusterOptions{
			Name:      "test-cluster",
			NodeImage: "kindest/node:v1.30.0",
		}

		require.Equal(t, "kindest/node:v1.30.0", options.GetNodeImageOrDefaultIfUnset())
	})
}

func TestCreateClusterOptions_GetNodeImageOrEmptyIfUnset(t *testing.T) {
	t.Run("returns empty string if unset", func(t *testing.T) {
		options := &kindparameteroptions.CreateClusterOptions{Name: "test-cluster"}

		require.Equal(t, "", options.GetNodeImageOrEmptyIfUnset())
	})

	t.Run("returns configured node image if set", func(t *testing.T) {
		options := &kindparameteroptions.CreateClusterOptions{
			Name:      "test-cluster",
			NodeImage: "kindest/node:v1.30.0",
		}

		require.Equal(t, "kindest/node:v1.30.0", options.GetNodeImageOrEmptyIfUnset())
	})
}
