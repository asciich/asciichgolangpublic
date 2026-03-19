package commandexecutorinstall_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/commandexecutorinstall"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_InstallFromSourceUrl(t *testing.T) {
	t.Run("only SrcUrl", func(t *testing.T) {
		ctx := getCtx()

		installDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

		installPath := filepath.Join(installDir, "installed")
		require.NoFileExists(t, installPath)

		commandExecutor := commandexecutorexecoo.Exec()

		const port int = 9123
		testWebserver, err := testwebserver.GetTestWebServer(port)
		require.NoError(t, err)
		err = testWebserver.StartInBackground(ctx)
		require.NoError(t, err)
		defer testWebserver.Stop(ctx)

		ctxInstall := contextutils.WithChangeIndicator(ctx)
		err = commandexecutorinstall.Install(ctxInstall, commandExecutor, &installoptions.InstallOptions{
			SrcUrl:      "http://localhost:9123/hello_world.txt",
			InstallPath: installPath,
			Mode:        "u=rwx",
		})
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(ctxInstall))

		modeString, err := commandexecutorfile.GetAccessPermissionsString(commandExecutor, installPath)
		require.NoError(t, err)
		require.EqualValues(t, "u=rwx,g=,o=", modeString)

		require.FileExists(t, installPath)
	})

	t.Run("only SrcUrl and valid checksum", func(t *testing.T) {
		ctx := getCtx()

		installDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

		installPath := filepath.Join(installDir, "installed")
		require.NoFileExists(t, installPath)

		commandExecutor := commandexecutorexecoo.Exec()

		const port int = 9123
		testWebserver, err := testwebserver.GetTestWebServer(port)
		require.NoError(t, err)
		err = testWebserver.StartInBackground(ctx)
		require.NoError(t, err)
		defer testWebserver.Stop(ctx)

		sha256sum := checksumutils.GetSha256SumFromString("hello world\n")

		// Since the file does not exist we expect the installation to happen:
		ctxInstall := contextutils.WithChangeIndicator(ctx)
		err = commandexecutorinstall.Install(ctxInstall, commandExecutor, &installoptions.InstallOptions{
			SrcUrl:      "http://localhost:9123/hello_world.txt",
			InstallPath: installPath,
			Mode:        "u=rwx",
			Sha256Sum:   sha256sum,
		})
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(ctxInstall))

		modeString, err := commandexecutorfile.GetAccessPermissionsString(commandExecutor, installPath)
		require.NoError(t, err)
		require.EqualValues(t, "u=rwx,g=,o=", modeString)

		// Since the file does now exist we expect the installation to do change:
		ctxInstall = contextutils.WithChangeIndicator(ctx)
		err = commandexecutorinstall.Install(ctxInstall, commandExecutor, &installoptions.InstallOptions{
			SrcUrl:      "http://localhost:9123/hello_world.txt",
			InstallPath: installPath,
			Mode:        "u=rwx",
			Sha256Sum:   sha256sum,
		})
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(ctxInstall))

		modeString, err = commandexecutorfile.GetAccessPermissionsString(commandExecutor, installPath)
		require.NoError(t, err)
		require.EqualValues(t, "u=rwx,g=,o=", modeString)

		require.FileExists(t, installPath)
	})

	t.Run("only SrcUrl and valid checksum ViaLocalTempDirectory", func(t *testing.T) {
		ctx := getCtx()

		installDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

		installPath := filepath.Join(installDir, "installed")
		require.NoFileExists(t, installPath)

		commandExecutor := commandexecutorexecoo.Exec()

		const port int = 9123
		testWebserver, err := testwebserver.GetTestWebServer(port)
		require.NoError(t, err)
		err = testWebserver.StartInBackground(ctx)
		require.NoError(t, err)
		defer testWebserver.Stop(ctx)

		sha256sum := checksumutils.GetSha256SumFromString("hello world\n")

		// Since the file does not exist we expect the installation to happen:
		ctxInstall := contextutils.WithChangeIndicator(ctx)
		err = commandexecutorinstall.Install(ctxInstall, commandExecutor, &installoptions.InstallOptions{
			SrcUrl:                "http://localhost:9123/hello_world.txt",
			InstallPath:           installPath,
			Mode:                  "u=rwx",
			Sha256Sum:             sha256sum,
			ViaLocalTempDirectory: true,
		})
		require.NoError(t, err)
		require.True(t, contextutils.IsChanged(ctxInstall))

		modeString, err := commandexecutorfile.GetAccessPermissionsString(commandExecutor, installPath)
		require.NoError(t, err)
		require.EqualValues(t, "u=rwx,g=,o=", modeString)

		// Since the file does now exist we expect the installation to do change:
		ctxInstall = contextutils.WithChangeIndicator(ctx)
		err = commandexecutorinstall.Install(ctxInstall, commandExecutor, &installoptions.InstallOptions{
			SrcUrl:      "http://localhost:9123/hello_world.txt",
			InstallPath: installPath,
			Mode:        "u=rwx",
			Sha256Sum:   sha256sum,
		})
		require.NoError(t, err)
		require.False(t, contextutils.IsChanged(ctxInstall))

		modeString, err = commandexecutorfile.GetAccessPermissionsString(commandExecutor, installPath)
		require.NoError(t, err)
		require.EqualValues(t, "u=rwx,g=,o=", modeString)

		require.FileExists(t, installPath)
	})

}
