package files

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// A CommandExecutorDirectory implements the functionality of a `Directory` by
// executing commands (like: test, stat, cat...).
//
// The benefit of this apporach is an easy way to access directories on any
// remote system like VMs, Containers, Hosts... while it easy to chain
// like inside Container on VM behind Jumphost...
//
// The downside of this is the poor performance and the possiblity to see
// in the process table which operations where done.
type CommandExecutorDirectory struct {
	DirectoryBase
	commandExecutor commandexecutorinterfaces.CommandExecutor
	dirPath         string
}

func GetCommandExecutorDirectoryByPath(commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (c *CommandExecutorDirectory, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	c, err = NewCommandExecutorDirectory(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = c.SetDirPath(path)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func GetLocalCommandExecutorDirectoryByPath(path string) (c *CommandExecutorDirectory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return GetCommandExecutorDirectoryByPath(commandexecutorbashoo.Bash(), path)
}

func NewCommandExecutorDirectory(commandExecutor commandexecutorinterfaces.CommandExecutor) (c *CommandExecutorDirectory, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	c = new(CommandExecutorDirectory)
	c.MustSetParentDirectoryForBaseClass(c)

	err = c.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CommandExecutorDirectory) Chmod(ctx context.Context, chmodOptions *parameteroptions.ChmodOptions) (err error) {
	if chmodOptions == nil {
		return tracederrors.TracedErrorNil("chmodOptions")
	}

	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return err
	}

	permissionString, err := chmodOptions.GetPermissionsString()
	if err != nil {
		return err
	}

	command := []string{"chmod", permissionString, dirPath}

	if chmodOptions.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Chmod '%s' for '%s' on '%s'", permissionString, dirPath, hostDescription)

	return nil
}

func (c *CommandExecutorDirectory) CopyContentToDirectory(destinationDir filesinterfaces.Directory, verbose bool) (err error) {
	if destinationDir == nil {
		return tracederrors.TracedErrorNil("destinationDir")
	}

	commandExecutor, srcDirPath, srcHostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return err
	}

	destHostDescription, err := destinationDir.GetHostDescription()
	if err != nil {
		return err
	}

	if srcHostDescription != destHostDescription {
		return tracederrors.TracedErrorf(
			"Copy from one host to another not imlemented. srcHostDescription='%s' != destHostDescription='%s'",
			srcHostDescription,
			destHostDescription,
		)
	}

	destDirPath, err := destinationDir.GetPath()
	if err != nil {
		return err
	}

	ctx := contextutils.GetVerbosityContextByBool(verbose)
	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"cp", "-r", "-v", srcDirPath, destDirPath},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf(
			"Copied directory '%s' on '%s' to '%s' on '%s'",
			srcDirPath,
			srcHostDescription,
			destDirPath,
			destHostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorDirectory) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	exists, err := c.Exists(ctx)
	if err != nil {
		return err
	}

	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Directory '%s' on '%s' already exists.", dirPath, hostDescription)
	} else {
		command := []string{"mkdir", "-p", dirPath}

		if options != nil && options.UseSudo {
			command = append([]string{"sudo"}, command...)
		}

		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Directory '%s' on '%s' created.", dirPath, hostDescription)
	}

	return nil
}

func (c *CommandExecutorDirectory) CreateSubDirectory(ctx context.Context, subDirectoryName string, options *filesoptions.CreateOptions) (createdSubDirectory filesinterfaces.Directory, err error) {
	if subDirectoryName == "" {
		return nil, tracederrors.TracedErrorEmptyString("subDirectoryName")
	}

	createdSubDirectory, err = c.GetSubDirectory(subDirectoryName)
	if err != nil {
		return nil, err
	}

	err = createdSubDirectory.Create(ctx, options)
	if err != nil {
		return nil, err
	}

	return createdSubDirectory, nil
}

