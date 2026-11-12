package commandexecutordocker

import (
	"context"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// BackupVolumes creates one gzip compressed multi volume archive at 'archivePath' containing all docker volumes in 'volumeNames'.
// Duplicate volume names are ignored.
// The archive contains a 'manifest.json' and the content of every volume below 'volumes/<volume-name>/'.
func (c *CommandExecutorDocker) BackupVolumes(ctx context.Context, volumeNames []string, archivePath string) error {
	names, err := dockergeneric.SortedUniqueVolumeNames(volumeNames)
	if err != nil {
		return err
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Backup '%d' docker volumes to '%s' started.", len(names), archivePath)

	executor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	archiveDir, archiveName, err := cmdSplitArchivePath(archivePath)
	if err != nil {
		return err
	}

	// Referencing a non-existing volume would silently create an empty one.
	// Fail early and name all missing volumes at once.
	existingVolumes, err := c.ListVolumes(ctx)
	if err != nil {
		return err
	}

	missing := []string{}
	for _, name := range names {
		if !slices.Contains(existingVolumes, name) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return tracederrors.TracedErrorf("Docker volumes do not exist. Unable to back them up: '%s'.", strings.Join(missing, "', '"))
	}

	manifest, err := dockergeneric.MarshalManifest(names)
	if err != nil {
		return err
	}

	cmd := []string{
		"docker", "run", "--rm",
		"-v", archiveDir + ":" + dockergeneric.MultiVolumeBackupDir,
	}
	for _, name := range names {
		mountPoint, err := dockergeneric.GetMultiVolumeMountPoint(name)
		if err != nil {
			return err
		}

		// The volumes are only read, so mount them read-only.
		cmd = append(cmd, "-v", name+":"+mountPoint+":ro")
	}

	// $0 = "sh", $1 = temporary archive name, $2 = owner, $3 = final archive name, $4 = manifest
	cmd = append(
		cmd,
		cmdVolumeHelperImage,
		"/bin/sh", "-c", dockergeneric.BackupVolumesScript,
		"sh", "."+archiveName+".tmp", cmdGetOwner(ctx, executor), archiveName, string(manifest),
	)

	_, err = executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to backup volumes '%s' to '%s': %w", strings.Join(names, "', '"), archivePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Backup '%d' docker volumes to '%s' finished.", len(names), archivePath)

	return nil
}

// RestoreVolumes restores all volumes contained in the multi volume archive 'archivePath' created by BackupVolumes.
// The volumes to restore are taken from the 'manifest.json' of the archive.
// Volumes which do not exist yet are created automatically.
// If the archive has no manifest, an error wrapping dockergeneric.ErrNoManifest is returned and nothing is restored.
func (c *CommandExecutorDocker) RestoreVolumes(ctx context.Context, archivePath string) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Restore all volumes from '%s' started.", archivePath)

	executor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	archiveDir, archiveName, err := cmdSplitArchivePath(archivePath)
	if err != nil {
		return err
	}

	// The archive directory is only read, so it is mounted read-only.
	archiveMount := archiveDir + ":" + dockergeneric.MultiVolumeBackupDir + ":ro"

	// Read the manifest first. Without it we don't know where to restore and nothing may be changed.
	// If the archive does not exist, tar fails and the error is returned.
	manifestOutput, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{
			"docker", "run", "--rm",
			"-v", archiveMount,
			cmdVolumeHelperImage,
			"/bin/sh", "-c", dockergeneric.ReadManifestScript, "sh", archiveName,
		},
	})
	if err != nil {
		return tracederrors.TracedErrorf("Reading the manifest of '%s' failed: %w", archivePath, err)
	}

	manifest, err := dockergeneric.ParseManifestOutput(manifestOutput)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to restore from '%s': %w", archivePath, err)
	}

	cmd := []string{
		"docker", "run", "--rm",
		"-v", archiveMount,
	}
	for _, name := range manifest.Volumes {
		mountPoint, err := dockergeneric.GetMultiVolumeMountPoint(name)
		if err != nil {
			return err
		}

		// A named volume is created automatically by docker if it does not exist.
		cmd = append(cmd, "-v", name+":"+mountPoint)
	}

	// $0 = "sh", $1 = archive name
	cmd = append(
		cmd,
		cmdVolumeHelperImage,
		"/bin/sh", "-c", dockergeneric.RestoreVolumesScript, "sh", archiveName,
	)

	_, err = executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to restore volumes '%s' from '%s': %w", strings.Join(manifest.Volumes, "', '"), archivePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Restore all volumes from '%s' finished. Restored '%d' volumes.", archivePath, len(manifest.Volumes))

	return nil
}
