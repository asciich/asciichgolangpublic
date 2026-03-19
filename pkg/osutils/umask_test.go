package osutils_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
)

func TestOperatingSystemGetUmask(t *testing.T) {
	ctx := getCtx()

	got, err := osutils.GetUmask(ctx)
	require.NoError(t, err)
	require.Greater(t, got, 0022)
}

func TestOperatingSystemGetProcessId(t *testing.T) {
	require.Greater(t, osutils.GetProcessId(), 0)
}

func TestGetProcessDefaultDirectoryModeAsInt(t *testing.T) {
	ctx := getCtx()
	got, err := osutils.GetProcessDefaultDirectoryModeAsInt(ctx)
	require.NoError(t, err)
	require.Greater(t, got, 0744)
}

func TestProcessDefaultDirectoryModeAsFsFileMode(t *testing.T) {
	ctx := getCtx()
	got, err := osutils.GetProcessDefaultDirectoryModeAsFsFileMode(ctx)
	require.NoError(t, err)
	require.Greater(t, got, fs.FileMode(0744))
}
