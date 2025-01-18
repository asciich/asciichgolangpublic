package files

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
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

func (f *FileInfo) MustGetPath() (path string) {
	path, err := f.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (f *FileInfo) MustGetPathAndSizeBytes() (path string, sizeBytes int64) {
	path, sizeBytes, err := f.GetPathAndSizeBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path, sizeBytes
}

func (f *FileInfo) MustGetSizeBytes() (sizeBytes int64) {
	sizeBytes, err := f.GetSizeBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sizeBytes
}

func (f *FileInfo) MustSetPath(path string) {
	err := f.SetPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileInfo) MustSetSizeBytes(sizeBytes int64) {
	err := f.SetSizeBytes(sizeBytes)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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
