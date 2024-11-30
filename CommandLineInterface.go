package asciichgolangpublic

import "strings"

type CommandLineInterfaceService struct {
}

func CommandLineInterface() (cli *CommandLineInterfaceService) {
	return NewCommandLineInterfaceService()
}

func NewCommandLineInterfaceService() (c *CommandLineInterfaceService) {
	return new(CommandLineInterfaceService)
}

func (c *CommandLineInterfaceService) IsLinePromptOnly(line string) (isPromptOnly bool) {
	stripped := strings.TrimSpace(line)

	if stripped == "" {
		return false
	}

	if strings.Contains(stripped, "\n") {
		return false
	}

	if !strings.HasSuffix(stripped, "$") {
		return false
	}

	return true
}
