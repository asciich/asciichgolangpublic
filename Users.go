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
