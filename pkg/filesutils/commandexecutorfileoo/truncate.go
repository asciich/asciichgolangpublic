package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func (f *File) Truncate(ctx context.Context, newSizeBytes int64) (err error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return err
	}

	return commandexecutorfile.Truncate(ctx, commandExecutor, filePath, newSizeBytes)
}

func (f *File) GetSizeBytes(ctx context.Context) (fileSize int64, err error) {
	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return 0, err
	}

	return commandexecutorfile.GetSizeBytes(ctx, commandExecutor, filePath)
}
