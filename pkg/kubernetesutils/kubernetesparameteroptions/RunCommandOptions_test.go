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

func TestGetContainerName(t *testing.T) {
	t.Run("explicit container name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			ContainerName: "mycontainer",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "mycontainer", name)
	})

	t.Run("fallback to pod name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			PodName: "mypod",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "mypod", name)
	})

	t.Run("fallback to ReplicaSet name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			ReplicaSetName: "myreplicaset",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "myreplicaset", name)
	})

	t.Run("fallback to Deployment name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			DeploymentName: "mydeployment",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "mydeployment", name)
	})

	t.Run("pod name takes precedence over ReplicaSet and Deployment name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			PodName:        "mypod",
			ReplicaSetName: "myreplicaset",
			DeploymentName: "mydeployment",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "mypod", name)
	})

	t.Run("ReplicaSet name takes precedence over Deployment name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			ReplicaSetName: "myreplicaset",
			DeploymentName: "mydeployment",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "myreplicaset", name)
	})

	t.Run("explicit container name takes precedence over pod name", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{
			ContainerName: "mycontainer",
			PodName:       "mypod",
		}
		name, err := options.GetContainerName()
		require.NoError(t, err)
		require.Equal(t, "mycontainer", name)
	})

	t.Run("error when no name available", func(t *testing.T) {
		options := &kubernetesparameteroptions.RunCommandOptions{}
		name, err := options.GetContainerName()
		require.Error(t, err)
		require.Empty(t, name)
		require.Contains(t, err.Error(), "ContainerName not set")
	})
}

