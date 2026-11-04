package nativedocker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Docker struct {
}

func NewDocker() dockerinterfaces.Docker {
	return new(Docker)
}

func RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error) {
	return NewDocker().RunContainer(ctx, options)
}

func GetContainerByName(name string) (containerinterfaces.Container, error) {
	return NewDocker().GetContainerByName(name)
}

func RemoveContainer(ctx context.Context, name string, options *dockeroptions.RemoveOptions) error {
	return NewDocker().RemoveContainer(ctx, name, options)
}

func (d *Docker) GetDeepCopyAsDocker() dockerinterfaces.Docker {
	return &Docker{}
}

func (d *Docker) GetContainerByName(name string) (containerinterfaces.Container, error) {
	return NewContainer(name)
}

func (d *Docker) GetHostDescription() (string, error) {
	return "localhost", nil
}

func (d *Docker) ListContainers(ctx context.Context) ([]containerinterfaces.Container, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (d *Docker) inspect(ctx context.Context, containerName string) (*client.ContainerInspectResult, error) {
	if containerName == "" {
		return nil, tracederrors.TracedErrorEmptyString("containerName")
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	inspect, err := cli.ContainerInspect(ctx, containerName, client.ContainerInspectOptions{})
	if err != nil {
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, tracederrors.TracedErrorf("Container inspect for '%s' failed: %w: %w", containerName, dockergeneric.ErrDockerContainerNotFound, err)
		}

		return nil, tracederrors.TracedErrorf("Container inspect for '%s' failed: %w", containerName, err)
	}

	return &inspect, nil
}

func (d *Docker) GetContainerId(ctx context.Context, containerName string) (string, error) {
	if containerName == "" {
		return "", tracederrors.TracedErrorEmptyString("containerName")
	}

	inspect, err := d.inspect(ctx, containerName)
	if err != nil {
		return "", err
	}

	containerId := inspect.Container.ID

	logging.LogInfoByCtxf(ctx, "Container '%s' has container ID '%s'.", containerName, containerId)

	return containerId, nil
}

func (d *Docker) RunContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (containerinterfaces.Container, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	name, err := options.GetName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run docker container '%s' started.", name)

	imageName, err := options.GetImageName()
	if err != nil {
		return nil, err
	}

	_, err = d.PullImage(ctx, imageName)
	if err != nil {
		return nil, err
	}

	command := options.GetCommandOrNil()
	entrypoint := options.GetEntryPointOrNil()

	autoremove := !options.KeepStoppedContainer

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create docker client: %w", err)
	}

	envVars := []string{}
	if len(options.AdditionalEnvVars) > 0 {
		envVars, err = environmentvariables.SetEnvVarsInStringSlice(envVars, options.AdditionalEnvVars)
		if err != nil {
			return nil, err
		}
	}

	portBindings := network.PortMap{}

	for _, p := range options.Ports {
		listenIpAddress := "127.0.0.1" // defaults to the local host.
		if strings.HasPrefix(p, "0.0.0.0:") {
			listenIpAddress = "0.0.0.0"
			p = strings.TrimPrefix(p, "0.0.0.0:")
		}
		if strings.HasPrefix(p, "127.0.0.1") {
			listenIpAddress = "127.0.0.1"
			p = strings.TrimPrefix(p, "127.0.0.1:")
		}

		splitted := strings.Split(p, ":")
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Unsupported port mapping '%s' to start docker container '%s'", p, name)
		}

		hostPortNumber, err := strconv.Atoi(splitted[0])
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse host port '%s': %w", splitted[0], err)
		}

		containerPortNumber, err := strconv.Atoi(splitted[1])
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse container port '%s': %w", splitted[1], err)
		}

		localhostIp, err := netip.ParseAddr(listenIpAddress)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse hostIP '%s': %w", listenIpAddress, err)
		}

		hostBinding := network.PortBinding{
			HostIP:   localhostIp,
			HostPort: strconv.Itoa(hostPortNumber),
		}

		containerPort, err := network.ParsePort(fmt.Sprintf("%d/tcp", containerPortNumber))
		if err != nil {
			return nil, err
		}

		portBindings[containerPort] = []network.PortBinding{hostBinding}
	}

	mounts := []mount.Mount{}
	for _, m := range options.Mounts {
		splitted := strings.Split(m, ":")
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Failed to process mount: '%s'.", m)
		}

		toAdd := mount.Mount{
			Type:     mount.TypeBind,
			Source:   splitted[0],
			Target:   splitted[1],
			ReadOnly: false,
		}

		mounts = append(mounts, toAdd)
	}

	exists, err := d.ContainerExists(ctx, name)
	if err != nil {
		return nil, err
	}

	var skipCreation bool
	if exists {
		if options.SkipIfAlreadyRunning {
			logging.LogInfoByCtxf(ctx, "Container creation will be skipped since '%s' is already running.", name)
			skipCreation = true
		}
	}

	if !skipCreation {

		createResult, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
			Name:  name,
			Image: imageName,
			Config: &container.Config{
				Env:        envVars,
				Cmd:        command,
				Entrypoint: entrypoint,
			},
			HostConfig: &container.HostConfig{
				AutoRemove:   autoremove,
				PortBindings: portBindings,
				Mounts:       mounts,
			},
		})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to create container '%s': %w", name, err)
		}

		_, err = cli.ContainerStart(ctx, createResult.ID, client.ContainerStartOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to start container '%s': %w", name, err)
		}

		ports, err := options.GetPortsOnHost()
		if err != nil {
			return nil, err
		}

		for _, port := range ports {
			err := netutils.WaitTcpPortOpen(ctx, "localhost", port, time.Minute*1)
			if err != nil {
				return nil, err
			}
		}
	}

	logging.LogInfoByCtxf(ctx, "Run docker container '%s' finished.", name)

	return d.GetContainerByName(name)
}

