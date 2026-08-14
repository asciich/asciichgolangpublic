package dockeroptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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

func Test_GetEntryPointOrNil(t *testing.T) {
	t.Run("nil when not set", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{}
		entrypoint := options.GetEntryPointOrNil()
		require.Nil(t, entrypoint)
	})

	t.Run("returns empty slice when explicitly set to empty", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			EntryPoint: []string{},
		}
		entrypoint := options.GetEntryPointOrNil()
		require.NotNil(t, entrypoint, "Empty slice should be returned as non-nil to allow overwriting entrypoint")
		require.EqualValues(t, []string{}, entrypoint)
	})

	t.Run("returns entrypoint when set", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			EntryPoint: []string{"/bin/sh", "-c"},
		}
		entrypoint := options.GetEntryPointOrNil()
		require.EqualValues(t, []string{"/bin/sh", "-c"}, entrypoint)
	})
}

func Test_SetEntryPoint(t *testing.T) {
	t.Run("nil entrypoint", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{}
		err := options.SetEntryPoint(nil)
		require.Error(t, err)
		require.True(t, tracederrors.IsTracedError(err))
	})

	t.Run("empty slice entrypoint allowed for overwrite", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{}
		err := options.SetEntryPoint([]string{})
		require.NoError(t, err)
		require.EqualValues(t, []string{}, options.EntryPoint)
	})

	t.Run("valid entrypoint", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{}
		err := options.SetEntryPoint([]string{"/bin/sh", "-c"})
		require.NoError(t, err)
		require.EqualValues(t, []string{"/bin/sh", "-c"}, options.EntryPoint)
	})
}

func Test_GetEntryPoint(t *testing.T) {
	t.Run("error when not set", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{}
		entrypoint, err := options.GetEntryPoint()
		require.Error(t, err)
		require.Nil(t, entrypoint)
		require.True(t, tracederrors.IsTracedError(err))
	})

	t.Run("returns empty slice when explicitly set", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			EntryPoint: []string{},
		}
		entrypoint, err := options.GetEntryPoint()
		require.NoError(t, err)
		require.EqualValues(t, []string{}, entrypoint)
	})

	t.Run("returns entrypoint when set", func(t *testing.T) {
		options := &dockeroptions.DockerRunContainerOptions{
			EntryPoint: []string{"/bin/sh", "-c"},
		}
		entrypoint, err := options.GetEntryPoint()
		require.NoError(t, err)
		require.EqualValues(t, []string{"/bin/sh", "-c"}, entrypoint)
	})
}

func Test_GetDeepCopy_includes_EntryPoint(t *testing.T) {
	original := &dockeroptions.DockerRunContainerOptions{
		EntryPoint: []string{"/bin/sh", "-c"},
	}

	copy := original.GetDeepCopy()

	// Verify the copy has the same entrypoint
	require.EqualValues(t, original.EntryPoint, copy.EntryPoint)

	// Verify it's a deep copy (modifying copy doesn't affect original)
	copy.EntryPoint[0] = "/bin/bash"
	require.EqualValues(t, []string{"/bin/sh", "-c"}, original.EntryPoint)
	require.EqualValues(t, []string{"/bin/bash", "-c"}, copy.EntryPoint)
}

func Test_OverwriteEntrypointBehavior(t *testing.T) {
	t.Run("nil vs empty slice distinction", func(t *testing.T) {
		// Test 1: nil EntryPoint should return nil from GetEntryPointOrNil
		optionsNil := &dockeroptions.DockerRunContainerOptions{}
		require.Nil(t, optionsNil.GetEntryPointOrNil(), "nil EntryPoint should return nil")

		// Test 2: empty slice EntryPoint should return non-nil empty slice
		optionsEmpty := &dockeroptions.DockerRunContainerOptions{
			EntryPoint: []string{},
		}
		entrypoint := optionsEmpty.GetEntryPointOrNil()
		require.NotNil(t, entrypoint, "Empty slice EntryPoint should return non-nil")
		require.EqualValues(t, []string{}, entrypoint, "Should return empty slice")

		// Test 3: SetEntryPoint should accept empty slice
		optionsSet := &dockeroptions.DockerRunContainerOptions{}
		err := optionsSet.SetEntryPoint([]string{})
		require.NoError(t, err)
		require.EqualValues(t, []string{}, optionsSet.EntryPoint)
		require.NotNil(t, optionsSet.GetEntryPointOrNil())

		// Test 4: GetEntryPoint should return empty slice without error when explicitly set
		entrypoint2, err := optionsSet.GetEntryPoint()
		require.NoError(t, err)
		require.EqualValues(t, []string{}, entrypoint2)
	})
}
