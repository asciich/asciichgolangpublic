package dockerutils

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/hosts"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorDocker struct {
	host hosts.Host
}

func GetCommandExecutorDocker(commandExecutor commandexecutorinterfaces.CommandExecutor) (docker dockerinterfaces.Docker, err error) {
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

func GetCommandExecutorDockerOnHost(host hosts.Host) (docker dockerinterfaces.Docker, err error) {
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

func GetLocalCommandExecutorDocker() (docker dockerinterfaces.Docker, err error) {
	return GetCommandExecutorDocker(commandexecutorbashoo.Bash())
}

func MustGetCommandExecutorDocker(commandExecutor commandexecutorinterfaces.CommandExecutor) (docker dockerinterfaces.Docker) {
	docker, err := GetCommandExecutorDocker(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func MustGetCommandExecutorDockerOnHost(host hosts.Host) (docker dockerinterfaces.Docker) {
	docker, err := GetCommandExecutorDockerOnHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func MustGetLocalCommandExecutorDocker() (docker dockerinterfaces.Docker) {
	docker, err := GetLocalCommandExecutorDocker()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func MustGetcommandExecutorDocker(commandExecutor commandexecutorinterfaces.CommandExecutor) (docker dockerinterfaces.Docker) {
	docker, err := GetCommandExecutorDocker(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return docker
}

func NewCommandExecutorDocker() (c *CommandExecutorDocker) {
	return new(CommandExecutorDocker)
}

func (c *CommandExecutorDocker) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
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

func (c *CommandExecutorDocker) GetContainerById(id string) (containerinterfaces.Container, error) {
	if id == "" {
		return nil, tracederrors.TracedErrorEmptyString("id")
	}

	toReturn := NewCommandExecutorDockerContainer()
	err := toReturn.SetId(id)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetDocker(c)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (c *CommandExecutorDocker) GetContainerByName(containerName string) (dockerContainer containerinterfaces.Container, err error) {
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

func (c *CommandExecutorDocker) KillContainerByName(ctx context.Context, name string) (err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	container, err := c.GetContainerByName(name)
	if err != nil {
		return err
	}

	err = container.Kill(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorDocker) RunCommand(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
	if runOptions == nil {
		return nil, tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommand(ctx, runOptions)
}

func (c *CommandExecutorDocker) RunCommandAndGetStdoutAsString(ctx context.Context, runOptions *parameteroptions.RunCommandOptions) (stdout string, err error) {
	if runOptions == nil {
		return "", tracederrors.TracedErrorNil("runOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.RunCommandAndGetStdoutAsString(ctx, runOptions)
}

func (c *CommandExecutorDocker) RunContainer(ctx context.Context, runOptions *DockerRunContainerOptions) (startedContainer containerinterfaces.Container, err error) {
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

	err = c.KillContainerByName(ctx, containerName)
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
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: startCommand,
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

func (c *CommandExecutorDocker) ListContainers(ctx context.Context) ([]containerinterfaces.Container, error) {
	executor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return nil, err
	}

	output, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"docker", "ps", "-a", "--no-trunc", "--format=json"},
	})
	if err != nil {
		return nil, err
	}

	type OutputEntry struct {
		Names string `json:"Names"`
		Id    string `json:"ID"`
	}

	parsed := []*OutputEntry{}

	for _, line := range stringsutils.SplitLines(output, true) {
		toAdd := new(OutputEntry)

		err = json.Unmarshal([]byte(line), toAdd)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to parse docker ps output: %w", err)
		}

		parsed = append(parsed, toAdd)
	}

	list := []containerinterfaces.Container{}
	for _, entry := range parsed {
		toAdd := NewCommandExecutorDockerContainer()

		err = toAdd.SetId(entry.Id)
		if err != nil {
			return nil, err
		}

		err = toAdd.SetCachedName(entry.Names)
		if err != nil {
			return nil, err
		}

		list = append(list, toAdd)
	}

	logging.LogInfoByCtxf(ctx, "Listed '%d' containers on host '%s'", len(list), hostDescription)

	return list, nil
}

func (c *CommandExecutorDocker) ListContainerNames(ctx context.Context) ([]string, error) {
	containers, err := c.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, c := range containers {
		cec, ok := c.(*CommandExecutorDockerContainer)
		if !ok {
			return nil, tracederrors.TracedErrorf("Unsupported type to get container name: %s", reflect.TypeOf(c))
		}

		name, err := cec.GetCachedName()
		if err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, nil
}
