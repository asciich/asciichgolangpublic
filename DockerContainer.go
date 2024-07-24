package asciichgolangpublic

import (
	"fmt"
	"strings"
	"time"
)

type DockerContainer struct {
	CommandExecutor
	dockerService *DockerService
	name          string
}

func NewDockerContainer() (dockerContainer *DockerContainer) {
	dockerContainer = new(DockerContainer)
	return dockerContainer
}

func (c *DockerContainer) GetDockerService() (dockerService *DockerService, err error) {
	if c.dockerService == nil {
		return nil, TracedError("dockerService not set")
	}

	return c.dockerService, nil
}

func (c *DockerContainer) GetName() (name string, err error) {
	if len(c.name) <= 0 {
		return "", TracedError("name not set")
	}

	return c.name, nil
}


func (c *DockerContainer) IsRunning(verbose bool) (isRunning bool, err error) {
	containerName, err := c.GetName()
	if err != nil {
		return false, err
	}

	dockerService, err := c.GetDockerService()
	if err != nil {
		return false, err
	}

	stdout, err := dockerService.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
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

	return false, TracedErrorf("Unexpected stdout to evaluate docker container running: '%s'", stdout)
}

func (c *DockerContainer) Kill(verbose bool) (err error) {
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
			LogInfof("Going to kill running container '%s'.", containerName)
		}

		dockerService, err := c.GetDockerService()
		if err != nil {
			return err
		}

		_, err = dockerService.RunCommand(
			&RunCommandOptions{
				Command: []string{"docker", "kill", containerName},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		sleepDuration := time.Second * 2
		if verbose {
			LogInfof(
				"Wait %v until delete of docker container '%s' is settled to avoid race condition.",
				sleepDuration,
				containerName,
			)
		}
		time.Sleep(sleepDuration)

		if verbose {
			LogChangedf("Killed container '%s'", containerName)
		}
	} else {
		if verbose {
			LogInfof("Container '%s' is already removed. Skip killing container.", containerName)
		}
	}

	return nil
}

func (c *DockerContainer) SetDockerService(dockerService *DockerService) (err error) {
	if dockerService == nil {
		return TracedError("dockerService is nil")
	}

	c.dockerService = dockerService

	return nil
}

func (c *DockerContainer) SetName(name string) (err error) {
	if len(name) <= 0 {
		return TracedError("name is empty string")
	}

	c.name = name

	return nil
}

func (d *DockerContainer) MustGetDockerService() (dockerService *DockerService) {
	dockerService, err := d.GetDockerService()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dockerService
}

func (d *DockerContainer) MustGetName() (name string) {
	name, err := d.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (d *DockerContainer) MustIsRunning(verbose bool) (isRunning bool) {
	isRunning, err := d.IsRunning(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isRunning
}

func (d *DockerContainer) MustKill(verbose bool) {
	err := d.Kill(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DockerContainer) MustSetDockerService(dockerService *DockerService) {
	err := d.SetDockerService(dockerService)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DockerContainer) MustSetName(name string) {
	err := d.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}
