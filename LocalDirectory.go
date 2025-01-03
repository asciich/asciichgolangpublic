package asciichgolangpublic

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type LocalDirectory struct {
	DirectoryBase
	localPath string
}

func GetLocalDirectoryByPath(path string) (l *LocalDirectory, err error) {
	if path == "" {
		return nil, TracedErrorEmptyString("path")
	}

	localDirectory := NewLocalDirectory()

	err = localDirectory.SetLocalPath(path)
	if err != nil {
		return nil, err
	}

	return localDirectory, nil
}

func MustGetLocalDirectoryByPath(path string) (l *LocalDirectory) {
	l, err := GetLocalDirectoryByPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return l
}

func NewLocalDirectory() (l *LocalDirectory) {
	l = new(LocalDirectory)

	// Allow usage of the base class functions:
	l.MustSetParentDirectoryForBaseClass(l)

	return l
}

func (l *LocalDirectory) Chmod(chmodOptions *ChmodOptions) (err error) {
	if chmodOptions == nil {
		return TracedErrorNil("chmodOptions")
	}

	chmodString, err := chmodOptions.GetPermissionsString()
	if err != nil {
		return err
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: []string{"chmod", chmodString, localPath},
			Verbose: chmodOptions.Verbose,
		},
	)
	if err != nil {
		return err
	}

	if chmodOptions.Verbose {
		LogChangedf("Chmod '%s' for local directory '%s'.", chmodString, localPath)
	}

	return nil
}

func (l *LocalDirectory) CopyContentToDirectory(destinationDir Directory, verbose bool) (err error) {
	if destinationDir == nil {
		return TracedError("destinationDir is empty string")
	}

	srcPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	destPath, err := destinationDir.GetLocalPath()
	if err != nil {
		return err
	}

	copyCommand := []string{
		"cp",
		"-rv",
		srcPath + "/.",
		destPath + "/.",
	}

	stdout, err := Bash().RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command:            copyCommand,
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Copied files from '%s' to '%s':\n%s",
			srcPath,
			destPath,
			stdout,
		)
	}

	return nil
}

func (l *LocalDirectory) CopyContentToLocalDirectory(destDirectory *LocalDirectory, verbose bool) (err error) {
	if destDirectory == nil {
		return TracedErrorNil("destDirectory")
	}

	destPath, err := destDirectory.GetLocalPath()
	if err != nil {
		return err
	}

	srcPath, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	exists, err := l.Exists(verbose)
	if err != nil {
		return err
	}

	if !exists {
		return TracedErrorf("Unable to copy content to local directory, '%s' does not exist.", srcPath)
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command:            []string{"cp", "-r", "-v", srcPath + "/.", destPath + "/."},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		LogChangedf("Copied files from '%s' to '%s'.", srcPath, destPath)
	}

	return nil
}

func (l *LocalDirectory) CopyFileToTemporaryFile(verbose bool, filePath ...string) (copy File, err error) {
	if len(filePath) <= 0 {
		return nil, TracedErrorEmptyString("filePath")
	}

	fileToCopy, err := l.GetFileInDirectory(filePath...)
	if err != nil {
		return nil, err
	}

	copy, err = TemporaryFiles().CreateTemporaryFileFromFile(fileToCopy, verbose)
	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (l *LocalDirectory) CopyFileToTemporaryFileAsLocalFile(verbose bool, filePath ...string) (copy *LocalFile, err error) {
	if len(filePath) <= 0 {
		return nil, TracedErrorEmptyString("filePath")
	}

	interfaceCopy, err := l.CopyFileToTemporaryFile(verbose, filePath...)
	if err != nil {
		return nil, err
	}

	copy, ok := interfaceCopy.(*LocalFile)
	if !ok {
		return nil, TracedErrorf("Internal error: Unable to convert to *LocalFile: '%v'", interfaceCopy)
	}

	return copy, nil
}

func (l *LocalDirectory) Create(verbose bool) (err error) {
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
			LogInfof("Local directory '%s' already exists. Skip create.", path)
		}
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				parentDirectoy, err := l.GetParentDirectory()
				if err != nil {
					return err
				}

				err = parentDirectoy.Create(verbose)
				if err != nil {
					return err
				}

				err = os.Mkdir(path, os.ModePerm)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		existsAfterCreate, err := l.Exists(verbose)
		if err != nil {
			return err
		}

		if !existsAfterCreate {
			return TracedErrorf("Local directory '%s' does not exist after creation.", path)
		}

		if verbose {
			LogChangedf("Created local directory '%s'", path)
		}
	}

	return nil
}

