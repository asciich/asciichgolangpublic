package tmux

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/commandlineinterface"
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tempfiles"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var ErrTmuxWindowCliPromptNotReady = errors.New("tmux window CLI promptnot ready")

type TmuxWindow struct {
	name    string
	session *TmuxSession
}

func NewTmuxWindow() (t *TmuxWindow) {
	return new(TmuxWindow)
}

// Default use case to send a command is using []string{"command to run", "enter"}. "enter" in this example is detected as enter key by tmux.
func (t *TmuxWindow) SendKeys(toSend []string, verbose bool) (err error) {
	if len(toSend) <= 0 {
		return tracederrors.TracedError("toSend has no elements")
	}

	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return err
	}

	for _, t := range toSend {
		commandToUse := []string{"tmux", "send-keys", "-t", sessionName + ":" + windowName}

		if IsTmuxKey(t) {
			commandToUse = append(commandToUse, t)
		} else {
			hexEncoded := stringsutils.ToHexStringSlice(t)
			commandToUse = append(commandToUse, "-H")
			commandToUse = append(commandToUse, hexEncoded...)
		}

		_, err = commandExecutor.RunCommand(
			&parameteroptions.RunCommandOptions{
				Command: commandToUse,
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}
	}

	if verbose {
		logging.LogChangedf(
			"Send keys to tmux window '%s' in session '%s'.",
			windowName,
			sessionName,
		)
	}

	return nil
}

// Delete the tmux session this window belongs to.
// Will implicitly also kill this window but also any other window in the session.
func (t *TmuxWindow) DeleteSession(verbose bool) (err error) {
	session, err := t.GetSession()
	if err != nil {
		return err
	}

	err = session.Delete(verbose)
	if err != nil {
		return err
	}

	return nil
}

// Since the latest line usually shows the command prompt this command can be used to receive the latest printed line.
func (t *TmuxWindow) GetSecondLatestPaneLine() (paneLine string, err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return "", err
	}

	lines, err := t.GetShownLines()
	if err != nil {
		return "", err
	}

	if len(lines) <= 0 {
		return "", tracederrors.TracedErrorf(
			"No lines for tmux window '%s' in session '%s' received.",
			windowName,
			sessionName,
		)
	}

	if len(lines) <= 1 {
		return "", tracederrors.TracedErrorf(
			"Only one line for tmux window '%s' in session '%s' received.",
			windowName,
			sessionName,
		)
	}

	paneLine = lines[len(lines)-2]

	return paneLine, nil
}

func (t *TmuxWindow) Create(verbose bool) (err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	exists, err := t.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		if verbose {
			logging.LogInfof(
				"Tmux window '%s' in session '%s' already exists. Skip create.",
				windowName,
				sessionName,
			)
		}
	} else {
		session, err := t.GetSession()
		if err != nil {
			return err
		}

		err = session.Create(verbose)
		if err != nil {
			return err
		}

		commandExecutor, err := t.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&parameteroptions.RunCommandOptions{
				Command: []string{"tmux", "new-window", "-t", sessionName, "-n", windowName},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Tmux window '%s' in session '%s' created.",
				windowName,
				sessionName,
			)
		}
	}

	return nil
}

func (t *TmuxWindow) Delete(verbose bool) (err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	exists, err := t.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		commandExecutor, err := t.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&parameteroptions.RunCommandOptions{
				Command: []string{"tmux", "kill-window", "-t", sessionName + ":" + windowName},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Tmux window '%s' of session '%s' is already deleted.",
			windowName,
			sessionName,
		)
	} else {
		logging.LogInfof(
			"Tmux window '%s' of session '%s' is already absent. Skip delete.",
			windowName,
			sessionName,
		)
	}

	return nil
}

func (t *TmuxWindow) Exists(verbose bool) (exists bool, err error) {
	windowName, err := t.GetName()
	if err != nil {
		return false, err
	}

	windowNames, err := t.ListWindowNames(verbose)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(windowNames, windowName)

	sessionName, err := t.GetSessionName()
	if err != nil {
		return false, err
	}

	if verbose {
		if exists {
			logging.LogInfof(
				"Window '%s' exists in tmux session '%s'.",
				windowName,
				sessionName,
			)
		} else {
			logging.LogInfof(
				"Window '%s' does not exist in tmux session '%s'.",
				windowName,
				sessionName,
			)
		}
	}

	return exists, nil
}

func (t *TmuxWindow) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {
	session, err := t.GetSession()
	if err != nil {
		return nil, err
	}

	commandExecutor, err = session.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor, nil
}

