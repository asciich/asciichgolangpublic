package asciichgolangpublic

import "strings"

type TmuxService struct {
	commandExecutor CommandExecutor
}

func GetTmuxOnLocalMachine() (tmux *TmuxService, err error) {
	tmux = NewTmuxService()

	err = tmux.SetCommandExecutor(Bash())
	if err != nil {
		return nil, err
	}

	return tmux, nil
}

func MustGetTmuxOnLocalMachine() (tmux *TmuxService) {
	tmux, err := GetTmuxOnLocalMachine()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tmux
}

func NewTmuxService() (t *TmuxService) {
	return new(TmuxService)
}

func (t *TmuxService) GetCommandExecutor() (commandExecutor CommandExecutor, err error) {

	return t.commandExecutor, nil
}

func (t *TmuxService) GetSessionByName(name string) (tmuxSession *TmuxSession, err error) {
	if name == "" {
		return nil, TracedErrorEmptyString("name")
	}

	tmuxSession = NewTmuxSession()

	err = tmuxSession.SetTmux(t)
	if err != nil {
		return nil, err
	}

	err = tmuxSession.SetName(name)
	if err != nil {
		return nil, err
	}

	return tmuxSession, err
}

func (t *TmuxService) GetWindowByNames(sessionName string, windowName string) (window *TmuxWindow, err error) {
	if sessionName == "" {
		return nil, TracedErrorEmptyString("sessionName")
	}

	if windowName == "" {
		return nil, TracedErrorEmptyString("windowName")
	}

	session, err := t.GetSessionByName(sessionName)
	if err != nil {
		return nil, err
	}

	window, err = session.GetWindowByName(windowName)
	if err != nil {
		return nil, err
	}

	return window, nil
}

func (t *TmuxService) ListSessionNames(verbose bool) (sessionNames []string, err error) {
	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	fullSessionLines, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		&RunCommandOptions{
			Command:            []string{"tmux", "ls"},
			LiveOutputOnStdout: verbose,
			Verbose:            verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	sessionNames = []string{}
	for _, line := range fullSessionLines {
		toAdd := strings.Split(line, ":")[0]
		toAdd = strings.TrimSpace(toAdd)

		if toAdd == "" {
			return nil, TracedErrorf(
				"toAdd is empty string after extracting session name from line='%s'",
				line,
			)
		}

		sessionNames = append(sessionNames, toAdd)
	}

	if verbose {
		LogInfof(
			"There are '%d' tmux sessions.",
			len(sessionNames),
		)
	}

	return sessionNames, nil
}

func (t *TmuxService) MustGetCommandExecutor() (commandExecutor CommandExecutor) {
	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (t *TmuxService) MustGetSessionByName(name string) (tmuxSession *TmuxSession) {
	tmuxSession, err := t.GetSessionByName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tmuxSession
}

func (t *TmuxService) MustGetWindowByNames(sessionName string, windowName string) (window *TmuxWindow) {
	window, err := t.GetWindowByNames(sessionName, windowName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return window
}

func (t *TmuxService) MustListSessionNames(verbose bool) (sessionNames []string) {
	sessionNames, err := t.ListSessionNames(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sessionNames
}

func (t *TmuxService) MustSetCommandExecutor(commandExecutor CommandExecutor) {
	err := t.SetCommandExecutor(commandExecutor)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (t *TmuxService) SetCommandExecutor(commandExecutor CommandExecutor) (err error) {
	t.commandExecutor = commandExecutor

	return nil
}
