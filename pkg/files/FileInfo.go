package files

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type FileInfo struct {
	Path      string
	SizeBytes int64
}

func NewFileInfo() (f *FileInfo) {
	return new(FileInfo)
}

func (f *FileInfo) GetPath() (path string, err error) {
	if f.Path == "" {
		return "", tracederrors.TracedErrorf("Path not set")
	}

	return f.Path, nil
}

func (f *FileInfo) GetPathAndSizeBytes() (path string, sizeBytes int64, err error) {
	path, err = f.GetPath()
	if err != nil {
		return "", -1, err
	}

	sizeBytes, err = f.GetSizeBytes()
	if err != nil {
		return "", -1, err
	}

	return path, sizeBytes, nil
}

func (f *FileInfo) GetSizeBytes() (sizeBytes int64, err error) {
	return f.SizeBytes, nil
}

func (f *FileInfo) SetPath(path string) (err error) {
	if path == "" {
		return tracederrors.TracedErrorf("path is empty string")
	}

	f.Path = path

	return nil
}

func (f *FileInfo) SetSizeBytes(sizeBytes int64) (err error) {
	f.SizeBytes = sizeBytes

	return nil
}
