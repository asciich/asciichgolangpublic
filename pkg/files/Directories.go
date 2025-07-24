package files

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
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

func (d *DirectoriesService) CreateLocalDirectoryByPath(path string, verbose bool) (l Directory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	dir, err := GetLocalDirectoryByPath(path)
	if err != nil {
		return nil, err
	}

	err = dir.Create(verbose)
	if err != nil {
		return nil, err
	}

	return dir, nil
}

func (d *DirectoriesService) MustCreateLocalDirectoryByPath(path string, verbose bool) (l Directory) {
	l, err := d.CreateLocalDirectoryByPath(path, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return l
}
