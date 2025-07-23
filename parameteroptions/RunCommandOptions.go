package parameteroptions

import (
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datetime"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type RunCommandOptions struct {
	Command           []string
	TimeoutString     string
	AllowAllExitCodes bool

	// If set this will be send to stdin of the command:
	StdinString string

	// Run as "root" user (or Administrator on Windows):
	RunAsRoot bool

	RemoveLastLineIfEmpty bool
}

func NewRunCommandOptions() (runCommandOptions *RunCommandOptions) {
	return new(RunCommandOptions)
}

func (o *RunCommandOptions) GetCommand() (command []string, err error) {
	if len(o.Command) <= 0 {
		return nil, tracederrors.TracedError("command not set")
	}

	command = slicesutils.GetDeepCopyOfStringsSlice(o.Command)

	if o.IsTimeoutSet() {
		timeout, err := o.GetTimeoutSecondsAsString()
		if err != nil {
			return nil, err
		}

		command = append([]string{"timeout", timeout}, command...)
	}

	return command, nil
}

func (o *RunCommandOptions) GetDeepCopy() (deepCopy *RunCommandOptions) {
	deepCopy = NewRunCommandOptions()
	*deepCopy = *o
	return deepCopy
}

func (o *RunCommandOptions) GetJoinedCommand() (joinedCommand string, err error) {
	command, err := o.GetCommand()
	if err != nil {
		return "", err
	}

	joinedCommand, err = shelllinehandler.Join(command)
	if err != nil {
		return "", err
	}

	return joinedCommand, nil
}

func (o *RunCommandOptions) GetTimeoutSecondsAsString() (timeoutSeconds string, err error) {
	if len(o.TimeoutString) <= 0 {
		return "", err
	}

	timeoutSeconds, err = datetime.DurationParser().ToSecondsAsString(o.TimeoutString)
	if err != nil {
		return "", err
	}

	return timeoutSeconds, nil
}

func (o *RunCommandOptions) IsStdinStringSet() (isSet bool) {
	return o.StdinString != ""
}

func (o *RunCommandOptions) IsTimeoutSet() (isSet bool) {
	return len(o.TimeoutString) > 0
}

func (r *RunCommandOptions) GetAllowAllExitCodes() (allowAllExitCodes bool, err error) {

	return r.AllowAllExitCodes, nil
}

func (r *RunCommandOptions) GetRemoveLastLineIfEmpty() (removeLastLineIfEmpty bool) {

	return r.RemoveLastLineIfEmpty
}

func (r *RunCommandOptions) GetRunAsRoot() (runAsRoot bool) {

	return r.RunAsRoot
}

func (r *RunCommandOptions) GetStdinString() (stdinString string, err error) {
	if r.StdinString == "" {
		return "", tracederrors.TracedErrorf("StdinString not set")
	}

	return r.StdinString, nil
}

func (r *RunCommandOptions) GetTimeoutString() (timeoutString string, err error) {
	if r.TimeoutString == "" {
		return "", tracederrors.TracedErrorf("TimeoutString not set")
	}

	return r.TimeoutString, nil
}

func (r *RunCommandOptions) SetAllowAllExitCodes(allowAllExitCodes bool) (err error) {
	r.AllowAllExitCodes = allowAllExitCodes

	return nil
}

func (r *RunCommandOptions) SetCommand(command []string) (err error) {
	if command == nil {
		return tracederrors.TracedErrorf("command is nil")
	}

	if len(command) <= 0 {
		return tracederrors.TracedErrorf("command has no elements")
	}

	r.Command = command

	return nil
}

func (r *RunCommandOptions) SetRemoveLastLineIfEmpty(removeLastLineIfEmpty bool) {
	r.RemoveLastLineIfEmpty = removeLastLineIfEmpty
}

func (r *RunCommandOptions) SetRunAsRoot(runAsRoot bool) {
	r.RunAsRoot = runAsRoot
}

func (r *RunCommandOptions) SetStdinString(stdinString string) (err error) {
	if stdinString == "" {
		return tracederrors.TracedErrorf("stdinString is empty string")
	}

	r.StdinString = stdinString

	return nil
}

func (r *RunCommandOptions) SetTimeoutString(timeoutString string) (err error) {
	if timeoutString == "" {
		return tracederrors.TracedErrorf("timeoutString is empty string")
	}

	r.TimeoutString = timeoutString

	return nil
}
