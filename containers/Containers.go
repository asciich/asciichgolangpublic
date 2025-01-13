package containers

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
)

// Returns true if running in a container like docker container.
func IsRunningInsideContainer(verbose bool) (isRunningInContainer bool, err error) {
	const procFilePath string = "/proc/1/cgroup"

	procFile, err := asciichgolangpublic.GetLocalFileByPath(procFilePath)
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
			return false, asciichgolangpublic.TracedErrorf("Unable to parse proc line '%s' from '%s'.", line, procFilePath)
		}

		pathToCheck := splittedLine[2]

		if !aslices.ContainsString([]string{"/", "/init.scope"}, pathToCheck) {
			if verbose {
				asciichgolangpublic.LogInfo("Currently running in a container")
			}
			return true, nil
		}
	}

	if verbose {
		asciichgolangpublic.LogInfo("Currently not running in a container")
	}

	return false, nil
}

// Returns true if running in a container like docker container.
func MustIsRunningInsideContainer(verbose bool) (isRunningInContainer bool) {
	isRunningInContainer, err := IsRunningInsideContainer(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return isRunningInContainer
}
