package dockerutils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorDockerContainer struct {
	docker Docker
	name   string
}

func NewCommandExecutorDockerContainer() (c *CommandExecutorDockerContainer) {
	return new(CommandExecutorDockerContainer)
}

func (c *CommandExecutorDockerContainer) GetCommandExecutor() (commandExectuor commandexecutor.CommandExecutor, err error) {
	docker, err := c.GetDocker()
	if err != nil {
		return nil, err
	}

	commandExecutorDocker, ok := docker.(*CommandExecutorDocker)
	if !ok {
		typeString, err := datatypes.GetTypeName(docker)
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"Only implemented for CommandExecutorDocker but got '%s'",
			typeString,
		)
	}

	return commandExecutorDocker.GetCommandExecutor()
}

func (c *CommandExecutorDockerContainer) GetName() (name string, err error) {
	if len(c.name) <= 0 {
		return "", tracederrors.TracedError("name not set")
	}

	return c.name, nil
}

func (c *CommandExecutorDockerContainer) IsRunning(ctx context.Context) (isRunning bool, err error) {
	containerName, err := c.GetName()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf("docker inspect '%s' &> /dev/null && echo yes || echo no", containerName),
			},
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)
	if stdout == "yes" {
		return true, nil
	}
	if stdout == "no" {
		return false, nil
	}

	return false, tracederrors.TracedErrorf("Unexpected stdout to evaluate docker container running: '%s'", stdout)
}

func (c *CommandExecutorDockerContainer) Kill(ctx context.Context) (err error) {
	isRunning, err := c.IsRunning(ctx)
	if err != nil {
		return err
	}

	containerName, err := c.GetName()
	if err != nil {
		return err
	}

	if isRunning {
		logging.LogInfoByCtxf(ctx, "Going to kill running container '%s'.", containerName)

		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"docker", "kill", containerName},
			},
		)
		if err != nil {
			return err
		}

		sleepDuration := time.Second * 2
		logging.LogInfoByCtxf(
			ctx,
			"Wait %v until delete of docker container '%s' is settled to avoid race condition.",
			sleepDuration,
			containerName,
		)
		time.Sleep(sleepDuration)

		logging.LogChangedByCtxf(ctx, "Killed container '%s'", containerName)
	} else {
		logging.LogInfoByCtxf(ctx, "Container '%s' is already removed. Skip killing container.", containerName)
	}

	return nil
}

func (c *CommandExecutorDockerContainer) RunCommand(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, runOptions)
}

func (c *CommandExecutorDockerContainer) RunCommandAndGetStdoutAsString(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (stdout string, err error) {
	if runOptions == nil {
		return "", tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.RunCommandAndGetStdoutAsString(ctx, runOptions)
}

func (c *CommandExecutorDockerContainer) SetName(name string) (err error) {
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	c.name = name

	return nil
}

func (d *CommandExecutorDockerContainer) GetDocker() (docker Docker, err error) {
	if d.docker == nil {
		return nil, tracederrors.TracedError("docker is not set")
	}
	return d.docker, nil
}

func (d *CommandExecutorDockerContainer) SetDocker(docker Docker) (err error) {
	if docker == nil {
		return tracederrors.TracedErrorNil("docker")
	}

	d.docker = docker

	return nil
}
