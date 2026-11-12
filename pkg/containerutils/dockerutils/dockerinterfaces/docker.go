package dockerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
)

type Docker interface {
	ContainerExists(ctx context.Context, name string) (bool, error)

	GetDeepCopyAsDocker() Docker

	GetContainerByName(name string) (containerinterfaces.Container, error)
	GetImageByName(name string) (containerinterfaces.Image, error)
	GetHostDescription() (string, error)

	ImageExists(ctx context.Context, name string) (bool, error)

	KillContainerByName(ctx context.Context, name string) error

	ListContainers(ctx context.Context) ([]containerinterfaces.Container, error)
	ListContainerNames(ctx context.Context) ([]string, error)
	ListVolumes(ctx context.Context) ([]string, error)

	BackupVolume(ctx context.Context, volumeName string, archivePath string) error
	ListVolumeFileNames(ctx context.Context, volumeName string) ([]string, error)
	RestoreVolume(ctx context.Context, volumeName string, archivePath string) error

	// BackupVolumes stores all given volumes in one gzip compressed multi volume archive.
	BackupVolumes(ctx context.Context, volumeNames []string, archivePath string) error
	// RestoreVolumes restores all volumes of a multi volume archive. The target volumes are taken from the manifest of the archive.
	RestoreVolumes(ctx context.Context, archivePath string) error

	PullImage(ctx context.Context, imageName string) (containerinterfaces.Image, error)

	RemoveImage(ctx context.Context, imageName string, options *dockeroptions.RemoveOptions) error
	RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error)
	BuildImage(ctx context.Context, options *dockeroptions.BuildContainerOptions) (containerinterfaces.Image, error)

	RemoveContainer(ctx context.Context, name string, options *dockeroptions.RemoveOptions) error

	RunCommandInTemporaryContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (*commandoutput.CommandOutput, error)
}
