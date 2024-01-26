package github.com/asciich/asciichgolangpublic

import (
	"fmt"
)

// A LocalFile represents a locally available file.
type LocalFile struct {
	path string
}

func MustNewLocalFileByPath(localPath string) (l *LocalFile) {
	l, err := NewLocalFileByPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func NewLocalFile() (l *LocalFile) {
	return new(LocalFile)
}

func NewLocalFileByPath(localPath string) (l *LocalFile, err error) {
	if localPath == "" {
		return nil, TracedErrorEmptyString("localPath")
	}

	l = NewLocalFile()

	err = l.SetPath(localPath)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *LocalFile) GetLocalPath() (path string, err error) {
	return l.GetPath()
}

func (l *LocalFile) GetPath() (path string, err error) {
	if l.path == "" {
		return "", fmt.Errorf("path not set")
	}

	return l.path, nil
}

func (l *LocalFile) GetUriAsString() (uri string, err error) {
	path, err := l.GetPath()
	if err != nil {
		return "", err
	}

	if Paths().IsRelativePath(path) {
		return "", TracedErrorf("Only implemeted for absolute paths but got '%s'", path)
	}

	uri = "file://" + path

	return uri, nil
}

func (l *LocalFile) MustExists() (exists bool) {
	exists, err := l.Exists()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalFile) MustGetLocalPath() (path string) {
	path, err := l.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path
}

func (l *LocalFile) MustGetPath() (path string) {
	path, err := l.GetPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path
}

func (l *LocalFile) MustGetUriAsString() (uri string) {
	uri, err := l.GetUriAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return uri
}

func (l *LocalFile) MustSetPath(path string) {
	err := l.SetPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalFile) SetPath(path string) (err error) {
	if path == "" {
		return TracedError("path is empty string")
	}

	l.path = path

	return nil
}

func (l LocalFile) Exists() (exists bool, err error) {
	return false, TracedError("not implemented")
}

func (l LocalFile) IsPathSet() (isSet bool) {
	return false
}
