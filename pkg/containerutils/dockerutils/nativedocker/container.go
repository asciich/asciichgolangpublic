package nativedocker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/client"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/ioutils"
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

func (c *Container) GetDeepCopy() commandexecutorinterfaces.CommandExecutor {
	ret := &Container{
		name: c.name,
	}

	// Without this the methods of the base class (e.g. RunCommandAndGetStdoutAsString)
	// would not work on the copy.
	ret.SetParentCommandExecutorForBaseClass(ret)

	return ret
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

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return false, err
	}

	_, err = c.inspect(ctx)
	if err != nil {
		if dockergeneric.IsErrorContainerNotFound(err) {
			logging.LogInfoByCtxf(ctx, "Docker container '%s' does not exist on '%s'.", name, hostDescription)
			return false, nil
		}

		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Docker container '%s' exists on '%s'.", name, hostDescription)
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

func (c *Container) WaitUntilRemoved(ctx context.Context) error {
	name, err := c.GetName()
	if err != nil {
		return err
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wait until container '%s' on '%s' is removed started.", name, hostDescription)

	maxTries := 20
	var exists bool
	for i := range maxTries {
		exists, err = c.Exists(ctx)
		if err != nil {
			return err
		}

		if exists {
			if i+1 >= maxTries {
				break
			}

			duration := time.Millisecond * 500
			logging.LogInfoByCtxf(ctx, "Container '%s' still present on host '%s'. Waiting another %v (%d/%d)", name, hostDescription, duration, i+1, maxTries)
			time.Sleep(duration)
			continue
		} else {
			logging.LogInfoByCtxf(ctx, "Container '%s' is now absent on '%s'.", name, hostDescription)
		}

		break
	}

	if exists {
		return tracederrors.TracedErrorf("Failed to wait for container '%s' to be deleted on '%s'. Container still available.", name, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Wait until container '%s' on '%s' is removed finished.", name, hostDescription)

	return nil
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
			if !IsRemovalAlreadyInProgressError(err) {
				return tracederrors.TracedErrorf("Failed to delete container '%s' on host '%s': %w", name, hostDescription, err)
			}
		}

		err = c.WaitUntilRemoved(ctx)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Docker container '%s' removed on host '%s'.", name, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Docker container '%s' is already absent on host '%s'. Skip removal of container.", name, hostDescription)
	}

	return nil
}

// getExecExitCode returns the exit code of the finished exec 'execId'.
//
// There is a race condition returning an empty inspect while the exec is still closing.
// Therefore the inspect is retried a few times.
func getExecExitCode(ctx context.Context, cli *client.Client, execId string, containerName string) (int, error) {
	for range 3 {
		inspect, err := cli.ExecInspect(ctx, execId, client.ExecInspectOptions{})
		if err != nil {
			return -1, tracederrors.TracedErrorf("Failed to exec inspect for exec id='%s' and container '%s': %w", execId, containerName, err)
		}

		if inspect.ID == "" || inspect.ContainerID == "" {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		return inspect.ExitCode, nil
	}

	return -1, tracederrors.TracedErrorf("Unable to get exit code for exec id='%s' in container '%s'", execId, containerName)
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

	var env []string
	if options.AdditionalEnvVars != nil {
		env, err = environmentvariables.SetEnvVarsInStringSlice(env, options.AdditionalEnvVars)
		if err != nil {
			return nil, err
		}
	}

	exec, err := cli.ExecCreate(ctx, name, client.ExecCreateOptions{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  isStdinSet,
		Cmd:          cmd,
		User:         options.RunAsUser,
		Env:          env,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to exec create to RunCommand in container '%s': %w", name, err)
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
			return nil, tracederrors.TracedErrorf("Failed to write to the stdin of container '%s' with exec id '%s': %w", name, execId, err)
		}

		if cw, ok := attach.Conn.(interface{ CloseWrite() error }); ok {
			err := cw.CloseWrite()
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to close stdin of container '%s' with exec id '%s': %w", name, execId, err)
			}
		} else {
			return nil, tracederrors.TracedErrorf("Unable to close stdin of container '%s' with exec id '%s': connection does not support CloseWrite", name, execId)
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

	exitCode, err := getExecExitCode(ctx, cli, execId, name)
	if err != nil {
		return nil, err
	}

	err = output.SetReturnCode(exitCode)
	if err != nil {
		return nil, err
	}

	if !options.AllowAllExitCodes {
		if !output.IsExitSuccess() {
			stderr, err := output.GetStderrAsString()
			if err != nil {
				return nil, err
			}

			return nil, tracederrors.TracedErrorf("Run command '%s' in docker container '%s' failed. Exit code is: %d, stderr is\n%s", cmdJoined, name, exitCode, stderr)
		}
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' in docker container '%s' finished with exit code=%d.", cmdJoined, name, exitCode)

	return output, nil
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

// RunCommandAndGetStdoutAsIoReadCloser runs the command and returns its stdout as io.ReadCloser.
//
//   - stdout is streamed to the returned reader, stderr is collected for error messages.
//   - After stdout reached EOF the exit code is validated. A non-zero exit code is
//     returned as error by Read (instead of io.EOF) including stderr, unless
//     options.AllowAllExitCodes is set.
//   - Close() can be called at any time. Closing before EOF stops reading and skips the
//     exit code validation (stopping early is not considered an error).
//   - Calling Close() multiple times is safe.
func (c *Container) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	cmdJoined, err := options.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' with stdout as io.ReadCloser in docker container '%s' started.", cmdJoined, name)

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}

	cmd, err := options.GetCommand()
	if err != nil {
		cli.Close()
		return nil, err
	}

	var env []string
	if options.AdditionalEnvVars != nil {
		env, err = environmentvariables.SetEnvVarsInStringSlice(env, options.AdditionalEnvVars)
		if err != nil {
			cli.Close()
			return nil, err
		}
	}

	exec, err := cli.ExecCreate(ctx, name, client.ExecCreateOptions{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
		User:         options.RunAsUser,
		Env:          env,
	})
	if err != nil {
		cli.Close()
		return nil, tracederrors.TracedErrorf("Failed to exec create to RunCommand in container '%s': %w", name, err)
	}

	execId := exec.ID

	attach, err := cli.ExecAttach(ctx, execId, client.ExecAttachOptions{})
	if err != nil {
		cli.Close()
		return nil, tracederrors.TracedErrorf("Failed to exec attach for id '%s' on container '%s': %w", execId, name, err)
	}

	pr, pw := io.Pipe()

	// Closed when the background goroutine finished. Close() waits for it before
	// closing the docker client, since the goroutine uses the client for ExecInspect.
	copyDone := make(chan struct{})

	go func() {
		defer close(copyDone)

		// stderr is only accessed inside this goroutine, so no locking is needed.
		var stderr bytes.Buffer
		_, err := stdcopy.StdCopy(pw, &stderr, attach.Reader)
		if err != nil {
			// Also happens if Close() was called before EOF: The reader is closed then
			// and nobody receives this error anymore.
			pw.CloseWithError(
				tracederrors.TracedErrorf("Failed to read stdout of command '%s' in docker container '%s' with exec id '%s': %w", cmdJoined, name, execId, err),
			)
			return
		}

		// stdout reached EOF: Validate the exit code before signaling EOF to the reader.
		_, err = WaitUntilExecFinished(ctx, execId)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		exitCode, err := getExecExitCode(ctx, cli, execId, name)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if exitCode != 0 && !options.AllowAllExitCodes {
			pw.CloseWithError(
				tracederrors.TracedErrorf(
					"Run command '%s' with stdout as io.ReadCloser in docker container '%s' failed. Exit code is: %d, stderr is\n%s",
					cmdJoined,
					name,
					exitCode,
					stderr.String(),
				),
			)
			return
		}

		logging.LogInfoByCtxf(ctx, "Command '%s' with stdout as io.ReadCloser in docker container '%s' finished with exit code=%d.", cmdJoined, name, exitCode)

		pw.Close() // signals EOF cleanly
	}()

	var closeOnce sync.Once
	var closeErr error

	ret := &ioutils.ReadCloser{
		CloseFunc: func() error {
			closeOnce.Do(func() {
				// Unblock the goroutine: Writes to pw fail after closing pr,
				// reads from attach.Reader fail after closing the connection.
				pr.Close()
				attach.HijackedResponse.Close()

				// Ensure the goroutine is finished before the client is closed.
				<-copyDone

				closeErr = cli.Close()
				if closeErr != nil {
					closeErr = tracederrors.TracedErrorf("Failed to close docker client: %w", closeErr)
				}
			})
			return closeErr
		},
		ReadFunc: func(p []byte) (n int, err error) {
			return pr.Read(p)
		},
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' with stdout as io.ReadCloser in docker container '%s' finished.", cmdJoined, name)

	return ret, nil
}

// RunCommandAndGetStdinAsIoWriteCloser runs the command and returns its stdin as io.WriteCloser.
//
// Close() must be called to finish the command. It:
//  1. Half-closes the connection (CloseWrite) so the command receives EOF on stdin
//     after all written data is transmitted.
//  2. Waits until stdout/stderr are fully consumed (EOF = command finished).
//     stdout/stderr are drained in the background from the very beginning. Otherwise
//     commands producing output (e.g. 'tee') could block, and closing the connection
//     with unread data pending can reset it and silently drop not yet transmitted stdin data.
//  3. Closes the connection and the client.
//  4. Validates the exit code. A non-zero exit code returns an error including stderr
//     unless options.AllowAllExitCodes is set.
//
// Calling Close() multiple times is safe, subsequent calls return the result of the first call.
func (c *Container) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	cmdJoined, err := options.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' with stdin as io.WriteCloser in docker container '%s' started.", cmdJoined, name)

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}

	cmd, err := options.GetCommand()
	if err != nil {
		cli.Close()
		return nil, err
	}

	var env []string
	if options.AdditionalEnvVars != nil {
		env, err = environmentvariables.SetEnvVarsInStringSlice(env, options.AdditionalEnvVars)
		if err != nil {
			cli.Close()
			return nil, err
		}
	}

	exec, err := cli.ExecCreate(ctx, name, client.ExecCreateOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		User:         options.RunAsUser,
		Env:          env,
	})
	if err != nil {
		cli.Close()
		return nil, tracederrors.TracedErrorf("Failed to exec create to RunCommand in container '%s': %w", name, err)
	}

	execId := exec.ID

	attach, err := cli.ExecAttach(ctx, execId, client.ExecAttachOptions{})
	if err != nil {
		cli.Close()
		return nil, tracederrors.TracedErrorf("Failed to exec attach for id '%s' on container '%s': %w", execId, name, err)
	}

	// Drain stdout/stderr in the background right from the start.
	// stdout is discarded (it's a stdin writer), stderr is kept for error messages.
	// The stderr buffer is only accessed after drainDone is received, so no locking is needed.
	var stderr bytes.Buffer
	drainDone := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(io.Discard, &stderr, attach.Reader)
		drainDone <- err
	}()

	closeAndWait := func() error {
		defer cli.Close()
		defer attach.HijackedResponse.Close()

		// Signal EOF on stdin. All data written so far is transmitted before.
		err := attach.HijackedResponse.CloseWrite()
		if err != nil {
			return tracederrors.TracedErrorf("Failed to close stdin of command '%s' in container '%s' with exec id '%s': %w", cmdJoined, name, execId, err)
		}

		// Wait until all output is consumed. EOF means the command has finished.
		select {
		case err = <-drainDone:
		case <-ctx.Done():
			return tracederrors.TracedErrorf("Context done while waiting for command '%s' in container '%s' to finish: %w", cmdJoined, name, ctx.Err())
		}
		if err != nil {
			return tracederrors.TracedErrorf("Failed to read stdout and stderr of exec id '%s' on container '%s': %w", execId, name, err)
		}

		_, err = WaitUntilExecFinished(ctx, execId)
		if err != nil {
			return err
		}

		exitCode, err := getExecExitCode(ctx, cli, execId, name)
		if err != nil {
			return err
		}

		if exitCode != 0 && !options.AllowAllExitCodes {
			return tracederrors.TracedErrorf(
				"Run command '%s' with stdin as io.WriteCloser in docker container '%s' failed. Exit code is: %d, stderr is\n%s",
				cmdJoined,
				name,
				exitCode,
				stderr.String(),
			)
		}

		logging.LogInfoByCtxf(ctx, "Command '%s' with stdin as io.WriteCloser in docker container '%s' finished with exit code=%d.", cmdJoined, name, exitCode)

		return nil
	}

	var closeOnce sync.Once
	var closeErr error

	ret := &ioutils.WriteCloser{
		CloseFunc: func() error {
			closeOnce.Do(func() {
				closeErr = closeAndWait()
			})
			return closeErr
		},
		WriteFunc: func(p []byte) (n int, err error) {
			return attach.HijackedResponse.Conn.Write(p)
		},
	}

	logging.LogInfoByCtxf(ctx, "Run command '%s' with stdin as io.WriteCloser in docker container '%s' finished.", cmdJoined, name)

	return ret, nil
}

