package asciichgolangpublic

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pathsutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A CommandExecutorFile implements the functionality of a `File` by
// executing commands (like: test, stat, cat...).
//
// The benefit of this apporach is an easy way to access files on any
// remote system like VMs, Containers, Hosts... while it easy to chain
// like inside Container on VM behind Jumphost...
//
// The downside of this is the poor performance and the possiblity to see
// in the process table which operations where done.
type CommandExecutorFile struct {
	FileBase
	commandExecutor CommandExecutor
	filePath        string
}

func NewCommandExecutorFile() (c *CommandExecutorFile) {
	c = new(CommandExecutorFile)
	c.SetParentFileForBaseClass(c)

	return c
}

func (c *CommandExecutorFile) AppendBytes(toWrite []byte, verbose bool) (err error) {
	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorFile) AppendString(toWrite string, verbose bool) (err error) {
	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"cat >> '%s'",
					filePath,
				),
			},
			Verbose:     false, // Would potentially expose the content to write.
			StdinString: toWrite,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorFile) Chmod(chmodOptions *ChmodOptions) (err error) {
	if chmodOptions == nil {
		return tracederrors.TracedErrorNil("chmodOptions")
	}

	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return err
	}

	permissionsString, err := chmodOptions.GetPermissionsString()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: []string{"chmod", permissionsString, filePath},
			Verbose: chmodOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandExecutorFile) Chown(options *ChownOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	userAndGroupForCommand, err := options.GetUserName()
	if err != nil {
		return err
	}

	if options.IsGroupNameSet() {
		groupName, err := options.GetGroupName()
		if err != nil {
			return err
		}

		userAndGroupForCommand += ":" + groupName
	}

	command := []string{"chown", userAndGroupForCommand, path}

	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return err
	}

	if options.Verbose {
		logging.LogChangedf(
			"Changed ownership of file '%s' to '%s' on host '%s'",
			path,
			userAndGroupForCommand,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorFile) CopyToFile(destFile File, verbose bool) (err error) {
	if destFile == nil {
		return tracederrors.TracedErrorNil("destFile")
	}

	return tracederrors.TracedErrorNotImplemented()
}

func (c *CommandExecutorFile) Create(verbose bool) (err error) {
	commandExecutor, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		return err
	}

	exists, err := c.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfof(
			"File '%s' on '%s' already exists.",
			filePath,
			hostDescription,
		)
	} else {
		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{"touch", filePath},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"File '%s' on '%s' created.",
			filePath,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorFile) Delete(verbose bool) (err error) {
	commandExecutor, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		return err
	}

	if !pathsutils.IsAbsolutePath(filePath) {
		return tracederrors.TracedErrorf(
			"For security reasons deleting a file is only implemented for absolute paths but got '%s'",
			filePath,
		)
	}

	exists, err := c.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{"rm", filePath},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"File '%s' on '%s' deleted.",
			filePath,
			hostDescription,
		)
	} else {
		logging.LogInfof(
			"File '%s' on '%s' already absent.",
			filePath,
			hostDescription,
		)
	}

	return nil
}

func (c *CommandExecutorFile) Exists(verbose bool) (exists bool, err error) {
	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return false, err
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"test -f '%s' && echo yes || echo no",
					filePath,
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
			"Unexpected output when checking for file to exist: '%s'",
			output,
		)
	}

	if verbose {
		hostDescription, err := c.GetHostDescription()
		if err != nil {
			return false, err
		}

		if exists {
			logging.LogInfof(
				"File '%s' on host '%s' exists.",
				filePath,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"File '%s' on host '%s' does not exist.",
				filePath,
				hostDescription,
			)
		}
	}

	return exists, nil
}

func (c *CommandExecutorFile) GetBaseName() (baseName string, err error) {
	filePath, err := c.GetFilePath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(filePath)

	return baseName, nil
}

func (c *CommandExecutorFile) GetCommandExecutor() (commandExecutor CommandExecutor, err error) {

	return c.commandExecutor, nil
}

func (c *CommandExecutorFile) GetCommandExecutorAndFilePath() (commandExecutor CommandExecutor, filePath string, err error) {
	commandExecutor, err = c.GetCommandExecutor()
	if err != nil {
		return nil, "", err
	}

	filePath, err = c.GetFilePath()
	if err != nil {
		return nil, "", err
	}

	return commandExecutor, filePath, nil
}

