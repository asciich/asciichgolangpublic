package dockeroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
)

func Test_GetPortsOnHost(t *testing.T) {
	t.Run("no ports specified returns empty list", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{}
		ports, err := options.GetPortsOnHost()
		require.NoError(t, err)
		require.EqualValues(t, []int{}, ports)
	})

	t.Run("one port", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			Ports: []string{"123:123"},
		}
		ports, err := options.GetPortsOnHost()
		require.NoError(t, err)
		require.EqualValues(t, []int{123}, ports)
	})

	t.Run("two ports", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			Ports: []string{"123:123", "345:345"},
		}
		ports, err := options.GetPortsOnHost()
		require.NoError(t, err)
		require.EqualValues(t, []int{123, 345}, ports)
	})

	t.Run("two ports unsorted", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			Ports: []string{"345:345", "123:123"},
		}
		ports, err := options.GetPortsOnHost()
		require.NoError(t, err)
		require.EqualValues(t, []int{123, 345}, ports)
	})

	t.Run("two ports unsorted and network reachable", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			Ports: []string{"345:345", "0.0.0.0:123:123"},
		}
		ports, err := options.GetPortsOnHost()
		require.NoError(t, err)
		require.EqualValues(t, []int{123, 345}, ports)
	})
}