func (c *Container) GetLogs(ctx context.Context) ([]byte, []byte, error) {
	name, err := c.GetName()
	if err != nil {
		return nil, nil, err
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	logs, err := cli.ContainerLogs(ctx, name, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to get logs for container '%s': %w", name, err)
	}
	defer logs.Close()

	// Docker logs are returned in a multiplexed format
	// Use stdcopy.StdCopy to separate stdout and stderr
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, logs)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Failed to demultiplex logs for container '%s': %w", name, err)
	}

	// Return empty byte slices if no data, never return nil for successful operations
	stdoutBytes := stdout.Bytes()
	stderrBytes := stderr.Bytes()

	if stdoutBytes == nil {
		stdoutBytes = []byte{}
	}
	if stderrBytes == nil {
		stderrBytes = []byte{}
	}

	return stdoutBytes, stderrBytes, nil
}

func (c *Container) WaitUntilFinished(ctx context.Context, timeout time.Duration) error {
	name, err := c.GetName()
	if err != nil {
		return err
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		isRunning, err := c.IsRunning(ctx)
		if err != nil {
			return err
		}

		if !isRunning {
			return nil
		}

		// Wait a bit before checking again
		time.Sleep(time.Millisecond * 100)
	}

	return tracederrors.TracedErrorf("Container '%s' did not finish within timeout %v", name, timeout)
}
