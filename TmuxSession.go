package asciichgolangpublic

import (
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type TmuxSession struct {
	name string
	tmux *TmuxService
}

func NewTmuxSession() (t *TmuxSession) {
	return new(TmuxSession)
}

func (t *TmuxSession) Create(verbose bool) (err error) {
	name, err := t.GetName()
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
				"Tmux session '%s' already exists. Skip creation.",
				name,
			)
		}
	} else {
		commandExecutor, err := t.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"tmux",
					"new-session",
					"-d",
					"-s",
					name,
				},
				Verbose:            verbose,
				LiveOutputOnStdout: verbose,
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Tmux session '%s' created.",
				name,
			)
		}
	}

	return nil
}

func (t *TmuxSession) Delete(verbose bool) (err error) {
	sessionExists, err := t.Exists(verbose)
	if err != nil {
		return err
	}

	sessionName, err := t.GetName()
	if err != nil {
		return err
	}

	if sessionExists {
		commandExecutor, err := t.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"tmux",
					"kill-session",
					"-t",
					sessionName,
				},
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Tmux session '%s' deleted.",
				sessionName,
			)
		}
	} else {
		if verbose {
			logging.LogInfof(
				"Session '%s' already absent. Skip delete.",
				sessionName,
			)
		}
	}

	return nil
}

func (t *TmuxSession) Exists(verbose bool) (exists bool, err error) {
	tmux, err := t.GetTmux()
	if err != nil {
		return false, err
	}

	name, err := t.GetName()
	if err != nil {
		return false, err
	}

	sessionNames, err := tmux.ListSessionNames(verbose)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(sessionNames, name)

	if exists {
		if verbose {
			logging.LogInfof(
				"Tmux session '%s' exists.",
				name,
			)
		}
	} else {
		if verbose {
			logging.LogInfof(
				"Tmux session '%s' does not exist.",
				name,
			)
		}
	}

	return exists, nil
}

func (t *TmuxSession) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {
	tmux, err := t.GetTmux()
	if err != nil {
		return nil, err
	}

	commandExecutor, err = tmux.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandExecutor, err
}

func (t *TmuxSession) GetName() (name string, err error) {
	if t.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return t.name, nil
}

func (t *TmuxSession) GetTmux() (tmux *TmuxService, err error) {
	if t.tmux == nil {
		return nil, tracederrors.TracedErrorf("tmux not set")
	}

	return t.tmux, nil
}

func (t *TmuxSession) GetWindowByName(windowName string) (window *TmuxWindow, err error) {
	if windowName == "" {
		return nil, tracederrors.TracedErrorEmptyString("windowName")
	}

	window = NewTmuxWindow()

	err = window.SetName(windowName)
	if err != nil {
		return nil, err
	}

	err = window.SetSession(t)
	if err != nil {
		return nil, err
	}

	return window, nil
}

func (t *TmuxSession) ListWindowNames(verbose bool) (windowsNames []string, err error) {
	name, err := t.GetName()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	lines, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		&parameteroptions.RunCommandOptions{
			Command: []string{"tmux", "list-windows", "-a"},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "\nno server running on ") {
			// no sever running means no sessions and windows available:
			lines = []string{}
		} else {
			return nil, err
		}
	}

	windowsNames = []string{}

	for _, l := range lines {
		if strings.HasPrefix(l, name+":") {
			splitted := strings.Split(l, ":")
			if len(splitted) < 3 {
				return nil, tracederrors.TracedErrorf("Unable to get window name out of line='%s'", l)
			}

			windowInfoString := strings.TrimSpace(splitted[2])

			toAdd := strings.Split(windowInfoString, " ")[0]
			toAdd = strings.TrimSuffix(toAdd, "*")

			windowsNames = append(windowsNames, toAdd)
		}
	}

	if verbose {
		logging.LogInfof(
			"Found '%d' windows in tmux session '%s'.",
			len(windowsNames),
			name,
		)
	}

	return windowsNames, nil
}

func (t *TmuxSession) MustCreate(verbose bool) {
	err := t.Create(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxSession) MustDelete(verbose bool) {
	err := t.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxSession) MustExists(verbose bool) (exists bool) {
	exists, err := t.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (t *TmuxSession) MustGetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor) {
	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (t *TmuxSession) MustGetName() (name string) {
	name, err := t.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (t *TmuxSession) MustGetTmux() (tmux *TmuxService) {
	tmux, err := t.GetTmux()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tmux
}

func (t *TmuxSession) MustGetWindowByName(windowName string) (window *TmuxWindow) {
	window, err := t.GetWindowByName(windowName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return window
}

func (t *TmuxSession) MustListWindowNames(verbose bool) (windowsNames []string) {
	windowsNames, err := t.ListWindowNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return windowsNames
}

func (t *TmuxSession) MustRecreate(verbose bool) {
	err := t.Recreate(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxSession) MustSetName(name string) {
	err := t.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxSession) MustSetTmux(tmux *TmuxService) {
	err := t.SetTmux(tmux)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TmuxSession) Recreate(verbose bool) (err error) {
	name, err := t.GetName()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Recreate tmux session '%s' started.", name)
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
		logging.LogInfof("Recreate tmux session '%s' finished.", name)
	}

	return nil
}

func (t *TmuxSession) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	t.name = name

	return nil
}

func (t *TmuxSession) SetTmux(tmux *TmuxService) (err error) {
	if tmux == nil {
		return tracederrors.TracedErrorf("tmux is nil")
	}

	t.tmux = tmux

	return nil
}
