package shelllinehandler

import (
	"strings"

	"github.com/google/shlex"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/stringsutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
)

func Join(command []string) (joinedCommand string, err error) {
	if len(command) == 1 {
		return command[0], nil
	}

	commandToJoin := []string{}
	for _, c := range command {
		c = strings.ReplaceAll(c, "'", "'\"'\"'")

		if len(c) <= 0 {
			c = "''"
		}

		if stringsutils.ContainsAtLeastOneSubstring(c, []string{" ", "\n", "\\n", "\""}) {
			c = "'" + c + "'"
		}

		commandToJoin = append(commandToJoin, c)
	}

	joinedCommand = strings.Join(commandToJoin, " ")
	return joinedCommand, nil
}

func MustJoin(command []string) (joinedCommand string) {
	joinedCommand, err := Join(command)
	if err != nil {
		logging.LogFatalf("shellLineHandler.Join failed: '%v'", err)
	}

	return joinedCommand
}

func MustSplit(command string) (splittedCommand []string) {
	splittedCommand, err := Split(command)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return splittedCommand
}

func Split(command string) (splittedCommand []string, err error) {
	splittedCommand, err = shlex.Split(command)
	if err != nil {
		return nil, err
	}
	return splittedCommand, nil
}
