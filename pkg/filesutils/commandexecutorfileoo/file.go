package commandexecutorfileoo

import (
	"path/filepath"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type File struct {
	files.FileBase
	commandExecutor commandexecutorinterfaces.CommandExecutor
	path            string
}

func New(commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (filesinterfaces.File, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExececutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	ret := &File{
		commandExecutor: commandExecutor,
		path:            path,
	}

	err := ret.SetParentFileForBaseClass(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (f *File) GetHostDescription() (string, error) {
	commandExecutor, err := f.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandExecutor.GetHostDescription()
}

func (f *File) GetBaseName() (string, error) {
	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	baseName := filepath.Base(path)

	if slices.Contains([]string{"", " ", ".", "/", "\\"}, baseName) {
		return "", tracederrors.TracedErrorf("Evaluated invalid baseName '%s' out of path '%s'.", baseName, path)
	}

	return baseName, nil
}

func (f *File) GetLocalPath() (localPath string, err error) {
	isLocalDirectory, err := f.IsLocalFile(contextutils.ContextSilent())
	if err != nil {
		return "", err
	}

	if isLocalDirectory {
		return f.GetPath()
	} else {
		hostDescription, err := f.GetHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf("File is on '%s', not on localhost", hostDescription)
	}
}
