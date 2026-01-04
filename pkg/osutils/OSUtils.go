package osutils

import (
	"os"
	"runtime"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetCurrentWorkingDirectoryAsString() (workingDirPath string, err error) {
	workingDirPath, err = os.Getwd()
	if err != nil {
		return "", tracederrors.TracedErrorf("Get working directory failed: %w", err)
	}

	if !pathsutils.IsAbsolutePath(workingDirPath) {
		return "", tracederrors.TracedErrorf(
			"Evaluated working directory path '%s' is not an absolute path after evaluation.",
			workingDirPath,
		)
	}

	return workingDirPath, nil
}

func IsRunningOnLinux() (isRunningOnLinux bool) {
	return runtime.GOOS == "linux"
}

func MustGetCurrentWorkingDirectoryAsString() (workingDirPath string) {
	workingDirPath, err := GetCurrentWorkingDirectoryAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return workingDirPath
}
