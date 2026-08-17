package commandexecutordocker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhost"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorDocker struct {
	commandexecutorgeneric.CommandExecutorBase
	host hostsutilsinterfaces.Host
}

func GetCommandExecutorDocker(commandExecutor commandexecutorinterfaces.CommandExecutor) (docker dockerinterfaces.Docker, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	ret := NewCommandExecutorDocker()

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

	host, err := hostsutils.GetLocalCommandExecutorHost()
	if err != nil {
		return nil, err
	}

	err = ret.SetHost(host)
	if err != nil {
		return nil, err
	}

	return ret, err
}

func GetCommandExecutorDockerOnHost(host hostsutilsinterfaces.Host) (docker dockerinterfaces.Docker, err error) {
	if host == nil {
		return nil, tracederrors.TracedErrorNil("host")
	}

	ret := NewCommandExecutorDocker()

	err = ret.SetHost(host)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetLocalCommandExecutorDocker() (docker dockerinterfaces.Docker, err error) {
	return GetCommandExecutorDocker(commandexecutorbashoo.Bash())
}

func NewCommandExecutorDocker() (c *CommandExecutorDocker) {
	return new(CommandExecutorDocker)
}

func (c *CommandExecutorDocker) GetDeepCopy() *CommandExecutorDocker {
	ret := NewCommandExecutorDocker()

	if c.host != nil {
		panic("use c.commandExecutor instead of c.host")
	}

	return ret
}

func (c *CommandExecutorDocker) GetDeepCopyAsCommandExecutor() commandexecutorinterfaces.CommandExecutor {
	return c.GetDeepCopy()
}

func (c *CommandExecutorDocker) GetDeepCopyAsDocker() dockerinterfaces.Docker {
	return c.GetDeepCopy()
}

func (c *CommandExecutorDocker) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	host, err := c.GetHost()
	if err != nil {
		return nil, err
	}

	commandExecutorHost, ok := host.(*commandexecutorhost.CommandExecutorHost)
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

	ret := NewCommandExecutorDockerContainer()
	err := ret.SetId(id)
	if err != nil {
		return nil, err
	}

	err = ret.SetDocker(c)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *CommandExecutorDocker) GetContainerByName(containerName string) (dockerContainer containerinterfaces.Container, err error) {
	if len(containerName) <= 0 {
		return nil, tracederrors.TracedError("containerName is empty string")
	}

	ret := NewCommandExecutorDockerContainer()
	err = ret.SetName(containerName)
	if err != nil {
		return nil, err
	}

	err = ret.SetDocker(c)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *CommandExecutorDocker) GetHost() (host hostsutilsinterfaces.Host, err error) {
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

func (c *CommandExecutorDocker) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommandAndGetStdoutAsIoReadCloser(ctx, options)
}

func (c *CommandExecutorDocker) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor.RunCommandAndGetStdinAsIoWriteCloser(ctx, options)
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

// appendEntryPointToCommand appends the --entrypoint flag to a docker command.
// It distinguishes between:
//   - nil: EntryPoint not set, don't add --entrypoint flag (use image default)
//   - empty slice ([]string{}): Explicitly overwrite entrypoint to empty string
//   - non-empty slice: Use the first element as entrypoint
func appendEntryPointToCommand(command []string, entryPoint []string) []string {
	if entryPoint == nil {
		// Not set, don't modify the entrypoint
		return command
	}

	if len(entryPoint) == 0 {
		// Explicitly overwrite entrypoint to empty
		command = append(command, "--entrypoint", "")
	} else {
		// Use specified entrypoint
		command = append(command, "--entrypoint", entryPoint[0])
		// Additional entrypoint args are not supported by docker --entrypoint flag,
		// they should be part of Command instead.
	}

	return command
}

