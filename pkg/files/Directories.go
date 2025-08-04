package files

import (
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
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

func (d *DirectoriesService) CreateLocalDirectoryByPath(path string, verbose bool) (l filesinterfaces.Directory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	dir, err := GetLocalDirectoryByPath(path)
	if err != nil {
		return nil, err
	}

	err = dir.Create(contextutils.GetVerbosityContextByBool(verbose))
	if err != nil {
		return nil, err
	}

	return dir, nil
}
