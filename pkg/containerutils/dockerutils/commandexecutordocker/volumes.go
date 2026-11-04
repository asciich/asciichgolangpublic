package commandexecutordocker

import (
	"context"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// Image used for the short living helper containers.
	cmdVolumeHelperImage = "alpine:latest"

	// Mount points inside the helper container.
	cmdHelperVolumeDir = "/volume"
	cmdHelperBackupDir = "/backup"
)

// BackupVolume creates a gzip compressed tar archive of the docker volume 'volumeName' at 'archivePath'.
func (c *CommandExecutorDocker) BackupVolume(ctx context.Context, volumeName string, archivePath string) error {
	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Backup docker volume '%s' to '%s' started.", volumeName, archivePath)

	executor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	archiveDir, archiveName, err := cmdSplitArchivePath(archivePath)
	if err != nil {
		return err
	}

	// Referencing a non-existing volume would silently create an empty one.
	// Fail early instead of backing up an empty, freshly created volume.
	volumes, err := c.ListVolumes(ctx)
	if err != nil {
		return err
	}
	if !slices.Contains(volumes, volumeName) {
		return tracederrors.TracedErrorf("Docker volume '%s' does not exist. Unable to back it up.", volumeName)
	}

	// Let the archive belong to the user running the commands instead of root.
	owner := cmdGetOwner(ctx, executor)

	// The archive is written to a temporary file first and moved to its final name on success.
	// The trap removes the temporary file if anything fails, so no corrupt archive is left behind.
	// The file names are passed as positional parameters to avoid any shell quoting problems:
	// $0 = "sh", $1 = temporary archive name, $2 = owner (may be empty), $3 = final archive name
	script := `set -e
trap 'rm -f "` + cmdHelperBackupDir + `/$1"' EXIT
tar czf "` + cmdHelperBackupDir + `/$1" -C ` + cmdHelperVolumeDir + ` .
if [ -n "$2" ]; then chown "$2" "` + cmdHelperBackupDir + `/$1"; fi
mv "` + cmdHelperBackupDir + `/$1" "` + cmdHelperBackupDir + `/$3"`

	cmd := []string{
		"docker", "run", "--rm",
		"-v", volumeName + ":" + cmdHelperVolumeDir,
		"-v", archiveDir + ":" + cmdHelperBackupDir,
		cmdVolumeHelperImage,
		"/bin/sh", "-c", script, "sh", "." + archiveName + ".tmp", owner, archiveName,
	}

	_, err = executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to backup volume '%s' to '%s': %w", volumeName, archivePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Backup docker volume '%s' to '%s' finished.", volumeName, archivePath)

	return nil
}

// RestoreVolume extracts the archive 'archivePath' created by BackupVolume into the docker volume 'volumeName'.
// The volume is created automatically if it does not exist yet.
func (c *CommandExecutorDocker) RestoreVolume(ctx context.Context, volumeName string, archivePath string) error {
	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Restore docker volume '%s' from '%s' started.", volumeName, archivePath)

	executor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	archiveDir, archiveName, err := cmdSplitArchivePath(archivePath)
	if err != nil {
		return err
	}

	// A named volume is created automatically by docker if it does not exist.
	// The archive directory is only read, so it is mounted read-only.
	// If the archive does not exist, tar fails and the error is returned.
	cmd := []string{
		"docker", "run", "--rm",
		"-v", volumeName + ":" + cmdHelperVolumeDir,
		"-v", archiveDir + ":" + cmdHelperBackupDir + ":ro",
		cmdVolumeHelperImage,
		"tar", "xzf", cmdHelperBackupDir + "/" + archiveName, "-C", cmdHelperVolumeDir,
	}

	_, err = executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to restore volume '%s' from '%s': %w", volumeName, archivePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Restore docker volume '%s' from '%s' finished.", volumeName, archivePath)

	return nil
}

// ListVolumes returns the names of all docker volumes in alphabetical order.
func (c *CommandExecutorDocker) ListVolumes(ctx context.Context) ([]string, error) {
	executor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return nil, err
	}

	output, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"docker", "volume", "ls", "--format", "{{.Name}}"},
	})
	if err != nil {
		return nil, err
	}

	names := stringsutils.SplitLines(output, true)

	sort.Strings(names)

	logging.LogInfoByCtxf(ctx, "Listed '%d' volumes on host '%s'", len(names), hostDescription)

	return names, nil
}

// cmdSplitArchivePath returns the absolute directory and the file name of 'archivePath'.
//
// The directory is mounted instead of the file itself, since the file may not exist yet.
// The path must be absolute, otherwise docker treats the value of '-v' as a named volume.
func cmdSplitArchivePath(archivePath string) (dir string, name string, err error) {
	absPath, err := filepath.Abs(archivePath)
	if err != nil {
		return "", "", tracederrors.TracedErrorf("Unable to get absolute path of '%s': %w", archivePath, err)
	}

	dir = filepath.Dir(absPath)
	name = filepath.Base(absPath)

	// ':' is the separator of the '-v' syntax and can't be escaped.
	if strings.Contains(dir, ":") {
		return "", "", tracederrors.TracedErrorf("The directory '%s' contains a ':' which is not supported in a docker bind mount.", dir)
	}

	return dir, name, nil
}

// cmdGetOwner returns "uid:gid" of the user running the commands on the host of the executor.
// An empty string is returned if it can't be determined. The archive then stays owned by root.
func cmdGetOwner(ctx context.Context, executor commandexecutorinterfaces.CommandExecutor) string {
	uid, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"id", "-u"},
	})
	if err != nil {
		logging.LogInfoByCtxf(ctx, "Unable to determine the user id, the archive will be owned by root: %v", err)
		return ""
	}

	gid, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"id", "-g"},
	})
	if err != nil {
		logging.LogInfoByCtxf(ctx, "Unable to determine the group id, the archive will be owned by root: %v", err)
		return ""
	}

	return strings.TrimSpace(uid) + ":" + strings.TrimSpace(gid)
}