func (c *CommandExecutorDirectory) Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error) {
	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return err
	}

	if !pathsutils.IsAbsolutePath(dirPath) {
		return tracederrors.TracedErrorf(
			"For security reasons deleting a directory is only implemented for absolute paths but got '%s'",
			dirPath,
		)
	}

	exists, err := c.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"rm", "-rf", dirPath},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Directory '%s' on '%s' deleted.",
			dirPath,
			hostDescription,
		)
	} else {
		logging.LogInfof(
			"Directory '%s' is already absent on '%s'.",
			dirPath,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorDirectory) Exists(ctx context.Context) (exists bool, err error) {
	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return false, err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"test -d '%s' && echo yes || echo no",
					dirPath,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)

	if output == "yes" {
		exists = true
	} else if output == "no" {
		exists = false
	} else {
		return false, tracederrors.TracedErrorf(
			"Unexpected output when evalution directory '%s' exists on '%s'",
			dirPath,
			hostDescription,
		)
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Directory '%s' exists on host '%s'.", dirPath, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Directory '%s' does not exist on host '%s'.", dirPath, hostDescription)
	}

	return exists, nil
}

func (c *CommandExecutorDirectory) GetBaseName() (baseName string, err error) {
	dirPath, err := c.GetDirPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(dirPath)

	if baseName == "" {
		return "", tracederrors.TracedError("baseName is nil after evaluation")
	}

	return baseName, nil
}

func (c *CommandExecutorDirectory) GetCommandExecutor() (commandExecutor commandexecutorinterfaces.CommandExecutor, err error) {
	return c.commandExecutor, nil
}

func (c *CommandExecutorDirectory) GetCommandExecutorAndDirPath() (commandExecutor commandexecutorinterfaces.CommandExecutor, dirPath string, err error) {
	commandExecutor, err = c.GetCommandExecutor()
	if err != nil {
		return nil, "", err
	}

	dirPath, err = c.GetDirPath()
	if err != nil {
		return nil, "", err
	}

	return commandExecutor, dirPath, nil
}

func (c *CommandExecutorDirectory) GetCommandExecutorAndDirPathAndHostDescription() (commandExecutor commandexecutorinterfaces.CommandExecutor, dirPath string, hostDescription string, err error) {
	commandExecutor, dirPath, err = c.GetCommandExecutorAndDirPath()
	if err != nil {
		return nil, "", "", err
	}

	hostDescription, err = commandExecutor.GetHostDescription()
	if err != nil {
		return nil, "", "", err
	}

	return commandExecutor, dirPath, hostDescription, nil
}

func (c *CommandExecutorDirectory) GetDirName() (parentPath string, err error) {
	path, err := c.GetPath()
	if err != nil {
		return "", err
	}

	return pathsutils.GetDirPath(path)
}

func (c *CommandExecutorDirectory) GetDirPath() (dirPath string, err error) {
	if c.dirPath == "" {
		return "", tracederrors.TracedErrorf("dirPath not set")
	}

	return c.dirPath, nil
}

func (c *CommandExecutorDirectory) GetFileInDirectory(pathToFile ...string) (file filesinterfaces.File, err error) {
	if len(pathToFile) <= 0 {
		return nil, tracederrors.TracedErrorNil("pathToFile")
	}

	commandExecutor, dirPath, err := c.GetCommandExecutorAndDirPath()
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

	f := NewCommandExecutorFile()

	err = f.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = f.SetFilePath(filePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (c *CommandExecutorDirectory) GetHostDescription() (hostDescription string, err error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	hostDescription, err = commandExecutor.GetHostDescription()
	if err != nil {
		return "", err
	}

	return hostDescription, nil
}

func (c *CommandExecutorDirectory) GetLocalPath() (localPath string, err error) {
	isLocalDirectory, err := c.IsLocalDirectory()
	if err != nil {
		return "", err
	}

	if isLocalDirectory {
		localPath, err = c.GetDirPath()
		if err != nil {
			return "", err
		}

		return localPath, nil
	} else {
		hostDescription, err := c.GetHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf("Directory is on '%s', not on localhost", hostDescription)
	}
}

func (c *CommandExecutorDirectory) GetParentDirectory() (parent filesinterfaces.Directory, err error) {
	parentPath, err := c.GetDirName()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return GetCommandExecutorDirectoryByPath(
		commandExecutor,
		parentPath,
	)
}

func (c *CommandExecutorDirectory) GetPath() (path string, err error) {
	path, err = c.GetDirPath()
	if err != nil {
		return "", err
	}

	if !pathsutils.IsAbsolutePath(path) {
		return "", tracederrors.TracedErrorf("path '%s' is not absolute.", path)
	}

	return path, nil
}

func (c *CommandExecutorDirectory) GetSubDirectory(path ...string) (subDirectory filesinterfaces.Directory, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedErrorNil("path")
	}

	commandExecutor, dirPath, err := c.GetCommandExecutorAndDirPath()
	if err != nil {
		return nil, err
	}

	subdir, err := NewCommandExecutorDirectory(commandExecutor)
	if err != nil {
		return nil, err
	}

	subDirPath := filepath.Join(append([]string{dirPath}, path...)...)

	toCheck := stringsutils.EnsureSuffix(dirPath, "/")
	if !strings.HasPrefix(
		subDirPath,
		toCheck,
	) {
		return nil, tracederrors.TracedErrorf(
			"subDirPath '%s' does not start with '%s' as expected.",
			subDirPath,
			toCheck,
		)
	}

	err = subdir.SetDirPath(subDirPath)
	if err != nil {
		return nil, err
	}

	return subdir, nil
}

func (c *CommandExecutorDirectory) IsLocalDirectory() (isLocalDirectory bool, err error) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return false, err
	}

	isLocalDirectory = hostDescription == "localhost"

	return isLocalDirectory, nil
}

