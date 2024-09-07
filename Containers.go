package asciichgolangpublic

import "strings"

type ContainersService struct{}

func Contaners() (c *ContainersService) {
	return NewContainersService()
}

// Returns true if running in a container like docker container.
func IsRunningInsideContainer(verbose bool) (isRunningInContainer bool, err error) {
	isRunningInContainer, err = Contaners().IsRunningInsideContainer(verbose)
	if err != nil {
		return false, err
	}

	return isRunningInContainer, nil
}

// Returns true if running in a container like docker container.
func MustIsRunningInsideContainer(verbose bool) (isRunningInContainer bool) {
	isRunningInContainer, err := IsRunningInsideContainer(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isRunningInContainer
}

func NewContainersService() (c *ContainersService) {
	return new(ContainersService)
}

// Returns true if running in a container like docker container.
func (c *ContainersService) IsRunningInsideContainer(verbose bool) (isRunningInContainer bool, err error) {
	const procFilePath string = "/proc/1/cgroup"

	procFile, err := GetLocalFileByPath(procFilePath)
	if err != nil {
		return false, err
	}

	procLines, err := procFile.ReadAsLines()
	if err != nil {
		return false, err
	}

	for _, line := range procLines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		splittedLine := strings.Split(line, ":")
		if len(splittedLine) != 3 {
			return false, TracedErrorf("Unable to parse proc line '%s' from '%s'.", line, procFilePath)
		}

		pathToCheck := splittedLine[2]

		if !Slices().ContainsString([]string{"/", "/init.scope"}, pathToCheck) {
			if verbose {
				LogInfo("Currently running in a container")
			}
			return true, nil
		}
	}

	if verbose {
		LogInfo("Currently not running in a container")
	}

	return false, nil
}

// Returns true if running in a container like docker container.
func (c *ContainersService) MustIsRunningInsideContainer(verbose bool) (isRunningInContainer bool) {
	isRunningInContainer, err := c.IsRunningInsideContainer(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isRunningInContainer
}
