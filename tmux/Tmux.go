package tmux

import (
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type TmuxService struct {
	commandExecutor commandexecutor.CommandExecutor
}

func GetTmuxOnLocalMachine() (tmux *TmuxService, err error) {
	tmux = NewTmuxService()

	err = tmux.SetCommandExecutor(commandexecutor.Bash())
	if err != nil {
		return nil, err
	}

	return tmux, nil
}

func MustGetTmuxOnLocalMachine() (tmux *TmuxService) {
	tmux, err := GetTmuxOnLocalMachine()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tmux
}

func NewTmuxService() (t *TmuxService) {
	return new(TmuxService)
}

func (t *TmuxService) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {

	return t.commandExecutor, nil
}

func (t *TmuxService) GetSessionByName(name string) (tmuxSession *TmuxSession, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
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
		return nil, tracederrors.TracedErrorEmptyString("sessionName")
	}

	if windowName == "" {
		return nil, tracederrors.TracedErrorEmptyString("windowName")
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
		&parameteroptions.RunCommandOptions{
			Command: []string{"tmux", "ls"},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "\nno server running on ") {
			// no sessions avaiable:
			fullSessionLines = []string{}
		} else {
			return nil, err
		}
	}

	sessionNames = []string{}
	for _, line := range fullSessionLines {
		toAdd := strings.Split(line, ":")[0]
		toAdd = strings.TrimSpace(toAdd)

		if toAdd == "" {
			return nil, tracederrors.TracedErrorf(
				"toAdd is empty string after extracting session name from line='%s'",
				line,
			)
		}

		sessionNames = append(sessionNames, toAdd)
	}

	if verbose {
		logging.LogInfof(
			"There are '%d' tmux sessions.",
			len(sessionNames),
		)
	}

	return sessionNames, nil
}

func (t *TmuxService) MustGetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor) {
	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (t *TmuxService) MustGetSessionByName(name string) (tmuxSession *TmuxSession) {
	tmuxSession, err := t.GetSessionByName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tmuxSession
}

func (t *TmuxService) MustGetWindowByNames(sessionName string, windowName string) (window *TmuxWindow) {
	window, err := t.GetWindowByNames(sessionName, windowName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return window
}

func (t *TmuxService) MustListSessionNames(verbose bool) (sessionNames []string) {
	sessionNames, err := t.ListSessionNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sessionNames
}

func (t *TmuxService) MustSetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) {
	err := t.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxService) SetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) (err error) {
	t.commandExecutor = commandExecutor

	return nil
}

// Returns true if input string is a valid tmux key like "enter".
// Returns false otherwise.
func IsTmuxKey(input string) (isKey bool) {
	knownKeys := []string{
		"enter", // Enter key
		"C-c",   // CTRL + C
		"C-d",   // CTRL + D
		"C-l",   // CTRL + L
	}

	return slices.Contains(knownKeys, input)
}
