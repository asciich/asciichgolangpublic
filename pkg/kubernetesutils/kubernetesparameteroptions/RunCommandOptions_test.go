package kubernetesparameteroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesparameteroptions"
)

func TestIsStdinDataAvailable(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{}
		require.False(t, options.IsStinDataAvailable())
	})

	t.Run("stdin bytes set", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{StdinBytes: []byte("example")}
		require.True(t, options.IsStinDataAvailable())
	})
} 
