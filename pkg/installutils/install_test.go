package installutils_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/commandexecutorinstall"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/nativeinstall"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func installByImplementation(ctx context.Context, implementationName string, options *installoptions.InstallOptions) error {
	if implementationName == "nativeinstall" {
		return nativeinstall.Install(ctx, options)
	}

	if implementationName == "commandexecutorinstall" {
		commandExecutor := commandexecutorexecoo.Exec()
		return commandexecutorinstall.Install(ctx, commandExecutor, options)
	}

	if implementationName == "convenience" {
		return installutils.Install(ctx, options)
	}

	panic("Unknown implementation name: " + implementationName)
}

func Test_InstallFromLocalFile(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_local_file", func(t *testing.T) {
					if tt.implementationName == "commandexecutorinstall" {
						t.Skip("commandexecutorinstall does not support SrcPath")
					}

					ctx := getCtx()

					sourceContent := "test file content for local installation\n"
					sourceFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, sourceContent)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, sourceFile, &filesoptions.DeleteOptions{})

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcPath:     sourceFile,
						InstallPath: installPath,
						Mode:        "u=rwx",
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, sourceContent, installedContent)
				})
			},
		)
	}
}

func Test_InstallFromURL(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL", func(t *testing.T) {
					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9124
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:      "http://localhost:9124/hello_world.txt",
						InstallPath: installPath,
						Mode:        "u=rwx",
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, "hello world\n", installedContent)
				})
			},
		)
	}
}

func Test_InstallFromURLWithChecksumValidation(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL_with_checksum_validation", func(t *testing.T) {
					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9125
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					sha256sum := checksumutils.GetSha256SumFromString("hello world\n")

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:      "http://localhost:9125/hello_world.txt",
						InstallPath: installPath,
						Mode:        "u=rwx",
						Sha256Sum:   sha256sum,
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					ctxInstall2 := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall2, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:      "http://localhost:9125/hello_world.txt",
						InstallPath: installPath,
						Mode:        "u=rwx",
						Sha256Sum:   sha256sum,
					})
					require.NoError(t, err)
					require.False(t, contextutils.IsChanged(ctxInstall2))
				})
			},
		)
	}
}

func Test_InstallFromLocalFileWithChecksumValidation(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_local_file_with_checksum_validation", func(t *testing.T) {
					if tt.implementationName == "commandexecutorinstall" {
						t.Skip("commandexecutorinstall does not support SrcPath")
					}

					ctx := getCtx()

					sourceContent := "test content with checksum\n"
					sourceFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, sourceContent)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, sourceFile, &filesoptions.DeleteOptions{})

					sha256sum, err := checksumutils.GetSha256SumFromFile(ctx, sourceFile)
					require.NoError(t, err)

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcPath:     sourceFile,
						InstallPath: installPath,
						Mode:        "u=rwx",
						Sha256Sum:   sha256sum,
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)
				})
			},
		)
	}
}

func Test_InstallFromURLWithViaLocalTempDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL_with_ViaLocalTempDirectory", func(t *testing.T) {
					if tt.implementationName != "commandexecutorinstall" {
						t.Skip("ViaLocalTempDirectory is only supported by commandexecutorinstall")
					}

					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9126
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					sha256sum := checksumutils.GetSha256SumFromString("hello world\n")

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:                "http://localhost:9126/hello_world.txt",
						InstallPath:           installPath,
						Mode:                  "u=rwx",
						Sha256Sum:             sha256sum,
						ViaLocalTempDirectory: true,
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, "hello world\n", installedContent)
				})
			},
		)
	}
}

func Test_InstallFromURLTarArchive(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL_tar_archive", func(t *testing.T) {
					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9127
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:         "http://localhost:9127/hello_world.tar",
						SrcArchivePath: "hello_world.txt",
						InstallPath:    installPath,
						Mode:           "u=rwx",
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, "hello world\n", installedContent)
				})
			},
		)
	}
}

func Test_InstallFromURLTarGzArchive(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL_tar_gz_archive", func(t *testing.T) {
					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9128
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:         "http://localhost:9128/hello_world.tar.gz",
						SrcArchivePath: "hello_world.txt",
						InstallPath:    installPath,
						Mode:           "u=rwx",
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, "hello world\n", installedContent)
				})
			},
		)
	}
}

func Test_InstallFromURLTarArchiveWithChecksumValidation(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL_tar_archive_with_checksum_validation", func(t *testing.T) {
					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9129
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					// Checksum refers to the extracted file, not the archive itself.
					sha256sum := checksumutils.GetSha256SumFromString("hello world\n")

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:         "http://localhost:9129/hello_world.tar",
						SrcArchivePath: "hello_world.txt",
						InstallPath:    installPath,
						Mode:           "u=rwx",
						Sha256Sum:      sha256sum,
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, "hello world\n", installedContent)

					// Second run must be idempotent (no change).
					ctxInstall2 := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall2, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:         "http://localhost:9129/hello_world.tar",
						SrcArchivePath: "hello_world.txt",
						InstallPath:    installPath,
						Mode:           "u=rwx",
						Sha256Sum:      sha256sum,
					})
					require.NoError(t, err)
					require.False(t, contextutils.IsChanged(ctxInstall2))
				})
			},
		)
	}
}

func Test_InstallFromURLTarGzArchiveWithChecksumValidation(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativeinstall"},
		{"commandexecutorinstall"},
		{"convenience"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				t.Run("Install_from_URL_tar_gz_archive_with_checksum_validation", func(t *testing.T) {
					ctx := getCtx()

					installDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, installDir, &filesoptions.DeleteOptions{})

					installPath := filepath.Join(installDir, "installed")
					require.NoFileExists(t, installPath)

					const port int = 9130
					testWebserver, err := testwebserver.GetTestWebServer(port)
					require.NoError(t, err)
					err = testWebserver.StartInBackground(ctx)
					require.NoError(t, err)
					defer testWebserver.Stop(ctx)

					// Checksum refers to the extracted file, not the archive itself.
					sha256sum := checksumutils.GetSha256SumFromString("hello world\n")

					ctxInstall := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:         "http://localhost:9130/hello_world.tar.gz",
						SrcArchivePath: "hello_world.txt",
						InstallPath:    installPath,
						Mode:           "u=rwx",
						Sha256Sum:      sha256sum,
					})
					require.NoError(t, err)
					require.True(t, contextutils.IsChanged(ctxInstall))
					require.FileExists(t, installPath)

					installedContent, err := nativefiles.ReadAsString(ctx, installPath, &filesoptions.ReadOptions{})
					require.NoError(t, err)
					require.Equal(t, "hello world\n", installedContent)

					// Second run must be idempotent (no change).
					ctxInstall2 := contextutils.WithChangeIndicator(ctx)
					err = installByImplementation(ctxInstall2, tt.implementationName, &installoptions.InstallOptions{
						SrcUrl:         "http://localhost:9130/hello_world.tar.gz",
						SrcArchivePath: "hello_world.txt",
						InstallPath:    installPath,
						Mode:           "u=rwx",
						Sha256Sum:      sha256sum,
					})
					require.NoError(t, err)
					require.False(t, contextutils.IsChanged(ctxInstall2))
				})
			},
		)
	}
}
