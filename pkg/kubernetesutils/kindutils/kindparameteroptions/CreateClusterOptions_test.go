package kindparameteroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kindutils/kindparameteroptions"
)

func TestCreateClusterOptions_GetNodeImage(t *testing.T) {
	t.Run("returns default node image if unset", func(t *testing.T) {
		options := &kindparameteroptions.CreateClusterOptions{Name: "test-cluster"}

		require.Equal(t, kindparameteroptions.DefaultNodeImage, options.GetNodeImage())
	})

	t.Run("returns configured node image if set", func(t *testing.T) {
		options := &kindparameteroptions.CreateClusterOptions{
			Name:      "test-cluster",
			NodeImage: "kindest/node:v1.30.0",
		}

		require.Equal(t, "kindest/node:v1.30.0", options.GetNodeImage())
	})
}
