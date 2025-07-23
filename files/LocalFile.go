package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A LocalFile represents a locally available file.
type LocalFile struct {
	FileBase
	path string
}

func GetLocalFileByFile(inputFile File) (localFile *LocalFile, err error) {
	if inputFile == nil {
		return nil, tracederrors.TracedErrorNil("inputFile")
	}

	localFile, ok := inputFile.(*LocalFile)
	if !ok {
		return nil, tracederrors.TracedError("inputFile is not a LocalFile")
	}

	return localFile, nil
}

func GetLocalFileByPath(localPath string) (l *LocalFile, err error) {
	return NewLocalFileByPath(localPath)
}

func MustGetLocalFileByFile(inputFile File) (localFile *LocalFile) {
	localFile, err := GetLocalFileByFile(inputFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localFile
}

func MustGetLocalFileByPath(localPath string) (l *LocalFile) {
	l, err := GetLocalFileByPath(localPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return l
}

func MustNewLocalFileByPath(localPath string) (l *LocalFile) {
	l, err := NewLocalFileByPath(localPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return l
}

func NewLocalFile() (l *LocalFile) {
	l = new(LocalFile)

	// Allow usage of the base class functions:
	l.MustSetParentFileForBaseClass(l)

	return l
}

func NewLocalFileByPath(localPath string) (l *LocalFile, err error) {
	if localPath == "" {
		return nil, tracederrors.TracedErrorEmptyString("localPath")
	}

	l = NewLocalFile()

	err = l.SetPath(localPath)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *LocalFile) String() string {
	if l.path == "" {
		return "<LocalFile.path NOT SET>"
	}

	return l.path
}

// Delete a file if it exists.
// If the file is already absent this function does nothing.
func (l *LocalFile) Delete(verbose bool) (err error) {
	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	exists, err := l.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		err = os.Remove(path)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delet localFile '%s': '%w'", path, err)
		}

		if verbose {
			logging.LogChangedf("Local file '%s' deleted.", path)
		}
	} else {
		if verbose {
			logging.LogInfof("Local file '%s' is already absent. Skip delete.", path)
		}
	}

	return nil
}

func (l *LocalFile) AppendBytes(toWrite []byte, verbose bool) (err error) {
	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	fileToWrite, err := os.OpenFile(localPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return tracederrors.TracedErrorf(
			"Unable to open file '%s' to append: '%w'",
			localPath,
			err,
		)
	}
	_, err = fileToWrite.Write(toWrite)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to append: '%w'", err)
	}

	err = fileToWrite.Close()
	if err != nil {
		return tracederrors.TracedErrorf("Unable to close file after append: '%w'", err)
	}

	if verbose {
		logging.LogChangedf("Appended data to localfile '%s'.", localPath)
	}

	return nil
}

func (l *LocalFile) AppendString(toWrite string, verbose bool) (err error) {
	err = l.AppendBytes([]byte(toWrite), verbose)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalFile) Chmod(chmodOptions *parameteroptions.ChmodOptions) (err error) {
	if chmodOptions == nil {
		return tracederrors.TracedErrorNil("chmodOptions")
	}

	chmodString, err := chmodOptions.GetPermissionsString()
	if err != nil {
		return err
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	_, err = commandexecutor.Bash().RunCommand(
		contextutils.GetVerbosityContextByBool(chmodOptions.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{"chmod", chmodString, localPath},
		},
	)
	if err != nil {
		return err
	}

	if chmodOptions.Verbose {
		logging.LogChangedf("Chmod '%s' for local file '%s'.", chmodString, localPath)
	}

	return nil
}

func (l *LocalFile) Chown(options *parameteroptions.ChownOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	path, hostDescription, err := l.GetPathAndHostDescription()
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

	_, err = commandexecutor.Bash().RunCommand(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
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

func (l *LocalFile) CopyToFile(destFile File, verbose bool) (err error) {
	if destFile == nil {
		return tracederrors.TracedErrorNil("destFile")
	}

	content, err := l.ReadAsBytes()
	if err != nil {
		return err
	}

	err = destFile.WriteBytes(content, verbose)
	if err != nil {
		return err
	}

	if verbose {
		srcPath, err := l.GetLocalPath()
		if err != nil {
			return err
		}

		destPath, err := destFile.GetLocalPath()
		if err != nil {
			return err
		}

		logging.LogChangedf("Copied '%s' to '%s'", srcPath, destPath)
	}

	return nil
}

func (l *LocalFile) Create(verbose bool) (err error) {
	exists, err := l.Exists(verbose)
	if err != nil {
		return err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if exists {
		if verbose {
			logging.LogInfof("Local file '%s' already exists", path)
		}
	} else {
		err = l.WriteString("", false)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf("Local file '%s' created", path)
		}
	}

	return nil
}

func (l *LocalFile) Exists(verbose bool) (exists bool, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return false, err
	}

	fileInfo, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, tracederrors.TracedErrorf("Unable to evaluate if local file exists: '%w'", err)
	}

	return !fileInfo.IsDir(), err
}

func (l *LocalFile) GetBaseName() (baseName string, err error) {
	path, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(path)

	if baseName == "" {
		return "", tracederrors.TracedErrorf(
			"Base name is empty string after evaluation of path='%s'",
			path,
		)
	}

	return baseName, nil
}

func (l *LocalFile) GetDeepCopy() (deepCopy File) {
	deepCopyLocalFile := NewLocalFile()
	deepCopyLocalFile.path = l.path

	deepCopy = deepCopyLocalFile

	return deepCopy
}

func (l *LocalFile) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
}