func (c *CommandExecutorDocker) RunContainer(ctx context.Context, runOptions *dockeroptions.DockerRunContainerOptions) (startedContainer containerinterfaces.Container, err error) {
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

	logging.LogInfoByCtxf(ctx,
		"Going to start container '%s' using image '%s'.",
		containerName,
		imageName,
	)

	_, err = c.PullImage(ctx, imageName)
	if err != nil {
		return nil, err
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

	for envName, envValue := range runOptions.AdditionalEnvVars {
		startCommand = append(startCommand, "-e", fmt.Sprintf("%s=%s", envName, envValue))
	}

	for _, port := range runOptions.Ports {
		startCommand = append(startCommand, "-p", port)
	}

	for _, mount := range runOptions.Mounts {
		startCommand = append(startCommand, "-v", mount)
	}

	// Add entrypoint if specified (nil = not set, empty = overwrite to empty, non-empty = use value)
	startCommand = appendEntryPointToCommand(startCommand, runOptions.EntryPoint)

	startCommand = append(startCommand, imageName)

	startCommand = append(startCommand, runOptions.Command...)

	logging.LogInfoByCtxf(ctx, "Going to start docker container using:\n%v", startCommand)

	stdout, err := c.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: startCommand,
		},
	)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Started container '%s':\n%s", containerName, stdout)

	startedContainer, err = c.GetContainerByName(containerName)
	if err != nil {
		return nil, err
	}

	return startedContainer, nil
}

func (c *CommandExecutorDocker) SetHost(host hostsutilsinterfaces.Host) (err error) {
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

func (c *CommandExecutorDocker) GetImageByName(imageName string) (containerinterfaces.Image, error) {
	if imageName == "" {
		return nil, tracederrors.TracedErrorEmptyString("imageName")
	}

	image := new(Image)

	err := image.SetName(imageName)
	if err != nil {
		return nil, err
	}

	err = image.SetDocker(c)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (c *CommandExecutorDocker) PullImage(ctx context.Context, imageName string) (containerinterfaces.Image, error) {
	if imageName == "" {
		return nil, tracederrors.TracedErrorEmptyString("imageName")
	}

	exists, err := c.ImageExists(ctx, imageName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' is already present. Skip pull", imageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pull docker image '%s' started.", imageName)

		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return nil, err
		}

		_, err = commandExecutor.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{"docker", "pull", imageName},
			},
		)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "Pulled docker image '%s'.", imageName)
	}

	return c.GetImageByName(imageName)
}

func (c *CommandExecutorDocker) ImageExists(ctx context.Context, imageName string) (bool, error) {
	if imageName == "" {
		return false, tracederrors.TracedErrorEmptyString("imageName")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	output, err := commandExecutor.RunCommand(
		contextutils.WithSilent(ctx),
		&parameteroptions.RunCommandOptions{
			Command:           []string{"docker", "image", "inspect", imageName},
			AllowAllExitCodes: true,
		},
	)
	if err != nil {
		return false, err
	}

	if output.IsExitSuccess() {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' exists.", imageName)
		return true, nil
	}

	stderr, err := output.GetStderrAsString()
	if err != nil {
		return false, err
	}

	if strings.Contains(stderr, "No such image:") {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' does not exist.", imageName)
		return false, err
	}

	return false, tracederrors.TracedErrorf("Unknown docker output on stderr: %w", err)
}

func (c *CommandExecutorDocker) RemoveImage(ctx context.Context, imageName string, options *dockeroptions.RemoveOptions) error {
	if imageName == "" {
		return tracederrors.TracedErrorEmptyString("imageName")
	}

	exists, err := c.ImageExists(ctx, imageName)
	if err != nil {
		return err
	}

	if exists {
		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		force := false
		if options != nil {
			force = options.Force
		}

		command := []string{"docker", "rmi"}
		if force {
			command = append(command, "--force")
		}
		command = append(command, imageName)

		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Docker image '%s' removed.", imageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' is already absent. Skip remove.", imageName)
	}

	return nil
}

func (c *CommandExecutorDocker) ContainerExists(ctx context.Context, containerName string) (bool, error) {
	if containerName == "" {
		return false, tracederrors.TracedErrorEmptyString("containerName")
	}

	container, err := c.GetContainerByName(containerName)
	if err != nil {
		return false, err
	}

	return container.Exists(ctx)
}

