package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegitoo"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetGitRepositoryByDirectory(directory filesinterfaces.Directory) (repository gitinterfaces.GitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	localDirectory, ok := directory.(*files.LocalDirectory)
	if ok {
		return GetLocalGitRepositoryFromDirectory(localDirectory)
	}

	nativeFilesDirectory, ok := directory.(*nativefilesoo.Directory)
	if ok {
		return nativegitoo.NewGitRepositoryFromDirectory(nativeFilesDirectory)
	}

	return commandexecutorgitoo.NewGitRepositoryFromDirectory(directory)
}
