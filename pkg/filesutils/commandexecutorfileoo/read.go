package commandexecutorfileoo

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func (f *File) ReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte, err error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	return commandexecutorfile.ReadFirstNBytes(commandExecutor, filePath, numberOfBytesToRead)
}

func (f *File) ReadAsBytes() (content []byte, err error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	return commandexecutorfile.ReadAsBytes(commandExecutor, filePath)
}
