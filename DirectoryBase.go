package asciichgolangpublic

import (
	"sort"
	"strconv"
	"strings"
)

type DirectoryBase struct {
	parentDirectoryForBaseClass Directory
}

func NewDirectoryBase() (d *DirectoryBase) {
	return new(DirectoryBase)
}

// TODO: Rename to WriteStringtoFile( to make it more generic.
// This renaming is needed to bring GitRepository and Directory together.
func (d *DirectoryBase) WriteStringToFileInDirectory(content string, verbose bool, path ...string) (writtenFile File, err error) {
	if len(path) <= 0 {
		return nil, TracedErrorf("Invalid path='%v'", path)
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	writtenFile, err = parent.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	err = writtenFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return writtenFile, nil
}

func (c *DirectoryBase) WriteBytesToFile(content []byte, verbose bool, path ...string) (writtenFile File, err error) {
	if content == nil {
		return nil, TracedErrorNil("content")
	}

	if len(path) <= 0 {
		return nil, TracedError("path is empty")
	}

	parent, err := c.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	file, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	err = file.WriteBytes(content, verbose)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d *DirectoryBase) CheckExists(verbose bool) (err error) {
	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return err
	}

	exists, err := parent.Exists(verbose)
	if err != nil {
		return err
	}

	if exists {
		return
	}

	path, err := parent.GetPath()
	if err != nil {
		return err
	}

	return TracedErrorf(
		"directory '%s' does not exist", path,
	)
}

