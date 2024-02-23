package asciichgolangpublic

type RunCommandOptions struct {
	Command            []string
	TimeoutString      string
	Verbose            bool
	AllowAllExitCodes  bool
	LiveOutputOnStdout bool

	// Run as "root" user (or Administrator on Windows):
	RunAsRoot bool
}

func NewRunCommandOptions() (runCommandOptions *RunCommandOptions) {
	return new(RunCommandOptions)
}

func (o *RunCommandOptions) GetCommand() (command []string, err error) {
	if len(o.Command) <= 0 {
		return nil, TracedError("command not set")
	}

	command = Slices().GetDeepCopyOfStringsSlice(o.Command)

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

func (r *RunCommandOptions) GetTimeoutString() (timeoutString string, err error) {
	if r.TimeoutString == "" {
		return "", TracedErrorf("TimeoutString not set")
	}

	return r.TimeoutString, nil
}

func (r *RunCommandOptions) GetVerbose() (verbose bool, err error) {

	return r.Verbose, nil
}

func (r *RunCommandOptions) MustGetAllowAllExitCodes() (allowAllExitCodes bool) {
	allowAllExitCodes, err := r.GetAllowAllExitCodes()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return allowAllExitCodes
}

func (r *RunCommandOptions) MustGetCommand() (command []string) {
	command, err := r.GetCommand()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return command
}

func (r *RunCommandOptions) MustGetJoinedCommand() (joinedCommand string) {
	joinedCommand, err := r.GetJoinedCommand()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return joinedCommand
}

func (r *RunCommandOptions) MustGetLiveOutputOnStdout() (liveOutputOnStdout bool) {
	liveOutputOnStdout, err := r.GetLiveOutputOnStdout()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return liveOutputOnStdout
}

func (r *RunCommandOptions) MustGetTimeoutSecondsAsString() (timeoutSeconds string) {
	timeoutSeconds, err := r.GetTimeoutSecondsAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return timeoutSeconds
}

func (r *RunCommandOptions) MustGetTimeoutString() (timeoutString string) {
	timeoutString, err := r.GetTimeoutString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return timeoutString
}

func (r *RunCommandOptions) MustGetVerbose() (verbose bool) {
	verbose, err := r.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (r *RunCommandOptions) MustSetAllowAllExitCodes(allowAllExitCodes bool) {
	err := r.SetAllowAllExitCodes(allowAllExitCodes)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetCommand(command []string) {
	err := r.SetCommand(command)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetLiveOutputOnStdout(liveOutputOnStdout bool) {
	err := r.SetLiveOutputOnStdout(liveOutputOnStdout)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetTimeoutString(timeoutString string) {
	err := r.SetTimeoutString(timeoutString)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) MustSetVerbose(verbose bool) {
	err := r.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (r *RunCommandOptions) SetAllowAllExitCodes(allowAllExitCodes bool) (err error) {
	r.AllowAllExitCodes = allowAllExitCodes

	return nil
}

func (r *RunCommandOptions) SetCommand(command []string) (err error) {
	if command == nil {
		return TracedErrorf("command is nil")
	}

	if len(command) <= 0 {
		return TracedErrorf("command has no elements")
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

func (r *RunCommandOptions) SetTimeoutString(timeoutString string) (err error) {
	if timeoutString == "" {
		return TracedErrorf("timeoutString is empty string")
	}

	r.TimeoutString = timeoutString

	return nil
}

func (r *RunCommandOptions) SetVerbose(verbose bool) (err error) {
	r.Verbose = verbose

	return nil
}
