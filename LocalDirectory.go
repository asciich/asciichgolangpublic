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

func (l *LocalDirectory) Create(verbose bool) (err error) {
	exists, err := l.Exists()
	if err != nil {
		return nil
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil
	}

	if exists {
		if verbose {
			LogInfof("Local directory '%s' already exists. Skip create.", path)
		}
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return TracedErrorf("Create directory '%s' failed: '%w'", path, err)
		}

		if verbose {
			LogChangedf("Created local directory '%s'", path)
		}
	}

	return nil
}

func (l *LocalDirectory) CreateFileInDirectory(path ...string) (createdFile File, err error) {
	createdFile, err = l.GetFileInDirectory(path...)
	if err != nil {
		return nil, err
	}

	const verbose = false
	err = createdFile.Create(verbose)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (l *LocalDirectory) CreateSubDirectory(subDirName string, verbose bool) (createdSubDir Directory, err error) {
	if subDirName == "" {
		return nil, TracedErrorEmptyString("subDirName")
	}

	subDirectory, subDirectoryPath, err := l.GetSubDirectoryAndLocalPath(subDirName)
	if err != nil {
		return nil, err
	}

	subDirExists, err := subDirectory.Exists()
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
	exists, err := l.Exists()
	if err != nil {
		return nil
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return nil
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

func (l *LocalDirectory) Exists() (exists bool, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return false, nil
	}

	dirInfo, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, TracedErrorf("Unable to evaluate if local directory exists: '%w'", err)
	}

	return dirInfo.IsDir(), nil
}

func (l *LocalDirectory) GetBaseName() (baseName string, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(localPath)

	return baseName, nil
}

func (l *LocalDirectory) GetDirName() (dirName string, err error) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		return "", err
	}

	dirName = filepath.Dir(localPath)

	return dirName, nil
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
	subDirectories, err := l.GetSubDirectories(&ListDirectoryOptions{
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

		isGitRepo, err := gitRepo.IsGitRepository()
		if err != nil {
			return nil, err
		}

		if !isGitRepo {
			continue
		}

		rootDirectoryPath, err := gitRepo.GetRootDirectoryPath()
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

func (l *LocalDirectory) GetLocalPath() (localPath string, err error) {
	if l.localPath == "" {
		return "", TracedErrorf("localPath not set")
	}

	return l.localPath, nil
}

func (l *LocalDirectory) GetSubDirectories(listDirectoryOptions *ListDirectoryOptions) (subDirectories []Directory, err error) {
	if listDirectoryOptions == nil {
		return nil, TracedErrorNil("listDirectoryOptions")
	}

	pathsToAdd, err := l.GetSubDirectoriesAsAbsolutePaths(listDirectoryOptions)
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

func (l *LocalDirectory) GetSubDirectoriesAsAbsolutePaths(listDirectoryOptions *ListDirectoryOptions) (subDirectoryPaths []string, err error) {
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

				subDirectoriesToAdd, err := subDirectory.GetSubDirectoriesAsAbsolutePaths(
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

func (l *LocalDirectory) GetSubDirectoryPaths(listOptions *ListDirectoryOptions) (paths []string, err error) {
	if listOptions == nil {
		return nil, TracedErrorNil("listOptions")
	}

	optionsToUse := &ListDirectoryOptions{
		Recursive: listOptions.Recursive,
		Verbose:   listOptions.Verbose,
	}

	absoultePaths, err := l.GetSubDirectoriesAsAbsolutePaths(optionsToUse)
	if err != nil {
		return nil, err
	}

	localPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	for _, absolutePath := range absoultePaths {
		toAdd := strings.TrimPrefix(absolutePath, localPath)
		toAdd = Strings().TrimAllPrefix(toAdd, "/")

		if toAdd == "" {
			return nil, TracedError("Internal error: toAdd is empty string after evaluation")
		}

		paths = append(paths, toAdd)
	}

	return paths, err
}

func (l *LocalDirectory) IsLocalDirectory() (isLocalDirectory bool) {
	return true
}

func (l *LocalDirectory) MustCreate(verbose bool) {
	err := l.Create(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCreateFileInDirectory(path ...string) (createdFile File) {
	createdFile, err := l.CreateFileInDirectory(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdFile
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

func (l *LocalDirectory) MustExists() (exists bool) {
	exists, err := l.Exists()
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

func (l *LocalDirectory) MustGetLocalPath() (localPath string) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return localPath
}

func (l *LocalDirectory) MustGetSubDirectories(listDirectoryOptions *ListDirectoryOptions) (subDirectories []Directory) {
	subDirectories, err := l.GetSubDirectories(listDirectoryOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectories
}

func (l *LocalDirectory) MustGetSubDirectoriesAsAbsolutePaths(listDirectoryOptions *ListDirectoryOptions) (subDirectoryPaths []string) {
	subDirectoryPaths, err := l.GetSubDirectoriesAsAbsolutePaths(listDirectoryOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return subDirectoryPaths
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

func (l *LocalDirectory) MustGetSubDirectoryPaths(listOptions *ListDirectoryOptions) (paths []string) {
	paths, err := l.GetSubDirectoryPaths(listOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return paths
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

func (l *LocalDirectory) MustWriteStringToFileInDirectory(content string, verbose bool, filePath ...string) (writtenFile *LocalFile) {
	writtenFile, err := l.WriteStringToFileInDirectory(content, verbose, filePath...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return writtenFile
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

	subDirExists, err = subDir.Exists()
	if err != nil {
		return false, err
	}

	return subDirExists, nil
}

func (l *LocalDirectory) WriteStringToFileInDirectory(content string, verbose bool, filePath ...string) (writtenFile *LocalFile, err error) {
	if len(filePath) <= 0 {
		return nil, TracedErrorNil("filePath")
	}

	writtenFile, err = l.GetFileInDirectoryAsLocalFile(filePath...)
	if err != nil {
		return nil, err
	}

	err = writtenFile.WriteString(content, verbose)
	if err != nil {
		return nil, err
	}

	return writtenFile, nil
}