func (l *LocalDirectory) CreateFileInDirectory(verbose bool, path ...string) (createdFile File, err error) {
	createdFile, err = l.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	parentDirectory, err := createdFile.GetParentDirectory()
	if err != nil {
		return nil, err
	}

	err = parentDirectory.Create(verbose)
	if err != nil {
		return nil, err
	}

	err = createdFile.Create(verbose)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (l *LocalDirectory) CreateFilesInDirectory(filesToCreate []string, verbose bool) (createdFiles []File, err error) {
	if filesToCreate == nil {
		return nil, TracedErrorNil("filesToCreate")
	}

	createdFiles = []File{}
	for _, fileName := range filesToCreate {
		toAdd, err := l.CreateFileInDirectory(verbose, fileName)
		if err != nil {
			return nil, err
		}

		createdFiles = append(createdFiles, toAdd)
	}

	if verbose {
		dirPath, err := l.GetLocalPath()
		if err != nil {
			return nil, err
		}

		LogInfof("Created '%d' files in directory '%s'.", len(createdFiles), dirPath)
	}

	return createdFiles, nil
}

func (l *LocalDirectory) CreateSubDirectory(subDirName string, verbose bool) (createdSubDir Directory, err error) {
	if subDirName == "" {
		return nil, TracedErrorEmptyString("subDirName")
	}

	subDirectory, subDirectoryPath, err := l.GetSubDirectoryAndLocalPath(subDirName)
	if err != nil {
		return nil, err
	}

	subDirExists, err := subDirectory.Exists(verbose)
	if err != nil {
		return nil, err
	}

	if subDirExists {
		if verbose {
			LogInfof("Sub directory '%s' already exists.", subDirectoryPath)
		}
	} else {
		err = subDirectory.Create(verbose)
		if err != nil {
			return nil, err
		}

		if verbose {
			LogChangedf("Sub directory '%s' already created.", subDirectoryPath)
		}
	}

	return subDirectory, nil
}

func (l *LocalDirectory) Delete(verbose bool) (err error) {
	exists, err := l.Exists(verbose)
	if err != nil {
		return err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if exists {
		err := os.RemoveAll(path)
		if err != nil {
			return TracedErrorf("Delete directory '%s' failed: '%w'", path, err)
		}

		if verbose {
			LogChangedf("Deleted local directory '%s'", path)
		}
	} else {
		if verbose {
			LogInfof("Local directory '%s' already absent. Skip delete.", path)
		}
	}

	return nil
}

func (l *LocalDirectory) Exists(verbose bool) (exists bool, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return false, err
	}

	dirInfo, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, TracedErrorf("Unable to evaluate if local directory exists: '%w'", err)
	}

	exists = dirInfo.IsDir()

	if verbose {
		if exists {
			LogInfof(
				"Local directory '%s' exists.",
				localPath,
			)
		} else {
			LogInfof(
				"Local directory '%s' does not exists.",
				localPath,
			)
		}
	}

	return exists, nil
}

func (l *LocalDirectory) GetBaseName() (baseName string, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(localPath)

	return baseName, nil
}

func (l *LocalDirectory) GetDirName() (parentPath string, err error) {
	path, err := l.GetPath()
	if err != nil {
		return "", err
	}

	return Paths().GetDirPath(path)
}

