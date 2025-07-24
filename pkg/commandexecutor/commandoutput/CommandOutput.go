package commandoutput

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/encodingutils/utf16utils"
	"github.com/asciich/asciichgolangpublic/pkg/exitcodes"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandOutput struct {
	ReturnCode  *int
	Stdout      *[]byte
	Stderr      *[]byte
	cmdRunError *error
}

func NewCommandOutput() (c *CommandOutput) {
	return new(CommandOutput)
}

func (c *CommandOutput) GetCmdRunError() (cmdRunError *error, err error) {
	if c.cmdRunError == nil {
		return nil, tracederrors.TracedErrorf("cmdRunError not set")
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
	if c.Stderr == nil {
		return nil, tracederrors.TracedErrorf("stderr not set")
	}

	return c.Stderr, nil
}

func (c *CommandOutput) GetStdout() (stdout *[]byte, err error) {
	if c.Stdout == nil {
		return nil, tracederrors.TracedErrorf("stdout not set")
	}

	return c.Stdout, nil
}

func (c *CommandOutput) GetStdoutAsFloat64() (stdout float64, err error) {
	stdoutString, err := c.GetStdoutAsString()
	if err != nil {
		return -1, err
	}

	stdoutString = strings.TrimSpace(stdoutString)

	stdout, err = strconv.ParseFloat(stdoutString, 64)
	if err != nil {
		return -1, tracederrors.TracedError(err)
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

	logging.LogInfo(stdout)

	return nil
}

func (o *CommandOutput) CheckExitSuccess(verbose bool) (err error) {
	if o.IsExitSuccess() {
		return nil
	} else {
		return tracederrors.TracedError(
			"Return code is not exit success",
		)
	}
}

func (o *CommandOutput) GetCmdRunErrorStringOrEmptyStringIfUnset() (cmdRunErrorString string) {
	if o.cmdRunError == nil {
		return ""
	}

	cmdRunErrorString = fmt.Sprintf("%v", *o.cmdRunError)
	if len(*o.Stderr) > 0 {
		cmdRunErrorString += "\n" + string(*o.Stderr)
	}

	return cmdRunErrorString
}

func (o *CommandOutput) GetReturnCode() (returnCode int, err error) {
	if o.ReturnCode == nil {
		return -1, tracederrors.TracedError("returnCode not set")
	}

	return *o.ReturnCode, nil
}

func (o *CommandOutput) GetStderrAsString() (stderr string, err error) {
	if o.Stderr == nil {
		return "", tracederrors.TracedError("stderr is not set")
	}

	if osutils.IsRunningOnWindows() {
		stderr, err = utf16utils.DecodeAsString(*o.Stderr)
		if err != nil {
			return "", err
		}

		stderr = strings.ReplaceAll(stderr, "\r\n", "\n")
		return stderr, nil
	} else {
		return string(*o.Stderr), nil
	}
}

func (o *CommandOutput) GetStderrAsStringOrEmptyIfUnset() (stderr string) {
	if o.Stderr == nil {
		return ""
	}

	return string(*o.Stderr)
}

func (o *CommandOutput) GetStdoutAsBytes() (stdout []byte, err error) {
	if o.Stdout == nil {
		return nil, tracederrors.TracedError("stdout is not set")
	}

	return *o.Stdout, nil
}

func (o *CommandOutput) GetStdoutAsLines(removeLastLineIfEmpty bool) (stdoutLines []string, err error) {
	stdoutString, err := o.GetStdoutAsString()
	if err != nil {
		return nil, err
	}

	stdoutLines = stringsutils.SplitLines(stdoutString, removeLastLineIfEmpty)

	stdoutLines = slicesutils.RemoveLastElementIfEmptyString(stdoutLines)
	return stdoutLines, nil
}

func (o *CommandOutput) GetStdoutAsString() (stdout string, err error) {
	if o.Stdout == nil {
		return "", tracederrors.TracedError("stdout is not set")
	}

	if osutils.IsRunningOnWindows() {
		stdout, err = utf16utils.DecodeAsString(*o.Stdout)
		if err != nil {
			return "", err
		}

		stdout = strings.ReplaceAll(stdout, "\r\n", "\n")
		return stdout, nil
	} else {
		return string(*o.Stdout), nil
	}
}

func (o *CommandOutput) IsExitSuccess() (isSuccess bool) {
	if o.ReturnCode == nil {
		return false
	}

	return *o.ReturnCode == exitcodes.EXIT_CODE_OK
}

func (o *CommandOutput) IsTimedOut() (IsTimedOut bool, err error) {
	returnCode, err := o.GetReturnCode()
	if err != nil {
		return false, err
	}

	return returnCode == exitcodes.EXIT_CODE_TIMEOUT, nil
}

func (o *CommandOutput) SetCmdRunError(err error) {
	errToAdd := err
	o.cmdRunError = &errToAdd
}

func (o *CommandOutput) SetReturnCode(returnCode int) (err error) {
	returnCodeToAdd := returnCode
	o.ReturnCode = &returnCodeToAdd

	return nil
}

func (o *CommandOutput) SetStderr(stderr []byte) (err error) {
	stderrToAdd := stderr
	o.Stderr = &stderrToAdd

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
	o.Stdout = &stdoutToAdd
	return nil
}

func (o *CommandOutput) SetStdoutByString(stdout string) (err error) {
	err = o.SetStdout([]byte(stdout))
	if err != nil {
		return err
	}

	return err
}
