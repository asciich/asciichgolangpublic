package containerutils

import (
	"context"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Returns true if running in a container like docker container.
func IsRunningInsideContainer(ctx context.Context) (isRunningInContainer bool, err error) {
	const procFilePath string = "/proc/1/cgroup"

	procFile, err := files.GetLocalFileByPath(procFilePath)
	if err != nil {
		return false, err
	}

	procLines, err := procFile.ReadAsLines(ctx)
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
			logging.LogInfoByCtx(ctx, "Currently running in a container")
			return true, nil
		}
	}

	logging.LogInfoByCtx(ctx, "Currently not running in a container")

	return false, nil
}
