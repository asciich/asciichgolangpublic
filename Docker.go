package asciichgolangpublic

import (
	"strings"
)

type DockerService struct {
	CommandExecutor
	host *Host
}

func Docker() (dockerService *DockerService) {
	dockerService = new(DockerService)
	return dockerService
}

func NewDockerService() (dockerService *DockerService) {
	dockerService = new(DockerService)
	return dockerService
}

func (d *DockerService) GetDockerContainerByName(containerName string) (dockerContainer *DockerContainer, err error) {
	if len(containerName) <= 0 {
		return nil, TracedError("containerName is empty string")
	}

	dockerContainer = NewDockerContainer()
	err = dockerContainer.SetName(containerName)
	if err != nil {
		return nil, err
	}

	err = dockerContainer.SetDockerService(d)
	if err != nil {
		return nil, err
	}

	return dockerContainer, nil
}

func (d *DockerService) GetHost() (host *Host, err error) {
	if d.host == nil {
		return nil, TracedError("host not set")
	}

	return d.host, nil
}


func (d *DockerService) IsHostSet() (isSet bool) {
	return d.host != nil
}

func (d *DockerService) KillContainerByName(name string, verbose bool) (err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return TracedError("name is empty string")
	}

	container, err := d.GetDockerContainerByName(name)
	if err != nil {
		return err
	}

	err = container.Kill(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (d *DockerService) MustGetDockerContainerByName(containerName string) (dockerContainer *DockerContainer) {
	dockerContainer, err := d.GetDockerContainerByName(containerName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dockerContainer
}

func (d *DockerService) MustGetHost() (host *Host) {
	host, err := d.GetHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return host
}

func (d *DockerService) MustKillContainerByName(name string, verbose bool) {
	err := d.KillContainerByName(name, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DockerService) MustRunContainer(runOptions *DockerRunContainerOptions) (startedContainer *DockerContainer) {
	startedContainer, err := d.RunContainer(runOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return startedContainer
}

func (d *DockerService) MustSetHost(host *Host) {
	err := d.SetHost(host)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DockerService) RunContainer(runOptions *DockerRunContainerOptions) (startedContainer *DockerContainer, err error) {
	if runOptions == nil {
		return nil, TracedError("runOptions is nil")
	}

	containerName, err := runOptions.GetName()
	if err != nil {
		return nil, err
	}

	imageName, err := runOptions.GetImageName()
	if err != nil {
		return nil, err
	}

	if runOptions.Verbose {
		LogInfof(
			"Going to start container '%s' using image '%s'.",
			containerName,
			imageName,
		)
	}

	err = d.KillContainerByName(containerName, runOptions.Verbose)
	if err != nil {
		return nil, err
	}

	startCommand := []string{
		"docker",
		"run",
	}

	if !runOptions.KeepStoppedContainer {
		startCommand = append(startCommand, "--rm")
	}

	startCommand = append(startCommand, "--detach", "--name", containerName)

	if runOptions.UseHostNet {
		startCommand = append(startCommand, "--net=host")
	}

	for _, port := range runOptions.Ports {
		startCommand = append(startCommand, "-p", port)
	}

	for _, mount := range runOptions.Mounts {
		startCommand = append(startCommand, "-v", mount)
	}

	startCommand = append(startCommand, imageName)

	startCommand = append(startCommand, runOptions.Command...)

	if runOptions.VerboseDockerRunCommand {
		LogInfof("Going to start docker container using:\n%v", startCommand)
	}

	stdout, err := d.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: startCommand,
			Verbose: runOptions.Verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	if runOptions.Verbose {
		LogChangedf("Started container '%s':\n%s", containerName, stdout)
	}

	startedContainer, err = d.GetDockerContainerByName(containerName)
	if err != nil {
		return nil, err
	}

	return startedContainer, nil
}

func (d *DockerService) SetHost(host *Host) (err error) {
	if host == nil {
		return TracedError("host not set")
	}

	d.host = host

	return nil
}
