package asciichgolangpublic

import (
	"fmt"
	"path/filepath"
	"strings"
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
	commandExecutor CommandExecutor
	dirPath         string
}

func GetCommandExecutorDirectoryByPath(commandExecutor CommandExecutor, path string) (c *CommandExecutorDirectory, err error) {
	if commandExecutor == nil {
		return nil, TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, TracedErrorEmptyString("path")
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
		return nil, TracedErrorEmptyString("path")
	}

	return GetCommandExecutorDirectoryByPath(Bash(), path)
}

func MustGetCommandExecutorDirectoryByPath(commandExecutor CommandExecutor, path string) (c *CommandExecutorDirectory) {
	c, err := GetCommandExecutorDirectoryByPath(commandExecutor, path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return c
}

func MustGetLocalCommandExecutorDirectoryByPath(path string) (c *CommandExecutorDirectory) {
	c, err := GetLocalCommandExecutorDirectoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return c
}

func MustNewCommandExecutorDirectory(commandExecutor CommandExecutor) (c *CommandExecutorDirectory) {
	c, err := NewCommandExecutorDirectory(commandExecutor)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return c
}

func NewCommandExecutorDirectory(commandExecutor CommandExecutor) (c *CommandExecutorDirectory, err error) {
	if commandExecutor == nil {
		return nil, TracedErrorNil("commandExecutor")
	}

	c = new(CommandExecutorDirectory)
	c.MustSetParentDirectoryForBaseClass(c)

	err = c.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CommandExecutorDirectory) Chmod(chmodOptions *ChmodOptions) (err error) {
	if chmodOptions == nil {
		return TracedErrorNil("chmodOptions")
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
		&RunCommandOptions{
			Command: command,
			Verbose: chmodOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	if chmodOptions.Verbose {
		LogChangedf(
			"Chmod '%s' for '%s' on '%s'",
			permissionString,
			dirPath,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorDirectory) CopyContentToDirectory(destinationDir Directory, verbose bool) (err error) {
	if destinationDir == nil {
		return TracedErrorNil("destinationDir")
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
		return TracedErrorf(
			"Copy from one host to another not imlemented. srcHostDescription='%s' != destHostDescription='%s'",
			srcHostDescription,
			destHostDescription,
		)
	}

	destDirPath, err := destinationDir.GetPath()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command:            []string{"cp", "-r", "-v", srcDirPath, destDirPath},
			LiveOutputOnStdout: verbose,
			Verbose:            verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf(
			"Copied directory '%s' on '%s' to '%s' on '%s'",
			srcDirPath,
			srcHostDescription,
			destDirPath,
			destHostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorDirectory) Create(verbose bool) (err error) {
	exists, err := c.Exists(verbose)
	if err != nil {
		return err
	}

	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return err
	}

	if exists {
		LogInfof(
			"Directory '%s' on '%s' already exists.",
			dirPath,
			hostDescription,
		)
	} else {
		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{"mkdir", "-p", dirPath},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		LogChangedf(
			"Directory '%s' on '%s' created.",
			dirPath,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorDirectory) CreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory, err error) {
	if subDirectoryName == "" {
		return nil, TracedErrorEmptyString("subDirectoryName")
	}

	createdSubDirectory, err = c.GetSubDirectory(subDirectoryName)
	if err != nil {
		return nil, err
	}

	err = createdSubDirectory.Create(verbose)
	if err != nil {
		return nil, err
	}

	return createdSubDirectory, nil
}

func (c *CommandExecutorDirectory) Delete(verbose bool) (err error) {
	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return err
	}

	if !Paths().IsAbsolutePath(dirPath) {
		return TracedErrorf(
			"For security reasons deleting a directory is only implemented for absolute paths but got '%s'",
			dirPath,
		)
	}

	exists, err := c.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{"rm", "-rf", dirPath},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		LogChangedf(
			"Directory '%s' on '%s' deleted.",
			dirPath,
			hostDescription,
		)
	} else {
		LogInfof(
			"Directory '%s' is already absent on '%s'.",
			dirPath,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorDirectory) Exists(verbose bool) (exists bool, err error) {
	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		return false, err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
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
		return false, TracedErrorf(
			"Unexpected output when evalution directory '%s' exists on '%s'",
			dirPath,
			hostDescription,
		)
	}

	if verbose {
		if exists {
			LogInfof(
				"Directory '%s' exists on host '%s'.",
				dirPath,
				hostDescription,
			)
		} else {
			LogInfof(
				"Directory '%s' exists on host '%s'.",
				dirPath,
				hostDescription,
			)
		}
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
		return "", TracedError("baseName is nil after evaluation")
	}

	return baseName, nil
}

func (c *CommandExecutorDirectory) GetCommandExecutor() (commandExecutor CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *CommandExecutorDirectory) GetCommandExecutorAndDirPath() (commandExecutor CommandExecutor, dirPath string, err error) {
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

func (c *CommandExecutorDirectory) GetCommandExecutorAndDirPathAndHostDescription() (commandExecutor CommandExecutor, dirPath string, hostDescription string, err error) {
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

	return Paths().GetDirPath(path)
}

func (c *CommandExecutorDirectory) GetDirPath() (dirPath string, err error) {
	if c.dirPath == "" {
		return "", TracedErrorf("dirPath not set")
	}

	return c.dirPath, nil
}

func (c *CommandExecutorDirectory) GetFileInDirectory(pathToFile ...string) (file File, err error) {
	if len(pathToFile) <= 0 {
		return nil, TracedErrorNil("pathToFile")
	}

	commandExecutor, dirPath, err := c.GetCommandExecutorAndDirPath()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(append([]string{dirPath}, pathToFile...)...)

	toCheck := Strings().EnsureSuffix(dirPath, "/")

	if !strings.HasPrefix(filePath, toCheck) {
		return nil, TracedErrorf(
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

		return "", TracedErrorf("Directory is on '%s', not on localhost", hostDescription)
	}
}

func (c *CommandExecutorDirectory) GetParentDirectory() (parent Directory, err error) {
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

	if !Paths().IsAbsolutePath(path) {
		return "", TracedErrorf("path '%s' is not absolute.", path)
	}

	return path, nil
}

func (c *CommandExecutorDirectory) GetSubDirectory(path ...string) (subDirectory Directory, err error) {
	if len(path) <= 0 {
		return nil, TracedErrorNil("path")
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

	toCheck := Strings().EnsureSuffix(dirPath, "/")
	if !strings.HasPrefix(
		subDirPath,
		toCheck,
	) {
		return nil, TracedErrorf(
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

func (c *CommandExecutorDirectory) ListFilePaths(listFileOptions *ListFileOptions) (filePaths []string, err error) {
	if listFileOptions == nil {
		return nil, TracedErrorNil("listFileOptions")
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
		&RunCommandOptions{
			Command: commandToUse,
		},
	)
	if err != nil {
		return nil, err
	}

	filePaths, err = Paths().FilterPaths(foundPaths, listFileOptions)
	if err != nil {
		return nil, err
	}

	if listFileOptions.ReturnRelativePaths {
		filePaths, err = Paths().GetRelativePathsTo(filePaths, dirPath)
		if err != nil {
			return nil, err
		}
	}

	filePaths = Slices().SortStringSlice(filePaths)

	return filePaths, nil
}

func (c *CommandExecutorDirectory) ListFiles(listFileOptions *ListFileOptions) (files []File, err error) {
	if listFileOptions == nil {
		return nil, TracedErrorNil("listFileOptions")
	}

	optionsToUse := listFileOptions.GetDeepCopy()

	optionsToUse.ReturnRelativePaths = true

	paths, err := c.ListFilePaths(optionsToUse)
	if err != nil {
		return nil, err
	}

	files = []File{}
	for _, path := range paths {
		toAdd, err := c.GetFileInDirectory(path)
		if err != nil {
			return nil, err
		}

		files = append(files, toAdd)
	}

	return files, nil
}

func (c *CommandExecutorDirectory) ListSubDirectories(options *ListDirectoryOptions) (subDirectories []Directory, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
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
		&RunCommandOptions{
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
		if Paths().IsAbsolutePath(pathToAdd) {
			pathToAdd, err = Paths().GetRelativePathTo(
				pathToAdd,
				path,
			)

			if err != nil {
				return nil, err
			}

			pathsToAdd = append(pathsToAdd, pathToAdd)
		}
	}

	pathsToAdd = Slices().SortStringSlice(pathsToAdd)

	subDirectories = []Directory{}
	for _, pathToAdd := range pathsToAdd {
		toAdd, err := c.GetSubDirectory(pathToAdd)
		if err != nil {
			return nil, err
		}

		subDirectories = append(subDirectories, toAdd)
	}

	return subDirectories, nil
}

func (c *CommandExecutorDirectory) MustChmod(chmodOptions *ChmodOptions) {
	err := c.Chmod(chmodOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDirectory) MustCopyContentToDirectory(destinationDir Directory, verbose bool) {
	err := c.CopyContentToDirectory(destinationDir, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDirectory) MustCreate(verbose bool) {
	err := c.Create(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDirectory) MustCreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory) {
	createdSubDirectory, err := c.CreateSubDirectory(subDirectoryName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdSubDirectory
}

func (c *CommandExecutorDirectory) MustDelete(verbose bool) {
	err := c.Delete(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDirectory) MustExists(verbose bool) (exists bool) {
	exists, err := c.Exists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (c *CommandExecutorDirectory) MustGetBaseName() (baseName string) {
	baseName, err := c.GetBaseName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return baseName
}

func (c *CommandExecutorDirectory) MustGetCommandExecutor() (commandExecutor CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorDirectory) MustGetCommandExecutorAndDirPath() (commandExecutor CommandExecutor, dirPath string) {
	commandExecutor, dirPath, err := c.GetCommandExecutorAndDirPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutor, dirPath
}

func (c *CommandExecutorDirectory) MustGetCommandExecutorAndDirPathAndHostDescription() (commandExecutor CommandExecutor, dirPath string, hostDescription string) {
	commandExecutor, dirPath, hostDescription, err := c.GetCommandExecutorAndDirPathAndHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandExecutor, dirPath, hostDescription
}

func (c *CommandExecutorDirectory) MustGetDirName() (dirName string) {
	dirName, err := c.GetDirName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dirName
}

func (c *CommandExecutorDirectory) MustGetDirPath() (dirPath string) {
	dirPath, err := c.GetDirPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dirPath
}

func (c *CommandExecutorDirectory) MustGetFileInDirectory(pathToFile ...string) (file File) {
	file, err := c.GetFileInDirectory(pathToFile...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return file
}

func (c *CommandExecutorDirectory) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostDescription
}

func (c *CommandExecutorDirectory) MustGetLocalPath() (localPath string) {
	localPath, err := c.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localPath
}

func (c *CommandExecutorDirectory) MustGetParentDirectory() (parent Directory) {
	parent, err := c.GetParentDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parent
}

func (c *CommandExecutorDirectory) MustGetPath() (path string) {
	path, err := c.GetPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path
}

func (c *CommandExecutorDirectory) MustGetSubDirectory(path ...string) (subDirectory Directory) {
	subDirectory, err := c.GetSubDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectory
}

func (c *CommandExecutorDirectory) MustIsLocalDirectory() (isLocalDirectory bool) {
	isLocalDirectory, err := c.IsLocalDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isLocalDirectory
}

func (c *CommandExecutorDirectory) MustListFilePaths(listFileOptions *ListFileOptions) (filePaths []string) {
	filePaths, err := c.ListFilePaths(listFileOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filePaths
}

func (c *CommandExecutorDirectory) MustListFiles(listFileOptions *ListFileOptions) (files []File) {
	files, err := c.ListFiles(listFileOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return files
}

func (c *CommandExecutorDirectory) MustListSubDirectories(options *ListDirectoryOptions) (subDirectories []Directory) {
	subDirectories, err := c.ListSubDirectories(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectories
}

func (c *CommandExecutorDirectory) MustSetCommandExecutor(commandExecutor CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDirectory) MustSetDirPath(dirPath string) {
	err := c.SetDirPath(dirPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorDirectory) SetCommandExecutor(commandExecutor CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorDirectory) SetDirPath(dirPath string) (err error) {
	if dirPath == "" {
		return TracedErrorf("dirPath is empty string")
	}

	c.dirPath = dirPath

	return nil
}
