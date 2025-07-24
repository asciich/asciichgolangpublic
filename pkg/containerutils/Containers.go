package containerutils

import (
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Returns true if running in a container like docker container.
func IsRunningInsideContainer(verbose bool) (isRunningInContainer bool, err error) {
	const procFilePath string = "/proc/1/cgroup"

	procFile, err := files.GetLocalFileByPath(procFilePath)
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
			return false, tracederrors.TracedErrorf("Unable to parse proc line '%s' from '%s'.", line, procFilePath)
		}

		pathToCheck := splittedLine[2]

		if !slices.Contains([]string{"/", "/init.scope"}, pathToCheck) {
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
