package asciichgolangpublic

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
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

func (e *ExecService) MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput) {
	commandOutput, err := e.RunCommand(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (e *ExecService) RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
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
		return nil, TracedError(err.Error())
	}
	cmd.Stderr = &stderr
	stdoutString := ""

	commandOutput = new(CommandOutput)

	cmd.Start()

	scanner := bufio.NewScanner(stdoutPipe)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()

		if OS().IsRunningOnWindows() {
			if len(m) > 0 {
				if []byte(m)[0] == 0x00 {
					m = string([]byte(m)[1:])
				}
			}

			if len(m) > 0 {
				if []byte(m)[len(m)-1] == '\r' {
					m = string([]byte(m)[:len(m)-2])
				}
			}

			m, err = Windows().DecodeStringAsString(m)
			if err != nil {
				return nil, err
			}
		}

		stdoutString += m + "\n"
		if options.LiveOutputOnStdout {
			fmt.Println(m)
		}
	}
	err = cmd.Wait()
	if err != nil {
		commandOutput.SetCmdRunError(err)
	}

	err = commandOutput.SetStdout([]byte(stdoutString))
	if err != nil {
		return nil, err
	}

	stderrBytes := stderr.Bytes()
	if OS().IsRunningOnWindows() {
		stderrBytes, err = Windows().DecodeAsBytes(stderrBytes)
		if err != nil {
			return nil, err
		}
	}

	err = commandOutput.SetStderr(stderrBytes)
	if err != nil {
		return nil, err
	}

	if cmd.ProcessState == nil {
		return nil, TracedErrorf(
			"unable to get exit code for failed command: '%v': '%v'",
			commandJoined,
			commandOutput.GetCmdRunErrorStringOrEmptyStringIfUnset(),
		)
	}

	err = commandOutput.SetReturnCode(cmd.ProcessState.ExitCode())
	if err != nil {
		return nil, err
	}

	if !commandOutput.IsExitSuccess() {
		if options.AllowAllExitCodes {
			if options.Verbose {
				LogInfof("Command '%v' has exit code != 0 but all exit codes are allowed by runOptions.AllowAllExitCodes.", commandJoined)
			}
		} else {
			errorMessage := fmt.Sprintf(
				"Command failed: '%v', %v\n%v",
				commandJoined,
				commandOutput.GetCmdRunErrorStringOrEmptyStringIfUnset(),
				commandOutput.GetStderrAsStringOrEmptyIfUnset(),
			)
			if options.Verbose {
				LogError(errorMessage)
			}
			return commandOutput, TracedError(errorMessage)
		}
	}

	return commandOutput, nil
}