func (l *LocalDirectory) GetFileInDirectory(path ...string) (file File, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no elements")
	}

	dirPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	filePath := dirPath
	for _, p := range path {
		filePath = filepath.Join(filePath, p)
	}

	file, err = GetLocalFileByPath(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *LocalDirectory) GetFileInDirectoryAsLocalFile(filePath ...string) (localFile *LocalFile, err error) {
	if len(filePath) <= 0 {
		return nil, TracedErrorNil("filePath")
	}

	fileInDir, err := l.GetFileInDirectory(filePath...)
	if err != nil {
		return nil, err
	}

	localFile, ok := fileInDir.(*LocalFile)
	if !ok {
		return nil, TracedError("Internal error: unable get as LocalFile")
	}

	return localFile, nil
}

func (l *LocalDirectory) GetGitRepositories(verbose bool) (gitRepos []GitRepository, err error) {
	localRepositories, err := l.GetGitRepositoriesAsLocalGitRepositories(verbose)
	if err != nil {
		return nil, err
	}

	for _, toAdd := range localRepositories {
		gitRepos = append(gitRepos, toAdd)
	}

	return gitRepos, nil
}

func (l *LocalDirectory) GetGitRepositoriesAsLocalGitRepositories(verbose bool) (gitRepos []*LocalGitRepository, err error) {
	subDirectories, err := l.ListSubDirectories(&ListDirectoryOptions{
		Recursive: true,
	})
	if err != nil {
		return nil, err
	}

	repoPaths := []string{}
	for _, subDir := range subDirectories {
		gitRepo, err := GetLocalGitReposioryFromDirectory(subDir)
		if err != nil {
			return nil, err
		}

		isGitRepo, err := gitRepo.IsGitRepository(verbose)
		if err != nil {
			return nil, err
		}

		if !isGitRepo {
			continue
		}

		rootDirectoryPath, err := gitRepo.GetRootDirectoryPath(verbose)
		if err != nil {
			return nil, err
		}

		if !Slices().ContainsString(repoPaths, rootDirectoryPath) {
			repoPaths = append(repoPaths, rootDirectoryPath)
		}
	}

	for _, toAdd := range repoPaths {
		gitRepo, err := GetLocalGitRepositoryByPath(toAdd)
		if err != nil {
			return nil, err
		}

		gitRepos = append(gitRepos, gitRepo)
	}

	if verbose {
		localPath, err := l.GetLocalPath()
		if err != nil {
			return nil, err
		}

		LogInfof("Found '%d' git repositories in '%s'.", len(gitRepos), localPath)
	}

	return gitRepos, nil
}

func (l *LocalDirectory) GetHostDescription() (hostDescription string, err error) {
	return "localhost", err
}

func (l *LocalDirectory) GetLocalPath() (localPath string, err error) {
	if l.localPath == "" {
		return "", TracedErrorf("localPath not set")
	}

	return l.localPath, nil
}

func (l *LocalDirectory) GetParentDirectory() (parentDirectory Directory, err error) {
	parentPath, err := l.GetDirName()
	if err != nil {
		return nil, err
	}

	parentDirectory, err = GetLocalDirectoryByPath(parentPath)
	if err != nil {
		return nil, err
	}

	return parentDirectory, err
}

func (l *LocalDirectory) GetPath() (dirPath string, err error) {
	dirPath, err = l.GetLocalPath()
	if err != nil {
		return "", err
	}

	return dirPath, nil
}

func (l *LocalDirectory) GetSubDirectory(path ...string) (subDirectory Directory, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no elements")
	}

	dirPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	directoryPath := dirPath
	for _, p := range path {
		directoryPath = filepath.Join(directoryPath, p)
	}

	subDirectory, err = GetLocalDirectoryByPath(directoryPath)
	if err != nil {
		return nil, err
	}

	return subDirectory, nil
}

