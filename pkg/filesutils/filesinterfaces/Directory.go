package filesinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

type Directory interface {
	Chmod(ctx context.Context, chmodOptions *parameteroptions.ChmodOptions) (err error)
	CopyContentToDirectory(destinationDir Directory, verbose bool) (err error)
	Create(ctx context.Context, options *filesoptions.CreateOptions) (err error)
	CreateSubDirectory(ctx context.Context, subDirectoryName string, options *filesoptions.CreateOptions) (createdSubDirectory Directory, err error)
	Delete(ctx context.Context, options *filesoptions.DeleteOptions) (err error)
	Exists(ctx context.Context) (exists bool, err error)
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

	// All methods below this line can be implemented by embedding the `DirectoryBase` struct:
	CheckExists(ctx context.Context) (err error)
	CreateFileInDirectory(ctx context.Context, path string, options *filesoptions.CreateOptions) (createdFile File, err error)
	GetFilePathInDirectory(path ...string) (filePath string, err error)
	GetPathAndHostDescription() (dirPath string, hostDescription string, err error)
	DeleteFilesMatching(ctx context.Context, listFileOptons *parameteroptions.ListFileOptions) (err error)
	FileInDirectoryExists(verbose bool, path ...string) (exists bool, err error)
	ListFilePaths(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (filePaths []string, err error)
	ListSubDirectoryPaths(options *parameteroptions.ListDirectoryOptions) (subDirectoryPaths []string, err error)
	ReadFileInDirectoryAsInt64(path ...string) (content int64, err error)
	ReadFileInDirectoryAsLines(path ...string) (content []string, err error)
	ReadFileInDirectoryAsString(path ...string) (content string, err error)
	ReadFirstLineOfFileInDirectoryAsString(path ...string) (firstLine string, err error)
	WriteStringToFile(ctx context.Context, path string, content string, options *filesoptions.WriteOptions) (writtenFile File, err error)
}
