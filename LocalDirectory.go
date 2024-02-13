package asciichgolangpublic

import (
	"errors"
	"os"
)

type LocalDirectory struct {
	localPath string
}

func GetLocalDirectoryByPath(path string) (directory *LocalDirectory, err error) {
	if path == "" {
		return nil, TracedErrorEmptyString("path")
	}

	directory = NewLocalDirectory()

	err = directory.SetLocalPath(path)
	if err != nil {
		return nil, err
	}

	return directory, nil
}

func MustGetLocalDirectoryByPath(path string) (directory *LocalDirectory) {
	directory, err := GetLocalDirectoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return directory
}

func NewLocalDirectory() (l *LocalDirectory) {
	return new(LocalDirectory)
}

func (l *LocalDirectory) Create(verbose bool) (err error) {
	exists, err := l.Exists()
	if err != nil {
		return nil
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil
	}

	if exists {
		if verbose {
			LogInfof("Local directory '%s' already exists. Skip create.", path)
		}
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return TracedErrorf("Create directory '%s' failed: '%w'", path, err)
		}

		if verbose {
			LogChangedf("Created local directory '%s'", path)
		}
	}

	return nil
}

func (l *LocalDirectory) Delete(verbose bool) (err error) {
	exists, err := l.Exists()
	if err != nil {
		return nil
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil
	}

	if exists {
		err := os.RemoveAll(path)
		if err != nil {
			return TracedErrorf("Delete directory '%s' failed: '%w'", path, err)
		}

		if verbose {
			LogChangedf("Deleted local directory '%s'", path)
		}
	} else {
		if verbose {
			LogInfof("Local directory '%s' already absent. Skip delete.", path)
		}
	}

	return nil
}

func (l *LocalDirectory) Exists() (exists bool, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return false, nil
	}

	dirInfo, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, TracedErrorf("Unable to evaluate if local directory exists: '%w'", err)
	}

	return dirInfo.IsDir(), nil
}

func (l *LocalDirectory) GetLocalPath() (localPath string, err error) {
	if l.localPath == "" {
		return "", TracedErrorf("localPath not set")
	}

	return l.localPath, nil
}

func (l *LocalDirectory) MustCreate(verbose bool) {
	err := l.Create(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustDelete(verbose bool) {
	err := l.Delete(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustExists() (exists bool) {
	exists, err := l.Exists()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalDirectory) MustGetLocalPath() (localPath string) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localPath
}

func (l *LocalDirectory) MustSetLocalPath(localPath string) {
	err := l.SetLocalPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) SetLocalPath(localPath string) (err error) {
	if localPath == "" {
		return TracedErrorf("localPath is empty string")
	}

	l.localPath = localPath

	return nil
}
