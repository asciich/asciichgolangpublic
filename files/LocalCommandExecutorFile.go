package files

import (
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func GetCommandExecutorFileByPath(commandExector commandexecutorinterfaces.CommandExecutor, path string) (commandExecutorFile *CommandExecutorFile, err error) {
	if commandExector == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
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
		return nil, tracederrors.TracedErrorEmptyString("file")
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
		return nil, tracederrors.TracedErrorEmptyString(localPath)
	}

	return GetCommandExecutorFileByPath(commandexecutorbashoo.Bash(), localPath)
}

func MustGetCommandExecutorFileByPath(commandExector commandexecutorinterfaces.CommandExecutor, path string) (commandExecutorFile *CommandExecutorFile) {
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
