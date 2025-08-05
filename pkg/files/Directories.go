package files

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DirectoriesService struct {
}

func Directories() (d *DirectoriesService) {
	return NewDirectoriesService()
}

func NewDirectoriesService() (d *DirectoriesService) {
	return new(DirectoriesService)
}

func (d *DirectoriesService) CreateLocalDirectoryByPath(ctx context.Context, path string, options *filesoptions.CreateOptions) (l filesinterfaces.Directory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	dir, err := GetLocalDirectoryByPath(path)
	if err != nil {
		return nil, err
	}

	err = dir.Create(ctx, options)
	if err != nil {
		return nil, err
	}

	return dir, nil
}
