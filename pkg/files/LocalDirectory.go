package files

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type LocalDirectory struct {
	DirectoryBase
	localPath string
}

func GetLocalDirectoryByPath(path string) (l *LocalDirectory, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
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
		logging.LogGoErrorFatal(err)
	}

	return l
}

func NewLocalDirectory() (l *LocalDirectory) {
	l = new(LocalDirectory)

	// Allow usage of the base class functions:
	l.MustSetParentDirectoryForBaseClass(l)

	return l
}

func (l *LocalDirectory) Chmod(chmodOptions *parameteroptions.ChmodOptions) (err error) {
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

	_, err = commandexecutorbashoo.Bash().RunCommand(
		contextutils.GetVerbosityContextByBool(chmodOptions.Verbose),
		&parameteroptions.RunCommandOptions{
			Command: []string{"chmod", chmodString, localPath},
		},
	)
	if err != nil {
		return err
	}

	if chmodOptions.Verbose {
		logging.LogChangedf("Chmod '%s' for local directory '%s'.", chmodString, localPath)
	}

	return nil
}

func (l *LocalDirectory) CopyContentToDirectory(destinationDir filesinterfaces.Directory, verbose bool) (err error) {
	if destinationDir == nil {
		return tracederrors.TracedError("destinationDir is empty string")
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

	ctx := contextutils.GetVerbosityContextByBool(verbose)
	stdout, err := commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: copyCommand,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
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
		return tracederrors.TracedErrorNil("destDirectory")
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
		return tracederrors.TracedErrorf("Unable to copy content to local directory, '%s' does not exist.", srcPath)
	}

	ctx := contextutils.GetVerbosityContextByBool(verbose)
	_, err = commandexecutorbashoo.Bash().RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"cp", "-r", "-v", srcPath + "/.", destPath + "/."},
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf("Copied files from '%s' to '%s'.", srcPath, destPath)
	}

	return nil
}

/* TODO move or remove
func (l *LocalDirectory) CopyFileToTemporaryFile(verbose bool, filePath ...string) (copy File, err error) {
	if len(filePath) <= 0 {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
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
*/

/* TODO move or remove
func (l *LocalDirectory) CopyFileToTemporaryFileAsLocalFile(verbose bool, filePath ...string) (copy *LocalFile, err error) {
	if len(filePath) <= 0 {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	interfaceCopy, err := l.CopyFileToTemporaryFile(verbose, filePath...)
	if err != nil {
		return nil, err
	}

	copy, ok := interfaceCopy.(*LocalFile)
	if !ok {
		return nil, tracederrors.TracedErrorf("Internal error: Unable to convert to *LocalFile: '%v'", interfaceCopy)
	}

	return copy, nil
}
*/

func (l *LocalDirectory) Create(ctx context.Context, options *filesoptions.CreateOptions) (err error) {
	exists, err := l.Exists(contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	path, err := l.GetLocalPath()
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Local directory '%s' already exists. Skip create.", path)
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				parentDirectoy, err := l.GetParentDirectory()
				if err != nil {
					return err
				}

				err = parentDirectoy.Create(ctx, options)
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

		existsAfterCreate, err := l.Exists(contextutils.GetVerboseFromContext(ctx))
		if err != nil {
			return err
		}

		if !existsAfterCreate {
			return tracederrors.TracedErrorf("Local directory '%s' does not exist after creation.", path)
		}

		logging.LogChangedByCtxf(ctx, "Created local directory '%s'", path)
	}

	return nil
}

