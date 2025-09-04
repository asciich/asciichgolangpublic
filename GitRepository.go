package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/datatypes"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetGitRepositoryByDirectory(directory filesinterfaces.Directory) (repository gitinterfaces.GitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	localDirectory, ok := directory.(*files.LocalDirectory)
	if ok {
		return GetLocalGitReposioryFromDirectory(localDirectory)
	}

	commandExecutorDirectory, ok := directory.(*files.CommandExecutorDirectory)
	if ok {
		return GetCommandExecutorGitRepositoryFromDirectory(commandExecutorDirectory)
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

func GitRepositoryDefaultCommitMessageForInitializeWithEmptyCommit() (msg string) {
	return "Initial empty commit during repo initialization"
}

func GitRepositryDefaultAuthorEmail() (email string) {
	return "asciichgolangpublic@example.net"
}

func GitRepositryDefaultAuthorName() (name string) {
	return "asciichgolangpublic git repo initializer"
}

func MustGetGitRepositoryByDirectory(directory filesinterfaces.Directory) (repository gitinterfaces.GitRepository) {
	repository, err := GetGitRepositoryByDirectory(directory)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repository
}
