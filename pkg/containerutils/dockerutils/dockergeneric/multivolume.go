package dockergeneric

import (
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// Directory inside the helper container the archive directory is mounted to.
	MultiVolumeBackupDir = "/backup"

	// Directory inside the helper container the archive content is built up and extracted.
	// Every volume is mounted to 'MultiVolumeRootDir/volumes/<volume-name>'.
	MultiVolumeRootDir = "/archive"

	// Printed by ReadManifestScript if the archive does not contain a manifest.
	NoManifestMarker = "ASCIICH_NO_MANIFEST"
)

// BackupVolumesScript creates the multi volume archive.
// The values are passed as positional parameters to avoid any shell quoting problems:
//
//	$0 = "sh"
//	$1 = temporary archive name (in MultiVolumeBackupDir)
//	$2 = owner "uid:gid" of the archive (may be empty)
//	$3 = final archive name (in MultiVolumeBackupDir)
//	$4 = content of the manifest.json
//
// The archive is written to the temporary name first and moved on success.
// The trap removes the temporary file if anything fails.
const BackupVolumesScript = `set -e
trap 'rm -f "` + MultiVolumeBackupDir + `/$1"' EXIT
printf '%s' "$4" > ` + MultiVolumeRootDir + `/manifest.json
tar czf "` + MultiVolumeBackupDir + `/$1" -C ` + MultiVolumeRootDir + ` .
if [ -n "$2" ]; then chown "$2" "` + MultiVolumeBackupDir + `/$1"; fi
mv "` + MultiVolumeBackupDir + `/$1" "` + MultiVolumeBackupDir + `/$3"`

// RestoreVolumesScript extracts the multi volume archive.
// Entries below 'volumes/<name>' end up in the mounted volumes.
// Everything else only lands in the throw away container file system.
//
//	$0 = "sh"
//	$1 = archive name (in MultiVolumeBackupDir)
const RestoreVolumesScript = `set -e
tar xzf "` + MultiVolumeBackupDir + `/$1" -C ` + MultiVolumeRootDir

// ReadManifestScript prints the manifest.json of the archive.
// If the archive has no manifest, NoManifestMarker is printed instead.
// A tar failure (e.g. not a gzip compressed archive) results in exit code 2.
//
//	$0 = "sh"
//	$1 = archive name (in MultiVolumeBackupDir)
const ReadManifestScript = `tar tzf "` + MultiVolumeBackupDir + `/$1" > /tmp/entries || exit 2
entry=$(grep -Ex '(\./)?manifest\.json' /tmp/entries | head -n 1)
if [ -z "$entry" ]; then echo "` + NoManifestMarker + `"; exit 0; fi
tar xzOf "` + MultiVolumeBackupDir + `/$1" "$entry"`

// GetMultiVolumeMountPoint returns the path inside the helper container
// where the volume has to be mounted.
func GetMultiVolumeMountPoint(volumeName string) (string, error) {
	err := CheckVolumeName(volumeName)
	if err != nil {
		return "", err
	}

	return MultiVolumeRootDir + "/" + GetVolumeDirectory(volumeName), nil
}

// SortedUniqueVolumeNames validates the volume names, removes duplicates and sorts them.
func SortedUniqueVolumeNames(volumeNames []string) ([]string, error) {
	if len(volumeNames) == 0 {
		return nil, tracederrors.TracedErrorf("No volume names specified.")
	}

	ret := slices.Clone(volumeNames)
	slices.Sort(ret)
	ret = slices.Compact(ret)

	for _, name := range ret {
		err := CheckVolumeName(name)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

// ParseManifestOutput parses the output of ReadManifestScript.
// Returns an error wrapping ErrNoManifest if the archive has no manifest.
func ParseManifestOutput(output string) (*Manifest, error) {
	output = strings.TrimSpace(output)

	if output == NoManifestMarker {
		return nil, tracederrors.TracedErrorf("%w", ErrNoManifest)
	}

	return ParseManifest([]byte(output))
}
