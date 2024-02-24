package asciichgolangpublic

import (
	"errors"
	"os"
	"path/filepath"
)

type LocalDirectory struct {
	DirectoryBase
	localPath string
}

func GetLocalDirectoryByPath(path string) (directory Directory, err error) {
	if path == "" {
		return nil, TracedErrorEmptyString("path")
	}

	localDirectory := NewLocalDirectory()

	err = localDirectory.SetLocalPath(path)
	if err != nil {
		return nil, err
	}

	return localDirectory, nil
}

func MustGetLocalDirectoryByPath(path string) (directory Directory) {
	directory, err := GetLocalDirectoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return directory
}

func NewLocalDirectory() (l *LocalDirectory) {
	l = new(LocalDirectory)

	// Allow usage of the base class functions:
	l.MustSetParentDirectoryForBaseClass(l)

	return l
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

func (l *LocalDirectory) GetFileInDirectory(path ...string) (file File, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no elements")
	}

	dirPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	filePath := dirPath
	for _, p := range path {
		filePath = filepath.Join(filePath, p)
	}

	file, err = GetLocalFileByPath(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *LocalDirectory) GetLocalPath() (localPath string, err error) {
	if l.localPath == "" {
		return "", TracedErrorf("localPath not set")
	}

	return l.localPath, nil
}

func (l *LocalDirectory) GetSubDirectory(path ...string) (subDirectory Directory, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no elements")
	}

	dirPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	directoryPath := dirPath
	for _, p := range path {
		directoryPath = filepath.Join(directoryPath, p)
	}

	subDirectory, err = GetLocalDirectoryByPath(directoryPath)
	if err != nil {
		return nil, err
	}

	return subDirectory, nil
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

func (l *LocalDirectory) MustGetFileInDirectory(path ...string) (file File) {
	file, err := l.GetFileInDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return file
}

func (l *LocalDirectory) MustGetLocalPath() (localPath string) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localPath
}

func (l *LocalDirectory) MustGetSubDirectory(path ...string) (subDirectory Directory) {
	subDirectory, err := l.GetSubDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectory
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
