package helm

import (
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type commandExecutorHelm struct {
	commandExecutor commandexecutor.CommandExecutor
}

func GetCommandExecutorHelm(executor commandexecutor.CommandExecutor) (helm Helm, err error) {
	if executor == nil {
		return nil, tracederrors.TracedErrorNil("executor")
	}

	toReturn := NewcommandExecutorHelm()

	err = toReturn.SetCommandExecutor(executor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorHelm() (helm Helm, err error) {
	return GetCommandExecutorHelm(commandexecutor.Bash())
}

func MustGetCommandExecutorHelm(executor commandexecutor.CommandExecutor) (helm Helm) {
	helm, err := GetCommandExecutorHelm(executor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return helm
}

func MustGetLocalCommandExecutorHelm() (helm Helm) {
	helm, err := GetLocalCommandExecutorHelm()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return helm
}

func NewcommandExecutorHelm() (c *commandExecutorHelm) {
	return new(commandExecutorHelm)
}

func (c *commandExecutorHelm) AddRepositoryByName(name string, url string, verbose bool) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	if url == "" {
		return tracederrors.TracedErrorEmptyString("url")
	}

	commandExecutor, hostDescription, err := c.GetCommandExecutorAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Add helm repository '%s' with url '%s' on host '%s' started.",
			name,
			url,
			hostDescription,
		)
	}

	_, err = commandExecutor.RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"helm",
				"repo",
				"add",
				name,
				url,
			},
			Verbose: verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf(
			"Added helm repository '%s' with url '%s' on host '%s'.",
			name,
			url,
			hostDescription,
		)
	}

	if verbose {
		logging.LogInfof(
			"Add helm repository '%s' with url '%s' on host '%s' finished.",
			name,
			url,
			hostDescription,
		)
	}

	return nil
}

func (c *commandExecutorHelm) GetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *commandExecutorHelm) GetCommandExecutorAndHostDescription() (commandExecutor commandexecutor.CommandExecutor, hostDescription string, err error) {
	commandExecutor, err = c.GetCommandExecutor()
	if err != nil {
		return nil, "", err
	}

	hostDescription, err = c.GetHostDescription()
	if err != nil {
		return nil, "", err
	}

	return commandExecutor, hostDescription, nil
}

func (c *commandExecutorHelm) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

func (c *commandExecutorHelm) MustAddRepositoryByName(name string, url string, verbose bool) {
	err := c.AddRepositoryByName(name, url, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *commandExecutorHelm) MustGetCommandExecutor() (commandExecutor commandexecutor.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *commandExecutorHelm) MustGetCommandExecutorAndHostDescription() (commandExecutor commandexecutor.CommandExecutor, hostDescription string) {
	commandExecutor, hostDescription, err := c.GetCommandExecutorAndHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor, hostDescription
}

func (c *commandExecutorHelm) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (c *commandExecutorHelm) MustSetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *commandExecutorHelm) SetCommandExecutor(commandExecutor commandexecutor.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}
