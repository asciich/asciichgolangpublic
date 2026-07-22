package gitutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewGitRepositoryFromDirectory(ctx context.Context, directory filesinterfaces.Directory) (gitinterfaces.GitRepository, error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	commandExecutorDirectory, ok := directory.(*commandexecutorfileoo.Directory)
	if ok {
		return commandexecutorgitoo.NewGitRepositoryFromDirectory(commandExecutorDirectory)
	}

	typeName, err := datatypes.GetTypeName(directory)
	if err != nil {
		return nil, err
	}
	return nil, tracederrors.TracedErrorf("Not implemented for '%s'", typeName)
}
