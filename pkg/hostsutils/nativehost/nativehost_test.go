package nativehost_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/nativehost"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// TestNativeHost_GetDirectoryByPath_ReturnsNativeFilesooDirectory validates the spec requirement
// that nativeHost returns a struct from the nativefilesoo package for directory operations.
// This test is in the nativehost package as required by nativehost.spec.md.
func TestNativeHost_GetDirectoryByPath_ReturnsNativeFilesooDirectory(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	ctx := getCtx()
	host := nativehost.NewNativeHost()

	directory, err := host.GetDirectoryByPath(ctx, "/home/")
	require.NoError(t, err)

	// Validate that the directory is a nativefilesoo.Directory as per spec
	_, ok := directory.(*nativefilesoo.Directory)
	require.True(t, ok, "GetDirectoryByPath should return a nativefilesoo.Directory for nativeHost")
}

// TestNativeHost_GetFileByPath_ReturnsNativeFilesooFile validates the spec requirement
// that nativeHost returns a struct from the nativefilesoo package for file operations.
// This test is in the nativehost package as required by nativehost.spec.md.
func TestNativeHost_GetFileByPath_ReturnsNativeFilesooFile(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	host := nativehost.NewNativeHost()

	// Cast to *NativeHost to access GetFileByPath which is not in the Host interface
	nativeHost, ok := host.(*nativehost.NativeHost)
	require.True(t, ok, "host should be *NativeHost")

	file, err := nativeHost.GetFileByPath("/etc/hosts")
	require.NoError(t, err)

	// Validate that the file is a nativefilesoo.File as per spec
	_, okFile := file.(*nativefilesoo.File)
	require.True(t, okFile, "GetFileByPath should return a nativefilesoo.File for nativeHost")
}
