package asciichgolangpublic

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/os/osutils"
	"github.com/asciich/asciichgolangpublic/os/windows"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type ExecService struct {
	CommandExecutorBase
}

func Exec() (e *ExecService) {
	return NewExec()
}

func NewExec() (e *ExecService) {
	e = new(ExecService)
	e.MustSetParentCommandExecutorForBaseClass(e)
	return e
}

func NewExecService() (e *ExecService) {
	return new(ExecService)
}

func (e *ExecService) GetDeepCopy() (deepCopy CommandExecutor) {
	d := NewExec()

	*d = *e

	deepCopy = d

	return deepCopy
}

func (e *ExecService) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
}

func (e *ExecService) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := e.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (e *ExecService) MustRunCommand(options *parameteroptions.RunCommandOptions) (commandOutput *CommandOutput) {
	commandOutput, err := e.RunCommand(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (e *ExecService) RunCommand(options *parameteroptions.RunCommandOptions) (commandOutput *CommandOutput, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	command, err := options.GetCommand()
	if err != nil {
		return nil, err
	}
	commandJoined, err := options.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(command[0])
	if len(options.Command) > 1 {
		cmd = exec.Command(command[0], command[1:]...)
	}

	var stderr bytes.Buffer

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}
	cmd.Stderr = &stderr

	commandOutput = new(CommandOutput)

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
			if options.LiveOutputOnStdout {
				mOutput := line

				if osutils.IsRunningOnWindows() {
					if len(mOutput) > 0 {
						if []byte(mOutput)[0] == 0x00 {
							mOutput = string([]byte(mOutput)[1:])
						}
					}

					mOutput, err = windows.DecodeStringAsString(mOutput)
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
	if osutils.IsRunningOnWindows() {
		stderrBytes, err = windows.DecodeAsBytes(stderrBytes)
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
			if options.Verbose {
				logging.LogInfof(
					"Command '%v' has exit code '%d' != 0 but all exit codes are allowed by runOptions.AllowAllExitCodes.",
					commandJoined,
					returnCode,
				)
			}
		} else {
			errorMessage := fmt.Sprintf(
				"Command failed: '%v', %v\n%v",
				commandJoined,
				commandOutput.GetCmdRunErrorStringOrEmptyStringIfUnset(),
				commandOutput.GetStderrAsStringOrEmptyIfUnset(),
			)
			if options.Verbose {
				logging.LogError(errorMessage)
			}
			return commandOutput, tracederrors.TracedError(errorMessage)
		}
	}

	return commandOutput, nil
}
