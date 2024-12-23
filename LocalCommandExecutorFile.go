package asciichgolangpublic

func GetCommandExecutorFileByPath(commandExector CommandExecutor, path string) (commandExecutorFile *CommandExecutorFile, err error) {
	if commandExector == nil {
		return nil, TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, TracedErrorEmptyString("path")
	}

	commandExecutorFile = NewCommandExecutorFile()

	err = commandExecutorFile.SetCommandExecutor(commandExector)
	if err != nil {
		return nil, err
	}

	err = commandExecutorFile.SetFilePath(path)
	if err != nil {
		return nil, err
	}

	return commandExecutorFile, nil
}

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

	return GetCommandExecutorFileByPath(Bash(), localPath)
}

func MustGetCommandExecutorFileByPath(commandExector CommandExecutor, path string) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetCommandExecutorFileByPath(commandExector, path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutorFile
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