func (d *DirectoryBase) CreateFileInDirectory(verbose bool, path ...string) (createdFile File, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	createdFile, err = parent.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	err = createdFile.Create(verbose)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (d *DirectoryBase) CreateFileInDirectoryFromString(content string, verbose bool, pathToCreate ...string) (createdFile File, err error) {
	if len(pathToCreate) <= 0 {
		return nil, TracedErrorf("Invalid pathToCreate='%v'", pathToCreate)
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	createdFile, err = parent.GetFileInDirectory(pathToCreate...)
	if err != nil {
		return nil, err
	}

	parentDir, err := createdFile.GetParentDirectory()
	if err != nil {
		return nil, err
	}

	err = parentDir.Create(verbose)
	if err != nil {
		return nil, err
	}

	err = createdFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (d *DirectoryBase) DeleteFilesMatching(listFileOptions *ListFileOptions) (err error) {
	if listFileOptions == nil {
		return TracedErrorNil("listFileOptions")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return err
	}

	toDelete, err := parent.ListFiles(listFileOptions)
	if err != nil {
		return err
	}

	for _, d := range toDelete {
		err = d.Delete(listFileOptions.Verbose)
		if err != nil {
			return err
		}
	}

	path, err := parent.GetPath()
	if err != nil {
		return err
	}

	hostDescription, err := parent.GetHostDescription()
	if err != nil {
		return err
	}

	if listFileOptions.Verbose {
		LogInfof(
			"Deleted '%d' in directoy '%s' on '%s'",
			len(toDelete),
			path,
			hostDescription,
		)
	}

	return err
}

func (d *DirectoryBase) FileInDirectoryExists(verbose bool, path ...string) (fileExists bool, err error) {
	if len(path) <= 0 {
		return false, TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return false, err
	}

	fileToCheck, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return false, err
	}

	fileExists, err = fileToCheck.Exists(verbose)
	if err != nil {
		return false, err
	}

	return fileExists, nil
}

func (d *DirectoryBase) GetFilePathInDirectory(path ...string) (filePath string, err error) {
	if len(path) <= 0 {
		return "", TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", err
	}

	localFile, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return "", err
	}

	filePath, err = localFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (d *DirectoryBase) GetParentDirectoryForBaseClass() (parentDirectoryForBaseClass Directory, err error) {
	if d.parentDirectoryForBaseClass == nil {
		return nil, TracedError("parentDirectoryForBaseClass not set")
	}
	return d.parentDirectoryForBaseClass, nil
}

func (d *DirectoryBase) GetPathAndHostDescription() (path string, hostDescription string, err error) {
	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", "", err
	}

	path, err = parent.GetPath()
	if err != nil {
		return "", "", err
	}

	hostDescription, err = parent.GetHostDescription()
	if err != nil {
		return "", "", err
	}

	return path, hostDescription, nil
}

func (d *DirectoryBase) ListFilePaths(listFileOptions *ListFileOptions) (filePaths []string, err error) {
	if listFileOptions == nil {
		return nil, TracedErrorNil("listFileOptions")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	files, err := parent.ListFiles(listFileOptions)
	if err != nil {
		return nil, err
	}

	directoryPath, err := parent.GetPath()
	if err != nil {
		return nil, err
	}

	filePaths = []string{}
	for _, f := range files {
		toAdd, err := f.GetPath()
		if err != nil {
			return nil, err
		}

		if listFileOptions.ReturnRelativePaths {
			toAdd = strings.TrimPrefix(toAdd, directoryPath+"/")
		}

		filePaths = append(filePaths, toAdd)
	}

	sort.Strings(filePaths)

	return filePaths, nil
}

func (d *DirectoryBase) ListSubDirectoryPaths(options *ListDirectoryOptions) (subDirectoryPaths []string, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	subDirs, err := parent.ListSubDirectories(options)
	if err != nil {
		return nil, err
	}

	dirPath, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	subDirectoryPaths = []string{}

	for _, subDir := range subDirs {
		path, err := subDir.GetPath()
		if err != nil {
			return nil, err
		}

		toAdd := path
		if options.ReturnRelativePaths {
			toAdd, err = Paths().GetRelativePathTo(
				toAdd,
				dirPath,
			)
			if err != nil {
				return nil, err
			}

			subDirectoryPaths = append(subDirectoryPaths, toAdd)
		}
	}

	sort.Strings(subDirectoryPaths)

	if options.Verbose {
		LogInfof(
			"Listed '%d' sub directory of directory '%s' on host '%s'.",
			len(subDirectoryPaths),
			dirPath,
			hostDescription,
		)
	}

	return subDirectoryPaths, nil
}

func (d *DirectoryBase) MustCheckExists(verbose bool) {
	err := d.CheckExists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DirectoryBase) MustCreateFileInDirectory(verbose bool, path ...string) (createdFile File) {
	createdFile, err := d.CreateFileInDirectory(verbose, path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFile
}

func (d *DirectoryBase) MustCreateFileInDirectoryFromString(content string, verbose bool, pathToCreate ...string) (createdFile File) {
	createdFile, err := d.CreateFileInDirectoryFromString(content, verbose, pathToCreate...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFile
}

func (d *DirectoryBase) MustDeleteFilesMatching(listFileOptions *ListFileOptions) {
	err := d.DeleteFilesMatching(listFileOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DirectoryBase) MustFileInDirectoryExists(verbose bool, path ...string) (fileExists bool) {
	fileExists, err := d.FileInDirectoryExists(verbose, path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return fileExists
}

func (d *DirectoryBase) MustGetFilePathInDirectory(path ...string) (filePath string) {
	filePath, err := d.GetFilePathInDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filePath
}

func (d *DirectoryBase) MustGetParentDirectoryForBaseClass() (parentDirectoryForBaseClass Directory) {
	parentDirectoryForBaseClass, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentDirectoryForBaseClass
}

func (d *DirectoryBase) MustGetPathAndHostDescription() (path string, hostDescription string) {
	path, hostDescription, err := d.GetPathAndHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path, hostDescription
}

func (d *DirectoryBase) MustListFilePaths(listFileOptions *ListFileOptions) (filePaths []string) {
	filePaths, err := d.ListFilePaths(listFileOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filePaths
}

func (d *DirectoryBase) MustListSubDirectoryPaths(options *ListDirectoryOptions) (subDirectoryPaths []string) {
	subDirectoryPaths, err := d.ListSubDirectoryPaths(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectoryPaths
}

func (d *DirectoryBase) MustReadFileInDirectoryAsInt64(path ...string) (value int64) {
	value, err := d.ReadFileInDirectoryAsInt64(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return value
}

func (d *DirectoryBase) MustReadFileInDirectoryAsLines(path ...string) (content []string) {
	content, err := d.ReadFileInDirectoryAsLines(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (d *DirectoryBase) MustReadFileInDirectoryAsString(path ...string) (content string) {
	content, err := d.ReadFileInDirectoryAsString(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (d *DirectoryBase) MustReadFirstLineOfFileInDirectoryAsString(path ...string) (firstLine string) {
	firstLine, err := d.ReadFirstLineOfFileInDirectoryAsString(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return firstLine
}

func (d *DirectoryBase) MustSetParentDirectoryForBaseClass(parentDirectoryForBaseClass Directory) {
	err := d.SetParentDirectoryForBaseClass(parentDirectoryForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (d *DirectoryBase) MustWriteBytesToFile(content []byte, verbose bool, path ...string) (writtenFile File) {
	writtenFile, err := d.WriteBytesToFile(content, verbose, path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return writtenFile
}

func (d *DirectoryBase) MustWriteStringToFileInDirectory(content string, verbose bool, path ...string) (writtenFile File) {
	writtenFile, err := d.WriteStringToFileInDirectory(content, verbose, path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return writtenFile
}

func (d *DirectoryBase) ReadFileInDirectoryAsInt64(path ...string) (value int64, err error) {
	if len(path) <= 0 {
		return -1, TracedError("path has no elements")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return -1, err
	}

	content, err := parent.ReadFileInDirectoryAsString(path...)
	if err != nil {
		return -1, err
	}

	content = strings.TrimSpace(content)

	value, err = strconv.ParseInt(content, 10, 64)
	if err != nil {
		return -1, TracedErrorf(
			"Failed to parse file content as int64: %w",
			err,
		)
	}

	return value, nil
}

func (d *DirectoryBase) ReadFileInDirectoryAsLines(path ...string) (content []string, err error) {
	if len(path) == 0 {
		return nil, TracedError("path is empty")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return nil, err
	}

	fileToRead, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	content, err = fileToRead.ReadAsLines()
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (d *DirectoryBase) ReadFileInDirectoryAsString(path ...string) (content string, err error) {
	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", err
	}

	fileToRead, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return "", err
	}

	content, err = fileToRead.ReadAsString()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (d *DirectoryBase) ReadFirstLineOfFileInDirectoryAsString(path ...string) (firstLine string, err error) {
	if len(path) <= 0 {
		return "", TracedError("No path given")
	}

	parent, err := d.GetParentDirectoryForBaseClass()
	if err != nil {
		return "", err
	}

	f, err := parent.GetFileInDirectory(path...)
	if err != nil {
		return "", err
	}

	return f.ReadFirstLine()
}

func (d *DirectoryBase) SetParentDirectoryForBaseClass(parentDirectoryForBaseClass Directory) (err error) {
	if parentDirectoryForBaseClass == nil {
		return TracedErrorNil("parentDirectoryForBaseClass")
	}

	d.parentDirectoryForBaseClass = parentDirectoryForBaseClass

	return nil
}