func (c *CommandExecutorDocker) RemoveContainer(ctx context.Context, containerName string, options *dockeroptions.RemoveOptions) error {
	if containerName == "" {
		return tracederrors.TracedErrorEmptyString("containerName")
	}

	container, err := c.GetContainerByName(containerName)
	if err != nil {
		return err
	}

	return container.Remove(ctx, options)
}

// RunCommandInTemporaryContainer runs a command in a temporary Docker container and returns the output.
//
// Implementation note — why we use a run → wait → logs → delete approach:
// This mirrors the Kubernetes implementation to provide a consistent API across container orchestration platforms.
// The container is created with --rm flag for automatic cleanup, but we also explicitly remove it
// to ensure deterministic cleanup and to have full control over the lifecycle.
//
// Steps:
//  1. `docker run`    — start the container with --rm flag without attaching
//  2. `docker wait`   — block until the container has completed
//  3. `docker logs`   — fetch the output exactly once
//  4. `docker rm`     — clean up the container (redundant with --rm but ensures cleanup)
func (c *CommandExecutorDocker) RunCommandInTemporaryContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	containerName, err := options.GetName()
	if err != nil {
		return nil, err
	}

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	commandToExecute, err := options.GetCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary container '%s' using container image '%s' started.", containerName, imageName)

	// First, ensure any existing container with this name is removed
	err = c.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	if err != nil {
		return nil, err
	}

	// Step 1: Start the container without attaching
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	runCommand := []string{
		"docker", "run",
		"--name", containerName,
		"--detach",
	}

	// Add environment variables
	for envName, envValue := range options.AdditionalEnvVars {
		runCommand = append(runCommand, "-e", fmt.Sprintf("%s=%s", envName, envValue))
	}

	// Add ports
	for _, port := range options.Ports {
		runCommand = append(runCommand, "-p", port)
	}

	// Add mounts
	for _, mount := range options.Mounts {
		runCommand = append(runCommand, "-v", mount)
	}

	// Add entrypoint if specified (nil = not set, empty = overwrite to empty, non-empty = use value)
	runCommand = appendEntryPointToCommand(runCommand, options.EntryPoint)

	runCommand = append(runCommand, imageName)
	runCommand = append(runCommand, commandToExecute...)

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: runCommand,
	})
	if err != nil {
		return nil, err
	}

	// Step 2: Wait for the container to complete and get the exit code
	waitCommand := []string{
		"docker", "wait", containerName,
	}

	waitOutput, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: waitCommand,
	})
	if err != nil {
		return nil, err
	}

	// Step 3: Fetch the logs exactly once
	logsCommand := []string{
		"docker", "logs", containerName,
	}

	output, err := commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: logsCommand,
	})
	if err != nil {
		return nil, err
	}

	// Set the return code from the wait command
	exitCodeStr, err := waitOutput.GetStdoutAsString()
	if err != nil {
		logging.LogErrorByCtxf(ctx, "Failed to get exit code from wait command: %v", err)
	} else {
		exitCodeStr = strings.TrimSpace(exitCodeStr)
		exitCode, err := strconv.Atoi(exitCodeStr)
		if err != nil {
			logging.LogErrorByCtxf(ctx, "Failed to parse exit code '%s': %v", exitCodeStr, err)
		} else {
			err = output.SetReturnCode(exitCode)
			if err != nil {
				logging.LogErrorByCtxf(ctx, "Failed to set return code: %v", err)
			}
		}
	}

	// Step 4: Delete the container to clean up
	deleteCommand := []string{
		"docker", "rm", "-f", containerName,
	}

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: deleteCommand,
	})
	if err != nil {
		// Log the error but don't fail the operation since we already got the output
		logging.LogErrorByCtxf(ctx, "Failed to remove temporary container '%s': %v", containerName, err)
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary container '%s' using container image '%s' finished.", containerName, imageName)

	return output, nil
}
