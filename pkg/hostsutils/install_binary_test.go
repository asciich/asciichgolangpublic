package hostsutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhost"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestCommandExecutorHost_InstallBinary_NilOptions(t *testing.T) {
	host := commandexecutorhost.NewCommandExecutorHost()

	_, err := host.InstallBinary(context.Background(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "installOptions")
}

func TestCommandExecutorHost_InstallBinary_SourcePathNotSet(t *testing.T) {
	tests := []struct {
		name    string
		options *parameteroptions.InstallOptions
		errMsg  string
	}{
		{
			name: "empty source path",
			options: &parameteroptions.InstallOptions{
				BinaryName: "myapp",
			},
			errMsg: "SourcePath not set",
		},
		{
			name: "empty binary name",
			options: &parameteroptions.InstallOptions{
				SourcePath: "/tmp/myapp",
			},
			errMsg: "BinaryName not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host := commandexecutorhost.NewCommandExecutorHost()

			// Set a command executor so GetHostName doesn't fail first
			bashExecutor := commandexecutorbashoo.Bash()
			err := host.SetCommandExecutor(bashExecutor)
			require.NoError(t, err)

			_, err = host.InstallBinary(context.Background(), tt.options)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestCommandExecutorHost_InstallBinary_NoCommandExecutor(t *testing.T) {
	host := commandexecutorhost.NewCommandExecutorHost()

	options := &parameteroptions.InstallOptions{
		SourcePath: "/tmp/source_binary",
		BinaryName: "myapp",
	}

	_, err := host.InstallBinary(context.Background(), options)
	require.Error(t, err)
}

func TestCommandExecutorHost_InstallBinary_SourceFileDoesNotExist(t *testing.T) {
	host := commandexecutorhost.NewCommandExecutorHost()

	bashExecutor := commandexecutorbashoo.Bash()
	err := host.SetCommandExecutor(bashExecutor)
	require.NoError(t, err)

	options := &parameteroptions.InstallOptions{
		SourcePath: "/tmp/nonexistent_binary_12345",
		BinaryName: "myapp",
	}

	_, err = host.InstallBinary(context.Background(), options)
	require.Error(t, err)
}

func TestCommandExecutorHost_InstallBinary_DefaultInstallationPath(t *testing.T) {
	// Verify that when InstallationPath is not set, it defaults to /bin/<binaryName>
	opts := parameteroptions.NewInstallOptions()
	err := opts.SetBinaryName("testbin")
	require.NoError(t, err)

	path, err := opts.GetInstallationPathOrDefaultIfUnset()
	require.NoError(t, err)
	assert.Equal(t, "/bin/testbin", path)
}

func TestCommandExecutorHost_InstallBinary_CustomInstallationPath(t *testing.T) {
	opts := parameteroptions.NewInstallOptions()
	err := opts.SetBinaryName("testbin")
	require.NoError(t, err)
	err = opts.SetInstallationPath("/usr/local/bin/testbin")
	require.NoError(t, err)

	path, err := opts.GetInstallationPathOrDefaultIfUnset()
	require.NoError(t, err)
	assert.Equal(t, "/usr/local/bin/testbin", path)
}

func TestCommandExecutorHost_InstallBinary_Integration(t *testing.T) {
	ctx := context.Background()

	host := commandexecutorhost.NewCommandExecutorHost()
	bashExecutor := commandexecutorbashoo.Bash()
	err := host.SetCommandExecutor(bashExecutor)
	require.NoError(t, err)

	// Create a temporary source file to act as the binary
	sourceFile := "/tmp/test_install_binary_source"
	_, err = bashExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"bash", "-c", "echo '#!/bin/bash\necho hello' > " + sourceFile + " && chmod +x " + sourceFile},
		},
	)
	require.NoError(t, err)

	// Clean up source after test
	defer func() {
		_, _ = bashExecutor.RunCommandAndGetStdoutAsString(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"rm", "-f", sourceFile},
			},
		)
	}()

	destPath := "/tmp/test_install_binary_dest"

	// Clean up destination after test
	defer func() {
		_, _ = bashExecutor.RunCommandAndGetStdoutAsString(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"rm", "-f", destPath},
			},
		)
	}()

	options := &parameteroptions.InstallOptions{
		SourcePath:       sourceFile,
		BinaryName:       "test_install_binary_dest",
		InstallationPath: destPath,
		UseSudoToInstall: false,
	}

	installedFile, err := host.InstallBinary(ctx, options)
	require.NoError(t, err)
	require.NotNil(t, installedFile)

	// Verify the file exists at the destination
	exists, err := installedFile.Exists(ctx)
	require.NoError(t, err)
	assert.True(t, exists)

	// Verify the installed path
	path, err := installedFile.GetPath()
	require.NoError(t, err)
	assert.Equal(t, destPath, path)
}
