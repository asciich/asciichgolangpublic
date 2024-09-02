package asciichgolangpublic

import (
	"os"
	"runtime"
)

type OsService struct{}

func NewOsService() (o *OsService) {
	return new(OsService)
}

func OS() (o *OsService) {
	return NewOsService()
}

func (o *OsService) GetCurrentWorkingDirectory() (workingDirectory *LocalDirectory, err error) {
	workingDirectoryPath, err := o.GetCurrentWorkingDirectoryAsString()
	if err != nil {
		return nil, err
	}

	workingDirectory, err = GetLocalDirectoryByPath(workingDirectoryPath)
	if err != nil {
		return nil, err
	}

	return workingDirectory, nil
}

func (o *OsService) GetCurrentWorkingDirectoryAsString() (workingDirPath string, err error) {
	workingDirPath, err = os.Getwd()
	if err != nil {
		return "", TracedErrorf("Get working directory failed: %w", err)
	}

	if !Paths().IsAbsolutePath(workingDirPath) {
		return "", TracedErrorf(
			"Evaluated working directory path '%s' is not an absolute path after evaluation.",
			workingDirPath,
		)
	}

	return workingDirPath, nil
}

func (o *OsService) IsRunningOnLinux() (isRunningOnLinux bool) {
	return runtime.GOOS == "linux"
}

func (o *OsService) IsRunningOnWindows() (isRunningOnWindows bool) {
	return runtime.GOOS == "windows"
}

func (o *OsService) MustGetCurrentWorkingDirectory() (workingDirectory *LocalDirectory) {
	workingDirectory, err := o.GetCurrentWorkingDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return workingDirectory
}

func (o *OsService) MustGetCurrentWorkingDirectoryAsString() (workingDirPath string) {
	workingDirPath, err := o.GetCurrentWorkingDirectoryAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return workingDirPath
}
