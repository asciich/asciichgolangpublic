package commandexecutordocker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorDockerContainer struct {
	commandexecutorgeneric.CommandExecutorBase

	docker dockerinterfaces.Docker
	name   string
	id     string

	// caching
	cachedName string
}

func NewCommandExecutorDockerContainer() (c *CommandExecutorDockerContainer) {
	return new(CommandExecutorDockerContainer)
}

func (c *CommandExecutorDockerContainer) SetCachedName(cachedName string) (err error) {
	c.cachedName = cachedName
	return nil
}

func (c *CommandExecutorDockerContainer) GetCommandExecutor() (commandExectuor commandexecutorinterfaces.CommandExecutor, err error) {
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

func (c *CommandExecutorDockerContainer) GetContainerStateStatus(ctx context.Context) (string, error) {
	containerName, err := c.GetName()
	if err != nil {
		return "", err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	var status string
	var stderr string
	for i := range 10 {
		output, err := commandExecutor.RunCommand(
			contextutils.WithSilent(ctx),
			&parameteroptions.RunCommandOptions{
				Command:           []string{"docker", "inspect", "--format", "{{.State.Status}}", containerName},
				AllowAllExitCodes: true,
			},
		)
		if err != nil {
			return "", err
		}

		if !output.IsExitSuccess() {
			stderr, err = output.GetStderrAsString()
			if err != nil {
				return "", err
			}

			// Important to ignore the case here:
			// Archlinux docker-ce has "error" in lowercase letters.
			// The docker used in the build system on Github has "Error" with a starting uppercase letter.
			if stringsutils.ContainsIgnoreCase(stderr, "error: no such object:") {
				return "", tracederrors.TracedErrorf("Unable to get docker container '%s' .State.Status: %w", containerName, dockergeneric.ErrDockerContainerNotFound)
			}
		}

		stdout, err := output.GetStdoutAsString()
		if err != nil {
			return "", err
		}

		status = strings.TrimSpace(stdout)
		if status == "" {
			retryDelay := time.Millisecond * 100 * time.Duration(i+1)
			logging.LogInfoByCtxf(ctx, "Empty .State.Status for docker container '%s' received. This indicates a race condition during startup. Retry in %v.", containerName, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		break
	}

	if status == "" {
		return "", tracederrors.TracedErrorf("Failed to evaluate .State.Status for docker container '%s'. status is empty string after evaluation. Stderr of last try is '%s'.", containerName, stderr)
	}

	logging.LogInfoByCtxf(ctx, "Container '%s' has .State.Status='%s'.", containerName, status)

	return status, nil
}

func (c *CommandExecutorDockerContainer) IsRunning(ctx context.Context) (bool, error) {
	containerName, err := c.GetName()
	if err != nil {
		return false, err
	}

	stateStatus, err := c.GetContainerStateStatus(ctx)
	if err != nil {
		if dockergeneric.IsErrorContainerNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Docker container '%s' does not exist and is therefore not running.", containerName)
			return false, nil
		}
		return false, err
	}

	isRunning := stateStatus == "running"
	if isRunning {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is running (container .State.Status is '%s').", containerName, stateStatus)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is not running (container .State.Status is '%s')..", containerName, stateStatus)
	}

	return isRunning, nil
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

func (c *CommandExecutorDockerContainer) RunCommand(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	command, err := runOptions.GetCommand()
	if err != nil {
		return nil, err
	}

	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	optionsToUse := runOptions.GetDeepCopy()
	newCommand := []string{"docker", "exec", name}
	newCommand = append(newCommand, command...)
	err = optionsToUse.SetCommand(newCommand)
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, optionsToUse)
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

func (c *CommandExecutorDockerContainer) SetId(id string) (err error) {
	if len(id) <= 0 {
		return tracederrors.TracedError("id is empty string")
	}

	c.id = id

	return nil
}

func (c *CommandExecutorDockerContainer) GetDocker() (docker dockerinterfaces.Docker, err error) {
	if c.docker == nil {
		return nil, tracederrors.TracedError("docker is not set")
	}
	return c.docker, nil
}

func (c *CommandExecutorDockerContainer) GetCachedName() (string, error) {
	if c.cachedName == "" {
		name, err := c.GetName()
		if err != nil {
			return "", err
		}

		c.cachedName = name
	}

	return c.cachedName, nil
}

func (c *CommandExecutorDockerContainer) SetDocker(docker dockerinterfaces.Docker) (err error) {
	if docker == nil {
		return tracederrors.TracedErrorNil("docker")
	}

	c.docker = docker

	return nil
}

func (c *CommandExecutorDockerContainer) Remove(ctx context.Context, options *dockeroptions.RemoveOptions) error {
	if options == nil {
		options = new(dockeroptions.RemoveOptions)
	}

	name, err := c.GetName()
	if err != nil {
		return err
	}

	force := options.Force
	if force {
		logging.LogInfoByCtxf(ctx, "Remove docker container '%s' started.", name)
	} else {
		logging.LogInfoByCtxf(ctx, "Force remove docker container '%s' started.", name)
	}

	exists, err := c.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		commandExectuor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		cmd := []string{"docker", "rm"}

		if force {
			cmd = append(cmd, "--force")
		}

		cmd = append(cmd, name)

		_, err = commandExectuor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: cmd,
		})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Docker container '%s' removed.", name)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' does not exists. Skip removal.", name)
	}

	logging.LogInfoByCtxf(ctx, "Remove docker container '%s' finished.", name)

	return nil
}

func (c *CommandExecutorDockerContainer) Run(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	name, err := c.GetName()
	if err != nil {
		return err
	}

	optionsToUse := options.GetDeepCopy()

	err = optionsToUse.SetName(name)
	if err != nil {
		return err
	}

	docker, err := c.GetDocker()
	if err != nil {
		return err
	}

	commandExecutorDocker, ok := docker.(*CommandExecutorDocker)
	if !ok {
		typeString, err := datatypes.GetTypeName(docker)
		if err != nil {
			return err
		}

		return tracederrors.TracedErrorf(
			"Only implemented for CommandExecutorDocker but got '%s'",
			typeString,
		)
	}

	_, err = commandExecutorDocker.RunContainer(ctx, optionsToUse)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorDockerContainer) Exists(ctx context.Context) (bool, error) {
	name, err := c.GetName()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	output, err := commandExecutor.RunCommand(
		contextutils.WithSilent(ctx),
		&parameteroptions.RunCommandOptions{
			Command:           []string{"docker", "inspect", name},
			AllowAllExitCodes: true,
		},
	)
	if err != nil {
		return false, err
	}

	exists := output.IsExitSuccess()
	if exists {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' exists.", name)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' does not exist", name)
	}

	return exists, nil
}

func (c *CommandExecutorDockerContainer) GetHostDescription() (string, error) {
	docker, err := c.GetDocker()
	if err != nil {
		return "", err
	}

	dockerHostDescription, err := docker.GetHostDescription()
	if err != nil {
		return "", err
	}

	name, err := c.GetName()
	if err != nil {
		return "", err
	}

	hostDescription := fmt.Sprintf("Docker container '%s' running on host '%s'.", name, dockerHostDescription)

	return hostDescription, nil
}
