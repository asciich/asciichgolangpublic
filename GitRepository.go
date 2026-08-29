package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
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

	commandExecutorDirectory, ok := directory.(*commandexecutorfileoo.Directory)
	if ok {
		return commandexecutorgitoo.NewGitRepositoryFromDirectory(commandExecutorDirectory)
	}

	nativeFilesDirectory, ok := directory.(*nativefilesoo.Directory)
	if ok {
		return nativegitoo.NewGitRepositoryFromDirectory(nativeFilesDirectory)
	}

	unknownTypeName, err := datatypes.GetTypeName(directory)
	if err != nil {
		return nil, err
	}

	return nil, tracederrors.TracedErrorf(
		"Unknown directory implementation '%s'. Unable to get GitRepository",
		unknownTypeName,
	)
}
