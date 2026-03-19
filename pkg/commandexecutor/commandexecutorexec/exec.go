package commandexecutorexec

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/ioutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/windowsutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) ([]byte, error) {
	output, err := RunCommand(ctx, options)
	if err != nil {
		return nil, err
	}

	stdout, err := output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

func RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (string, error) {
	output, err := RunCommand(ctx, options)
	if err != nil {
		return "", err
	}

	stdout, err := output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}

func RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (*commandoutput.CommandOutput, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	command, err := options.GetFullCommand()
	if err != nil {
		return nil, err
	}

	commandJoined, err := options.GetJoinedFullCommand()
	if err != nil {
		return nil, err
	}

	hostDescription := "localhost"

	logging.LogInfoByCtxf(ctx, "Exec command '%s' on '%s' started.", commandJoined, hostDescription)

	const avoidExecEnvVarName = "ASCIICHGOLANGPUBLIC_AVOID_EXEC"
	const trueValue = "1"
	if os.Getenv(avoidExecEnvVarName) == trueValue {
		return nil, tracederrors.TracedErrorf(
			"env var '%s' is set to '%s'. The command exec is therefore blocked. The blocked command is '%s'",
			avoidExecEnvVarName,
			trueValue,
			commandJoined,
		)
	}

	cmd := exec.Command(command[0])
	if len(command) > 1 {
		cmd = exec.Command(command[0], command[1:]...)
	}

	var stderr bytes.Buffer

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}
	cmd.Stderr = &stderr

	commandOutput := new(commandoutput.CommandOutput)

	if len(options.AdditionalEnvVars) > 0 {
		envVars, err := environmentvariables.SetEnvVarsInStringSlice(os.Environ(), options.AdditionalEnvVars)
		if err != nil {
			return nil, err
		}

		cmd.Env = envVars
	}

	writeStdin := options.IsStdinStringSet()

	var stdin io.WriteCloser

	if writeStdin {
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		defer stdin.Close()
	}

	cmd.Start()

	if writeStdin {
		bytesToWrite := []byte(options.StdinString)
		nBytesToWrite := len(bytesToWrite)

		nWrittenBytes, err := stdin.Write([]byte(options.StdinString))
		if err != nil {
			return nil, err
		}

		if nBytesToWrite != nWrittenBytes {
			return nil, tracederrors.TracedErrorf(
				"Writing to stdin of command '%v' failed. Expected '%d' bytes to write but '%d' got written",
				command,
				nBytesToWrite,
				nWrittenBytes,
			)
		}

		stdin.Close()
	}

	scanner := bufio.NewScanner(stdoutPipe)

	scanner.Split(bufio.ScanBytes)
	stdoutBytes := []byte{}
	goOn := true
	lastProcessedByteWasNewLine := false
	for {
		lastProcessedByteWasNewLine = false
		line := ""
		for {
			goOn = scanner.Scan()
			if !goOn {
				break
			}

			b := scanner.Text()
			if b == "\n" {
				lastProcessedByteWasNewLine = true
				break
			} else {
				lastProcessedByteWasNewLine = false
			}

			line += b
		}

		if goOn {
			if commandexecutorgeneric.IsLiveOutputOnStdoutEnabled(ctx) {
				mOutput := line

				if windowsutils.IsRunningOnWindows() {
					if len(mOutput) > 0 {
						if []byte(mOutput)[0] == 0x00 {
							mOutput = string([]byte(mOutput)[1:])
						}
					}

					mOutput, err = windowsutils.DecodeStringAsString(mOutput)
					if err != nil {
						return nil, err
					}
				}

				fmt.Println(mOutput)
			}
		}

		stdoutBytes = append(stdoutBytes, []byte(line)...)
		if lastProcessedByteWasNewLine {
			stdoutBytes = append(stdoutBytes, byte('\n'))
		}

		if !goOn {
			break
		}
	}

	err = cmd.Wait()
	if err != nil {
		commandOutput.SetCmdRunError(err)
	}

	err = commandOutput.SetStdout(stdoutBytes)
	if err != nil {
		return nil, err
	}

	stderrBytes := stderr.Bytes()
	if windowsutils.IsRunningOnWindows() {
		stderrBytes, err = windowsutils.DecodeAsBytes(stderrBytes)
		if err != nil {
			return nil, err
		}
	}

	err = commandOutput.SetStderr(stderrBytes)
	if err != nil {
		return nil, err
	}

	if cmd.ProcessState == nil {
		return nil, tracederrors.TracedErrorf(
			"unable to get exit code for failed command: '%v': '%v'",
			commandJoined,
			commandOutput.GetCmdRunErrorStringOrEmptyStringIfUnset(),
		)
	}

	err = commandOutput.SetReturnCode(cmd.ProcessState.ExitCode())
	if err != nil {
		return nil, err
	}

	returnCode, err := commandOutput.GetReturnCode()
	if err != nil {
		return nil, err
	}

	if !commandOutput.IsExitSuccess() {
		if options.AllowAllExitCodes {
			logging.LogInfoByCtxf(
				ctx,
				"Command '%v' has exit code '%d' != 0 but all exit codes are allowed by runOptions.AllowAllExitCodes.",
				commandJoined,
				returnCode,
			)
		} else {
			errorMessage := fmt.Sprintf(
				"Command failed: '%v', %v\n%v",
				commandJoined,
				commandOutput.GetCmdRunErrorStringOrEmptyStringIfUnset(),
				commandOutput.GetStderrAsStringOrEmptyIfUnset(),
			)
			return commandOutput, tracederrors.TracedError(errorMessage)
		}
	}

	logging.LogInfoByCtxf(ctx, "Exec command '%s' on '%s' finished.", commandJoined, hostDescription)

	return commandOutput, nil
}

func RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	fullCommand, err := options.GetFullCommand()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(fullCommand[0])
	if len(fullCommand) > 0 {
		cmd = exec.Command(fullCommand[0], fullCommand[1:]...)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create stdout pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		err = stdout.Close()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Close stdout reader after start of command failed: %w", err)
		}
		return nil, tracederrors.TracedErrorf("Failed to start command: %w", err)
	}

	ret := &ioutils.ReadCloser{
		CloseFunc: func() error {
			err := stdout.Close()
			if err != nil {
				return tracederrors.TracedErrorf("Failed to close stdout: %w", err)
			}

			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
			}()

			select {
			case <-time.After(10 * time.Second):
				err := cmd.Process.Kill()
				if err != nil {
					if !errors.Is(err, os.ErrProcessDone) {
						return tracederrors.TracedErrorf("Failed to kill command in ReadCloser: %w", err)
					}
				}
				logging.LogInfoByCtxf(ctx, "Killed ReadCloser process.")
			case err := <-done:
				// Process finished before the timeout
				if err != nil {
					logging.LogErrorByCtxf(ctx, "ReadCloser process finished with error: %v\n", err)
				} else {
					logging.LogInfoByCtxf(ctx, "ReadCloser process finished successfully")
				}
			}

			return nil
		},
		ReadFunc: func(p []byte) (n int, err error) {
			return stdout.Read(p)
		},
	}

	return ret, nil
}

func RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	fullCommand, err := options.GetFullCommand()
	if err != nil {
		return nil, err
	}

	fullCommandJoined, err := shelllinehandler.Join(fullCommand)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(fullCommand[0])
	if len(fullCommand) > 0 {
		cmd = exec.Command(fullCommand[0], fullCommand[1:]...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create stdin pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		err = stdin.Close()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Close stdin writer after start of command failed: %w", err)
		}
		return nil, tracederrors.TracedErrorf("Failed to start command: %w", err)
	}

	ret := &ioutils.WriteCloser{
		CloseFunc: func() error {
			err := stdin.Close()
			if err != nil {
				return tracederrors.TracedErrorf("Failed to close stdin: %w", err)
			}

			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
			}()

			select {
			case <-time.After(10 * time.Second):
				err := cmd.Process.Kill()
				if err != nil {
					if !errors.Is(err, os.ErrProcessDone) {
						return tracederrors.TracedErrorf("Failed to kill command in WriteCloser: %w", err)
					}
				}
				logging.LogInfoByCtxf(ctx, "Killed WriteCloser process '%s' .", fullCommandJoined)
			case err := <-done:
				// Process finished before the timeout
				if err != nil {
					logging.LogErrorByCtxf(ctx, "WriteCloser command '%s' finished with error: %v\n", fullCommandJoined, err)
				} else {
					logging.LogInfoByCtxf(ctx, "WriteCloser command '%s' finished successfully", fullCommandJoined)
				}
			}

			return nil
		},
		WriteFunc: func(p []byte) (n int, err error) {
			return stdin.Write(p)
		},
	}

	return ret, nil
}