func (l *LocalDirectory) GetSubDirectoryAndLocalPath(path ...string) (subDirectory Directory, subDirectoryPath string, err error) {
	subDirectory, err = l.GetSubDirectory(path...)
	if err != nil {
		return nil, "", err
	}

	subDirectoryPath, err = subDirectory.GetLocalPath()
	if err != nil {
		return nil, "", err
	}

	return subDirectory, subDirectoryPath, nil
}

func (l *LocalDirectory) IsEmptyDirectory(verbose bool) (isEmpty bool, err error) {
	subDirs, err := l.ListSubDirectories(
		&ListDirectoryOptions{
			Verbose: verbose,
		},
	)
	if err != nil {
		return false, err
	}

	if len(subDirs) > 0 {
		return false, nil
	}

	files, err := l.ListFiles(
		&ListFileOptions{
			Verbose:                       verbose,
			AllowEmptyListIfNoFileIsFound: true,
		},
	)
	if err != nil {
		return false, err
	}

	if len(files) > 0 {
		return false, nil
	}

	return true, nil
}

func (l *LocalDirectory) IsLocalDirectory() (isLocalDirectory bool, err error) {
	return true, nil
}

func (l *LocalDirectory) ListFilePaths(listOptions *ListFileOptions) (filePathList []string, err error) {
	if listOptions == nil {
		return nil, TracedError("listOptions is nil")
	}

	listOptions = listOptions.GetDeepCopy()
	listOptions.OnlyFiles = true

	directoryPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	filePathList = []string{}
	err = filepath.Walk(
		directoryPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			filePathList = append(filePathList, path)
			return nil
		},
	)
	if err != nil {
		return nil, TracedErrorf("Unable to filepath.Walk: '%w'", err)
	}

	filePathList = Slices().RemoveEmptyStrings(filePathList)

	filePathList, err = Paths().FilterPaths(filePathList, listOptions)
	if err != nil {
		return nil, err
	}

	if listOptions.ReturnRelativePaths {
		filePathList, err = Paths().GetRelativePathsTo(filePathList, directoryPath)
		if err != nil {
			return nil, err
		}
	}

	filePathList = Slices().SortStringSliceAndRemoveEmpty(filePathList)

	if len(filePathList) <= 0 {
		if !listOptions.AllowEmptyListIfNoFileIsFound {
			return nil, TracedErrorf("No files in '%s' found", directoryPath)
		}
	}

	return filePathList, nil
}

func (l *LocalDirectory) ListFiles(options *ListFileOptions) (files []File, err error) {
	if options == nil {
		return nil, TracedError("options is nil")
	}

	optionsToUse := options.GetDeepCopy()
	optionsToUse.ReturnRelativePaths = true

	filePathList, err := l.ListFilePaths(optionsToUse)
	if err != nil {
		return nil, err
	}

	files = []File{}
	for _, name := range filePathList {
		fileToAdd, err := l.GetFileInDirectory(name)
		if err != nil {
			return nil, err
		}

		files = append(files, fileToAdd)
	}

	return files, nil
}

func (l *LocalDirectory) ListSubDirectories(listDirectoryOptions *ListDirectoryOptions) (subDirectories []Directory, err error) {
	if listDirectoryOptions == nil {
		return nil, TracedErrorNil("listDirectoryOptions")
	}

	pathsToAdd, err := l.ListSubDirectoriesAsAbsolutePaths(listDirectoryOptions)
	if err != nil {
		return nil, err
	}

	for _, pathToAdd := range pathsToAdd {
		toAdd, err := GetLocalDirectoryByPath(pathToAdd)
		if err != nil {
			return nil, err
		}

		subDirectories = append(subDirectories, toAdd)
	}

	return subDirectories, nil
}

