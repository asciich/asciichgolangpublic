package asciichgolangpublic

type Directory interface {
	Create(verbose bool) (err error)
	CreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory, err error)
	Delete(verbose bool) (err error)
	Exists() (exists bool, err error)
	GetBaseName() (baseName string, err error)
	GetDirName() (dirName string, err error)
	GetFileInDirectory(pathToFile ...string) (file File, err error)
	GetLocalPath() (localPath string, err error)
	GetSubDirectories(options *ListDirectoryOptions) (subDirectories []Directory, err error)
	GetSubDirectory(path ...string) (subDirectory Directory, err error)
	IsLocalDirectory() (isLocalDirectory bool)
	MustCreate(verbose bool)
	MustCreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory)
	MustDelete(verbose bool)
	MustExists() (exists bool)
	MustGetBaseName() (baseName string)
	MustGetDirName() (dirName string)
	MustGetSubDirectory(path ...string) (subDirectory Directory)
	MustGetSubDirectories(options *ListDirectoryOptions) (subDirectories []Directory)
	MustGetFileInDirectory(pathToFile ...string) (file File)
	MustGetLocalPath() (localPath string)

	// All methods below this line can be implemented by embedding the `DirectoryBase` struct:
	GetFilePathInDirectory(path ...string) (filePath string, err error)
	FileInDirectoryExists(path ...string) (exists bool, err error)
	MustGetFilePathInDirectory(path ...string) (filePath string)
	MustFileInDirectoryExists(path ...string) (exists bool)
	MustReadFileInDirectoryAsString(path ...string) (content string)
	ReadFileInDirectoryAsString(path ...string) (content string, err error)
}
