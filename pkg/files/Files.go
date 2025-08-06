package files

import (
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/osutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type FilesService struct {
}

func Files() (f *FilesService) {
	return NewFilesService()
}

func DeleteFileByPath(path string, verbose bool) (err error) {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	toDelete, err := GetLocalFileByPath(path)
	if err != nil {
		return err
	}

	err = toDelete.Delete(contextutils.GetVerbosityContextByBool(verbose), &filesoptions.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func NewFilesService() (f *FilesService) {
	return new(FilesService)
}

func ReadFileAsString(path string) (content string, err error) {
	return Files().ReadAsString(path)
}

func (f *FilesService) ReadAsString(path string) (content string, err error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString(path)
	}

	localFile, err := GetLocalFileByPath(path)
	if err != nil {
		return "", err
	}

	content, err = localFile.ReadAsString()
	if err != nil {
		return "", err
	}

	return content, nil
}

func GetCurrentWorkingDirectory() (workingDirectory *LocalDirectory, err error) {
	workingDirectoryPath, err := osutils.GetCurrentWorkingDirectoryAsString()
	if err != nil {
		return nil, err
	}

	workingDirectory, err = GetLocalDirectoryByPath(workingDirectoryPath)
	if err != nil {
		return nil, err
	}

	return workingDirectory, nil
}