func (l *LocalFile) GetLocalPath() (path string, err error) {
	return l.GetPath()
}

func (l *LocalFile) GetLocalPathOrEmptyStringIfUnset() (localPath string, err error) {
	return l.path, nil
}

func (l *LocalFile) GetParentDirectory() (parentDirectory Directory, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	localDirPath := filepath.Dir(localPath)

	parentDirectory, err = GetLocalDirectoryByPath(localDirPath)
	if err != nil {
		return nil, err
	}

	return parentDirectory, nil
}

func (l *LocalFile) GetParentFileForBaseClassAsLocalFile() (parentAsLocalFile *LocalFile, err error) {
	parent, err := l.GetParentFileForBaseClass()
	if err != nil {
		return nil, err
	}

	parentAsLocalFile, err = GetLocalFileByFile(parent)
	if err != nil {
		return nil, err
	}

	return parentAsLocalFile, nil
}

func (l *LocalFile) GetPath() (path string, err error) {
	if l.path == "" {
		return "", fmt.Errorf("path not set")
	}

	return l.path, nil
}

func (l *LocalFile) GetSizeBytes() (fileSizeBytes int64, err error) {
	path, err := l.GetPath()
	if err != nil {
		return -1, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return -1, tracederrors.TracedError(err)
	}
	fileSizeBytes = fi.Size()

	return fileSizeBytes, nil
}

func (l *LocalFile) GetUriAsString() (uri string, err error) {
	path, err := l.GetPath()
	if err != nil {
		return "", err
	}

	if pathsutils.IsRelativePath(path) {
		return "", tracederrors.TracedErrorf("Only implemeted for absolute paths but got '%s'", path)
	}

	uri = "file://" + path

	return uri, nil
}

func (l *LocalFile) IsPathSet() (isSet bool) {
	return false
}

func (l *LocalFile) MoveToPath(path string, useSudo bool, verbose bool) (movedFile File, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString(path)
	}

	srcPath, hostDescription, err := l.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	if useSudo {
		_, err = commandexecutor.Bash().RunCommand(
			contextutils.ContextSilent(),
			&parameteroptions.RunCommandOptions{
				Command: []string{"sudo", "mv", srcPath, path},
			},
		)
		if err != nil {
			return nil, err
		}
	} else {
		err = os.Rename(srcPath, path)
		if err != nil {
			return nil, tracederrors.TracedErrorf(
				"Move '%s' to '%s' on host '%s' failed: %w",
				srcPath,
				path,
				hostDescription,
				err,
			)
		}
	}

	if verbose {
		logging.LogChangedf(
			"Moved '%s' to '%s' on host '%s'.",
			srcPath,
			path,
			hostDescription,
		)
	}

	return GetLocalFileByPath(path)
}

