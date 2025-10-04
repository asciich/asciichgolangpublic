package files

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
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

func GetLocalCommandExecutorFileByFile(ctx context.Context, file filesinterfaces.File) (commandExecutorFile *CommandExecutorFile, err error) {
	if file == nil {
		return nil, tracederrors.TracedErrorEmptyString("file")
	}

	err = file.CheckIsLocalFile(ctx)
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
