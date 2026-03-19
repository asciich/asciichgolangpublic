package commandexecutorfileoo

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func (f *File) ReadFirstNBytes(ctx context.Context, numberOfBytesToRead int) (firstBytes []byte, err error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	return commandexecutorfile.ReadFirstNBytes(ctx, commandExecutor, filePath, numberOfBytesToRead)
}

func (f *File) ReadAsBytes(ctx context.Context) (content []byte, err error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	return commandexecutorfile.ReadAsBytes(commandExecutor, filePath)
}

func (f *File) OpenAsReadCloser(ctx context.Context) (io.ReadCloser, error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	return commandexecutorfile.OpenAsReadCloser(ctx, commandExecutor, filePath)
}
