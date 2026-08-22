package nativefilesoo

import (
	"context"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (f *File) AppendBytes(ctx context.Context, toWrite []byte) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.AppendBytes(ctx, path, toWrite)
}

func (f *File) AppendString(ctx context.Context, toWrite string) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.AppendString(ctx, path, toWrite)
}

func (f *File) Chown(ctx context.Context, options *parameteroptions.ChownOptions) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Chown(ctx, path, options)
}

func (f *File) GetBaseName() (baseName string, err error) {
	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	return nativefiles.GetBaseName(path)
}

func (f *File) GetHostDescription() (hostDescription string, err error) {
	return "localhost", nil
}

func (f *File) GetLocalPath() (localPath string, err error) {
	return f.GetPath()
}

func (f *File) GetLocalPathOrEmptyStringIfUnset() (localPath string, err error) {
	if f.path == "" {
		return "", nil
	}
	return f.GetPath()
}

func (f *File) GetParentDirectory(ctx context.Context) (parentDirectory filesinterfaces.Directory, err error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	parentPath, err := nativefiles.GetParentDirectoryPath(path)
	if err != nil {
		return nil, err
	}

	return NewDirectoryByPath(parentPath)
}

func (f *File) GetUriAsString() (uri string, err error) {
	path, err := f.GetPath()
	if err != nil {
		return "", err
	}

	return path, nil
}

func (f *File) MoveToPath(ctx context.Context, destPath string, useSudo bool) (movedFile filesinterfaces.File, err error) {
	path, err := f.GetPath()
	if err != nil {
		return nil, err
	}

	err = nativefiles.MoveToPath(ctx, path, destPath, useSudo)
	if err != nil {
		return nil, err
	}

	return NewFileByPath(destPath)
}

func (f *File) SecurelyDelete(ctx context.Context) (err error) {
	path, err := f.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.SecureDelete(ctx, path)
}

func (f *File) String() (path string) {
	return f.path
}

func (d *Directory) Chmod(ctx context.Context, chmodOptions *filesoptions.ChmodOptions) error {
	path, err := d.GetPath()
	if err != nil {
		return err
	}

	return nativefiles.Chmod(ctx, path, chmodOptions)
}

func (d *Directory) CopyContentToDirectory(ctx context.Context, destinationDir filesinterfaces.Directory) (err error) {
	// List all files in source directory
	files, err := d.ListFiles(ctx, &parameteroptions.ListFileOptions{})
	if err != nil {
		return err
	}

	// Copy each file to destination
	for _, file := range files {
		srcFile := file
		// Get just the base name of the file
		baseName, err := srcFile.GetBaseName()
		if err != nil {
			return err
		}

		destFile, err := destinationDir.GetFileInDirectory(baseName)
		if err != nil {
			return err
		}

		err = srcFile.CopyToFile(ctx, destFile, &filesoptions.CopyOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Directory) GetBaseName() (baseName string, err error) {
	path, err := d.GetPath()
	if err != nil {
		return "", err
	}

	return nativefiles.GetBaseName(path)
}

func (d *Directory) GetDirName() (dirName string, err error) {
	return d.GetBaseName()
}

func (d *Directory) GetFileInDirectory(pathToFile ...string) (file filesinterfaces.File, err error) {
	if len(pathToFile) <= 0 {
		return nil, tracederrors.TracedError("pathToFile is empty")
	}

	dirPath, err := d.GetPath()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(dirPath, filepath.Join(pathToFile...))

	return NewFileByPath(fullPath)
}

func (d *Directory) IsLocalDirectory() (isLocalDirectory bool, err error) {
	hostDescription, err := d.GetHostDescription()
	if err != nil {
		return false, err
	}

	return hostDescription == "localhost", nil
}

func (d *Directory) ListFiles(ctx context.Context, listFileOptions *parameteroptions.ListFileOptions) (files []filesinterfaces.File, err error) {
	path, err := d.GetPath()
	if err != nil {
		return nil, err
	}

	if listFileOptions == nil {
		listFileOptions = &parameteroptions.ListFileOptions{}
	}

	filePaths, err := nativefiles.ListFiles(ctx, path, listFileOptions)
	if err != nil {
		return nil, err
	}

	files = make([]filesinterfaces.File, 0, len(filePaths))
	for _, filePath := range filePaths {
		file, err := NewFileByPath(filePath)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (d *Directory) ListSubDirectories(ctx context.Context, options *parameteroptions.ListDirectoryOptions) (subDirectories []filesinterfaces.Directory, err error) {
	path, err := d.GetPath()
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = &parameteroptions.ListDirectoryOptions{}
	}

	// Read directory entries directly
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	subDirectories = make([]filesinterfaces.Directory, 0)
	for _, entry := range entries {
		// Check if it's a directory
		if entry.IsDir() {
			entryPath := filepath.Join(path, entry.Name())
			subDir, err := NewDirectoryByPath(entryPath)
			if err != nil {
				return nil, err
			}
			subDirectories = append(subDirectories, subDir)
		}
	}

	return subDirectories, nil
}