func (c *CommandExecutorFile) GetCommandExecutorAndFilePathAndHostDescription() (commandExecutor CommandExecutor, filePath string, hostDescription string, err error) {
	commandExecutor, filePath, err = c.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, "", "", err
	}

	hostDescription, err = commandExecutor.GetHostDescription()
	if err != nil {
		return nil, "", "", err
	}

	return commandExecutor, filePath, hostDescription, nil
}

func (c *CommandExecutorFile) GetDeepCopy() (deepCopy File) {
	d := NewCommandExecutorFile()

	*d = *c

	if c.commandExecutor != nil {
		d.commandExecutor = MustGetDeepCopyOfCommandExecutor(c.commandExecutor)
	}

	return d
}

func (c *CommandExecutorFile) GetFilePath() (filePath string, err error) {
	if c.filePath == "" {
		return "", tracederrors.TracedErrorf("filePath not set")
	}

	return c.filePath, nil
}

func (c *CommandExecutorFile) GetHostDescription() (hostDescription string, err error) {
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

func (c *CommandExecutorFile) GetLocalPath() (localPath string, err error) {
	isRunningOnLocalhost, err := c.IsRunningOnLocalhost()
	if err != nil {
		return "", err
	}

	_, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		return "", err
	}

	if !isRunningOnLocalhost {
		return "", tracederrors.TracedErrorf(
			"File '%s' is not local. It is on '%s'.",
			filePath,
			hostDescription,
		)
	}

	if !pathsutils.IsAbsolutePath(filePath) {
		return "", tracederrors.TracedErrorf(
			"File path '%s' is not absolute.",
			filePath,
		)
	}

	return filePath, nil
}

func (c *CommandExecutorFile) GetLocalPathOrEmptyStringIfUnset() (localPath string, err error) {
	isRunningOnLocalHost, err := c.IsRunningOnLocalhost()
	if err != nil {
		return "", err
	}

	if isRunningOnLocalHost {
		return c.filePath, nil
	}

	return "", tracederrors.TracedError("Not running on local host.")
}

func (c *CommandExecutorFile) GetParentDirectory() (parentDirectory Directory, err error) {
	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	p, err := NewCommandExecutorDirectory(commandExecutor)
	if err != nil {
		return nil, err
	}

	dirName := filepath.Dir(filePath)

	err = p.SetDirPath(dirName)
	if err != nil {
		return nil, err
	}

	parentDirectory = p

	return parentDirectory, nil
}

func (c *CommandExecutorFile) GetPath() (path string, err error) {
	isRunningOnLocalhost, err := c.IsRunningOnLocalhost()
	if err != nil {
		return "", err
	}

	if isRunningOnLocalhost {
		path, err := c.GetLocalPath()
		if err == nil {
			return path, nil
		}
	} else {
		return "", tracederrors.TracedErrorNotImplemented()
	}

	if path == "" {
		return "", tracederrors.TracedError("path is empty string after evaluation.")
	}

	return path, nil
}

func (c *CommandExecutorFile) GetPathAndHostDescription() (path string, hostDescription string, err error) {
	path, err = c.GetPath()
	if err != nil {
		return "", "", err
	}

	hostDescription, err = c.GetHostDescription()
	if err != nil {
		return "", "", err
	}

	return path, hostDescription, nil
}

func (c *CommandExecutorFile) GetSizeBytes() (fileSize int64, err error) {
	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return -1, err
	}

	fileSize, err = commandExecutor.RunCommandAndGetStdoutAsInt64(
		&RunCommandOptions{
			Command: []string{
				"stat", "--printf=%s", filePath,
			},
			Verbose: false,
		},
	)
	if err != nil {
		return -1, err
	}

	return fileSize, nil
}

