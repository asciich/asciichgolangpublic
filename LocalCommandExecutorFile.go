package asciichgolangpublic

func GetLocalCommandExecutorFileByPath(localPath string) (commandExecutorFile *CommandExecutorFile, err error) {
	if localPath == "" {
		return nil, TracedErrorEmptyString(localPath)
	}

	commandExecutorFile = NewCommandExecutorFile()

	err = commandExecutorFile.SetCommandExecutor(Bash())
	if err != nil {
		return nil, err
	}

	err = commandExecutorFile.SetFilePath(localPath)
	if err != nil {
		return nil, err
	}

	return commandExecutorFile, nil
}

func MustGetLocalCommandExecutorFileByPath(localPath string) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetLocalCommandExecutorFileByPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutorFile
}
