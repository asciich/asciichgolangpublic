package helm

import (
	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type commandExecutorHelm struct {
	commandExecutor asciichgolangpublic.CommandExecutor
}

func GetCommandExecutorHelm(executor asciichgolangpublic.CommandExecutor) (helm Helm, err error) {
	if executor == nil {
		return nil, errors.TracedErrorNil("executor")
	}

	toReturn := NewcommandExecutorHelm()

	err = toReturn.SetCommandExecutor(executor)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func GetLocalCommandExecutorHelm() (helm Helm, err error) {
	return GetCommandExecutorHelm(asciichgolangpublic.Bash())
}

func MustGetCommandExecutorHelm(executor asciichgolangpublic.CommandExecutor) (helm Helm) {
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
		return errors.TracedErrorEmptyString("name")
	}

	if url == "" {
		return errors.TracedErrorEmptyString("url")
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
		&asciichgolangpublic.RunCommandOptions{
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

func (c *commandExecutorHelm) GetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *commandExecutorHelm) GetCommandExecutorAndHostDescription() (commandExecutor asciichgolangpublic.CommandExecutor, hostDescription string, err error) {
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

func (c *commandExecutorHelm) MustGetCommandExecutor() (commandExecutor asciichgolangpublic.CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *commandExecutorHelm) MustGetCommandExecutorAndHostDescription() (commandExecutor asciichgolangpublic.CommandExecutor, hostDescription string) {
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

func (c *commandExecutorHelm) MustSetCommandExecutor(commandExecutor asciichgolangpublic.CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *commandExecutorHelm) SetCommandExecutor(commandExecutor asciichgolangpublic.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}
