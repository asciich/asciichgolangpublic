package dockerutils

import "context"

// Backup docker volume on localhost to a tar archive.
func BackupVolume(ctx context.Context, volumeName string, archivePath string) error {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return err
	}

	return docker.BackupVolume(ctx, volumeName, archivePath)
}

// Restore docker volume on localhost from a tar archive.
func RestoreVolume(ctx context.Context, volumeName string, archivePath string) error {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return err
	}

	return docker.RestoreVolume(ctx, volumeName, archivePath)
}

// Backup multiple docker volumes on localhost into one multi volume tar archive.
func BackupVolumes(ctx context.Context, volumeNames []string, archivePath string) error {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return err
	}

	return docker.BackupVolumes(ctx, volumeNames, archivePath)
}

// Restore all docker volumes on localhost from a multi volume tar archive.
// The target volumes are taken from the manifest of the archive.
func RestoreVolumes(ctx context.Context, archivePath string) error {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return err
	}

	return docker.RestoreVolumes(ctx, archivePath)
}

// List docker volumes on localhost.
func ListVolumesNames(ctx context.Context) ([]string, error) {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return nil, err
	}

	return docker.ListVolumes(ctx)
}

// List the names of all files and directories inside the docker volume
// 'volumeName' on localhost. Hidden files are included and the result
// is sorted alphabetically.
func ListVolumeFileNames(ctx context.Context, volumeName string) ([]string, error) {
	docker, err := GetDockerOnLocalHost()
	if err != nil {
		return nil, err
	}

	return docker.ListVolumeFileNames(ctx, volumeName)
}