func (t *TmuxWindow) GetLatestPaneLine() (paneLine string, err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return "", err
	}

	lines, err := t.GetShownLines()
	if err != nil {
		return "", err
	}

	if len(lines) <= 0 {
		return "", tracederrors.TracedErrorf(
			"No lines for tmux window '%s' in session '%s' received.",
			windowName,
			sessionName,
		)
	}

	paneLine = lines[len(lines)-1]

	return paneLine, nil
}

func (t *TmuxWindow) GetName() (name string, err error) {
	if t.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return t.name, nil
}

func (t *TmuxWindow) GetSession() (session *TmuxSession, err error) {
	if t.session == nil {
		return nil, tracederrors.TracedErrorf("session not set")
	}

	return t.session, nil
}

func (t *TmuxWindow) GetSessionAndWindowName() (sessionName string, windowName string, err error) {
	sessionName, err = t.GetSessionName()
	if err != nil {
		return "", "", err
	}

	windowName, err = t.GetName()
	if err != nil {
		return "", "", err
	}

	return sessionName, windowName, nil
}

func (t *TmuxWindow) GetSessionName() (sessionName string, err error) {
	session, err := t.GetSession()
	if err != nil {
		return "", err
	}

	sessionName, err = session.GetName()
	if err != nil {
		return "", err
	}

	return sessionName, nil
}

func (t *TmuxWindow) GetShownOutput() (output string, err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return "", err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	output, err = commandExecutor.RunCommandAndGetStdoutAsString(
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"tmux",
				"capture-pane",
				"-t",
				sessionName + ":" + windowName,
				"-J",
				"-p",
			},
		},
	)
	if err != nil {
		return "", err
	}

	return output, nil
}

func (t *TmuxWindow) GetShownLines() (lines []string, err error) {
	output, err := t.GetShownOutput()
	if err != nil {
		return nil, err
	}

	lines = stringsutils.SplitLines(output, true)
	lines = slicesutils.RemoveEmptyStringsAtEnd(lines)

	return lines, nil
}

func (t *TmuxWindow) ListWindowNames(verbose bool) (windowNames []string, err error) {
	session, err := t.GetSession()
	if err != nil {
		return nil, err
	}

	windowNames, err = session.ListWindowNames(verbose)
	if err != nil {
		return nil, err
	}

	return windowNames, nil
}

