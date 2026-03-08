package commandexecutorfileoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func (f *File) GetSha256Sum(ctx context.Context) (string, error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	return commandexecutorfile.GetSha256Sum(ctx, commandExecutor, path)
}