func (c *CommandExecutorFile) GetUriAsString() (uri string, err error) {
	isRunningOnLocalHost, err := c.IsRunningOnLocalhost()
	if err != nil {
		return "", err
	}

	if isRunningOnLocalHost {
		filePath, err := c.GetFilePath()
		if err != nil {
			return "", err
		}

		if pathsutils.IsRelativePath(filePath) {
			return "", tracederrors.TracedErrorf("Only implemeted for absolute paths but got '%s'", filePath)
		}

		uri = "file://" + filePath

		return uri, nil
	}

	_, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		return "", err
	}

	return "", tracederrors.TracedErrorf("not implemented for '%s' on '%s'", filePath, hostDescription)
}

func (c *CommandExecutorFile) IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return false, err
	}

	isRunningOnLocalhost = hostDescription == "localhost"

	return isRunningOnLocalhost, nil
}

func (c *CommandExecutorFile) MoveToPath(path string, useSudo bool, verbose bool) (movedFile File, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	srcPath, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	commandToUse := []string{"mv", srcPath, path}
	if useSudo {
		commandToUse = append([]string{"sudo"}, commandToUse...)
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: commandToUse,
		},
	)
	if err != nil {
		return nil, err
	}

	if verbose {
		logging.LogChangedf(
			"Moved file '%s' to '%s' on host '%s'.",
			srcPath,
			path,
			hostDescription,
		)
	}

	return GetCommandExecutorFileByPath(commandExecutor, path)
}

