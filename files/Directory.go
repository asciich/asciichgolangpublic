package files

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

type Directory interface {
	Chmod(chmodOptions *parameteroptions.ChmodOptions) (err error)
	CopyContentToDirectory(destinationDir Directory, verbose bool) (err error)
	Create(verbose bool) (err error)
	CreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory, err error)
	Delete(verbose bool) (err error)
	Exists(verbose bool) (exists bool, err error)
	GetBaseName() (baseName string, err error)
	GetDirName() (dirName string, err error)
	GetFileInDirectory(pathToFile ...string) (file File, err error)
	GetHostDescription() (hostDescription string, err error)
	// Returns the path on the local machine. If the path is not available locally an error is returned.
	GetLocalPath() (localPath string, err error)
	GetParentDirectory() (parentDirectory Directory, err error)
	// Returns the absolute path to the file without any indication of the host.
	GetPath() (dirPath string, err error)
	// TODO rename GetSubDirectory with GetDirectoryByPath to make it consistent.
	GetSubDirectory(path ...string) (subDirectory Directory, err error)
	IsLocalDirectory() (isLocalDirectory bool, err error)
	ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []File, err error)
	ListSubDirectories(options *parameteroptions.ListDirectoryOptions) (subDirectories []Directory, err error)
	MustChmod(chmodOptions *parameteroptions.ChmodOptions)
	MustCopyContentToDirectory(destinationDir Directory, verbose bool)
	MustCreate(verbose bool)
	MustCreateSubDirectory(subDirectoryName string, verbose bool) (createdSubDirectory Directory)
	MustDelete(verbose bool)
	MustExists(verbose bool) (exists bool)
	MustGetBaseName() (baseName string)
	MustGetDirName() (dirName string)
	MustGetFileInDirectory(pathToFile ...string) (file File)
	MustGetHostDescription() (hostDescription string)
	// Returns the path on the local machine. If the path is not available locally an error is returned.
	MustGetLocalPath() (localPath string)
	MustGetParentDirectory() (parentDirectory Directory)
	// Returns the absolute path to the file without any indication of the host.
	MustGetPath() (dirPath string)
	MustGetSubDirectory(path ...string) (subDirectory Directory)
	MustIsLocalDirectory() (isLocalDirectory bool)

	// All methods below this line can be implemented by embedding the `DirectoryBase` struct:
	CheckExists(ctx context.Context) (err error)
	CreateFileInDirectory(verbose bool, path ...string) (createdFile File, err error)
	GetFilePathInDirectory(path ...string) (filePath string, err error)
	GetPathAndHostDescription() (dirPath string, hostDescription string, err error)
	DeleteFilesMatching(ctx context.Context, listFileOptons *parameteroptions.ListFileOptions) (err error)
	FileInDirectoryExists(verbose bool, path ...string) (exists bool, err error)
	ListFilePaths(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (filePaths []string, err error)
	ListSubDirectoryPaths(options *parameteroptions.ListDirectoryOptions) (subDirectoryPaths []string, err error)
	MustCheckExists(ctx context.Context)
	MustCreateFileInDirectory(verbose bool, path ...string) (createdFile File)
	MustGetFilePathInDirectory(path ...string) (filePath string)
	MustGetPathAndHostDescription() (pathString string, hostDescription string)
	MustFileInDirectoryExists(verbose bool, path ...string) (exists bool)
	MustListSubDirectoryPaths(options *parameteroptions.ListDirectoryOptions) (subDirectoryPaths []string)
	MustReadFileInDirectoryAsInt64(path ...string) (content int64)
	MustReadFileInDirectoryAsLines(path ...string) (content []string)
	MustReadFileInDirectoryAsString(path ...string) (content string)
	MustReadFirstLineOfFileInDirectoryAsString(path ...string) (firstLine string)
	MustWriteStringToFileInDirectory(content string, verbose bool, path ...string) (writtenFile File)
	ReadFileInDirectoryAsInt64(path ...string) (content int64, err error)
	ReadFileInDirectoryAsLines(path ...string) (content []string, err error)
	ReadFileInDirectoryAsString(path ...string) (content string, err error)
	ReadFirstLineOfFileInDirectoryAsString(path ...string) (firstLine string, err error)
	WriteStringToFileInDirectory(content string, verbose bool, path ...string) (writtenFile File, err error)
}
