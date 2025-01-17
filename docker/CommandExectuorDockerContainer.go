package docker

import (
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic"
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

func (c *CommandExecutorDockerContainer) GetCommandExecutor() (commandExectuor asciichgolangpublic.CommandExecutor, err error) {
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

func (c *CommandExecutorDockerContainer) IsRunning(verbose bool) (isRunning bool, err error) {
	containerName, err := c.GetName()
	if err != nil {
		return false, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf("docker inspect '%s' &> /dev/null && echo yes || echo no", containerName),
			},
			Verbose: verbose,
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

func (c *CommandExecutorDockerContainer) Kill(verbose bool) (err error) {
	isRunning, err := c.IsRunning(verbose)
	if err != nil {
		return err
	}

	containerName, err := c.GetName()
	if err != nil {
		return err
	}

	if isRunning {
		if verbose {
			logging.LogInfof("Going to kill running container '%s'.", containerName)
		}

		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&parameteroptions.RunCommandOptions{
				Command: []string{"docker", "kill", containerName},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		sleepDuration := time.Second * 2
		if verbose {
			logging.LogInfof(
				"Wait %v until delete of docker container '%s' is settled to avoid race condition.",
				sleepDuration,
				containerName,
			)
		}
		time.Sleep(sleepDuration)

		if verbose {
			logging.LogChangedf("Killed container '%s'", containerName)
		}
	} else {
		if verbose {
			logging.LogInfof("Container '%s' is already removed. Skip killing container.", containerName)
		}
	}

	return nil
}

func (c *CommandExecutorDockerContainer) MustGetCommandExecutor() (commandExectuor asciichgolangpublic.CommandExecutor) {
	commandExectuor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExectuor
}

func (c *CommandExecutorDockerContainer) MustRunCommand(runOptions *parameteroptions.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput) {
	commandOutput, err := c.RunCommand(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (c *CommandExecutorDockerContainer) MustRunCommandAndGetStdoutAsString(runOptions *parameteroptions.RunCommandOptions) (stdout string) {
	stdout, err := c.RunCommandAndGetStdoutAsString(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorDockerContainer) RunCommand(runOptions *parameteroptions.RunCommandOptions) (commandOutput *asciichgolangpublic.CommandOutput, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(runOptions)
}

func (c *CommandExecutorDockerContainer) RunCommandAndGetStdoutAsString(runOptions *parameteroptions.RunCommandOptions) (stdout string, err error) {
	if runOptions == nil {
		return "", tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.RunCommandAndGetStdoutAsString(runOptions)
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

func (d *CommandExecutorDockerContainer) MustGetDocker() (docker Docker) {
	docker, err := d.GetDocker()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func (d *CommandExecutorDockerContainer) MustGetName() (name string) {
	name, err := d.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (d *CommandExecutorDockerContainer) MustIsRunning(verbose bool) (isRunning bool) {
	isRunning, err := d.IsRunning(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunning
}

func (d *CommandExecutorDockerContainer) MustKill(verbose bool) {
	err := d.Kill(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *CommandExecutorDockerContainer) MustSetDocker(docker Docker) {
	err := d.SetDocker(docker)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *CommandExecutorDockerContainer) MustSetName(name string) {
	err := d.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (d *CommandExecutorDockerContainer) SetDocker(docker Docker) (err error) {
	if docker == nil {
		return tracederrors.TracedErrorNil("docker")
	}

	d.docker = docker

	return nil
}