func (l *LocalFile) MustAppendBytes(toWrite []byte, verbose bool) {
	err := l.AppendBytes(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustAppendString(toWrite string, verbose bool) {
	err := l.AppendString(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustChmod(chmodOptions *parameteroptions.ChmodOptions) {
	err := l.Chmod(chmodOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustChown(options *parameteroptions.ChownOptions) {
	err := l.Chown(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustCopyToFile(destFile File, verbose bool) {
	err := l.CopyToFile(destFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustCreate(verbose bool) {
	err := l.Create(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustDelete(verbose bool) {
	err := l.Delete(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustExists(verbose bool) (exists bool) {
	exists, err := l.Exists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalFile) MustGetBaseName() (baseName string) {
	baseName, err := l.GetBaseName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return baseName
}

func (l *LocalFile) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := l.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (l *LocalFile) MustGetLocalPath() (path string) {
	path, err := l.GetLocalPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (l *LocalFile) MustGetLocalPathOrEmptyStringIfUnset() (localPath string) {
	localPath, err := l.GetLocalPathOrEmptyStringIfUnset()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localPath
}

func (l *LocalFile) MustGetParentDirectory() (parentDirectory Directory) {
	parentDirectory, err := l.GetParentDirectory()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentDirectory
}

func (l *LocalFile) MustGetParentFileForBaseClassAsLocalFile() (parentAsLocalFile *LocalFile) {
	parentAsLocalFile, err := l.GetParentFileForBaseClassAsLocalFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentAsLocalFile
}

func (l *LocalFile) MustGetPath() (path string) {
	path, err := l.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path
}

func (l *LocalFile) MustGetSizeBytes() (fileSizeBytes int64) {
	fileSizeBytes, err := l.GetSizeBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileSizeBytes
}

func (l *LocalFile) MustGetUriAsString() (uri string) {
	uri, err := l.GetUriAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return uri
}

func (l *LocalFile) MustMoveToPath(path string, useSudo bool, verbose bool) (movedFile File) {
	movedFile, err := l.MoveToPath(path, useSudo, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return movedFile
}

func (l *LocalFile) MustReadAsBytes() (content []byte) {
	content, err := l.ReadAsBytes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (l *LocalFile) MustReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte) {
	firstBytes, err := l.ReadFirstNBytes(numberOfBytesToRead)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return firstBytes
}

func (l *LocalFile) MustSetPath(path string) {
	err := l.SetPath(path)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustTruncate(newSizeBytes int64, verbose bool) {
	err := l.Truncate(newSizeBytes, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) MustWriteBytes(toWrite []byte, verbose bool) {
	err := l.WriteBytes(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalFile) ReadAsBytes() (content []byte, err error) {
	path, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	content, err = os.ReadFile(path)
	if err != nil {
		return nil, tracederrors.TracedError(err)
	}

	return content, err
}

func (l *LocalFile) ReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte, err error) {
	if numberOfBytesToRead <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid numberOfBytesToRead: '%d'", numberOfBytesToRead)
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}

	defer fd.Close()

	firstBytes = make([]byte, numberOfBytesToRead)
	readBytes, err := fd.Read(firstBytes)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, tracederrors.TracedError(err.Error())
		}
	}

	firstBytes = firstBytes[:readBytes]

	return firstBytes, nil
}

func (l *LocalFile) SecurelyDelete(ctx context.Context) (err error) {
	pathToDelete, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if !pathsutils.IsAbsolutePath(pathToDelete) {
		return tracederrors.TracedErrorf("pathToDelete='%v' is not absolute", pathToDelete)
	}

	return filesutils.SecureDelete(ctx, pathToDelete)
}

func (l *LocalFile) SetPath(path string) (err error) {
	if path == "" {
		return tracederrors.TracedError("path is empty string")
	}

	path, err = pathsutils.GetAbsolutePath(path)
	if err != nil {
		return err
	}

	if !pathsutils.IsAbsolutePath(path) {
		return tracederrors.TracedErrorf(
			"Path '%s' is not absolute. Beware this is an internal issue since the code before this line should fix that.",
			path,
		)
	}

	l.path = path

	return nil
}

func (l *LocalFile) Truncate(newSizeBytes int64, verbose bool) (err error) {
	if newSizeBytes < 0 {
		return tracederrors.TracedErrorf("Invalid newSizeBytes='%d'", newSizeBytes)
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	currentSize, err := l.GetSizeBytes()
	if err != nil {
		return err
	}

	if currentSize == newSizeBytes {
		logging.LogInfof(
			"Local file '%s' is already of size '%d' bytes. Skip truncate.",
			localPath,
			newSizeBytes,
		)
	} else {
		fileToTruncate, err := os.OpenFile(localPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return tracederrors.TracedErrorf(
				"Unable to open file '%s' to truncate: '%w'",
				localPath,
				err,
			)
		}
		defer fileToTruncate.Close()

		err = fileToTruncate.Truncate(newSizeBytes)
		if err != nil {
			return tracederrors.TracedErrorf(
				"Unable to truncate file '%s': '%w'",
				localPath,
				err,
			)
		}

		if verbose {
			logging.LogChangedf(
				"Truncated local file '%s' to new size '%d'.",
				localPath,
				newSizeBytes,
			)
		}
	}

	return nil
}

func (l *LocalFile) WriteBytes(toWrite []byte, verbose bool) (err error) {
	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(localPath, toWrite, 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to write file '%s': %w", localPath, err)
	}

	if verbose {
		logging.LogInfof("Wrote data to '%s'", localPath)
	}

	return nil
}

func (l *LocalFile) MustGetAccessPermissionsString() (permissionString string) {
	permissionString, err := l.GetAccessPermissionsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permissionString
}

func (l *LocalFile) MustGetAccessPermissions() (permissions int) {
	permissions, err := l.GetAccessPermissions()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return permissions
}

func (l *LocalFile) GetAccessPermissions() (permissions int, err error) {
	path, err := l.GetPath()
	if err != nil {
		return 0, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, tracederrors.TracedErrorf("Unable to get fileInfo of '%s': %w", path, err)
	}

	perm := fileInfo.Mode().Perm()

	return int(perm), nil
}

func (l *LocalFile) GetAccessPermissionsString() (permissionsString string, err error) {
	permissions, err := l.GetAccessPermissions()
	if err != nil {
		return "", err
	}

	permissionsString, err = unixfilepermissionsutils.GetPermissionString(permissions)
	if err != nil {
		return "", err
	}

	return permissionsString, nil
}
