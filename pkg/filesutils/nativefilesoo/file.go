package nativefilesoo

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type File struct {
	filesgeneric.FileBase
	path string
}

func NewFileByPath(path string) (filesinterfaces.File, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	ret := &File{
		path: path,
	}

	err := ret.SetParentFileForBaseClass(ret)
	if err != nil {
		panic(err)
	}

	return ret, nil
}

func (f *File) GetDeepCopy() filesinterfaces.File {
	copy := new(File)
	*copy = *f

	err := copy.SetParentFileForBaseClass(copy)
	if err != nil {
		panic(err)
	}

	return copy
}

func (f *File) GetPath() (path string, err error) {
	if f.path == "" {
		return "", tracederrors.TracedError("Path not set")
	}

	return f.path, nil
}

func (f *File) ReadFirstNBytes(ctx context.Context, numberOfBytesToRead int) (firstBytes []byte, err error) {
	if numberOfBytesToRead <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid numberOfBytesToRead: '%d'", numberOfBytesToRead)
	}

	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, tracederrors.TracedError(err)
	}

	defer fd.Close()

	firstBytes = make([]byte, numberOfBytesToRead)
	readBytes, err := fd.Read(firstBytes)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, tracederrors.TracedError(err)
		}
	}

	firstBytes = firstBytes[:readBytes]

	return firstBytes, nil
}
