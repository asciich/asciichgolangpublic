package asciichgolangpublic

func GetLocalCommandExecutorFileByFile(file File, verbose bool) (commandExecutorFile *CommandExecutorFile, err error) {
	if file == nil {
		return nil, TracedErrorEmptyString("file")
	}

	err = file.CheckIsLocalFile(verbose)
	if err != nil {
		return nil, err
	}

	pathToUse, err := file.GetLocalPath()
	if err != nil {
		return nil, err
	}

	commandExecutorFile, err = GetLocalCommandExecutorFileByPath(pathToUse)
	if err != nil {
		return nil, err
	}

	return commandExecutorFile, nil
}

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

func MustGetLocalCommandExecutorFileByFile(file File, verbose bool) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetLocalCommandExecutorFileByFile(file, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutorFile
}

func MustGetLocalCommandExecutorFileByPath(localPath string) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetLocalCommandExecutorFileByPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutorFile
}
