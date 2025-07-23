package tmux

import (
	"context"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type TmuxSession struct {
	name string
	tmux *TmuxService
}

func NewTmuxSession() (t *TmuxSession) {
	return new(TmuxSession)
}

func (t *TmuxSession) Create(ctx context.Context) (err error) {
	name, err := t.GetName()
	if err != nil {
		return err
	}

	exists, err := t.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Tmux session '%s' already exists. Skip creation.", name)
	} else {
		commandExecutor, err := t.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			commandexecutor.WithLiveOutputOnStdoutIfVerbose(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{
					"tmux",
					"new-session",
					"-d",
					"-s",
					name,
				},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Tmux session '%s' created.", name)
	}

	return nil
}

func (t *TmuxSession) Delete(ctx context.Context) (err error) {
	sessionExists, err := t.Exists(ctx)
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
			ctx,
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

		logging.LogChangedByCtxf(ctx, "Tmux session '%s' deleted.", sessionName)
	} else {
		logging.LogInfoByCtxf(ctx, "Session '%s' already absent. Skip delete.", sessionName)
	}

	return nil
}

func (t *TmuxSession) Exists(ctx context.Context) (exists bool, err error) {
	tmux, err := t.GetTmux()
	if err != nil {
		return false, err
	}

	name, err := t.GetName()
	if err != nil {
		return false, err
	}

	sessionNames, err := tmux.ListSessionNames(ctx)
	if err != nil {
		return false, err
	}

	exists = slices.Contains(sessionNames, name)

	if exists {
		logging.LogInfoByCtxf(ctx, "Tmux session '%s' exists.", name)
	} else {
		logging.LogInfoByCtxf(ctx, "Tmux session '%s' does not exist.", name)
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

func (t *TmuxSession) ListWindowNames(ctx context.Context) (windowsNames []string, err error) {
	name, err := t.GetName()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := t.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	lines, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		ctx,
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

	logging.LogInfoByCtxf(
		ctx,
		"Found '%d' windows in tmux session '%s'.",
		len(windowsNames),
		name,
	)

	return windowsNames, nil
}

func (t *TmuxSession) Recreate(ctx context.Context) (err error) {
	name, err := t.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Recreate tmux session '%s' started.", name)

	err = t.Delete(ctx)
	if err != nil {
		return err
	}

	err = t.Create(ctx)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Recreate tmux session '%s' finished.", name)

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
