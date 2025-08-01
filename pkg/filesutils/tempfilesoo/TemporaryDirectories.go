package tempfilesoo

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/files"
)

func CreateEmptyTemporaryDirectory(ctx context.Context) (temporaryDirectory *files.LocalDirectory, err error) {
	dirPath, err := os.MkdirTemp("", "empty")
	if err != nil {
		return nil, err
	}

	temporaryDirectory, err = files.GetLocalDirectoryByPath(dirPath)
	if err != nil {
		return nil, err
	}

	return temporaryDirectory, nil
}

func CreateEmptyTemporaryDirectoryAndGetPath(ctx context.Context) (TemporaryDirectoryPath string, err error) {
	TemporaryDirectory, err := CreateEmptyTemporaryDirectory(ctx)
	if err != nil {
		return "", err
	}

	TemporaryDirectoryPath, err = TemporaryDirectory.GetLocalPath()
	if err != nil {
		return "", err
	}

	return TemporaryDirectoryPath, nil
}
