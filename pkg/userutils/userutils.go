package userutils

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetCurrentUserName(verbose bool) (currentUserName string, err error) {
	nativeUser, err := GetNativeUser()
	if err != nil {
		return "", err
	}

	currentUserName = nativeUser.Username

	if currentUserName == "" {
		return "", tracederrors.TracedError("currentUserName is empty string after evaluation")
	}

	if verbose {
		logging.LogInfof("Current username is '%s'.", currentUserName)
	}

	return currentUserName, nil
}

func GetDirectoryInHomeDirectory(path ...string) (fileInUnsersHome files.Directory, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no length")
	}

	usersHome, err := GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	fileInUnsersHome, err = usersHome.GetSubDirectory(path...)
	if err != nil {
		return nil, err
	}

	return fileInUnsersHome, nil
}

func GetFileInHomeDirectory(path ...string) (fileInUnsersHome files.File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no length")
	}

	usersHome, err := GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	fileInUnsersHome, err = usersHome.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	return fileInUnsersHome, nil
}

func GetFileInHomeDirectoryAsLocalFile(path ...string) (localFile *files.LocalFile, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path is empty")
	}

	fileToReturn, err := GetFileInHomeDirectory(path...)
	if err != nil {
		return nil, err
	}

	localFile, ok := fileToReturn.(*files.LocalFile)
	if !ok {
		return nil, tracederrors.TracedError("Unable to convert to local file")
	}

	return localFile, nil
}

func GetFilePathInHomeDirectory(path ...string) (string, error) {
	if len(path) <= 0 {
		return "", tracederrors.TracedError("path has no elements.")
	}

	homePath, err := GetHomeDirectoryPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(append([]string{homePath}, path...)...), nil
}

func GetHomeDirectory() (homeDir files.Directory, err error) {
	homeDirPath, err := GetHomeDirectoryPath()
	if err != nil {
		return nil, err
	}

	homeDir, err = files.GetLocalDirectoryByPath(homeDirPath)
	if err != nil {
		return nil, err
	}

	return homeDir, nil
}

func GetHomeDirectoryPath() (homeDirPath string, err error) {
	homeDirPath, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDirPath, nil
}

func GetNativeUser() (nativeUser *user.User, err error) {
	nativeUser, err = user.Current()
	if err != nil {
		return nil, err
	}

	return nativeUser, nil
}

func IsRunningAsRoot(verbose bool) (isRunningAsRoot bool, err error) {
	userName, err := GetCurrentUserName(verbose)
	if err != nil {
		return false, err
	}

	isRunningAsRoot = userName == "root"

	if verbose {
		if isRunningAsRoot {
			logging.LogInfof("Running as root since current user name is '%s'.", userName)
		} else {
			logging.LogInfof("Not running as root, current user name is '%s'.", userName)
		}
	}

	return isRunningAsRoot, nil
}

func WhoAmI(verbose bool) (userName string, err error) {
	userName, err = GetCurrentUserName(verbose)
	if err != nil {
		return "", err
	}

	return userName, nil
}