func (d *Docker) ListContainerNames(ctx context.Context) ([]string, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	nativeList, err := cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list containers")
	}

	names := []string{}
	for _, container := range nativeList.Items {
		for _, name := range container.Names {
			names = append(names, name)
		}
	}

	return names, nil
}

func (d *Docker) KillContainerByName(ctx context.Context, name string) error {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	_, err = cli.ContainerKill(ctx, name, client.ContainerKillOptions{})
	if err != nil {
		return tracederrors.TracedErrorf("Unable to kill container '%s': %w", name, err)
	}

	logging.LogChangedByCtxf(ctx, "Killed container '%s'.", name)

	return err
}

func (d *Docker) imageInspect(ctx context.Context, imageName string) (*client.ImageInspectResult, error) {
	if imageName == "" {
		return nil, tracederrors.TracedErrorEmptyString("imageName")
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	inspect, err := cli.ImageInspect(ctx, imageName, client.ImageInspectWithManifests(false))
	if err != nil {
		if errors.Is(err, errdefs.ErrNotFound) {
			return nil, tracederrors.TracedErrorf("Docker image inspect for '%s' failed: %w: %w", imageName, dockergeneric.ErrDockerImageNotFound, err)
		}
		return nil, tracederrors.TracedErrorf("ImageInspect failed for image '%s': %w", imageName, err)
	}

	return &inspect, nil
}

func (d *Docker) ImageExists(ctx context.Context, imageName string) (bool, error) {
	if imageName == "" {
		return false, tracederrors.TracedErrorEmptyString("imageName")
	}

	_, err := d.imageInspect(ctx, imageName)
	if err != nil {
		if dockergeneric.IsErrorImageNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Docker image '%s' does not exist.", imageName)
			return false, nil
		}

		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Docker image '%s' exists.", imageName)
	return true, nil
}

func (d *Docker) GetImageByName(imageName string) (containerinterfaces.Image, error) {
	image := NewImage()

	err := image.SetName(imageName)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (d *Docker) PullImage(ctx context.Context, imageName string) (containerinterfaces.Image, error) {
	if imageName == "" {
		return nil, tracederrors.TracedErrorEmptyString("imageName")
	}

	exists, err := d.ImageExists(ctx, imageName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' already exists, skip pull.", imageName)
	} else {
		logging.LogInfoByCtxf(ctx, "Pull docker image '%s' started.", imageName)

		cli, err := client.New(client.FromEnv)
		if err != nil {
			return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
		}
		defer cli.Close()

		out, err := cli.ImagePull(ctx, imageName, client.ImagePullOptions{})
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to pull image '%s': %w", imageName, err)
		}
		defer out.Close()
		io.Copy(os.Stderr, out)

		logging.LogChangedByCtxf(ctx, "Pulled docker image '%s'.", imageName)
	}

	return d.GetImageByName(imageName)
}

func (d *Docker) RemoveImage(ctx context.Context, imageName string, options *dockeroptions.RemoveOptions) error {
	if imageName == "" {
		return tracederrors.TracedErrorEmptyString("imageName")
	}

	exists, err := d.ImageExists(ctx, imageName)
	if err != nil {
		return err
	}

	if exists {
		cli, err := client.New(client.FromEnv)
		if err != nil {
			return tracederrors.TracedErrorf("unable to create docker client: %w", err)
		}
		defer cli.Close()

		force := false
		if options != nil {
			force = options.Force
		}

		_, err = cli.ImageRemove(ctx, imageName, client.ImageRemoveOptions{
			Force: force,
		})
		if err != nil {
			return tracederrors.TracedErrorf("Unable to remove image '%s': %w", imageName, err)
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Docker image '%s' is already absents. Skip remove.", imageName)
	}

	return nil
}

func (d *Docker) ContainerExists(ctx context.Context, containerName string) (bool, error) {
	if containerName == "" {
		return false, tracederrors.TracedErrorEmptyString("containerName")
	}

	container, err := d.GetContainerByName(containerName)
	if err != nil {
		return false, err
	}

	return container.Exists(ctx)
}

func (d *Docker) RemoveContainer(ctx context.Context, containerName string, options *dockeroptions.RemoveOptions) error {
	if containerName == "" {
		return tracederrors.TracedErrorEmptyString("containerName")
	}

	container, err := d.GetContainerByName(containerName)
	if err != nil {
		return err
	}

	return container.Remove(ctx, options)
}

// Waits until the execId finished and returns its exit code.
func WaitUntilExecFinished(ctx context.Context, execId string) (int, error) {
	if execId == "" {
		return 0, tracederrors.TracedErrorEmptyString("execId")
	}

	logging.LogInfoByCtxf(ctx, "Wait until exec '%s' finished started.", execId)

	var inspectResult *client.ExecInspectResult
	var maxTries = 50
	for i := range maxTries {
		cli, err := client.New(client.FromEnv)
		if err != nil {
			return 0, tracederrors.TracedErrorf("Failed to get native docker client: %w", err)
		}
		defer cli.Close()

		inspect, err := cli.ExecInspect(ctx, execId, client.ExecInspectOptions{})
		if err != nil {
			return 0, tracederrors.TracedErrorf("Attach with execId '%s' failed: %w", execId, err)
		}

		if inspect.Running {
			delay := time.Millisecond * 100
			logging.LogInfoByCtxf(ctx, "Exec '%s' is still running in docker container as PID=%d. Wait anoter '%s' (%d/%d).", execId, inspect.PID, delay, i+1, maxTries)
			time.Sleep(delay)
			continue
		}

		inspectResult = &inspect
		break
	}

	if inspectResult == nil {
		return 0, tracederrors.TracedErrorf("Container exec '%s' still running.", execId)
	}

	exitCode := inspectResult.ExitCode

	logging.LogInfoByCtxf(ctx, "Wait until exec '%s' finished finished. Exit code is %d.", execId, exitCode)

	return exitCode, nil
}

// RunCommandInTemporaryContainer runs a command in a temporary Docker container and returns the output.
//
// Implementation note — why we use a run → wait → logs → delete approach:
// This mirrors the Kubernetes implementation to provide a consistent API across container orchestration platforms.
// The container is created with AutoRemove=false for explicit control, and we manually remove it
// to ensure deterministic cleanup.
//
// Steps:
//  1. `ContainerCreate` — create and start the container without attaching
//  2. `ContainerWait`   — block until the container has completed
//  3. `ContainerLogs`   — fetch the output exactly once
//  4. `ContainerRemove` — clean up the container
func (d *Docker) RunCommandInTemporaryContainer(ctx context.Context, options *dockeroptions.DockerRunContainerOptions) (*commandoutput.CommandOutput, error) {
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

	command := options.GetCommandOrNil()
	entrypoint := options.GetEntryPointOrNil()

	logging.LogInfoByCtxf(ctx, "Run command in temporary container '%s' using container image '%s' started.", containerName, imageName)

	// First, ensure any existing container with this name is removed
	err = d.RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true})
	if err != nil {
		return nil, err
	}

	_, err = d.PullImage(ctx, imageName)
	if err != nil {
		return nil, err
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create docker client: %w", err)
	}
	defer cli.Close()

	envVars := []string{}
	if len(options.AdditionalEnvVars) > 0 {
		envVars, err = environmentvariables.SetEnvVarsInStringSlice(envVars, options.AdditionalEnvVars)
		if err != nil {
			return nil, err
		}
	}

	portBindings := network.PortMap{}
	for _, p := range options.Ports {
		listenIpAddress := "127.0.0.1"
		if strings.HasPrefix(p, "0.0.0.0:") {
			listenIpAddress = "0.0.0.0"
			p = strings.TrimPrefix(p, "0.0.0.0:")
		}
		if strings.HasPrefix(p, "127.0.0.1:") {
			listenIpAddress = "127.0.0.1"
			p = strings.TrimPrefix(p, "127.0.0.1:")
		}

		splitted := strings.Split(p, ":")
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Unsupported port mapping '%s' to start docker container '%s'", p, containerName)
		}

		hostPortNumber, err := strconv.Atoi(splitted[0])
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse host port '%s': %w", splitted[0], err)
		}

		containerPortNumber, err := strconv.Atoi(splitted[1])
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse container port '%s': %w", splitted[1], err)
		}

		localhostIp, err := netip.ParseAddr(listenIpAddress)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse hostIP '%s': %w", listenIpAddress, err)
		}

		hostBinding := network.PortBinding{
			HostIP:   localhostIp,
			HostPort: strconv.Itoa(hostPortNumber),
		}

		containerPort, err := network.ParsePort(fmt.Sprintf("%d/tcp", containerPortNumber))
		if err != nil {
			return nil, err
		}

		portBindings[containerPort] = []network.PortBinding{hostBinding}
	}

	mounts := []mount.Mount{}
	for _, m := range options.Mounts {
		splitted := strings.Split(m, ":")
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Failed to process mount: '%s'.", m)
		}

		toAdd := mount.Mount{
			Type:     mount.TypeBind,
			Source:   splitted[0],
			Target:   splitted[1],
			ReadOnly: false,
		}

		mounts = append(mounts, toAdd)
	}

	// Step 1: Create and start the container
	createResult, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:  containerName,
		Image: imageName,
		Config: &container.Config{
			Env:        envVars,
			Cmd:        command,
			Entrypoint: entrypoint,
		},
		HostConfig: &container.HostConfig{
			AutoRemove:   false, // We handle removal explicitly
			PortBindings: portBindings,
			Mounts:       mounts,
		},
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create container '%s': %w", containerName, err)
	}

	_, err = cli.ContainerStart(ctx, createResult.ID, client.ContainerStartOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to start container '%s': %w", containerName, err)
	}

	// Step 2: Wait for the container to complete
	waitResult := cli.ContainerWait(ctx, createResult.ID, client.ContainerWaitOptions{})

	// Wait for the result or error
	var exitCode int64
	select {
	case err := <-waitResult.Error:
		return nil, tracederrors.TracedErrorf("Error waiting for container '%s': %w", containerName, err)
	case result := <-waitResult.Result:
		exitCode = result.StatusCode
	}

	// Step 3: Fetch the logs exactly once
	logs, err := cli.ContainerLogs(ctx, createResult.ID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to get logs for container '%s': %w", containerName, err)
	}
	defer logs.Close()

	// Docker logs are returned in a multiplexed format
	// Use stdcopy.StdCopy to separate stdout and stderr
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, logs)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to demultiplex logs for container '%s': %w", containerName, err)
	}

	// Create CommandOutput from the logs
	output := &commandoutput.CommandOutput{}
	err = output.SetStdout(stdout.Bytes())
	if err != nil {
		return nil, err
	}
	err = output.SetStderr(stderr.Bytes())
	if err != nil {
		return nil, err
	}
	err = output.SetReturnCode(int(exitCode))
	if err != nil {
		return nil, err
	}

	// Step 4: Delete the container to clean up
	removeOptions := client.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: false,
	}
	_, err = cli.ContainerRemove(ctx, createResult.ID, removeOptions)
	if err != nil {
		// Log the error but don't fail the operation since we already got the output
		logging.LogErrorByCtxf(ctx, "Failed to remove temporary container '%s': %v", containerName, err)
	}

	logging.LogInfoByCtxf(ctx, "Run command in temporary container '%s' using container image '%s' finished.", containerName, imageName)

	return output, nil
}