func (l *LocalDirectory) ListSubDirectoriesAsAbsolutePaths(listDirectoryOptions *ListDirectoryOptions) (subDirectoryPaths []string, err error) {
	if listDirectoryOptions == nil {
		return nil, TracedErrorNil("listDirectoryOptions")
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	allEntries, err := os.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range allEntries {
		if entry.IsDir() {
			pathToAdd := filepath.Join(localPath, entry.Name())
			if err != nil {
				return nil, err
			}

			subDirectoryPaths = append(subDirectoryPaths, pathToAdd)

			if listDirectoryOptions.Recursive {
				subDirectory, err := GetLocalDirectoryByPath(pathToAdd)
				if err != nil {
					return nil, err
				}

				subDirectoriesToAdd, err := subDirectory.ListSubDirectoriesAsAbsolutePaths(
					&ListDirectoryOptions{
						Recursive: true,
					},
				)
				if err != nil {
					return nil, err
				}

				subDirectoryPaths = append(subDirectoryPaths, subDirectoriesToAdd...)
			}
		}
	}

	subDirectoryPaths = Slices().SortStringSlice(subDirectoryPaths)

	return subDirectoryPaths, nil
}

func (l *LocalDirectory) MustChmod(chmodOptions *ChmodOptions) {
	err := l.Chmod(chmodOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCopyContentToDirectory(destinationDir Directory, verbose bool) {
	err := l.CopyContentToDirectory(destinationDir, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCopyContentToLocalDirectory(destDirectory *LocalDirectory, verbose bool) {
	err := l.CopyContentToLocalDirectory(destDirectory, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCopyFileToTemporaryFile(verbose bool, filePath ...string) (copy File) {
	copy, err := l.CopyFileToTemporaryFile(verbose, filePath...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return copy
}

func (l *LocalDirectory) MustCopyFileToTemporaryFileAsLocalFile(verbose bool, filePath ...string) (copy *LocalFile) {
	copy, err := l.CopyFileToTemporaryFileAsLocalFile(verbose, filePath...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return copy
}

func (l *LocalDirectory) MustCreate(verbose bool) {
	err := l.Create(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCreateFileInDirectory(verbose bool, path ...string) (createdFile File) {
	createdFile, err := l.CreateFileInDirectory(verbose, path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFile
}

func (l *LocalDirectory) MustCreateFilesInDirectory(filesToCreate []string, verbose bool) (createdFiles []File) {
	createdFiles, err := l.CreateFilesInDirectory(filesToCreate, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFiles
}

func (l *LocalDirectory) MustCreateSubDirectory(subDirName string, verbose bool) (createdSubDir Directory) {
	createdSubDir, err := l.CreateSubDirectory(subDirName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdSubDir
}

func (l *LocalDirectory) MustDelete(verbose bool) {
	err := l.Delete(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustExists(verbose bool) (exists bool) {
	exists, err := l.Exists(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (l *LocalDirectory) MustGetBaseName() (baseName string) {
	baseName, err := l.GetBaseName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return baseName
}

func (l *LocalDirectory) MustGetDirName() (dirName string) {
	dirName, err := l.GetDirName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dirName
}

func (l *LocalDirectory) MustGetFileInDirectory(path ...string) (file File) {
	file, err := l.GetFileInDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return file
}

func (l *LocalDirectory) MustGetFileInDirectoryAsLocalFile(filePath ...string) (localFile *LocalFile) {
	localFile, err := l.GetFileInDirectoryAsLocalFile(filePath...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localFile
}

func (l *LocalDirectory) MustGetGitRepositories(verbose bool) (gitRepos []GitRepository) {
	gitRepos, err := l.GetGitRepositories(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitRepos
}

func (l *LocalDirectory) MustGetGitRepositoriesAsLocalGitRepositories(verbose bool) (gitRepos []*LocalGitRepository) {
	gitRepos, err := l.GetGitRepositoriesAsLocalGitRepositories(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitRepos
}

func (l *LocalDirectory) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := l.GetHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hostDescription
}

func (l *LocalDirectory) MustGetLocalPath() (localPath string) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localPath
}

func (l *LocalDirectory) MustGetParentDirectory() (parentDirectory Directory) {
	parentDirectory, err := l.GetParentDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentDirectory
}

func (l *LocalDirectory) MustGetPath() (dirPath string) {
	dirPath, err := l.GetPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return dirPath
}

func (l *LocalDirectory) MustGetSubDirectory(path ...string) (subDirectory Directory) {
	subDirectory, err := l.GetSubDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectory
}

func (l *LocalDirectory) MustGetSubDirectoryAndLocalPath(path ...string) (subDirectory Directory, subDirectoryPath string) {
	subDirectory, subDirectoryPath, err := l.GetSubDirectoryAndLocalPath(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectory, subDirectoryPath
}

func (l *LocalDirectory) MustIsEmptyDirectory(verbose bool) (isEmpty bool) {
	isEmpty, err := l.IsEmptyDirectory(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isEmpty
}

func (l *LocalDirectory) MustIsLocalDirectory() (isLocalDirectory bool) {
	isLocalDirectory, err := l.IsLocalDirectory()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isLocalDirectory
}

func (l *LocalDirectory) MustListFilePaths(listOptions *ListFileOptions) (filePathList []string) {
	filePathList, err := l.ListFilePaths(listOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filePathList
}

func (l *LocalDirectory) MustListFilePathsns(listOptions *ListFileOptions) (filePathList []string) {
	filePathList, err := l.ListFilePaths(listOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return filePathList
}

func (l *LocalDirectory) MustListFiles(options *ListFileOptions) (files []File) {
	files, err := l.ListFiles(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return files
}

func (l *LocalDirectory) MustListSubDirectories(listDirectoryOptions *ListDirectoryOptions) (subDirectories []Directory) {
	subDirectories, err := l.ListSubDirectories(listDirectoryOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectories
}

func (l *LocalDirectory) MustListSubDirectoriesAsAbsolutePaths(listDirectoryOptions *ListDirectoryOptions) (subDirectoryPaths []string) {
	subDirectoryPaths, err := l.ListSubDirectoriesAsAbsolutePaths(listDirectoryOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectoryPaths
}

func (l *LocalDirectory) MustReplaceBetweenMarkers(verbose bool) {
	err := l.ReplaceBetweenMarkers(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustSetLocalPath(localPath string) {
	err := l.SetLocalPath(localPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustSubDirectoryExists(subDirName string, verbose bool) (subDirExists bool) {
	subDirExists, err := l.SubDirectoryExists(subDirName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirExists
}

func (l *LocalDirectory) ReplaceBetweenMarkers(verbose bool) (err error) {
	files, err := l.ListFiles(
		&ListFileOptions{
			Verbose: verbose,
		},
	)
	if err != nil {
		return err
	}

	for _, f := range files {
		err = f.ReplaceBetweenMarkers(verbose)
		if err != nil {
			return err
		}
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if verbose {
		LogInfof(
			"Replaces between markers in '%d' files in '%s'.",
			len(files),
			path,
		)
	}

	return nil
}

func (l *LocalDirectory) SetLocalPath(localPath string) (err error) {
	if localPath == "" {
		return TracedErrorf("localPath is empty string")
	}

	localPath, err = Paths().GetAbsolutePath(localPath)
	if err != nil {
		return err
	}

	if !Paths().IsAbsolutePath(localPath) {
		return TracedErrorf(
			"Path '%s' is not absolute. Beware this is an internal issue since the code before this line should fix that.",
			localPath,
		)
	}

	l.localPath = localPath

	return nil
}

func (l *LocalDirectory) SubDirectoryExists(subDirName string, verbose bool) (subDirExists bool, err error) {
	if subDirName == "" {
		return false, TracedErrorEmptyString("subDirName")
	}

	subDir, err := l.GetSubDirectory(subDirName)
	if err != nil {
		return false, err
	}

	subDirExists, err = subDir.Exists(verbose)
	if err != nil {
		return false, err
	}

	return subDirExists, nil
}
