package nativedocker

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// Image used for the short living helper containers.
	volumeHelperImage = "alpine:latest"

	// Mount points inside the helper container.
	helperVolumeDir = "/volume"
	helperBackupDir = "/backup"

	// Maximum time granted to remove a helper container, even if the caller's ctx is already cancelled.
	helperCleanupTimeout = 30 * time.Second
)

// BackupVolume creates a gzip compressed tar archive of the docker volume 'volumeName' at 'archivePath'.
func (d *Docker) BackupVolume(ctx context.Context, volumeName string, archivePath string) error {
	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Backup docker volume '%s' to '%s' started.", volumeName, archivePath)

	archiveDir, archiveName, err := splitArchivePath(archivePath)
	if err != nil {
		return err
	}

	// Referencing a non-existing volume in a mount would silently create an empty one.
	// Fail early instead of backing up an empty, freshly created volume.
	volumes, err := d.ListVolumes(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(volumes, volumeName) {
		return tracederrors.TracedErrorf("Docker volume '%s' does not exist. Unable to back it up.", volumeName)
	}

	err = os.MkdirAll(archiveDir, 0o755)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to create directory '%s': %w", archiveDir, err)
	}

	// The archive is written to a temporary file first and renamed on success.
	// This way a failed backup never leaves a corrupt archive at the final path.
	tmpName := archiveName + ".tmp"
	tmpPath := filepath.Join(archiveDir, tmpName)
	defer os.Remove(tmpPath) // No-op after a successful rename.

	// Let the archive belong to the calling user instead of root.
	// os.Getuid() returns -1 on Windows, in this case the chown is skipped.
	owner := ""
	if uid := os.Getuid(); uid >= 0 {
		owner = strconv.Itoa(uid) + ":" + strconv.Itoa(os.Getgid())
	}

	// The file names are passed as positional parameters to avoid any shell quoting problems.
	// $0 = "sh", $1 = temporary archive name, $2 = owner (may be empty)
	script := `tar czf "` + helperBackupDir + `/$1" -C ` + helperVolumeDir + ` . && { [ -z "$2" ] || chown "$2" "` + helperBackupDir + `/$1"; }`

	err = runHelperContainer(
		ctx,
		"backup",
		[]string{"/bin/sh", "-c", script, "sh", tmpName, owner},
		[]mount.Mount{
			{Type: mount.TypeVolume, Source: volumeName, Target: helperVolumeDir},
			{Type: mount.TypeBind, Source: archiveDir, Target: helperBackupDir},
		},
	)
	if err != nil {
		return tracederrors.TracedErrorf("Backup of docker volume '%s' failed: %w", volumeName, err)
	}

	err = os.Rename(tmpPath, filepath.Join(archiveDir, archiveName))
	if err != nil {
		return tracederrors.TracedErrorf("Unable to move temporary archive '%s' to its final name '%s': %w", tmpPath, archiveName, err)
	}

	logging.LogInfoByCtxf(ctx, "Backup docker volume '%s' to '%s' finished.", volumeName, archivePath)

	return nil
}

// RestoreVolume extracts the archive 'archivePath' created by BackupVolume into the docker volume 'volumeName'.
// The volume is created automatically if it does not exist yet.
func (d *Docker) RestoreVolume(ctx context.Context, volumeName string, archivePath string) error {
	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Restore docker volume '%s' from '%s' started.", volumeName, archivePath)

	archiveDir, archiveName, err := splitArchivePath(archivePath)
	if err != nil {
		return err
	}

	info, err := os.Stat(filepath.Join(archiveDir, archiveName))
	if err != nil {
		return tracederrors.TracedErrorf("Unable to access archive '%s': %w", archivePath, err)
	}
	if !info.Mode().IsRegular() {
		return tracederrors.TracedErrorf("Archive '%s' is not a regular file.", archivePath)
	}

	err = runHelperContainer(
		ctx,
		"restore",
		[]string{"tar", "xzf", helperBackupDir + "/" + archiveName, "-C", helperVolumeDir},
		[]mount.Mount{
			// A named volume is created automatically by docker if it does not exist.
			{Type: mount.TypeVolume, Source: volumeName, Target: helperVolumeDir},
			// The archive directory is only read, so mount it read-only.
			{Type: mount.TypeBind, Source: archiveDir, Target: helperBackupDir, ReadOnly: true},
		},
	)
	if err != nil {
		return tracederrors.TracedErrorf("Restore of docker volume '%s' failed: %w", volumeName, err)
	}

	logging.LogInfoByCtxf(ctx, "Restore docker volume '%s' from '%s' finished.", volumeName, archivePath)

	return nil
}

