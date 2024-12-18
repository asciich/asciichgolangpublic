package asciichgolangpublic

func GetLocalCommandExecutorGitRepositoryByDirectory(directory Directory) (gitRepo *CommandExecutorGitRepository, err error) {
	if directory == nil {
		return nil, TracedErrorNil("directory")
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
		return nil, TracedErrorf(
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
		return nil, TracedErrorEmptyString("path")
	}

	gitRepo, err = NewCommandExecutorGitRepository(Bash())
	if err != nil {
		return nil, err
	}

	err = gitRepo.SetDirPath(path)
	if err != nil {
		return nil, err
	}

	return gitRepo, nil
}

func MustGetLocalCommandExecutorGitRepositoryByDirectory(directory Directory) (gitRepo *CommandExecutorGitRepository) {
	gitRepo, err := GetLocalCommandExecutorGitRepositoryByDirectory(directory)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitRepo
}

func MustGetLocalCommandExecutorGitRepositoryByPath(path string) (gitRepo *CommandExecutorGitRepository) {
	gitRepo, err := GetLocalCommandExecutorGitRepositoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitRepo
}
