package asciichgolangpublic

import (
	"fmt"
	"strconv"
	"strings"

	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
)

type CommandOutput struct {
	returnCode  *int
	stdout      *[]byte
	stderr      *[]byte
	cmdRunError *error
}

func NewCommandOutput() (c *CommandOutput) {
	return new(CommandOutput)
}

func (c *CommandOutput) GetCmdRunError() (cmdRunError *error, err error) {
	if c.cmdRunError == nil {
		return nil, TracedErrorf("cmdRunError not set")
	}

	return c.cmdRunError, nil
}

func (c *CommandOutput) GetFirstLineOfStdoutAsString() (firstLine string, err error) {
	lines, err := c.GetStdoutAsLines(false)
	if err != nil {
		return "", err
	}

	firstLine = lines[0]
	return firstLine, nil
}

func (c *CommandOutput) GetStderr() (stderr *[]byte, err error) {
	if c.stderr == nil {
		return nil, TracedErrorf("stderr not set")
	}

	return c.stderr, nil
}

func (c *CommandOutput) GetStdout() (stdout *[]byte, err error) {
	if c.stdout == nil {
		return nil, TracedErrorf("stdout not set")
	}

	return c.stdout, nil
}

func (c *CommandOutput) GetStdoutAsFloat64() (stdout float64, err error) {
	stdoutString, err := c.GetStdoutAsString()
	if err != nil {
		return -1, err
	}

	stdoutString = strings.TrimSpace(stdoutString)

	stdout, err = strconv.ParseFloat(stdoutString, 64)
	if err != nil {
		return -1, TracedError(err)
	}

	return stdout, err
}

func (c *CommandOutput) IsStderrEmpty() (isEmpty bool, err error) {
	stderr, err := c.GetStderrAsString()
	if err != nil {
		return false, err
	}

	return stderr == "", nil
}

func (c *CommandOutput) IsStdoutAndStderrEmpty() (isEmpty bool, err error) {
	isEmpty, err = c.IsStdoutEmpty()
	if err != nil {
		return false, err
	}

	if !isEmpty {
		return false, err
	}

	return c.IsStderrEmpty()
}

func (c *CommandOutput) IsStdoutEmpty() (isEmpty bool, err error) {
	stdout, err := c.GetStdoutAsString()
	if err != nil {
		return false, err
	}

	return stdout == "", nil
}

func (c *CommandOutput) LogStdoutAsInfo() (err error) {
	stdout, err := c.GetStdoutAsString()
	if err != nil {
		return err
	}

	LogInfo(stdout)

	return nil
}

