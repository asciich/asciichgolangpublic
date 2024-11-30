package asciichgolangpublic

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
		return TracedError("toSend has no elements")
	}

	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return err
	}

	commandToUse := append([]string{"tmux", "send-keys", "-t", sessionName + ":" + windowName}, toSend...)

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: commandToUse,
			Verbose: verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf(
			"Send keys to tmux window '%s' in session '%s'.",
			windowName,
			sessionName,
		)
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
		return "", TracedErrorf(
			"No lines for tmux window '%s' in session '%s' received.",
			windowName,
			sessionName,
		)
	}

	if len(lines) <= 1 {
		return "", TracedErrorf(
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
			LogInfof(
				"Tmux window '%s' in session '%s' already exists. Skip create.",
				windowName,
				sessionName,
			)
		}
	} else {
		commandExecutor, err := t.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{"tmux", "new-window", "-t", sessionName, "-n", windowName},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			LogChangedf(
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
			&RunCommandOptions{
				Command: []string{"tmux", "kill-window", "-t", sessionName + ":" + windowName},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		LogChangedf(
			"Tmux window '%s' of session '%s' is already deleted.",
			windowName,
			sessionName,
		)
	} else {
		LogInfof(
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

	exists = Slices().ContainsString(windowNames, windowName)

	sessionName, err := t.GetSessionName()
	if err != nil {
		return false, err
	}

	if verbose {
		if exists {
			LogInfof(
				"Window '%s' exists in tmux session '%s'.",
				windowName,
				sessionName,
			)
		} else {
			LogInfof(
				"Window '%s' does not exist in tmux session '%s'.",
				windowName,
				sessionName,
			)
		}
	}

	return exists, nil
}

func (t *TmuxWindow) GetCommandExecutor() (commandExecutor CommandExecutor, err error) {
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
		return "", TracedErrorf(
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
		return "", TracedErrorf("name not set")
	}

	return t.name, nil
}

func (t *TmuxWindow) GetSession() (session *TmuxSession, err error) {
	if t.session == nil {
		return nil, TracedErrorf("session not set")
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

func (t *TmuxWindow) GetShownLines() (lines []string, err error) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	lines, err = commandExecutor.RunCommandAndGetStdoutAsLines(
		&RunCommandOptions{
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
		return nil, err
	}

	lines = Slices().RemoveEmptyStringsAtEnd(lines)

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
		LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustDelete(verbose bool) {
	err := t.Delete(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustExists(verbose bool) (exists bool) {
	exists, err := t.Exists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (t *TmuxWindow) MustGetCommandExecutor() (commandExecutor CommandExecutor) {
	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (t *TmuxWindow) MustGetLatestPaneLine() (paneLine string) {
	paneLine, err := t.GetLatestPaneLine()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return paneLine
}

func (t *TmuxWindow) MustGetName() (name string) {
	name, err := t.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (t *TmuxWindow) MustGetSecondLatestPaneLine() (paneLine string) {
	paneLine, err := t.GetSecondLatestPaneLine()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return paneLine
}

func (t *TmuxWindow) MustGetSession() (session *TmuxSession) {
	session, err := t.GetSession()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return session
}

func (t *TmuxWindow) MustGetSessionAndWindowName() (sessionName string, windowName string) {
	sessionName, windowName, err := t.GetSessionAndWindowName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sessionName, windowName
}

func (t *TmuxWindow) MustGetSessionName() (sessionName string) {
	sessionName, err := t.GetSessionName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sessionName
}

func (t *TmuxWindow) MustGetShownLines() (lines []string) {
	lines, err := t.GetShownLines()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return lines
}

func (t *TmuxWindow) MustListWindowNames(verbose bool) (windowNames []string) {
	windowNames, err := t.ListWindowNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return windowNames
}

func (t *TmuxWindow) MustSendKeys(toSend []string, verbose bool) {
	err := t.SendKeys(toSend, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustSetName(name string) {
	err := t.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) MustSetSession(session *TmuxSession) {
	err := t.SetSession(session)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (t *TmuxWindow) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	t.name = name

	return nil
}

func (t *TmuxWindow) SetSession(session *TmuxSession) (err error) {
	if session == nil {
		return TracedErrorf("session is nil")
	}

	t.session = session

	return nil
}