func (t *TmuxWindow) MustCreate(verbose bool) {
	err := t.Create(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustDelete(verbose bool) {
	err := t.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustDeleteSession(verbose bool) {
	err := t.DeleteSession(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustExists(verbose bool) (exists bool) {
	exists, err := t.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (t *TmuxWindow) MustGetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor) {
	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (t *TmuxWindow) MustGetLatestPaneLine() (paneLine string) {
	paneLine, err := t.GetLatestPaneLine()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return paneLine
}

func (t *TmuxWindow) MustGetName() (name string) {
	name, err := t.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (t *TmuxWindow) MustGetSecondLatestPaneLine() (paneLine string) {
	paneLine, err := t.GetSecondLatestPaneLine()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return paneLine
}

func (t *TmuxWindow) MustGetSession() (session *TmuxSession) {
	session, err := t.GetSession()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return session
}

func (t *TmuxWindow) MustGetSessionAndWindowName() (sessionName string, windowName string) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sessionName, windowName
}

func (t *TmuxWindow) MustGetSessionName() (sessionName string) {
	sessionName, err := t.GetSessionName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sessionName
}

func (t *TmuxWindow) MustGetShownLines() (lines []string) {
	lines, err := t.GetShownLines()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return lines
}

func (t *TmuxWindow) MustListWindowNames(verbose bool) (windowNames []string) {
	windowNames, err := t.ListWindowNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return windowNames
}

func (t *TmuxWindow) MustRecreate(verbose bool) {
	err := t.Recreate(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustRunCommand(runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput) {
	commandOutput, err := t.RunCommand(runCommandOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (t *TmuxWindow) WaitUntilOutputMatchesRegex(regex string, timeout time.Duration, verbose bool) (err error) {
	if regex == "" {
		return tracederrors.TracedErrorEmptyString("regex")
	}

	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	tStart := time.Now()

	for {
		if time.Since(tStart) > timeout {
			return tracederrors.TracedErrorf(
				"Timeout (%v) while waiting for tmux terminal output of '%s:%s' matches regex '%s'.",
				timeout,
				sessionName,
				windowName,
				regex,
			)
		}

		output, err := t.GetShownOutput()
		if err != nil {
			return err
		}

		matches, err := stringsutils.MatchesRegex(output, regex)
		if err != nil {
			return err
		}

		if matches {
			break
		}

		retryDelay := time.Millisecond * 100

		if verbose {
			logging.LogInfof(
				"Tmux output of '%s:%s' does not match regex '%s'. Going to retry in '%s'",
				sessionName,
				windowName,
				regex,
				retryDelay,
			)
		}

		time.Sleep(retryDelay)
	}

	if verbose {
		logging.LogInfof(
			"Tmux terminal output of '%s:%s' now matches regex '%s'.",
			sessionName,
			windowName,
			regex,
		)
	}

	return nil
}

func (t *TmuxWindow) IsOutputMatchingRegex(regex string, verbose bool) (isMatching bool, err error) {
	if regex == "" {
		return false, tracederrors.TracedErrorEmptyString("regex")
	}

	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return false, err
	}

	output, err := t.GetShownOutput()
	if err != nil {
		return false, err
	}

	isMatching, err = stringsutils.MatchesRegex(output, regex)
	if err != nil {
		return false, err
	}

	if verbose {
		if isMatching {
			logging.LogInfof(
				"Output of tmux window '%s:%s' matches regex '%s'.",
				sessionName,
				windowName,
				regex,
			)
		} else {
			logging.LogInfof(
				"Output of tmux window '%s:%s' does not match regex '%s'.",
				sessionName,
				windowName,
				regex,
			)
		}
	}

	return isMatching, nil
}

func (t *TmuxWindow) MustIsOutputMatchingRegex(regex string, verbose bool) (isMatching bool) {
	isMatching, err := t.IsOutputMatchingRegex(regex, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isMatching
}

func (t *TmuxWindow) MustWaitUntilOutputMatchesRegex(regex string, timeout time.Duration, verbose bool) {
	err := t.WaitUntilOutputMatchesRegex(regex, timeout, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustSendKeys(toSend []string, verbose bool) {
	err := t.SendKeys(toSend, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustSetName(name string) {
	err := t.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustSetSession(session *TmuxSession) {
	err := t.SetSession(session)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustWaitUntilCliPromptReady(verbose bool) {
	err := t.WaitUntilCliPromptReady(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) Recreate(verbose bool) (err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	err = t.Delete(verbose)
	if err != nil {
		return err
	}

	err = t.Create(verbose)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf(
			"Tmux window '%s' in session '%s' recreated.",
			windowName,
			sessionName,
		)
	}

	return nil
}

func (t *TmuxWindow) RunCommand(runCommandOptions *parameteroptions.RunCommandOptions) (commandOutput *commandexecutor.CommandOutput, err error) {
	if runCommandOptions == nil {
		return nil, tracederrors.TracedErrorNil("runCommandOptions")
	}

	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return nil, err
	}

	if runCommandOptions.Verbose {
		logging.LogInfof(
			"Run command in tmux window '%s' of tmux session '%s' started.",
			windowName,
			sessionName,
		)
	}

	err = t.Create(runCommandOptions.Verbose)
	if err != nil {
		return nil, err
	}

	captureFile, err := tempfiles.CreateEmptyTemporaryFile(runCommandOptions.Verbose)
	if err != nil {
		return nil, err
	}

	captureFilePath, err := captureFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	err = t.WaitUntilCliPromptReady(runCommandOptions.Verbose)
	if err != nil {
		return nil, err
	}

	// start output capture
	_, err = commandExecutor.RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"tmux",
				"pipe-pane",
				"-t",
				sessionName + ":" + windowName,
				"cat > '" + captureFilePath + "'",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if runCommandOptions.Verbose {
		logging.LogInfof("'%s' will be used to capture tmux output for command to run.", captureFilePath)
	}

	commandToSend, err := runCommandOptions.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	const endCommandMarkerPrefix = "Last command ended exited with status code"
	err = t.SendKeys(
		[]string{" " + commandToSend + "; echo -en \"\\n" + endCommandMarkerPrefix + " $?\\n\"", "enter"},
		runCommandOptions.Verbose,
	)
	if err != nil {
		return nil, err
	}

	// Wait for command to finish
	for {
		lines, err := t.GetShownLines()
		if err != nil {
			return nil, err
		}

		if len(lines) > 0 {
			if strings.HasPrefix(lines[len(lines)-1], endCommandMarkerPrefix) {
				if runCommandOptions.Verbose {
					logging.LogInfo("Found endCommandMarkerPrefix in latest line. Command is finished.")
				}
				break
			}
		}

		if len(lines) > 1 {
			if strings.HasPrefix(lines[len(lines)-2], endCommandMarkerPrefix) {
				if runCommandOptions.Verbose {
					logging.LogInfo("Found endCommandMarkerPrefix in second last latest line. Command is finished.")
				}
				break
			}
		}

		waitTime := time.Millisecond * 200
		if runCommandOptions.Verbose {
			logging.LogInfof("Wait another '%v' until command is finished.", waitTime)
		}

		time.Sleep(waitTime)
	}

	// stop output capture
	_, err = commandExecutor.RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"tmux",
				"pipe-pane",
				"-t",
				sessionName + ":" + windowName,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	allOutputLines, err := captureFile.ReadAsLines()
	if err != nil {
		return nil, err
	}

	allOutputLines = slicesutils.RemoveEmptyStringsAtEnd(allOutputLines)

	if len(allOutputLines) < 2 {
		return nil, tracederrors.TracedErrorf(
			"Unable to parse tmux command output: allOutputLines is '%v'",
			allOutputLines,
		)
	}

	// Remove first line as it only contains the command input.
	outputLines := allOutputLines[1:]

	if len(outputLines) < 1 {
		return nil, tracederrors.TracedErrorf(
			"Unable to parse tmux command output: outputLines is '%v'",
			outputLines,
		)
	}

	// extract exit code
	exitCodeLine := outputLines[len(outputLines)-1]
	splitted := strings.Split(exitCodeLine, " ")
	if len(splitted) <= 2 {
		return nil, tracederrors.TracedErrorf(
			"Unable to parse tmux command output: splitted is '%v'",
			splitted,
		)
	}

	exitCodeString := splitted[len(splitted)-1]

	exitCode, err := strconv.Atoi(exitCodeString)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to parse exitCodeString='%s' %w",
			exitCodeString,
			err,
		)
	}

	outputLines = outputLines[:len(outputLines)-1]

	// remove escape sequence in first line of command output:
	if len(splitted) > 0 {
		firstLine := outputLines[0]

		if strings.HasPrefix(firstLine, "\x1b") {
			splitted := strings.SplitN(firstLine, "\r", 2)
			firstLine = splitted[len(splitted)-1]
		}

		outputLines[0] = firstLine
	}

	stdout := strings.Join(outputLines, "\n")

	commandOutput = commandexecutor.NewCommandOutput()

	err = commandOutput.SetStdoutByString(stdout)
	if err != nil {
		return nil, err
	}

	err = commandOutput.SetReturnCode(exitCode)
	if err != nil {
		return nil, err
	}

	err = commandOutput.CheckExitSuccess(false)
	if err != nil {
		return nil, err
	}

	if runCommandOptions.Verbose {
		logging.LogInfof(
			"Run command in tmux window '%s' of tmux session '%s' finished.",
			windowName,
			sessionName,
		)
	}

	return commandOutput, nil
}

func (t *TmuxWindow) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	t.name = name

	return nil
}

func (t *TmuxWindow) SetSession(session *TmuxSession) (err error) {
	if session == nil {
		return tracederrors.TracedErrorf("session is nil")
	}

	t.session = session

	return nil
}

func (t *TmuxWindow) WaitUntilCliPromptReady(verbose bool) (err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	nTries := 30
	for i := 0; i < nTries; i++ {
		lines, err := t.GetShownLines()
		if err != nil {
			return err
		}

		if len(lines) > 0 {
			lastLine := lines[len(lines)-1]

			if commandlineinterface.IsLinePromptOnly(lastLine) {
				if verbose {
					logging.LogInfof(
						"Tmux window '%s' in session '%s' shows CLI prompt and is ready to use.",
						windowName,
						sessionName,
					)
				}

				return nil
			}
		}

		delayTime := 100 * time.Millisecond

		if verbose {
			logging.LogInfof(
				"Wait '%v' before tmux window '%s' in session '%s' becomes ready (%d/%d).",
				delayTime,
				windowName,
				sessionName,
				i+1,
				nTries,
			)
		}

		time.Sleep(delayTime)
	}

	return tracederrors.TracedErrorf(
		"%w: tmux window '%s' in session '%s' is not ready",
		ErrTmuxWindowCliPromptNotReady,
		windowName,
		sessionName,
	)
}
