package nativedocker

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Docker struct {
}

func NewDocker() dockerinterfaces.Docker {
	return new(Docker)
}

func (d *Docker) GetContainerByName(name string) (containerinterfaces.Container, error) {
	return NewContainer(name)
}

func (d *Docker) GetHostDescription() (string, error) {
	return "localhost", nil
}

func (d *Docker) ListContainers(ctx context.Context) ([]containerinterfaces.Container, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (d *Docker) inspect(ctx context.Context, containerName string) (*client.ContainerInspectResult, error) {
	if containerName == "" {
		return nil, tracederrors.TracedErrorEmptyString("containerName")
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	inspect, err := cli.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
	if err != nil {
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, tracederrors.TracedErrorf("Container inspect for '%s' failed: %w: %w", containerName, dockergeneric.ErrDockerContainerNotFound, err)
		}

		return nil, tracederrors.TracedErrorf("Container inspect for '%s' failed: %w", containerName, err)
	}

	return &inspect, nil
}

func (d *Docker) GetContainerId(ctx context.Context, containerName string) (string, error) {
	if containerName == "" {
		return "", tracederrors.TracedErrorEmptyString("containerName")
	}

	inspect, err := d.inspect(ctx, containerName)
	if err != nil {
		return "", err
	}

	containerId := inspect.Container.ID

	logging.LogInfoByCtxf(ctx, "Container '%s' has container ID '%s'.", containerName, containerId)

	return containerId, nil
}

func (d *Docker) RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	name, err := options.GetName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run docker container '%s' started.", name)

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	_, err = d.PullImage(ctx, imageName)
	if err != nil {
		return nil, err
	}

	command, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	autoremove := !options.KeepStoppedContainer

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create docker client: %w", err)
	}

	createResult, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:  name,
		Image: imageName,
		Config: &container.Config{
			Cmd: command,
		},
		HostConfig: &container.HostConfig{
			AutoRemove: autoremove,
		},
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create container '%s': %w", name, err)
	}

	_, err = cli.ContainerStart(ctx, createResult.ID, client.ContainerStartOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to start container '%s': %w", name, err)
	}

	logging.LogInfoByCtxf(ctx, "Run docker container '%s' finished.", name)

	return d.GetContainerByName(name)
}

func (d *Docker) ListContainerNames(ctx context.Context) ([]string, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	nativeList, err := cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list containers")
	}

	names := []string{}
	for _, container := range nativeList.Items {
		for _, name := range container.Names {
			names = append(names, name)
		}
	}

	return names, nil
}

func (d *Docker) KillContainerByName(ctx context.Context, name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	_, err = cli.ContainerKill(ctx, name, client.ContainerKillOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("Unable to kill container '%s': %w", name, err)
	}

	logging.LogChangedByCtxf(ctx, "Killed container '%s'.", name)

	return err
}

func (d *Docker) imageInspect(ctx context.Context, imageName string) (*client.ImageInspectResult, error) {
	if imageName == "" {
		return nil, tracederrors.TracedErrorEmptyString("imageName")
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	inspect, err := cli.ImageInspect(ctx, imageName, client.ImageInspectWithManifests(false))
	if err != nil {
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, tracederrors.TracedErrorf("Docker image inspect for '%s' failed: %w: %w", imageName, dockergeneric.ErrDockerImageNotFound, err)
		}
		return nil, tracederrors.TracedErrorf("ImageInspect failed for image '%s': %w", imageName, err)
	}

	return &inspect, nil
}

func (d *Docker) ImageExists(ctx context.Context, imageName string) (bool, error) {
	if imageName == "" {
		return false, tracederrors.TracedErrorEmptyString("imageName")
	}

	_, err := d.imageInspect(ctx, imageName)
	if err != nil {
		if dockergeneric.IsErrorImageNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Docker image '%s' does not exist.", imageName)
			return false, nil
		}

		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Docker image '%s' exists.", imageName)
	return true, nil
}

func (d *Docker) GetImageByName(imageName string) (containerinterfaces.Image, error) {
	image := NewImage()

	err := image.SetName(imageName)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (d *Docker) PullImage(ctx context.Context, imageName string) (containerinterfaces.Image, error) {
	if imageName == "" {
		return nil, tracederrors.TracedErrorEmptyString("imageName")
	}

	exists, err := d.ImageExists(ctx, imageName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' already exists, skip pull.", imageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pull docker image '%s' started.", imageName)

		cli, err := client.New(client.FromEnv)
		if err != nil {
			return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
		}
		defer cli.Close()

		out, err := cli.ImagePull(ctx, imageName, client.ImagePullOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to pull image '%s': %w", imageName, err)
		}
		defer out.Close()
		io.Copy(os.Stderr, out)

		logging.LogChangedByCtxf(ctx, "Pulled docker image '%s'.", imageName)
	}

	return d.GetImageByName(imageName)
}

func (d *Docker) RemoveImage(ctx context.Context, imageName string) error {
	if imageName == "" {
		return tracederrors.TracedErrorEmptyString("imageName")
	}

	exists, err := d.ImageExists(ctx, imageName)
	if err != nil {
		return err
	}

	if exists {
		cli, err := client.New(client.FromEnv)
		if err != nil {
			return tracederrors.TracedErrorf("unable to create docker client: %w", err)
		}
		defer cli.Close()

		_, err = cli.ImageRemove(ctx, imageName, client.ImageRemoveOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Unable to remove image '%s': %w", imageName, err)
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' is already absents. Skip remove.", imageName)
	}

	return nil
}

func (d *Docker) ContainerExists(ctx context.Context, containerName string) (bool, error) {
	if containerName == "" {
		return false, tracederrors.TracedErrorEmptyString("containerName")
	}

	container, err := d.GetContainerByName(containerName)
	if err != nil {
		return false, err
	}

	return container.Exists(ctx)
}

func (d *Docker) RemoveContainer(ctx context.Context, containerName string, options *dockeroptions.RemoveOptions) error {
	if containerName == "" {
		return tracederrors.TracedErrorEmptyString("containerName")
	}

	container, err := d.GetContainerByName(containerName)
	if err != nil {
		return err
	}

	return container.Remove(ctx, options)
}