func (c *CommandExecutorDirectory) ListFilePaths(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (filePaths []string, err error) {
	if listFileOptions == nil {
		return nil, tracederrors.TracedErrorNil("listFileOptions")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	dirPath, err := c.GetPath()
	if err != nil {
		return nil, err
	}

	commandToUse := []string{"find", dirPath, "-type", "f"}
	if listFileOptions.NonRecursive {
		commandToUse = []string{"find", dirPath, "-type", "f", "-maxdepth", "1"}
	}

	foundPaths, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: commandToUse,
		},
	)
	if err != nil {
		return nil, err
	}

	filePaths, err = pathsutils.FilterPaths(foundPaths, listFileOptions)
	if err != nil {
		return nil, err
	}

	if listFileOptions.ReturnRelativePaths {
		filePaths, err = pathsutils.GetRelativePathsTo(filePaths, dirPath)
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(filePaths)

	return filePaths, nil
}

func (c *CommandExecutorDirectory) ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []filesinterfaces.File, err error) {
	if listFileOptions == nil {
		return nil, tracederrors.TracedErrorNil("listFileOptions")
	}

	optionsToUse := listFileOptions.GetDeepCopy()

	optionsToUse.ReturnRelativePaths = true

	paths, err := c.ListFilePaths(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	files = []filesinterfaces.File{}
	for _, path := range paths {
		toAdd, err := c.GetFileInDirectory(path)
		if err != nil {
			return nil, err
		}

		files = append(files, toAdd)
	}

	return files, nil
}

func (c *CommandExecutorDirectory) ListSubDirectories(options *parameteroptions.ListDirectoryOptions) (subDirectories []filesinterfaces.Directory, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	path, err := c.GetPath()
	if err != nil {
		return nil, err
	}

	findCommand := []string{"find", path, "-type", "d"}
	findCommand = append(findCommand, "-mindepth", "1") // do not list the current directory itself.

	if !options.Recursive {
		findCommand = append(findCommand, "-maxdepth", "1")
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	stdoutLines, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		contextutils.GetVerbosityContextByBool(options.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: findCommand,
		},
	)
	if err != nil {
		return nil, err
	}

	pathsToAdd := []string{}
	for _, line := range stdoutLines {
		if line == "" {
			continue
		}

		pathToAdd := strings.TrimPrefix(line, "./")
		if pathsutils.IsAbsolutePath(pathToAdd) {
			pathToAdd, err = pathsutils.GetRelativePathTo(
				pathToAdd,
				path,
			)

			if err != nil {
				return nil, err
			}

			pathsToAdd = append(pathsToAdd, pathToAdd)
		}
	}

	sort.Strings(pathsToAdd)

	subDirectories = []filesinterfaces.Directory{}
	for _, pathToAdd := range pathsToAdd {
		toAdd, err := c.GetSubDirectory(pathToAdd)
		if err != nil {
			return nil, err
		}

		subDirectories = append(subDirectories, toAdd)
	}

	return subDirectories, nil
}

func (c *CommandExecutorDirectory) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorDirectory) SetDirPath(dirPath string) (err error) {
	if dirPath == "" {
		return tracederrors.TracedErrorf("dirPath is empty string")
	}

	c.dirPath = dirPath

	return nil
}
