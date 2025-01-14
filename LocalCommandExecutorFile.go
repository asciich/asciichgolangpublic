package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

func GetCommandExecutorFileByPath(commandExector CommandExecutor, path string) (commandExecutorFile *CommandExecutorFile, err error) {
	if commandExector == nil {
		return nil, errors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, errors.TracedErrorEmptyString("path")
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
		return nil, errors.TracedErrorEmptyString("file")
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
		return nil, errors.TracedErrorEmptyString(localPath)
	}

	return GetCommandExecutorFileByPath(Bash(), localPath)
}

func MustGetCommandExecutorFileByPath(commandExector CommandExecutor, path string) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetCommandExecutorFileByPath(commandExector, path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutorFile
}

func MustGetLocalCommandExecutorFileByFile(file File, verbose bool) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetLocalCommandExecutorFileByFile(file, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutorFile
}

func MustGetLocalCommandExecutorFileByPath(localPath string) (commandExecutorFile *CommandExecutorFile) {
	commandExecutorFile, err := GetLocalCommandExecutorFileByPath(localPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutorFile
}
