package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetCommandExecutorGitRepositoryByPath(commandExecutor commandexecutor.CommandExecutor, path string) (gitRepo *CommandExecutorGitRepository, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecturo")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	gitRepo, err = NewCommandExecutorGitRepository(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = gitRepo.SetDirPath(path)
	if err != nil {
		return nil, err
	}

	return gitRepo, nil
}

func GetLocalCommandExecutorGitRepositoryByDirectory(directory Directory) (gitRepo *CommandExecutorGitRepository, err error) {
	if directory == nil {
		return nil, tracederrors.TracedErrorNil("directory")
	}

	isLocalDir, err := directory.IsLocalDirectory()
	if err != nil {
		return nil, err
	}

	path, err := directory.GetPath()
	if err != nil {
		return nil, err
	}

	if !isLocalDir {
		return nil, tracederrors.TracedErrorf(
			"Unable to get LocalCommandExecutorGitRepository for non local path '%s'",
			path,
		)
	}

	gitRepo, err = GetLocalCommandExecutorGitRepositoryByPath(path)
	if err != nil {
		return nil, err
	}

	return gitRepo, nil
}

func GetLocalCommandExecutorGitRepositoryByPath(path string) (gitRepo *CommandExecutorGitRepository, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return GetCommandExecutorGitRepositoryByPath(
		commandexecutor.Bash(),
		path,
	)
}

func MustGetCommandExecutorGitRepositoryByPath(commandExecutor commandexecutor.CommandExecutor, path string) (gitRepo *CommandExecutorGitRepository) {
	gitRepo, err := GetCommandExecutorGitRepositoryByPath(commandExecutor, path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepo
}

func MustGetLocalCommandExecutorGitRepositoryByDirectory(directory Directory) (gitRepo *CommandExecutorGitRepository) {
	gitRepo, err := GetLocalCommandExecutorGitRepositoryByDirectory(directory)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepo
}

func MustGetLocalCommandExecutorGitRepositoryByPath(path string) (gitRepo *CommandExecutorGitRepository) {
	gitRepo, err := GetLocalCommandExecutorGitRepositoryByPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepo
}
