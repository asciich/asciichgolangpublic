package nativedocker

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/moby/moby/api/types/mount"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// BackupVolumes creates one gzip compressed multi volume archive at 'archivePath' containing all docker volumes in 'volumeNames'.
// Duplicate volume names are ignored.
// The archive contains a 'manifest.json' and the content of every volume below 'volumes/<volume-name>/'.
func (d *Docker) BackupVolumes(ctx context.Context, volumeNames []string, archivePath string) error {
	names, err := dockergeneric.SortedUniqueVolumeNames(volumeNames)
	if err != nil {
		return err
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Backup '%d' docker volumes to '%s' started.", len(names), archivePath)

	archiveDir, archiveName, err := splitArchivePath(archivePath)
	if err != nil {
		return err
	}

	// Referencing a non-existing volume in a mount would silently create an empty one.
	// Fail early and name all missing volumes at once.
	existingVolumes, err := d.ListVolumes(ctx)
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

	err = os.MkdirAll(archiveDir, 0o755)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to create directory '%s': %w", archiveDir, err)
	}

	// The script writes to the temporary name and moves it to the final name on success.
	// Removing the temporary file here is a second safety net, e.g. if the container is killed.
	tmpName := "." + archiveName + ".tmp"
	defer os.Remove(filepath.Join(archiveDir, tmpName))

	mounts := []mount.Mount{
		{Type: mount.TypeBind, Source: archiveDir, Target: dockergeneric.MultiVolumeBackupDir},
	}
	for _, name := range names {
		mountPoint, err := dockergeneric.GetMultiVolumeMountPoint(name)
		if err != nil {
			return err
		}

		// The volumes are only read, so mount them read-only.
		mounts = append(mounts, mount.Mount{Type: mount.TypeVolume, Source: name, Target: mountPoint, ReadOnly: true})
	}

	// $0 = "sh", $1 = temporary archive name, $2 = owner, $3 = final archive name, $4 = manifest
	_, err = runHelperContainer(
		ctx,
		"backup volumes",
		[]string{"/bin/sh", "-c", dockergeneric.BackupVolumesScript, "sh", tmpName, getArchiveOwner(), archiveName, string(manifest)},
		mounts,
	)
	if err != nil {
		return tracederrors.TracedErrorf("Backup of docker volumes '%s' failed: %w", strings.Join(names, "', '"), err)
	}

	logging.LogInfoByCtxf(ctx, "Backup '%d' docker volumes to '%s' finished.", len(names), archivePath)

	return nil
}

// RestoreVolumes restores all volumes contained in the multi volume archive 'archivePath' created by BackupVolumes.
// The volumes to restore are taken from the 'manifest.json' of the archive.
// Volumes which do not exist yet are created automatically.
// If the archive has no manifest, an error wrapping dockergeneric.ErrNoManifest is returned and nothing is restored.
func (d *Docker) RestoreVolumes(ctx context.Context, archivePath string) error {
	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	logging.LogInfoByCtxf(ctx, "Restore all volumes from '%s' started.", archivePath)

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

	// The archive directory is only read, so mount it read-only.
	archiveMount := mount.Mount{Type: mount.TypeBind, Source: archiveDir, Target: dockergeneric.MultiVolumeBackupDir, ReadOnly: true}

	// Read the manifest first. Without it we don't know where to restore and nothing may be changed.
	manifestOutput, err := runHelperContainer(
		ctx,
		"read manifest",
		[]string{"/bin/sh", "-c", dockergeneric.ReadManifestScript, "sh", archiveName},
		[]mount.Mount{archiveMount},
	)
	if err != nil {
		return tracederrors.TracedErrorf("Reading the manifest of '%s' failed: %w", archivePath, err)
	}

	manifest, err := dockergeneric.ParseManifestOutput(manifestOutput)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to restore from '%s': %w", archivePath, err)
	}

	mounts := []mount.Mount{archiveMount}
	for _, name := range manifest.Volumes {
		mountPoint, err := dockergeneric.GetMultiVolumeMountPoint(name)
		if err != nil {
			return err
		}

		// A named volume is created automatically by docker if it does not exist.
		mounts = append(mounts, mount.Mount{Type: mount.TypeVolume, Source: name, Target: mountPoint})
	}

	// $0 = "sh", $1 = archive name
	_, err = runHelperContainer(
		ctx,
		"restore volumes",
		[]string{"/bin/sh", "-c", dockergeneric.RestoreVolumesScript, "sh", archiveName},
		mounts,
	)
	if err != nil {
		return tracederrors.TracedErrorf("Restore of docker volumes '%s' failed: %w", strings.Join(manifest.Volumes, "', '"), err)
	}

	logging.LogInfoByCtxf(ctx, "Restore all volumes from '%s' finished. Restored '%d' volumes.", archivePath, len(manifest.Volumes))

	return nil
}

// getArchiveOwner returns "uid:gid" of the calling user so the archive does not belong to root.
// An empty string is returned if it can't be determined (os.Getuid() returns -1 on Windows).
func getArchiveOwner() string {
	uid := os.Getuid()
	if uid < 0 {
		return ""
	}

	return strconv.Itoa(uid) + ":" + strconv.Itoa(os.Getgid())
}
