package gitutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegitoo"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetGitRepositoryByDirectory(ctx context.Context, directory filesinterfaces.Directory) (repository gitinterfaces.GitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	hostDescription, err := directory.GetHostDescription()
	if err != nil {
		return nil, err
	}

	if hostDescription == "localhost" {
		return nativegitoo.NewGitRepositoryFromDirectory(directory)
	}

	return commandexecutorgitoo.NewGitRepositoryFromDirectory(directory)
}
