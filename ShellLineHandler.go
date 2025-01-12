package asciichgolangpublic

import (
	"strings"

	"github.com/google/shlex"

	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
)

type ShellLineHandlerService struct {
}

func NewShellLineHandlerService() (s *ShellLineHandlerService) {
	return new(ShellLineHandlerService)
}

func ShellLineHandler() (shellLineHandler *ShellLineHandlerService) {
	return new(ShellLineHandlerService)
}

func (s *ShellLineHandlerService) Join(command []string) (joinedCommand string, err error) {
	if len(command) == 1 {
		return command[0], nil
	}

	commandToJoin := []string{}
	for _, c := range command {
		c = strings.ReplaceAll(c, "'", "'\"'\"'")

		if len(c) <= 0 {
			c = "''"
		}

		if astrings.ContainsAtLeastOneSubstring(c, []string{" ", "\n", "\\n", "\""}) {
			c = "'" + c + "'"
		}

		commandToJoin = append(commandToJoin, c)
	}

	joinedCommand = strings.Join(commandToJoin, " ")
	return joinedCommand, nil
}

func (s *ShellLineHandlerService) MustJoin(command []string) (joinedCommand string) {
	joinedCommand, err := s.Join(command)
	if err != nil {
		LogFatalf("shellLineHandler.Join failed: '%v'", err)
	}

	return joinedCommand
}

func (s *ShellLineHandlerService) MustSplit(command string) (splittedCommand []string) {
	splittedCommand, err := s.Split(command)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return splittedCommand
}

func (s *ShellLineHandlerService) Split(command string) (splittedCommand []string, err error) {
	splittedCommand, err = shlex.Split(command)
	if err != nil {
		return nil, err
	}
	return splittedCommand, nil
}
