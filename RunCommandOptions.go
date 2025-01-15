package asciichgolangpublic

import (
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type RunCommandOptions struct {
	Command            []string
	TimeoutString      string
	Verbose            bool
	AllowAllExitCodes  bool
	LiveOutputOnStdout bool

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

	command = aslices.GetDeepCopyOfStringsSlice(o.Command)

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

	joinedCommand, err = ShellLineHandler().Join(command)
	if err != nil {
		return "", err
	}

	return joinedCommand, nil
}

func (o *RunCommandOptions) GetTimeoutSecondsAsString() (timeoutSeconds string, err error) {
	if len(o.TimeoutString) <= 0 {
		return "", err
	}

	timeoutSeconds, err = DurationParser().ToSecondsAsString(o.TimeoutString)
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

func (r *RunCommandOptions) GetLiveOutputOnStdout() (liveOutputOnStdout bool, err error) {

	return r.LiveOutputOnStdout, nil
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

func (r *RunCommandOptions) GetVerbose() (verbose bool, err error) {

	return r.Verbose, nil
}

func (r *RunCommandOptions) MustGetAllowAllExitCodes() (allowAllExitCodes bool) {
	allowAllExitCodes, err := r.GetAllowAllExitCodes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return allowAllExitCodes
}

func (r *RunCommandOptions) MustGetCommand() (command []string) {
	command, err := r.GetCommand()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return command
}

func (r *RunCommandOptions) MustGetJoinedCommand() (joinedCommand string) {
	joinedCommand, err := r.GetJoinedCommand()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return joinedCommand
}

func (r *RunCommandOptions) MustGetLiveOutputOnStdout() (liveOutputOnStdout bool) {
	liveOutputOnStdout, err := r.GetLiveOutputOnStdout()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return liveOutputOnStdout
}

func (r *RunCommandOptions) MustGetStdinString() (stdinString string) {
	stdinString, err := r.GetStdinString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdinString
}

func (r *RunCommandOptions) MustGetTimeoutSecondsAsString() (timeoutSeconds string) {
	timeoutSeconds, err := r.GetTimeoutSecondsAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return timeoutSeconds
}

func (r *RunCommandOptions) MustGetTimeoutString() (timeoutString string) {
	timeoutString, err := r.GetTimeoutString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return timeoutString
}

func (r *RunCommandOptions) MustGetVerbose() (verbose bool) {
	verbose, err := r.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (r *RunCommandOptions) MustSetAllowAllExitCodes(allowAllExitCodes bool) {
	err := r.SetAllowAllExitCodes(allowAllExitCodes)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetCommand(command []string) {
	err := r.SetCommand(command)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetLiveOutputOnStdout(liveOutputOnStdout bool) {
	err := r.SetLiveOutputOnStdout(liveOutputOnStdout)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetStdinString(stdinString string) {
	err := r.SetStdinString(stdinString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetTimeoutString(timeoutString string) {
	err := r.SetTimeoutString(timeoutString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetVerbose(verbose bool) {
	err := r.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (r *RunCommandOptions) SetLiveOutputOnStdout(liveOutputOnStdout bool) (err error) {
	r.LiveOutputOnStdout = liveOutputOnStdout

	return nil
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

func (r *RunCommandOptions) SetVerbose(verbose bool) (err error) {
	r.Verbose = verbose

	return nil
}
