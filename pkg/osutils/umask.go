package osutils

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
)

// Get umask of the current process
func GetUmask(ctx context.Context) (umask int, err error) {
	umask, err = GetProcessProcStatusFileValueAsInt(ctx, "Umask")
	if err != nil {
		return -1, err
	}

	return umask, nil
}

func GetProcessProcStatusFileValueAsInt(ctx context.Context, valueName string) (value int, err error) {
	statusFile, err := GetProcessProcStatusFile(ctx)
	if err != nil {
		return -1, err
	}

	value, err = statusFile.GetValueAsInt(ctx, valueName)
	if err != nil {
		return -1, err
	}

	return value, err
}

func  GetProcessDefaultDirectoryModeAsFsFileMode(ctx context.Context) (defaultMode fs.FileMode, err error) {
	defaultModeInt, err := GetProcessDefaultDirectoryModeAsInt(ctx)
	if err != nil {
		return 0, err
	}

	defaultMode = fs.FileMode(defaultModeInt)
	return defaultMode, nil
}

func  GetProcessDefaultDirectoryModeAsInt(ctx context.Context) (defaultMode int, err error) {
	umask, err := GetUmask(ctx)
	if err != nil {
		return -1, err
	}

	defaultMode = 0777 - umask

	return defaultMode, nil
}

func  GetProcessId() (processId int) {
	processId = os.Getpid()
	return processId
}

func  GetProcessProcDirectory(ctx context.Context) (procDirectory filesinterfaces.Directory, err error) {
	pid := GetProcessId()
	procDirectory, err = files.GetLocalDirectoryByPath(ctx, fmt.Sprintf("/proc/%d", pid))
	if err != nil {
		return nil, err
	}

	return procDirectory, nil
}

func GetProcessProcStatusFile(ctx context.Context) (statusFile filesinterfaces.File, err error) {
	pidDir, err := GetProcessProcDirectory(ctx)
	if err != nil {
		return nil, err
	}

	statusFile, err = pidDir.GetFileInDirectory("status")
	if err != nil {
		return nil, err
	}

	return statusFile, nil
}