// ListVolumes returns the names of all docker volumes in alphabetical order.
func (d *Docker) ListVolumes(ctx context.Context) ([]string, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	volumeList, err := cli.VolumeList(ctx, client.VolumeListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list volumes: %w", err)
	}

	names := make([]string, 0, len(volumeList.Items))
	for _, v := range volumeList.Items {
		names = append(names, v.Name)
	}

	sort.Strings(names)

	return names, nil
}

// runHelperContainer runs 'cmd' in a short living container with the given mounts
// and returns an error if the command does not exit with code 0.
// 'description' is only used for log and error messages.
func runHelperContainer(ctx context.Context, description string, cmd []string, mounts []mount.Mount) error {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	// No fixed container name: docker generates a unique one.
	// This avoids name collisions between parallel runs and stale containers.
	createResult, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: volumeHelperImage,
		Config: &container.Config{
			// No entrypoint override: 'Cmd' is executed directly.
			Cmd: cmd,
			// With a TTY the logs are a plain stream and not multiplexed,
			// so they can be read directly into an error message.
			Tty: true,
		},
		HostConfig: &container.HostConfig{
			AutoRemove: false, // Removed explicitly below, so the logs stay readable on failure.
			Mounts:     mounts,
		},
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to create %s container: %w", description, err)
	}

	defer func() {
		// Cleanup must also work when ctx is already cancelled.
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), helperCleanupTimeout)
		defer cancel()

		_, err := cli.ContainerRemove(cleanupCtx, createResult.ID, client.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: false, // Never remove the volume we are working on.
		})
		if err != nil {
			logging.LogErrorByCtxf(ctx, "Failed to remove %s container '%s': %v", description, createResult.ID, err)
		}
	}()

	_, err = cli.ContainerStart(ctx, createResult.ID, client.ContainerStartOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to start %s container: %w", description, err)
	}

	// The container is not auto-removed, so waiting after the start also works
	// for containers that already exited.
	waitResult := cli.ContainerWait(ctx, createResult.ID, client.ContainerWaitOptions{})

	select {
	case <-ctx.Done():
		return tracederrors.TracedErrorf("Context done while waiting for %s container: %w", description, ctx.Err())
	case err := <-waitResult.Error:
		return tracederrors.TracedErrorf("Error waiting for %s container: %w", description, err)
	case result := <-waitResult.Result:
		if result.StatusCode != 0 {
			logs, logErr := cli.ContainerLogs(ctx, createResult.ID, client.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
			})
			if logErr == nil {
				defer logs.Close()
				logBytes, _ := io.ReadAll(logs)
				return tracederrors.TracedErrorf(
					"The %s container exited with code %d. Logs: %s",
					description, result.StatusCode, strings.TrimSpace(string(logBytes)),
				)
			}
			return tracederrors.TracedErrorf("The %s container exited with code %d", description, result.StatusCode)
		}
	}

	return nil
}

// splitArchivePath returns the absolute directory and the file name of 'archivePath'.
// Docker bind mounts require absolute paths and a directory is mounted
// instead of the file itself, since the file may not exist yet.
func splitArchivePath(archivePath string) (dir string, name string, err error) {
	absPath, err := filepath.Abs(archivePath)
	if err != nil {
		return "", "", tracederrors.TracedErrorf("Unable to get absolute path of '%s': %w", archivePath, err)
	}

	return filepath.Dir(absPath), filepath.Base(absPath), nil
}
