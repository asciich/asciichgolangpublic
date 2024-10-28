package asciichgolangpublic

type Directory interface {
	Chmod(chmodOptions *ChmodOptions) (err error)
	CopyContentToDirectory(destinationDir Directory, verbose bool) (err error)
	Create(verbose bool) (err error)
	CreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory, err error)
	Delete(verbose bool) (err error)
	Exists() (exists bool, err error)
	GetBaseName() (baseName string, err error)
	GetDirName() (dirName string, err error)
	GetFileInDirectory(pathToFile ...string) (file File, err error)
	GetHostDescription() (hostDescription string, err error)
	// Returns the path on the local machine. If the path is not available locally an error is returned.
	GetLocalPath() (localPath string, err error)
	// Returns the absolute path to the file without any indication of the host.
	GetPath() (dirPath string, err error)
	GetSubDirectories(options *ListDirectoryOptions) (subDirectories []Directory, err error)
	GetSubDirectory(path ...string) (subDirectory Directory, err error)
	IsLocalDirectory() (isLocalDirectory bool, err error)
	MustChmod(chmodOptions *ChmodOptions)
	MustCopyContentToDirectory(destinationDir Directory, verbose bool)
	MustCreate(verbose bool)
	MustCreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory)
	MustDelete(verbose bool)
	MustExists() (exists bool)
	MustGetBaseName() (baseName string)
	MustGetDirName() (dirName string)
	MustGetFileInDirectory(pathToFile ...string) (file File)
	MustGetHostDescription() (hostDescription string)
	// Returns the path on the local machine. If the path is not available locally an error is returned.
	MustGetLocalPath() (localPath string)
	// Returns the absolute path to the file without any indication of the host.
	MustGetPath() (dirPath string)
	MustGetSubDirectory(path ...string) (subDirectory Directory)
	MustGetSubDirectories(options *ListDirectoryOptions) (subDirectories []Directory)
	MustIsLocalDirectory() (isLocalDirectory bool)

	// All methods below this line can be implemented by embedding the `DirectoryBase` struct:
	GetFilePathInDirectory(path ...string) (filePath string, err error)
	FileInDirectoryExists(path ...string) (exists bool, err error)
	MustGetFilePathInDirectory(path ...string) (filePath string)
	MustFileInDirectoryExists(path ...string) (exists bool)
	MustReadFileInDirectoryAsLines(path ...string) (content []string)
	MustReadFileInDirectoryAsString(path ...string) (content string)
	MustWriteStringToFileInDirectory(content string, verbose bool, path ...string) (writtenFile File)
	ReadFileInDirectoryAsLines(path ...string) (content []string, err error)
	ReadFileInDirectoryAsString(path ...string) (content string, err error)
	WriteStringToFileInDirectory(content string, verbose bool, path ...string) (writtenFile File, err error)
}
