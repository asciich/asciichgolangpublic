package asciichgolangpublic

type Directory interface {
	Create(verbose bool) (err error)
	Delete(verbose bool) (err error)
	Exists() (exists bool, err error)
	GetFileInDirectory(pathToFile ...string) (file File, err error)
	GetLocalPath() (localPath string, err error)
	GetSubDirectory(path ...string) (subDirectory Directory, err error)
	IsLocalDirectory() (isLocalDirectory bool)
	MustCreate(verbose bool)
	MustDelete(verbose bool)
	MustExists() (exists bool)
	MustGetSubDirectory(path ...string) (subDirectory Directory)
	MustGetFileInDirectory(pathToFile ...string) (file File)
	MustGetLocalPath() (localPath string)

	// All methods below this line can be implemented by embedding the `DirectoryBase` struct:
	GetFilePathInDirectory(path ...string) (filePath string, err error)
	MustGetFilePathInDirectory(path ...string) (filePath string)
}
