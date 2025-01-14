package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

func GetCommandExecutorGitRepositoryByPath(commandExecutor CommandExecutor, path string) (gitRepo *CommandExecutorGitRepository, err error) {
	if commandExecutor == nil {
		return nil, errors.TracedErrorNil("commandExecturo")
	}

	if path == "" {
		return nil, errors.TracedErrorEmptyString("path")
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
		return nil, errors.TracedErrorNil("directory")
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
		return nil, errors.TracedErrorf(
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
		return nil, errors.TracedErrorEmptyString("path")
	}

	return GetCommandExecutorGitRepositoryByPath(
		Bash(),
		path,
	)
}

func MustGetCommandExecutorGitRepositoryByPath(commandExecutor CommandExecutor, path string) (gitRepo *CommandExecutorGitRepository) {
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