func (c *CommandOutput) MustCheckExitSuccess(verbose bool) {
	err := c.CheckExitSuccess(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandOutput) MustGetCmdRunError() (cmdRunError *error) {
	cmdRunError, err := c.GetCmdRunError()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return cmdRunError
}

func (c *CommandOutput) MustGetFirstLineOfStdoutAsString() (firstLine string) {
	firstLine, err := c.GetFirstLineOfStdoutAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return firstLine
}

func (c *CommandOutput) MustGetReturnCode() (returnCode int) {
	returnCode, err := c.GetReturnCode()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return returnCode
}

func (c *CommandOutput) MustGetStderr() (stderr *[]byte) {
	stderr, err := c.GetStderr()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stderr
}

func (c *CommandOutput) MustGetStderrAsString() (stdout string) {
	stdout, err := c.GetStderrAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandOutput) MustGetStdout() (stdout *[]byte) {
	stdout, err := c.GetStdout()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandOutput) MustGetStdoutAsBytes() (stdout []byte) {
	stdout, err := c.GetStdoutAsBytes()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandOutput) MustGetStdoutAsFloat64() (stdout float64) {
	stdout, err := c.GetStdoutAsFloat64()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandOutput) MustGetStdoutAsLines(removeLastLineIfEmpty bool) (stdoutLines []string) {
	stdoutLines, err := c.GetStdoutAsLines(removeLastLineIfEmpty)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdoutLines
}

func (c *CommandOutput) MustGetStdoutAsString() (stdout string) {
	stdout, err := c.GetStdoutAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandOutput) MustIsStderrEmpty() (isEmpty bool) {
	isEmpty, err := c.IsStderrEmpty()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isEmpty
}

func (c *CommandOutput) MustIsStdoutAndStderrEmpty() (isEmpty bool) {
	isEmpty, err := c.IsStdoutAndStderrEmpty()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isEmpty
}

func (c *CommandOutput) MustIsStdoutEmpty() (isEmpty bool) {
	isEmpty, err := c.IsStdoutEmpty()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isEmpty
}

func (c *CommandOutput) MustIsTimedOut() (IsTimedOut bool) {
	IsTimedOut, err := c.IsTimedOut()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return IsTimedOut
}

func (c *CommandOutput) MustLogStdoutAsInfo() {
	err := c.LogStdoutAsInfo()
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandOutput) MustSetReturnCode(returnCode int) {
	err := c.SetReturnCode(returnCode)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandOutput) MustSetStderr(stderr []byte) {
	err := c.SetStderr(stderr)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandOutput) MustSetStderrByString(stderr string) {
	err := c.SetStderrByString(stderr)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandOutput) MustSetStdout(stdout []byte) {
	err := c.SetStdout(stdout)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandOutput) MustSetStdoutByString(stdout string) {
	err := c.SetStdoutByString(stdout)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (o *CommandOutput) CheckExitSuccess(verbose bool) (err error) {
	if o.IsExitSuccess() {
		return nil
	} else {
		return TracedError(
			"Return code is not exit success",
		)
	}
}

func (o *CommandOutput) GetCmdRunErrorStringOrEmptyStringIfUnset() (cmdRunErrorString string) {
	if o.cmdRunError == nil {
		return ""
	}

	cmdRunErrorString = fmt.Sprintf("%v", *o.cmdRunError)
	if len(*o.stderr) > 0 {
		cmdRunErrorString += "\n" + string(*o.stderr)
	}

	return cmdRunErrorString
}

func (o *CommandOutput) GetReturnCode() (returnCode int, err error) {
	if o.returnCode == nil {
		return -1, TracedError("returnCode not set")
	}

	return *o.returnCode, nil
}

func (o *CommandOutput) GetStderrAsString() (stderr string, err error) {
	if o.stderr == nil {
		return "", TracedError("stderr is not set")
	}

	if OS().IsRunningOnWindows() {
		stderr, err = UTF16().DecodeAsString(*o.stderr)
		if err != nil {
			return "", err
		}

		stderr = strings.ReplaceAll(stderr, "\r\n", "\n")
		return stderr, nil
	} else {
		return string(*o.stderr), nil
	}
}

func (o *CommandOutput) GetStderrAsStringOrEmptyIfUnset() (stderr string) {
	if o.stderr == nil {
		return ""
	}

	return string(*o.stderr)
}

func (o *CommandOutput) GetStdoutAsBytes() (stdout []byte, err error) {
	if o.stdout == nil {
		return nil, TracedError("stdout is not set")
	}

	return *o.stdout, nil
}

func (o *CommandOutput) GetStdoutAsLines(removeLastLineIfEmpty bool) (stdoutLines []string, err error) {
	stdoutString, err := o.GetStdoutAsString()
	if err != nil {
		return nil, err
	}

	stdoutLines = astrings.SplitLines(stdoutString, removeLastLineIfEmpty)

	stdoutLines = aslices.RemoveLastElementIfEmptyString(stdoutLines)
	return stdoutLines, nil
}

func (o *CommandOutput) GetStdoutAsString() (stdout string, err error) {
	if o.stdout == nil {
		return "", TracedError("stdout is not set")
	}

	if OS().IsRunningOnWindows() {
		stdout, err = UTF16().DecodeAsString(*o.stdout)
		if err != nil {
			return "", err
		}

		stdout = strings.ReplaceAll(stdout, "\r\n", "\n")
		return stdout, nil
	} else {
		return string(*o.stdout), nil
	}
}

func (o *CommandOutput) IsExitSuccess() (isSuccess bool) {
	if o.returnCode == nil {
		return false
	}

	return *o.returnCode == ExitCodes().ExitCodeOK()
}

func (o *CommandOutput) IsTimedOut() (IsTimedOut bool, err error) {
	returnCode, err := o.GetReturnCode()
	if err != nil {
		return false, err
	}

	return returnCode == ExitCodes().ExitCodeTimeout(), nil
}

func (o *CommandOutput) SetCmdRunError(err error) {
	errToAdd := err
	o.cmdRunError = &errToAdd
}

func (o *CommandOutput) SetReturnCode(returnCode int) (err error) {
	returnCodeToAdd := returnCode
	o.returnCode = &returnCodeToAdd

	return nil
}

func (o *CommandOutput) SetStderr(stderr []byte) (err error) {
	stderrToAdd := stderr
	o.stderr = &stderrToAdd

	return nil
}

func (o *CommandOutput) SetStderrByString(stderr string) (err error) {
	err = o.SetStderr([]byte(stderr))
	if err != nil {
		return err
	}

	return err
}

func (o *CommandOutput) SetStdout(stdout []byte) (err error) {
	stdoutToAdd := stdout
	o.stdout = &stdoutToAdd
	return nil
}

func (o *CommandOutput) SetStdoutByString(stdout string) (err error) {
	err = o.SetStdout([]byte(stdout))
	if err != nil {
		return err
	}

	return err
}
