package commandexecutorfileoo

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Directory struct {
	files.DirectoryBase
	commandExecutor commandexecutorinterfaces.CommandExecutor
	path            string
}

func NewDirectory(commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (filesinterfaces.Directory, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExececutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	ret := &Directory{
		commandExecutor: commandExecutor,
		path:            path,
	}

	err := ret.SetParentDirectoryForBaseClass(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (d *Directory) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if d.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExectuor not set")
	}

	return d.commandExecutor, nil
}

func (d *Directory) GetCommandExecutorAndDirectoryPath() (commandexecutorinterfaces.CommandExecutor, string, error) {
	ce, err := d.GetCommandExecutor()
	if err != nil {
		return nil, "", err
	}

	path, err := d.GetPath()
	if err != nil {
		return nil, "", err
	}

	return ce, path, nil
}

func (d *Directory) GetPath() (dirPath string, err error) {
	if d.path == "" {
		return "", tracederrors.TracedError("path not set")
	}

	return d.path, nil
}

func (d *Directory) SetPath(path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	d.path = path

	return nil
}

func (d *Directory) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	d.commandExecutor = commandExecutor

	return nil
}

func (d *Directory) GetFileInDirectory(pathToFile ...string) (file filesinterfaces.File, err error) {
	if len(pathToFile) <= 0 {
		return nil, tracederrors.TracedErrorNil("pathToFile")
	}

	commandExecutor, dirPath, err := d.GetCommandExecutorAndDirectoryPath()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(append([]string{dirPath}, pathToFile...)...)

	toCheck := stringsutils.EnsureSuffix(dirPath, "/")

	if !strings.HasPrefix(filePath, toCheck) {
		return nil, tracederrors.TracedErrorf(
			"filePath '%s' does not start with dirPath '%s' as expected.",
			filePath,
			dirPath,
		)
	}

	return New(commandExecutor, filePath)
}

func (d *Directory) CreateSubDirectory(ctx context.Context, subDirectoryName string, options *filesoptions.CreateOptions) (createdSubDirectory filesinterfaces.Directory, err error) {
	if subDirectoryName == "" {
		return nil, tracederrors.TracedErrorEmptyString("subDirectoryName")
	}

	createdSubDirectory, err = d.GetSubDirectory(subDirectoryName)
	if err != nil {
		return nil, err
	}

	err = createdSubDirectory.Create(ctx, options)
	if err != nil {
		return nil, err
	}

	return createdSubDirectory, nil
}

func (d *Directory) GetSubDirectory(path ...string) (subDirectory filesinterfaces.Directory, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedErrorNil("path")
	}

	commandExecutor, dirPath, err := d.GetCommandExecutorAndDirectoryPath()
	if err != nil {
		return nil, err
	}

	subDirPath := filepath.Join(append([]string{dirPath}, path...)...)

	toCheck := stringsutils.EnsureSuffix(dirPath, "/")

	if !strings.HasPrefix(subDirPath, toCheck) {
		return nil, tracederrors.TracedErrorf(
			"subDirPath '%s' does not start with '%s' as expected.",
			subDirPath,
			toCheck,
		)
	}

	return NewDirectory(commandExecutor, subDirPath)
}

func (d *Directory) IsLocalDirectory() (isLocalDirectory bool, err error) {
	hostDescription, err := d.GetHostDescription()
	if err != nil {
		return false, err
	}

	isLocalDirectory = hostDescription == "localhost"

	return isLocalDirectory, nil
}

func (d *Directory) GetLocalPath() (localPath string, err error) {
	isLocalDirectory, err := d.IsLocalDirectory()
	if err != nil {
		return "", err
	}

	if isLocalDirectory {
		return d.GetPath()
	} else {
		hostDescription, err := d.GetHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf("Directory is on '%s', not on localhost", hostDescription)
	}
}
