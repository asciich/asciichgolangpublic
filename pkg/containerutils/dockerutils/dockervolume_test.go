package dockerutils_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/moby/moby/client"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Test_BackupAndRestoreVolume(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				const volumeName = "test-backup-restore-volume"
				const testContent = "Hello from test volume!\n"
				const testFilename = "testfile.txt"

				// Create a temporary directory for the backup archive
				tempDir, err := os.MkdirTemp("", "volume-backup-test-")
				require.NoError(t, err)
				defer os.RemoveAll(tempDir)

				archivePath := filepath.Join(tempDir, "backup.tar.gz")

				// Clean up any existing volume with the same name
				err = runCommandToCreateVolume(ctx, volumeName)
				require.NoError(t, err)
				defer func() {
					// Clean up the volume after the test
					err := runCommandToRemoveVolume(ctx, volumeName)
					if err != nil {
						logging.LogErrorByCtxf(ctx, "Failed to cleanup volume '%s': %v", volumeName, err)
					}
				}()

				// Write test content to the volume using a temporary container
				err = writeContentToVolume(ctx, docker, volumeName, testFilename, testContent)
				require.NoError(t, err)

				// Verify the content was written
				content, err := readContentFromVolume(ctx, docker, volumeName, testFilename)
				require.NoError(t, err)
				require.EqualValues(t, testContent, content)

				// Backup the volume
				err = docker.BackupVolume(ctx, volumeName, archivePath)
				require.NoError(t, err)

				// Verify the backup archive was created
				_, err = os.Stat(archivePath)
				require.NoError(t, err, "Backup archive should exist")

				// Remove the volume
				err = runCommandToRemoveVolume(ctx, volumeName)
				require.NoError(t, err)

				// Verify the volume is gone
				volumes, err := docker.ListVolumes(ctx)
				require.NoError(t, err)
				require.NotContains(t, volumes, volumeName, "Volume should be removed")

				// Restore the volume from backup
				err = docker.RestoreVolume(ctx, volumeName, archivePath)
				require.NoError(t, err)

				// Verify the volume exists again
				volumes, err = docker.ListVolumes(ctx)
				require.NoError(t, err)
				require.Contains(t, volumes, volumeName, "Volume should exist after restore")

				// Verify the content was restored
				content, err = readContentFromVolume(ctx, docker, volumeName, testFilename)
				require.NoError(t, err)
				require.EqualValues(t, testContent, content, "Content should be restored from backup")
			},
		)
	}
}

func runCommandToCreateVolume(ctx context.Context, volumeName string) error {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	_, err = cli.VolumeCreate(ctx, client.VolumeCreateOptions{
		Name: volumeName,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create volume '%s': %w", volumeName, err)
	}

	return nil
}

func runCommandToRemoveVolume(ctx context.Context, volumeName string) error {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	_, err = cli.VolumeRemove(ctx, volumeName, client.VolumeRemoveOptions{Force: true})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to remove volume '%s': %w", volumeName, err)
	}

	return nil
}

func writeContentToVolume(ctx context.Context, docker dockerinterfaces.Docker, volumeName string, filename string, content string) error {
	cmd := []string{
		"docker", "run", "--rm",
		"-v", volumeName + ":/volume",
		"alpine:latest",
		"/bin/sh", "-c", `printf '%s' "$1" > "/volume/$2"`,
		"sh", content, filename, // $0, $1, $2
	}

	executor := commandexecutorexecoo.Exec()
	_, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write content to volume '%s': %w", volumeName, err)
	}

	return nil
}


func readContentFromVolume(ctx context.Context, docker dockerinterfaces.Docker, volumeName string, filename string) (string, error) {
	// Use docker run with --rm to read content from the volume
	cmd := []string{
		"docker", "run", "--rm",
		"-v", volumeName + ":/volume",
		"alpine:latest",
		"cat", "/volume/" + filename,
	}

	executor := commandexecutorexecoo.Exec()
	stdout, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to read content from volume '%s': %w", volumeName, err)
	}

	return stdout, nil
}
