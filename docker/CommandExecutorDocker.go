package docker

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/containers"
	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/hosts"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorDocker struct {
	host hosts.Host
}

func GetCommandExecutorDocker(commandExecutor commandexecutor.CommandExecutor) (docker Docker, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	toReturn := NewCommandExecutorDocker()

	isRunningOnLocalhost, err := commandExecutor.IsRunningOnLocalhost()
	if err != nil {
		return nil, err
	}

	if !isRunningOnLocalhost {
		hostDescription, err := commandExecutor.GetHostDescription()
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"Not implemented for command executor running on '%s'.",
			hostDescription,
		)
	}

	host, err := hosts.GetLocalCommandExecutorHost()
	if err != nil {
		return nil, err
	}

	err = toReturn.SetHost(host)
	if err != nil {
		return nil, err
	}

	return toReturn, err
}

func GetCommandExecutorDockerOnHost(host hosts.Host) (docker Docker, err error) {
	if host == nil {
		return nil, tracederrors.TracedErrorNil("host")
	}

	toReturn := NewCommandExecutorDocker()

	err = toReturn.SetHost(host)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorDocker() (docker Docker, err error) {
	return GetCommandExecutorDocker(commandexecutor.Bash())
}

func MustGetCommandExecutorDocker(commandExecutor commandexecutor.CommandExecutor) (docker Docker) {
	docker, err := GetCommandExecutorDocker(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func MustGetCommandExecutorDockerOnHost(host hosts.Host) (docker Docker) {
	docker, err := GetCommandExecutorDockerOnHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func MustGetLocalCommandExecutorDocker() (docker Docker) {
	docker, err := GetLocalCommandExecutorDocker()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func MustGetcommandExecutorDocker(commandExecutor commandexecutor.CommandExecutor) (docker Docker) {
	docker, err := GetCommandExecutorDocker(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func NewCommandExecutorDocker() (c *CommandExecutorDocker) {
	return new(CommandExecutorDocker)
}

func (c *CommandExecutorDocker) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {
	host, err := c.GetHost()
	if err != nil {
		return nil, err
	}

	commandExecutorHost, ok := host.(*hosts.CommandExecutorHost)
	if !ok {
		typeString, err := datatypes.GetTypeName(host)
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"Only available for commandExecutorHost but got '%s'",
			typeString,
		)
	}

	return commandExecutorHost, nil
}

func (c *CommandExecutorDocker) GetContainerByName(containerName string) (dockerContainer containers.Container, err error) {
	if len(containerName) <= 0 {
		return nil, tracederrors.TracedError("containerName is empty string")
	}

	toReturn := NewCommandExecutorDockerContainer()
	err = toReturn.SetName(containerName)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetDocker(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorDocker) GetHost() (host hosts.Host, err error) {
	if c.host == nil {
		return nil, tracederrors.TracedError("host not set")
	}

	return c.host, nil
}

func (c *CommandExecutorDocker) GetHostDescription() (hostDescription string, err error) {
	host, err := c.GetHost()
	if err != nil {
		return "", err
	}

	return host.GetHostDescription()
}

func (c *CommandExecutorDocker) IsHostSet() (isSet bool) {
	return c.host != nil
}

func (c *CommandExecutorDocker) KillContainerByName(name string, verbose bool) (err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	container, err := c.GetContainerByName(name)
	if err != nil {
		return err
	}

	err = container.Kill(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorDocker) MustGetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorDocker) MustGetContainerByName(containerName string) (dockerContainer containers.Container) {
	dockerContainer, err := c.GetContainerByName(containerName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dockerContainer
}

func (c *CommandExecutorDocker) MustGetHost() (host hosts.Host) {
	host, err := c.GetHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return host
}

func (c *CommandExecutorDocker) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (c *CommandExecutorDocker) MustKillContainerByName(name string, verbose bool) {
	err := c.KillContainerByName(name, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDocker) MustRunCommand(runOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput) {
	commandOutput, err := c.RunCommand(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (c *CommandExecutorDocker) MustRunCommandAndGetStdoutAsString(runOptions *parameteroptions.RunCommandOptions) (stdout string) {
	stdout, err := c.RunCommandAndGetStdoutAsString(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorDocker) MustRunContainer(runOptions *DockerRunContainerOptions) (startedContainer containers.Container) {
	startedContainer, err := c.RunContainer(runOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return startedContainer
}

func (c *CommandExecutorDocker) MustSetHost(host hosts.Host) {
	err := c.SetHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDocker) RunCommand(runOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(runOptions)
}

func (c *CommandExecutorDocker) RunCommandAndGetStdoutAsString(runOptions *parameteroptions.RunCommandOptions) (stdout string, err error) {
	if runOptions == nil {
		return "", tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.RunCommandAndGetStdoutAsString(runOptions)
}

func (c *CommandExecutorDocker) RunContainer(runOptions *DockerRunContainerOptions) (startedContainer containers.Container, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedError("runOptions is nil")
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
		logging.LogInfof(
			"Going to start container '%s' using image '%s'.",
			containerName,
			imageName,
		)
	}

	err = c.KillContainerByName(containerName, runOptions.Verbose)
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
		logging.LogInfof("Going to start docker container using:\n%v", startCommand)
	}

	stdout, err := c.RunCommandAndGetStdoutAsString(
		&parameteroptions.RunCommandOptions{
			Command: startCommand,
			Verbose: runOptions.Verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	if runOptions.Verbose {
		logging.LogChangedf("Started container '%s':\n%s", containerName, stdout)
	}

	startedContainer, err = c.GetContainerByName(containerName)
	if err != nil {
		return nil, err
	}

	return startedContainer, nil
}

func (c *CommandExecutorDocker) SetHost(host hosts.Host) (err error) {
	if host == nil {
		return tracederrors.TracedError("host not set")
	}

	c.host = host

	return nil
}
