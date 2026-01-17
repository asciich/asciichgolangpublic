package nativedocker

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Container struct {
	commandexecutorgeneric.CommandExecutorBase
	name string
}

func NewContainer(name string) (*Container, error) {
	ret := new(Container)

	ret.SetParentCommandExecutorForBaseClass(ret)

	err := ret.SetName(name)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Container) GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor {
	return &Container{
		name: c.name,
	}
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
	dockerHostDescription, err := NewDocker().GetHostDescription()
	if err != nil {
		return "", err
	}

	name, err := c.GetName()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Docker container '%s' running on host '%s'.", name, dockerHostDescription), nil
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

func (c *Container) Remove(ctx context.Context, options *dockeroptions.RemoveOptions) error {
	if options == nil {
		options = new(dockeroptions.RemoveOptions)
	}

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
		force := options.Force

		if force {
			logging.LogInfoByCtxf(ctx, "Remove docker container '%s' started.", name)
		} else {
			logging.LogInfoByCtxf(ctx, "Force remove docker container '%s' started.", name)
		}

		cli, err := client.New(client.FromEnv)
		if err != nil {
			return tracederrors.TracedErrorf("unable to create docker client: %w", err)
		}
		defer cli.Close()

		clientOptions := client.ContainerRemoveOptions{
			Force:         force,
			RemoveVolumes: false,
		}
		_, err = cli.ContainerRemove(ctx, name, clientOptions)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete container '%s' on host '%s': %w", name, hostDescription, err)
		}

		logging.LogChangedByCtxf(ctx, "Docker container '%s' removed on host '%s'.", name, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is already absent on host '%s'. Skip removal of container.", name, hostDescription)
	}

	return nil
}

func (c *Container) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	output := new(commandoutput.CommandOutput)

	cmdJoined, err := options.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' in docker container '%s' started.", cmdJoined, name)

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	cmd, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	isStdinSet := len(options.StdinString) > 0

	exec, err := cli.ExecCreate(ctx, name, client.ExecCreateOptions{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  isStdinSet,
		Cmd:          cmd,
		User:         options.RunAsUser,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to exec crate to RunCommand in container '%s': %w", name, err)
	}

	execId := exec.ID

	attach, err := cli.ExecAttach(ctx, execId, client.ExecAttachOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to exec attach for id '%s' on container '%s': %w", execId, name, err)
	}
	defer attach.HijackedResponse.Close()

	if isStdinSet {
		_, err = attach.HijackedResponse.Conn.Write([]byte(options.StdinString))
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to write to the stdin of container '%s' with exect id '%s': %w", name, execId, err)
		}

		if cw, ok := attach.Conn.(interface{ CloseWrite() error }); ok {
			err := cw.CloseWrite()
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to close stdin of container '%s' with exect id '%s': %w", name, execId, err)
			}
		} else {
			return nil, tracederrors.TracedErrorf("Unable to close stdin of container '%s' with exect id '%s': %w", name, execId, err)
		}
	}

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, attach.HijackedResponse.Reader)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read stdout and stderr of execid '%s' on container '%s': %w", execId, name, err)
	}

	err = output.SetStdout(stdout.Bytes())
	if err != nil {
		return nil, err
	}

	err = output.SetStderr(stderr.Bytes())
	if err != nil {
		return nil, err
	}

	for range 3 {
		inspect, err := cli.ExecInspect(ctx, execId, client.ExecInspectOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to exec inspect for exec id='%s' and container '%s': %w", execId, name, err)
		}

		if inspect.ID == "" || inspect.ContainerID == "" {
			// There is a race condition returing an empty inspect while the container is still closing
			time.Sleep(time.Millisecond * 100)
			continue
		}

		err = output.SetReturnCode(inspect.ExitCode)
		if err != nil {
			return nil, err
		}

		break
	}

	if !output.IsReturnCodeSet() {
		return nil, tracederrors.TracedError("Unable to set return code for docker exec")
	}

	if !options.AllowAllExitCodes {
		if !output.IsExitSuccess() {
			exitCode, err := output.GetReturnCode()
			if err != nil {
				return nil, err
			}

			stderr, err := output.GetStderrAsString()
			if err != nil {
				return nil, err
			}

			return nil, tracederrors.TracedErrorf("Run command '%s' in docker container '%s' failed. Exit code is: %d, stderr is\n%s", cmdJoined, name, exitCode, stderr)
		}
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' in docker container '%s' finished.", cmdJoined, name)

	return output, err
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
