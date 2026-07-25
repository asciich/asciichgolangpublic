package openhandsutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBindAddress_LocalhostDefault(t *testing.T) {
	opts := &StartContainerOptions{
		ContainerName:            "test-container",
		Port:                     3000,
		ReachableByOtherMachines: false,
	}

	address, err := opts.GetBindAddress()
	require.NoError(t, err)
	require.Equal(t, "127.0.0.1:3000", address)
}

func TestGetBindAddress_ReachableByOtherMachines(t *testing.T) {
	opts := &StartContainerOptions{
		ContainerName:            "test-container",
		Port:                     3000,
		ReachableByOtherMachines: true,
	}

	address, err := opts.GetBindAddress()
	require.NoError(t, err)
	require.Equal(t, "0.0.0.0:3000", address)
}

func TestGetBindAddress_DifferentPorts(t *testing.T) {
	tests := []struct {
		name      string
		port      int
		reachable bool
		expected  string
	}{
		{"port 8080 localhost", 8080, false, "127.0.0.1:8080"},
		{"port 8080 all interfaces", 8080, true, "0.0.0.0:8080"},
		{"port 443 localhost", 443, false, "127.0.0.1:443"},
		{"port 443 all interfaces", 443, true, "0.0.0.0:443"},
		{"port 1 localhost", 1, false, "127.0.0.1:1"},
		{"port 65535 all interfaces", 65535, true, "0.0.0.0:65535"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &StartContainerOptions{
				ContainerName:            "test-container",
				Port:                     tt.port,
				ReachableByOtherMachines: tt.reachable,
			}

			address, err := opts.GetBindAddress()
			require.NoError(t, err)
			require.Equal(t, tt.expected, address)
		})
	}
}

func TestGetBindAddress_PortNotSet(t *testing.T) {
	opts := &StartContainerOptions{
		ContainerName: "test-container",
		Port:          0,
	}

	address, err := opts.GetBindAddress()
	require.Error(t, err)
	require.Empty(t, address)
}

func TestGetBindAddress_NegativePort(t *testing.T) {
	opts := &StartContainerOptions{
		ContainerName: "test-container",
		Port:          -1,
	}

	address, err := opts.GetBindAddress()
	require.Error(t, err)
	require.Empty(t, address)
}

func TestGetWorkspacePath_AbsolutePath(t *testing.T) {
	opts := &StartContainerOptions{
		WorkspacePath: "/tmp/workspace",
	}

	path, err := opts.GetWorkspacePath()
	require.NoError(t, err)
	require.Equal(t, "/tmp/workspace", path)
}

func TestGetWorkspacePath_RelativePath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	opts := &StartContainerOptions{
		WorkspacePath: "my-workspace",
	}

	path, err := opts.GetWorkspacePath()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(cwd, "my-workspace"), path)
}

func TestGetWorkspacePath_DotPath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	opts := &StartContainerOptions{
		WorkspacePath: ".",
	}

	path, err := opts.GetWorkspacePath()
	require.NoError(t, err)
	require.Equal(t, cwd, path)
}

func TestGetWorkspacePath_ParentDirectory(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	opts := &StartContainerOptions{
		WorkspacePath: "..",
	}

	path, err := opts.GetWorkspacePath()
	require.NoError(t, err)
	require.Equal(t, filepath.Dir(cwd), path)
}

func TestGetWorkspacePath_Empty(t *testing.T) {
	opts := &StartContainerOptions{
		WorkspacePath: "",
	}

	path, err := opts.GetWorkspacePath()
	require.Error(t, err)
	require.Empty(t, path)
}
