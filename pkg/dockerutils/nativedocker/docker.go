package nativedocker

import (
	"context"
	"errors"

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

	image, err := options.GetImageName()
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
		Image: image,
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
