package containers

import (
	"strings"

	"github.com/asciich/asciichgolangpublic"
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
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
			return false, errors.TracedErrorf("Unable to parse proc line '%s' from '%s'.", line, procFilePath)
		}

		pathToCheck := splittedLine[2]

		if !aslices.ContainsString([]string{"/", "/init.scope"}, pathToCheck) {
			if verbose {
				logging.LogInfo("Currently running in a container")
			}
			return true, nil
		}
	}

	if verbose {
		logging.LogInfo("Currently not running in a container")
	}

	return false, nil
}

// Returns true if running in a container like docker container.
func MustIsRunningInsideContainer(verbose bool) (isRunningInContainer bool) {
	isRunningInContainer, err := IsRunningInsideContainer(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunningInContainer
}