func (c *CommandExecutorFile) MustAppendBytes(toWrite []byte, verbose bool) {
	err := c.AppendBytes(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustAppendString(toWrite string, verbose bool) {
	err := c.AppendString(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustChmod(chmodOptions *ChmodOptions) {
	err := c.Chmod(chmodOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustChown(options *ChownOptions) {
	err := c.Chown(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustCopyToFile(destFile File, verbose bool) {
	err := c.CopyToFile(destFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustCreate(verbose bool) {
	err := c.Create(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustDelete(verbose bool) {
	err := c.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustExists(verbose bool) (exist bool) {
	exist, err := c.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exist
}

func (c *CommandExecutorFile) MustGetBaseName() (baseName string) {
	baseName, err := c.GetBaseName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return baseName
}

func (c *CommandExecutorFile) MustGetCommandExecutor() (commandExecutor CommandExecutor) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor
}

func (c *CommandExecutorFile) MustGetCommandExecutorAndFilePath() (commandExecutor CommandExecutor, filePath string) {
	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor, filePath
}

func (c *CommandExecutorFile) MustGetCommandExecutorAndFilePathAndHostDescription() (commandExecutor CommandExecutor, filePath string, hostDescription string) {
	commandExecutor, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandExecutor, filePath, hostDescription
}

func (c *CommandExecutorFile) MustGetFilePath() (filePath string) {
	filePath, err := c.GetFilePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return filePath
}

func (c *CommandExecutorFile) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := c.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (c *CommandExecutorFile) MustGetLocalPath() (localPath string) {
	localPath, err := c.GetLocalPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localPath
}

func (c *CommandExecutorFile) MustGetLocalPathOrEmptyStringIfUnset() (localPath string) {
	localPath, err := c.GetLocalPathOrEmptyStringIfUnset()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localPath
}

func (c *CommandExecutorFile) MustGetParentDirectory() (parentDirectory Directory) {
	parentDirectory, err := c.GetParentDirectory()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentDirectory
}

func (c *CommandExecutorFile) MustGetPath() (path string) {
	path, err := c.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (c *CommandExecutorFile) MustGetPathAndHostDescription() (path string, hostDescription string) {
	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path, hostDescription
}

func (c *CommandExecutorFile) MustGetSizeBytes() (fileSize int64) {
	fileSize, err := c.GetSizeBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileSize
}

func (c *CommandExecutorFile) MustGetUriAsString() (uri string) {
	uri, err := c.GetUriAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return uri
}

func (c *CommandExecutorFile) MustIsRunningOnLocalhost() (isRunningOnLocalhost bool) {
	isRunningOnLocalhost, err := c.IsRunningOnLocalhost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunningOnLocalhost
}

func (c *CommandExecutorFile) MustMoveToPath(path string, useSudo bool, verbose bool) (movedFile File) {
	movedFile, err := c.MoveToPath(path, useSudo, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return movedFile
}

func (c *CommandExecutorFile) MustReadAsBytes() (content []byte) {
	content, err := c.ReadAsBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (c *CommandExecutorFile) MustReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte) {
	firstBytes, err := c.ReadFirstNBytes(numberOfBytesToRead)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return firstBytes
}

func (c *CommandExecutorFile) MustSecurelyDelete(verbose bool) {
	err := c.SecurelyDelete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustSetCommandExecutor(commandExecutor CommandExecutor) {
	err := c.SetCommandExecutor(commandExecutor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustSetFilePath(filePath string) {
	err := c.SetFilePath(filePath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustTruncate(newSizeBytes int64, verbose bool) {
	err := c.Truncate(newSizeBytes, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) MustWriteBytes(toWrite []byte, verbose bool) {
	err := c.WriteBytes(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorFile) ReadAsBytes() (content []byte, err error) {
	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	content, err = commandExecutor.RunCommandAndGetStdoutAsBytes(
		&RunCommandOptions{
			Command: []string{"cat", filePath},
			Verbose: false,
		},
	)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (c *CommandExecutorFile) ReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte, err error) {
	if numberOfBytesToRead < 0 {
		return nil, tracederrors.TracedErrorf("Invalid number of bytest to read: '%d'", numberOfBytesToRead)
	}

	commandExecutor, filePath, err := c.GetCommandExecutorAndFilePath()
	if err != nil {
		return nil, err
	}

	firstBytes, err = commandExecutor.RunCommandAndGetStdoutAsBytes(
		&RunCommandOptions{
			Command: []string{
				"head",
				fmt.Sprintf(
					"--bytes=%d",
					numberOfBytesToRead,
				),
				filePath,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return firstBytes, nil
}

func (c *CommandExecutorFile) SecurelyDelete(verbose bool) (err error) {
	commandExecutor, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		return err
	}

	exits, err := c.Exists(verbose)
	if err != nil {
		return err
	}

	if exits {
		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{"shred", "-u", filePath},
				Verbose: verbose,
			},
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Securely deleted file '%s' on '%s'.",
				filePath,
				hostDescription,
			)
		}
	} else {
		if verbose {
			logging.LogInfof(
				"File '%s' on '%s' is alreay absent.",
				filePath,
				hostDescription,
			)
		}
	}

	return nil
}

func (c *CommandExecutorFile) SetCommandExecutor(commandExecutor CommandExecutor) (err error) {
	c.commandExecutor = commandExecutor

	return nil
}

func (c *CommandExecutorFile) SetFilePath(filePath string) (err error) {
	if filePath == "" {
		return tracederrors.TracedErrorf("filePath is empty string")
	}

	c.filePath = filePath

	return nil
}

func (c *CommandExecutorFile) Truncate(newSizeBytes int64, verbose bool) (err error) {
	if newSizeBytes < 0 {
		return tracederrors.TracedErrorf(
			"Invalid size for truncating: newSizeBytes='%d'",
			newSizeBytes,
		)
	}

	currentSize, err := c.GetSizeBytes()
	if err != nil {
		return err
	}

	path, hostDescription, err := c.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if currentSize == newSizeBytes {
		logging.LogInfof(
			"File '%s' on host '%s' is already of size '%d' bytes. Skip truncate.",
			path,
			hostDescription,
			newSizeBytes,
		)
	} else {
		commandExecutor, err := c.GetCommandExecutor()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(
			&RunCommandOptions{
				Command: []string{
					"truncate",
					fmt.Sprintf("-s%d", newSizeBytes),
					path,
				},
				Verbose:            verbose,
				LiveOutputOnStdout: verbose,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"File '%s' on host '%s' is truncated to '%d' bytes.",
			path,
			hostDescription,
			newSizeBytes,
		)
	}

	return nil
}

func (c *CommandExecutorFile) WriteBytes(toWrite []byte, verbose bool) (err error) {
	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	commandExecutor, filePath, hostDescription, err := c.GetCommandExecutorAndFilePathAndHostDescription()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		&RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"cat > '%s'",
					filePath,
				),
			},
			StdinString: string(toWrite),
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf(
			"Wrote '%d' bytes to file '%s' on '%s'",
			len(toWrite),
			filePath,
			hostDescription,
		)
	}

	return nil
}
