package asciichgolangpublic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// A LocalFile represents a locally available file.
type LocalFile struct {
	FileBase
	path string
}

func GetLocalFileByPath(localPath string) (l *LocalFile, err error) {
	return NewLocalFileByPath(localPath)
}

func MustGetLocalFileByPath(localPath string) (l *LocalFile) {
	l, err := GetLocalFileByPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func MustNewLocalFileByPath(localPath string) (l *LocalFile) {
	l, err := NewLocalFileByPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func NewLocalFile() (l *LocalFile) {
	l = new(LocalFile)

	// Allow usage of the base class functions:
	l.MustSetParentFileForBaseClass(l)

	return l
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

// Delete a file if it exists.
// If the file is already absent this function does nothing.
func (l *LocalFile) Delete(verbose bool) (err error) {
	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	exists, err := l.Exists()
	if err != nil {
		return err
	}

	if exists {
		err = os.Remove(path)
		if err != nil {
			return TracedErrorf("Failed to delet localFile '%s': '%w'", path, err)
		}

		if verbose {
			LogChangedf("Local file '%s' deleted.", path)
		}
	} else {
		if verbose {
			LogChangedf("Local file '%s' is already absent. Skip delete.", path)
		}
	}

	return nil
}

func (l *LocalFile) Create(verbose bool) (err error) {
	exists, err := l.Exists()
	if err != nil {
		return err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if exists {
		if verbose {
			LogInfof("Local file '%s' already exists", path)
		}
	} else {
		err = l.WriteString("", false)
		if err != nil {
			return err
		}

		if verbose {
			LogInfof("Local file '%s' created", path)
		}
	}

	return nil
}

func (l *LocalFile) Exists() (exists bool, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return false, err
	}

	fileInfo, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, TracedErrorf("Unable to evaluate if local file exists: '%w'", err)
	}

	return !fileInfo.IsDir(), err
}

func (l *LocalFile) GetBaseName() (baseName string, err error) {
	path, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(path)

	if baseName == "" {
		return "", TracedErrorf(
			"Base name is empty string after evaluation of path='%s'",
			path,
		)
	}

	return baseName, nil
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

func (l *LocalFile) IsPathSet() (isSet bool) {
	return false
}

func (l *LocalFile) MustCreate(verbose bool) {
	err := l.Create(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustDelete(verbose bool) {
	err := l.Delete(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustExists() (exists bool) {
	exists, err := l.Exists()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalFile) MustGetBaseName() (baseName string) {
	baseName, err := l.GetBaseName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return baseName
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

func (l *LocalFile) MustReadAsBytes() (content []byte) {
	content, err := l.ReadAsBytes()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (l *LocalFile) MustSetPath(path string) {
	err := l.SetPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustWriteBytes(toWrite []byte, verbose bool) {
	err := l.WriteBytes(toWrite, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalFile) ReadAsBytes() (content []byte, err error) {
	path, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	content, err = os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, err
}

func (l *LocalFile) SetPath(path string) (err error) {
	if path == "" {
		return TracedError("path is empty string")
	}

	l.path = path

	return nil
}

func (l *LocalFile) WriteBytes(toWrite []byte, verbose bool) (err error) {
	if toWrite == nil {
		return TracedErrorNil("toWrite")
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(path, toWrite, 0644)
	if err != nil {
		return TracedErrorf("Unable to write file '%s': %w", path, err)
	}

	return nil
}
