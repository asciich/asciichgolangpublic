package asciichgolangpublic

import "os"

type UsersService struct {
}

func NewUsersService() (u *UsersService) {
	return new(UsersService)
}

func Users() (u *UsersService) {
	return NewUsersService()
}

func (u *UsersService) GetDirectoryInHomeDirectory(path ...string) (fileInUnsersHome Directory, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no length")
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
		return nil, TracedError("path has no length")
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
		return nil, TracedError("path is empty")
	}

	fileToReturn, err := u.GetFileInHomeDirectory(path...)
	if err != nil {
		return nil, err
	}

	localFile, ok := fileToReturn.(*LocalFile)
	if !ok {
		return nil, TracedError("Unable to convert to local file")
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

func (u *UsersService) MustGetDirectoryInHomeDirectory(path ...string) (fileInUnsersHome Directory) {
	fileInUnsersHome, err := u.GetDirectoryInHomeDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fileInUnsersHome
}

func (u *UsersService) MustGetFileInHomeDirectory(path ...string) (fileInUnsersHome File) {
	fileInUnsersHome, err := u.GetFileInHomeDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fileInUnsersHome
}

func (u *UsersService) MustGetFileInHomeDirectoryAsLocalFile(path ...string) (localFile *LocalFile) {
	localFile, err := u.GetFileInHomeDirectoryAsLocalFile(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localFile
}

func (u *UsersService) MustGetHomeDirectory() (homeDir Directory) {
	homeDir, err := u.GetHomeDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return homeDir
}

func (u *UsersService) MustGetHomeDirectoryAsString() (homeDirPath string) {
	homeDirPath, err := u.GetHomeDirectoryAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return homeDirPath
}
