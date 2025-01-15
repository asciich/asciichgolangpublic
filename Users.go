package asciichgolangpublic

import (
	"os"
	"os/user"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type UsersService struct {
}

func GetDirectoryInHomeDirectory(path ...string) (directoryInHome Directory, err error) {
	directoryInHome, err = Users().GetDirectoryInHomeDirectory(path...)
	if err != nil {
		return nil, err
	}

	return directoryInHome, nil
}

func GetFileInHomeDirectory(path ...string) (fileInHome File, err error) {
	fileInHome, err = Users().GetFileInHomeDirectory(path...)
	if err != nil {
		return nil, err
	}

	return fileInHome, nil
}

func IsRunningAsRoot(verbose bool) (isRunningAsRoot bool, err error) {
	isRunningAsRoot, err = Users().IsRunningAsRoot(verbose)
	if err != nil {
		return false, err
	}

	return isRunningAsRoot, nil
}

func MustGetDirectoryInHomeDirectory(path ...string) (directoryInHome Directory) {
	directoryInHome, err := GetDirectoryInHomeDirectory(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return directoryInHome
}

func MustGetFileInHomeDirectory(path ...string) (fileInHome File) {
	fileInHome, err := Users().GetFileInHomeDirectory(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileInHome
}

func MustIsRunningAsRoot(verbose bool) (isRunningAsRoot bool) {
	isRunningAsRoot, err := IsRunningAsRoot(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunningAsRoot
}

func MustWhoAmI(verbose bool) (userName string) {
	userName, err := WhoAmI(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userName
}

func NewUsersService() (u *UsersService) {
	return new(UsersService)
}

func Users() (u *UsersService) {
	return NewUsersService()
}

func WhoAmI(verbose bool) (userName string, err error) {
	userName, err = Users().WhoAmI(verbose)
	if err != nil {
		return "", err
	}

	return userName, err
}

func (u *UsersService) GetCurrentUserName(verbose bool) (currentUserName string, err error) {
	nativeUser, err := u.GetNativeUser()
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

func (u *UsersService) GetDirectoryInHomeDirectory(path ...string) (fileInUnsersHome Directory, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no length")
	}

	usersHome, err := u.GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	fileInUnsersHome, err = usersHome.GetSubDirectory(path...)
	if err != nil {
		return nil, err
	}

	return fileInUnsersHome, nil
}

func (u *UsersService) GetFileInHomeDirectory(path ...string) (fileInUnsersHome File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no length")
	}

	usersHome, err := u.GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	fileInUnsersHome, err = usersHome.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	return fileInUnsersHome, nil
}

func (u *UsersService) GetFileInHomeDirectoryAsLocalFile(path ...string) (localFile *LocalFile, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path is empty")
	}

	fileToReturn, err := u.GetFileInHomeDirectory(path...)
	if err != nil {
		return nil, err
	}

	localFile, ok := fileToReturn.(*LocalFile)
	if !ok {
		return nil, tracederrors.TracedError("Unable to convert to local file")
	}

	return localFile, nil
}

func (u *UsersService) GetHomeDirectory() (homeDir Directory, err error) {
	homeDirPath, err := u.GetHomeDirectoryAsString()
	if err != nil {
		return nil, err
	}

	homeDir, err = GetLocalDirectoryByPath(homeDirPath)
	if err != nil {
		return nil, err
	}

	return homeDir, nil
}

func (u *UsersService) GetHomeDirectoryAsString() (homeDirPath string, err error) {
	homeDirPath, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDirPath, nil
}

func (u *UsersService) GetNativeUser() (nativeUser *user.User, err error) {
	nativeUser, err = user.Current()
	if err != nil {
		return nil, err
	}

	return nativeUser, nil
}

func (u *UsersService) IsRunningAsRoot(verbose bool) (isRunningAsRoot bool, err error) {
	userName, err := u.GetCurrentUserName(verbose)
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

func (u *UsersService) MustGetCurrentUserName(verbose bool) (currentUserName string) {
	currentUserName, err := u.GetCurrentUserName(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return currentUserName
}

func (u *UsersService) MustGetDirectoryInHomeDirectory(path ...string) (fileInUnsersHome Directory) {
	fileInUnsersHome, err := u.GetDirectoryInHomeDirectory(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileInUnsersHome
}

func (u *UsersService) MustGetFileInHomeDirectory(path ...string) (fileInUnsersHome File) {
	fileInUnsersHome, err := u.GetFileInHomeDirectory(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileInUnsersHome
}

func (u *UsersService) MustGetFileInHomeDirectoryAsLocalFile(path ...string) (localFile *LocalFile) {
	localFile, err := u.GetFileInHomeDirectoryAsLocalFile(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localFile
}

func (u *UsersService) MustGetHomeDirectory() (homeDir Directory) {
	homeDir, err := u.GetHomeDirectory()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return homeDir
}

func (u *UsersService) MustGetHomeDirectoryAsString() (homeDirPath string) {
	homeDirPath, err := u.GetHomeDirectoryAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return homeDirPath
}

func (u *UsersService) MustGetNativeUser() (nativeUser *user.User) {
	nativeUser, err := u.GetNativeUser()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nativeUser
}

func (u *UsersService) MustIsRunningAsRoot(verbose bool) (isRunningAsRoot bool) {
	isRunningAsRoot, err := u.IsRunningAsRoot(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunningAsRoot
}

func (u *UsersService) MustWhoAmI(verbose bool) (userName string) {
	userName, err := u.WhoAmI(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userName
}

func (u *UsersService) WhoAmI(verbose bool) (userName string, err error) {
	userName, err = u.GetCurrentUserName(verbose)
	if err != nil {
		return "", err
	}

	return userName, nil
}
