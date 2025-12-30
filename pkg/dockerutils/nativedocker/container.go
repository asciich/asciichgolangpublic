package nativedocker

import (
	"context"

	"github.com/moby/moby/client"

	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Container struct {
	name string
}

func NewContainer(name string) (*Container, error) {
	ret := new(Container)
	err := ret.SetName(name)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Container) SetName(name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	c.name = name

	return nil
}

func (c *Container) GetName() (string, error) {
	if c.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return c.name, nil
}

func (c *Container) GetHostDescription() (string, error) {
	return NewDocker().GetHostDescription()
}

func (c *Container) Exists(ctx context.Context) (bool, error) {
	name, err := c.GetName()
	if err != nil {
		return false, err
	}

	_, err = c.inspect(ctx)
	if err != nil {
		if dockergeneric.IsErrorContainerNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Docker container '%s' does not exist.", name)
			return false, nil
		}

		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Docker container '%s' exists.", name)
	return true, nil
}

func (c *Container) inspect(ctx context.Context) (*client.ContainerInspectResult, error) {
	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	return new(Docker).inspect(ctx, name)
}

func (c *Container) IsRunning(ctx context.Context) (bool, error) {
	containerName, err := c.GetName()
	if err != nil {
		return false, err
	}

	inspect, err := c.inspect(ctx)
	if err != nil {
		if dockergeneric.IsErrorContainerNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Docker container '%s' does not exist and is therefore not running.", containerName)
			return false, nil
		}

		return false, err
	}

	isRunning := inspect.Container.State.Status == "running"

	if isRunning {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is running.", containerName)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is not running.", containerName)
	}

	return isRunning, nil
}

func (c *Container) Kill(ctx context.Context) error {
	containerName, err := c.GetName()
	if err != nil {
		return err
	}

	return NewDocker().KillContainerByName(ctx, containerName)
}

func (c *Container) Remove(ctx context.Context) error {
	name, err := c.GetName()
	if err != nil {
		return err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return err
	}

	exists, err := c.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return tracederrors.TracedErrorf("unable to create docker client: %w", err)
		}
		defer cli.Close()

		options := client.ContainerRemoveOptions{
			Force:         false,
			RemoveVolumes: false,
		}
		_, err = cli.ContainerRemove(ctx, name, options)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete container '%s' on host '%s': %w", name, hostDescription, err)
		}

		logging.LogChangedByCtxf(ctx, "Docker container '%s' removed on host '%s'.", name, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is already absent on host '%s'. Skip removal of container.", name, hostDescription)
	}

	return nil
}

func (c *Container) Run(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	containerName, err := c.GetName()
	if err != nil {
		return err
	}

	optionsToUse := options.GetDeepCopy()
	err = optionsToUse.SetName(containerName)
	if err != nil {
		return err
	}

	_, err = NewDocker().RunContainer(ctx, optionsToUse)
	if err != nil {
		return err
	}

	return nil
}