func (l *LocalDirectory) CreateFileInDirectory(ctx context.Context, path string, options *filesoptions.CreateOptions) (createdFile filesinterfaces.File, err error) {
	createdFile, err = l.GetFileInDirectory(path)
	if err != nil {
		return nil, err
	}

	parentDirectory, err := createdFile.GetParentDirectory()
	if err != nil {
		return nil, err
	}

	err = parentDirectory.Create(ctx, options)
	if err != nil {
		return nil, err
	}

	err = createdFile.Create(ctx, options)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (l *LocalDirectory) CreateFilesInDirectory(ctx context.Context, filesToCreate []string, options *filesoptions.CreateOptions) (createdFiles []filesinterfaces.File, err error) {
	if filesToCreate == nil {
		return nil, tracederrors.TracedErrorNil("filesToCreate")
	}

	createdFiles = []filesinterfaces.File{}
	for _, fileName := range filesToCreate {
		toAdd, err := l.CreateFileInDirectory(ctx, fileName, options)
		if err != nil {
			return nil, err
		}

		createdFiles = append(createdFiles, toAdd)
	}

	dirPath, err := l.GetLocalPath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Created '%d' files in directory '%s'.", len(createdFiles), dirPath)

	return createdFiles, nil
}

func (l *LocalDirectory) CreateSubDirectory(ctx context.Context, subDirName string, options *filesoptions.CreateOptions) (createdSubDir filesinterfaces.Directory, err error) {
	if subDirName == "" {
		return nil, tracederrors.TracedErrorEmptyString("subDirName")
	}

	subDirectory, subDirectoryPath, err := l.GetSubDirectoryAndLocalPath(subDirName)
	if err != nil {
		return nil, err
	}

	subDirExists, err := subDirectory.Exists(contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return nil, err
	}

	if subDirExists {
		logging.LogInfoByCtxf(ctx, "Sub directory '%s' already exists.", subDirectoryPath)
	} else {
		err = subDirectory.Create(ctx, options)
		if err != nil {
			return nil, err
		}
		logging.LogChangedByCtxf(ctx, "Sub directory '%s' already created.", subDirectoryPath)
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
			return tracederrors.TracedErrorf("Delete directory '%s' failed: '%w'", path, err)
		}

		if verbose {
			logging.LogChangedf("Deleted local directory '%s'", path)
		}
	} else {
		if verbose {
			logging.LogInfof("Local directory '%s' already absent. Skip delete.", path)
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

		return false, tracederrors.TracedErrorf("Unable to evaluate if local directory exists: '%w'", err)
	}

	exists = dirInfo.IsDir()

	if verbose {
		if exists {
			logging.LogInfof(
				"Local directory '%s' exists.",
				localPath,
			)
		} else {
			logging.LogInfof(
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

	return pathsutils.GetDirPath(path)
}

func (l *LocalDirectory) GetFileInDirectory(path ...string) (file filesinterfaces.File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no elements")
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
		return nil, tracederrors.TracedErrorNil("filePath")
	}

	fileInDir, err := l.GetFileInDirectory(filePath...)
	if err != nil {
		return nil, err
	}

	localFile, ok := fileInDir.(*LocalFile)
	if !ok {
		return nil, tracederrors.TracedError("Internal error: unable get as LocalFile")
	}

	return localFile, nil
}

/* TODO remove or move
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
*/

/* TODO remove or move
func (l *LocalDirectory) GetGitRepositoriesAsLocalGitRepositories(verbose bool) (gitRepos []*LocalGitRepository, err error) {
	subDirectories, err := l.ListSubDirectories(&parameteroptions.ListDirectoryOptions{
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

		if !slicesutils.ContainsString(repoPaths, rootDirectoryPath) {
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

		logging.LogInfof("Found '%d' git repositories in '%s'.", len(gitRepos), localPath)
	}

	return gitRepos, nil
}
*/

func (l *LocalDirectory) GetHostDescription() (hostDescription string, err error) {
	return "localhost", err
}

func (l *LocalDirectory) GetLocalPath() (localPath string, err error) {
	if l.localPath == "" {
		return "", tracederrors.TracedErrorf("localPath not set")
	}

	return l.localPath, nil
}

func (l *LocalDirectory) GetParentDirectory() (parentDirectory filesinterfaces.Directory, err error) {
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

func (l *LocalDirectory) GetSubDirectory(path ...string) (subDirectory filesinterfaces.Directory, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no elements")
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

func (l *LocalDirectory) GetSubDirectoryAndLocalPath(path ...string) (subDirectory filesinterfaces.Directory, subDirectoryPath string, err error) {
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
		&parameteroptions.ListDirectoryOptions{
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
		contextutils.GetVerbosityContextByBool(verbose),
		&parameteroptions.ListFileOptions{
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

func (l *LocalDirectory) ListFilePaths(ctx context.Context, listOptions *parameteroptions.ListFileOptions) (filePathList []string, err error) {
	if listOptions == nil {
		return nil, tracederrors.TracedError("listOptions is nil")
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
		return nil, tracederrors.TracedErrorf("Unable to filepath.Walk: '%w'", err)
	}

	filePathList = slicesutils.RemoveEmptyStrings(filePathList)

	filePathList, err = pathsutils.FilterPaths(filePathList, listOptions)
	if err != nil {
		return nil, err
	}

	if listOptions.ReturnRelativePaths {
		filePathList, err = pathsutils.GetRelativePathsTo(filePathList, directoryPath)
		if err != nil {
			return nil, err
		}
	}

	filePathList = slicesutils.SortStringSliceAndRemoveEmpty(filePathList)

	if len(filePathList) <= 0 {
		if !listOptions.AllowEmptyListIfNoFileIsFound {
			return nil, tracederrors.TracedErrorf("No files in '%s' found", directoryPath)
		}
	}

	return filePathList, nil
}

func (l *LocalDirectory) ListFiles(ctx context.Context, options *parameteroptions.ListFileOptions) (files []filesinterfaces.File, err error) {
	if options == nil {
		return nil, tracederrors.TracedError("options is nil")
	}

	optionsToUse := options.GetDeepCopy()
	optionsToUse.ReturnRelativePaths = true

	filePathList, err := l.ListFilePaths(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	files = []filesinterfaces.File{}
	for _, name := range filePathList {
		fileToAdd, err := l.GetFileInDirectory(name)
		if err != nil {
			return nil, err
		}

		files = append(files, fileToAdd)
	}

	return files, nil
}

func (l *LocalDirectory) ListSubDirectories(listDirectoryOptions *parameteroptions.ListDirectoryOptions) (subDirectories []filesinterfaces.Directory, err error) {
	if listDirectoryOptions == nil {
		return nil, tracederrors.TracedErrorNil("listDirectoryOptions")
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

func (l *LocalDirectory) ListSubDirectoriesAsAbsolutePaths(listDirectoryOptions *parameteroptions.ListDirectoryOptions) (subDirectoryPaths []string, err error) {
	if listDirectoryOptions == nil {
		return nil, tracederrors.TracedErrorNil("listDirectoryOptions")
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

			subDirectoryPaths = append(subDirectoryPaths, pathToAdd)

			if listDirectoryOptions.Recursive {
				subDirectory, err := GetLocalDirectoryByPath(pathToAdd)
				if err != nil {
					return nil, err
				}

				subDirectoriesToAdd, err := subDirectory.ListSubDirectoriesAsAbsolutePaths(
					&parameteroptions.ListDirectoryOptions{
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

	sort.Strings(subDirectoryPaths)

	return subDirectoryPaths, nil
}

func (l *LocalDirectory) MustChmod(chmodOptions *parameteroptions.ChmodOptions) {
	err := l.Chmod(chmodOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCopyContentToDirectory(destinationDir filesinterfaces.Directory, verbose bool) {
	err := l.CopyContentToDirectory(destinationDir, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustCopyContentToLocalDirectory(destDirectory *LocalDirectory, verbose bool) {
	err := l.CopyContentToLocalDirectory(destDirectory, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustGetBaseName() (baseName string) {
	baseName, err := l.GetBaseName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return baseName
}

func (l *LocalDirectory) MustGetDirName() (dirName string) {
	dirName, err := l.GetDirName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dirName
}

func (l *LocalDirectory) MustGetFileInDirectory(path ...string) (file filesinterfaces.File) {
	file, err := l.GetFileInDirectory(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return file
}

func (l *LocalDirectory) MustGetFileInDirectoryAsLocalFile(filePath ...string) (localFile *LocalFile) {
	localFile, err := l.GetFileInDirectoryAsLocalFile(filePath...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localFile
}

/* TODO remove or move
func (l *LocalDirectory) MustGetGitRepositories(verbose bool) (gitRepos []GitRepository) {
	gitRepos, err := l.GetGitRepositories(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepos
}
*/

/* TODO remove or move
func (l *LocalDirectory) MustGetGitRepositoriesAsLocalGitRepositories(verbose bool) (gitRepos []*LocalGitRepository) {
	gitRepos, err := l.GetGitRepositoriesAsLocalGitRepositories(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepos
}
*/

func (l *LocalDirectory) MustGetHostDescription() (hostDescription string) {
	hostDescription, err := l.GetHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostDescription
}

func (l *LocalDirectory) MustGetLocalPath() (localPath string) {
	localPath, err := l.GetLocalPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return localPath
}

func (l *LocalDirectory) MustGetParentDirectory() (parentDirectory filesinterfaces.Directory) {
	parentDirectory, err := l.GetParentDirectory()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentDirectory
}

func (l *LocalDirectory) MustGetPath() (dirPath string) {
	dirPath, err := l.GetPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dirPath
}

func (l *LocalDirectory) MustGetSubDirectory(path ...string) (subDirectory filesinterfaces.Directory) {
	subDirectory, err := l.GetSubDirectory(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subDirectory
}

func (l *LocalDirectory) MustGetSubDirectoryAndLocalPath(path ...string) (subDirectory filesinterfaces.Directory, subDirectoryPath string) {
	subDirectory, subDirectoryPath, err := l.GetSubDirectoryAndLocalPath(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subDirectory, subDirectoryPath
}

func (l *LocalDirectory) MustIsEmptyDirectory(verbose bool) (isEmpty bool) {
	isEmpty, err := l.IsEmptyDirectory(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isEmpty
}

func (l *LocalDirectory) MustIsLocalDirectory() (isLocalDirectory bool) {
	isLocalDirectory, err := l.IsLocalDirectory()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isLocalDirectory
}

func (l *LocalDirectory) MustListSubDirectories(listDirectoryOptions *parameteroptions.ListDirectoryOptions) (subDirectories []filesinterfaces.Directory) {
	subDirectories, err := l.ListSubDirectories(listDirectoryOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subDirectories
}

func (l *LocalDirectory) MustListSubDirectoriesAsAbsolutePaths(listDirectoryOptions *parameteroptions.ListDirectoryOptions) (subDirectoryPaths []string) {
	subDirectoryPaths, err := l.ListSubDirectoriesAsAbsolutePaths(listDirectoryOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subDirectoryPaths
}

func (l *LocalDirectory) MustSetLocalPath(localPath string) {
	err := l.SetLocalPath(localPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (l *LocalDirectory) MustSubDirectoryExists(subDirName string, verbose bool) (subDirExists bool) {
	subDirExists, err := l.SubDirectoryExists(subDirName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return subDirExists
}

/* TODO remove or move
func (l *LocalDirectory) ReplaceBetweenMarkers(verbose bool) (err error) {
	files, err := l.ListFiles(
		&parameteroptions.ListFileOptions{
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
		logging.LogInfof(
			"Replaces between markers in '%d' files in '%s'.",
			len(files),
			path,
		)
	}

	return nil
}
*/

func (l *LocalDirectory) SetLocalPath(localPath string) (err error) {
	if localPath == "" {
		return tracederrors.TracedErrorf("localPath is empty string")
	}

	localPath, err = pathsutils.GetAbsolutePath(localPath)
	if err != nil {
		return err
	}

	if !pathsutils.IsAbsolutePath(localPath) {
		return tracederrors.TracedErrorf(
			"Path '%s' is not absolute. Beware this is an internal issue since the code before this line should fix that.",
			localPath,
		)
	}

	l.localPath = localPath

	return nil
}

func (l *LocalDirectory) SubDirectoryExists(subDirName string, verbose bool) (subDirExists bool, err error) {
	if subDirName == "" {
		return false, tracederrors.TracedErrorEmptyString("subDirName")
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
