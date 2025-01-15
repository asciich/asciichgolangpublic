package asciichgolangpublic

import (
	"slices"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pathsutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitRepositoryBase struct {
	parentRepositoryForBaseClass GitRepository
}

func NewGitRepositoryBase() (g *GitRepositoryBase) {
	return new(GitRepositoryBase)
}

func (g *GitRepositoryBase) AddFilesByPath(pathsToAdd []string, verbose bool) (err error) {
	if len(pathsToAdd) <= 0 {
		return tracederrors.TracedError("pathToAdd has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	for _, p := range pathsToAdd {
		err = parent.AddFileByPath(p, verbose)
		if err != nil {
			return err
		}
	}

	if verbose {
		path, hostDescription, err := parent.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Added '%d' files to git repository '%s' on host '%s'.",
			len(pathsToAdd),
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GitRepositoryBase) BranchByNameExists(branchName string, verbose bool) (branchExists bool, err error) {
	if branchName == "" {
		return false, tracederrors.TracedErrorEmptyString("branchName")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	branchNames, err := parent.ListBranchNames(false)
	if err != nil {
		return false, err
	}

	branchExists = slices.Contains(branchNames, branchName)

	if verbose {
		path, hostDescription, err := parent.GetPathAndHostDescription()
		if err != nil {
			return false, err
		}

		if branchExists {
			logging.LogInfof(
				"Branch '%s' in git repository '%s' on host '%s' exists.",
				branchName,
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"Branch '%s' in git repository '%s' on host '%s' does not exist.",
				branchName,
				path,
				hostDescription,
			)
		}
	}

	return branchExists, nil
}

func (g *GitRepositoryBase) CheckHasNoUncommittedChanges(verbose bool) (err error) {
	hasNoUncommitedChanges, err := g.HasNoUncommittedChanges(verbose)
	if err != nil {
		return err
	}

	if !hasNoUncommitedChanges {
		parent, err := g.GetParentRepositoryForBaseClass()
		if err != nil {
			return err
		}

		path, hostDescription, err := parent.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		return tracederrors.TracedErrorf(
			"There are uncommited changes in git repository '%s' on host '%s'",
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GitRepositoryBase) CheckIsGolangApplication(verbose bool) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	ok, err := parent.IsGolangApplication(verbose)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	return tracederrors.TracedErrorf(
		"git repository '%s' on host '%s' is not a golang application",
		path,
		hostDescription,
	)
}

func (g *GitRepositoryBase) CheckIsGolangPackage(verbose bool) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	isGolangPackage, err := parent.IsGolangPackage(verbose)
	if err != nil {
		return err
	}

	if !isGolangPackage {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		return tracederrors.TracedErrorf(
			"git repository '%s' on host '%s' is not a golang package",
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GitRepositoryBase) CheckIsOnLocalhost(verbose bool) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	isOnLocalhost, err := parent.IsOnLocalhost(verbose)
	if err != nil {
		return err
	}

	if !isOnLocalhost {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return err
		}

		return tracederrors.TracedErrorf(
			"git repository '%s' is not on localhost. Host is '%s'",
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GitRepositoryBase) CheckIsPreCommitRepository(verbose bool) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	isPreCommitRepository, err := parent.IsPreCommitRepository(verbose)
	if err != nil {
		return err
	}

	path, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if !isPreCommitRepository {
		return tracederrors.TracedErrorf(
			"Repository '%s' on host '%s' is not a pre-commit repository.",
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GitRepositoryBase) CommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	createdCommit, err = parent.Commit(commitOptions)
	if err != nil {
		return nil, err
	}

	err = parent.Push(commitOptions.Verbose)
	if err != nil {
		return nil, err
	}

	return createdCommit, nil
}

func (g *GitRepositoryBase) CommitIfUncommittedChanges(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	hasUncommitedChanges, err := parent.HasUncommittedChanges(commitOptions.Verbose)
	if err != nil {
		return nil, err
	}

	path, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	if hasUncommitedChanges {
		optionsToUse := commitOptions.GetDeepCopy()
		optionsToUse.CommitAllChanges = true

		createdCommit, err = parent.Commit(optionsToUse)
		if err != nil {
			return nil, err
		}

		createdHash, err := createdCommit.GetHash()
		if err != nil {
			return nil, err
		}

		if commitOptions.Verbose {
			logging.LogInfof(
				"Commited all uncommited changes in git repository '%s' on host '%s' as commit '%s'.",
				path,
				hostDescription,
				createdHash,
			)
		}
	} else {
		createdCommit, err = parent.GetCurrentCommit(false)
		if err != nil {
			return nil, err
		}

		if commitOptions.Verbose {
			logging.LogInfof(
				"No uncommited changes to commit in git repository '%s' on host '%s'.",
				path,
				hostDescription,
			)
		}
	}

	return createdCommit, nil
}

func (g *GitRepositoryBase) ContainsGoSourceFileOfMainPackageWithMainFunction(verbose bool) (mainFound bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	goFiles, err := parent.ListFiles(
		&parameteroptions.ListFileOptions{
			NonRecursive:                  true,
			MatchBasenamePattern:          []string{".*.go"},
			AllowEmptyListIfNoFileIsFound: true,
		},
	)
	if err != nil {
		return false, err
	}

	path, err := parent.GetPath()
	if err != nil {
		return false, err
	}

	const packageMainString string = "package main"
	const funcMainString string = "func main() {"

	for _, goFile := range goFiles {
		firstLine, err := goFile.ReadFirstLineAndTrimSpace()
		if err != nil {
			return false, err
		}

		if firstLine != packageMainString {
			continue
		}

		containsLine, err := goFile.ContainsLine(funcMainString)
		if err != nil {
			return false, err
		}

		filePath, err := goFile.GetLocalPath()
		if err != nil {
			return false, err
		}

		if containsLine {
			if verbose {
				logging.LogInfof(
					"Found '%s' and '%s' in '%s'",
					packageMainString,
					funcMainString,
					filePath,
				)
			}

			return true, nil
		}
	}

	if verbose {
		logging.LogInfof(
			"No file containing '%s' and '%s' found in '%s'.",
			packageMainString,
			funcMainString,
			path,
		)
	}

	return false, nil
}

func (g *GitRepositoryBase) CreateAndInit(createOptions *CreateRepositoryOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	err = parent.Create(createOptions.Verbose)
	if err != nil {
		return err
	}

	err = parent.Init(createOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitRepositoryBase) DirectoryByPathExists(verbose bool, path ...string) (exists bool, err error) {
	if len(path) <= 0 {
		return false, tracederrors.TracedError("path has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	subDir, err := parent.GetDirectoryByPath(path...)
	if err != nil {
		return false, err
	}

	exists, err = subDir.Exists(false)
	if err != nil {
		return false, err
	}

	if verbose {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return false, err
		}

		subDirPath, err := subDir.GetPath()
		if err != nil {
			return false, err
		}

		relativeSubDirPath, err := pathsutils.GetRelativePathTo(subDirPath, path)
		if err != nil {
			return false, err
		}

		if exists {
			logging.LogInfof(
				"Directory '%s' in git repository '%s' on host '%s' exists.",
				relativeSubDirPath,
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"Directory '%s' in git repository '%s' on host '%s' does not exist.",
				relativeSubDirPath,
				path,
				hostDescription,
			)
		}
	}

	return exists, err
}

func (g *GitRepositoryBase) EnsureMainReadmeMdExists(verbose bool) (err error) {
	const fileName string = "README.md"

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	_, err = parent.CreateFileInDirectory(verbose, "README.md")
	if err != nil {
		return err
	}

	err = parent.AddFileByPath(fileName, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitRepositoryBase) GetCurrentCommitMessage(verbose bool) (currentCommitMessage string, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return "", err
	}

	currentCommit, err := parent.GetCurrentCommit(verbose)
	if err != nil {
		return "", err
	}

	return currentCommit.GetCommitMessage()
}

func (g *GitRepositoryBase) GetCurrentCommitsNewestVersion(verbose bool) (newestVersion Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	currentCommit, err := parent.GetCurrentCommit(verbose)
	if err != nil {
		return nil, err
	}

	return currentCommit.GetNewestTagVersion(verbose)
}

func (g *GitRepositoryBase) GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	currentCommit, err := parent.GetCurrentCommit(verbose)
	if err != nil {
		return nil, err
	}

	return currentCommit.GetNewestTagVersionOrNilIfUnset(verbose)
}

func (g *GitRepositoryBase) GetFileByPath(path ...string) (file File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	rootDir, err := parent.GetRootDirectory(false)
	if err != nil {
		return nil, err
	}

	return rootDir.GetFileInDirectory(path...)
}

func (g *GitRepositoryBase) GetLatestTagVersion(verbose bool) (latestTagVersion Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	latestTagVersion, err = parent.GetLatestTagVersionOrNilIfNotFound(verbose)
	if err != nil {
		return nil, err
	}

	if latestTagVersion == nil {
		path, hostDescription, err := parent.GetPathAndHostDescription()
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"No version tag in git repository '%s' on host '%s' found.",
			path,
			hostDescription,
		)
	}

	return latestTagVersion, nil
}

func (g *GitRepositoryBase) GetLatestTagVersionAsString(verbose bool) (latestTagVersion string, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return "", err
	}

	version, err := parent.GetLatestTagVersion(verbose)
	if err != nil {
		return "", err
	}

	return version.GetAsString()
}

func (g *GitRepositoryBase) GetLatestTagVersionOrNilIfNotFound(verbose bool) (latestTagVersion Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	versionTags, err := parent.ListVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	for _, tag := range versionTags {
		toCheck, err := tag.GetVersion()
		if err != nil {
			return nil, err
		}

		if latestTagVersion == nil {
			latestTagVersion = toCheck
		}

		latestTagVersion, err = Versions().ReturnNewerVersion(latestTagVersion, toCheck)
		if err != nil {
			return nil, err
		}
	}

	return latestTagVersion, nil
}

func (g *GitRepositoryBase) GetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository, err error) {
	if g.parentRepositoryForBaseClass == nil {
		return nil, tracederrors.TracedErrorf("parentRepositoryForBaseClass not set")
	}

	return g.parentRepositoryForBaseClass, nil
}

func (g *GitRepositoryBase) GetPathAndHostDescription() (path string, hostDescription string, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
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

func (g *GitRepositoryBase) HasNoUncommittedChanges(verbose bool) (hasNoUncommittedChanges bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	hasUncommittedChanges, err := parent.HasUncommittedChanges(false)
	if err != nil {
		return false, err
	}

	path, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	hasNoUncommittedChanges = !hasUncommittedChanges

	if verbose {
		if hasNoUncommittedChanges {
			logging.LogInfof(
				"Git repository '%s' on host '%s' has no uncommitted changes.",
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"Git repository '%s' on host '%s' has uncommitted changes.",
				path,
				hostDescription,
			)
		}
	}

	return hasNoUncommittedChanges, nil
}

func (g *GitRepositoryBase) IsGolangApplication(verbose bool) (isGolangApplication bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	repoPath, err := parent.GetPath()
	if err != nil {
		return false, err
	}

	goModExists, err := parent.FileByPathExists("go.mod", verbose)
	if err != nil {
		return false, err
	}

	if !goModExists {
		if verbose {
			logging.LogInfof(
				"'%s' has no 'go.mod' present and is therefore not a golang application.",
				repoPath,
			)
		}
		return false, nil
	}

	mainGoExists, err := parent.FileByPathExists("main.go", verbose)
	if err != nil {
		return false, err
	}

	if mainGoExists {
		if verbose {
			logging.LogInfof("'%s' contains a go application since 'main.go' was found.", repoPath)
		}
	}

	isMainFuncPresent, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(verbose)
	if err != nil {
		return false, err
	}

	if isMainFuncPresent {
		if verbose {
			logging.LogInfof(
				"'%s' contains a go application since 'main' function was found.",
				repoPath,
			)
		}

		return true, nil
	}

	if verbose {
		logging.LogInfof(
			"'%s' does not contain a go application.",
			repoPath,
		)
	}

	return false, nil
}

func (g *GitRepositoryBase) IsGolangPackage(verbose bool) (isGolangPackage bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	repoPath, err := parent.GetPath()
	if err != nil {
		return false, err
	}

	goModExists, err := parent.FileByPathExists("go.mod", verbose)
	if err != nil {
		return false, err
	}

	if !goModExists {
		if verbose {
			logging.LogInfof("'%s' has no 'go.mod' present is not a golang package.", repoPath)
		}
		return false, nil
	}

	isMainFuncPresent, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(verbose)
	if err != nil {
		return false, err
	}

	if isMainFuncPresent {
		if verbose {
			logging.LogInfof("'%s' contains not a go package since 'main' function was found.", repoPath)
		}

		return false, nil
	}

	if verbose {
		logging.LogInfof("'%s' contains a go package.", repoPath)
	}

	return true, nil
}

func (g *GitRepositoryBase) IsOnLocalhost(verbose bool) (isOnLocalhost bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	hostDescription, err := parent.GetHostDescription()
	if err != nil {
		return false, err
	}

	isOnLocalhost = hostDescription == "localhost"

	if verbose {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return false, err
		}

		if isOnLocalhost {
			logging.LogInfof(
				"Git repository '%s' is on localhost",
				path,
			)
		} else {
			logging.LogInfof(
				"Git repository '%s' is not on localhost. Gost is '%s'.",
				path,
				hostDescription,
			)
		}
	}

	return isOnLocalhost, nil
}

func (g *GitRepositoryBase) IsPreCommitRepository(verbose bool) (isPreCommitRepository bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	isPreCommitRepository, err = parent.DirectoryByPathExists(false, "pre_commit_hooks")
	if err != nil {
		return false, err
	}

	if verbose {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return false, err
		}

		if isPreCommitRepository {
			logging.LogInfof(
				"Git reposiotry '%s' on host '%s' is a pre-commit repository.",
				path,
				hostDescription,
			)
		} else {
			logging.LogInfof(
				"Git reposiotry '%s' on host '%s' is not a pre-commit repository.",
				path,
				hostDescription,
			)
		}
	}

	return isPreCommitRepository, nil
}

func (g *GitRepositoryBase) ListVersionTags(verbose bool) (versionTags []GitTag, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	allTags, err := parent.ListTags(verbose)
	if err != nil {
		return nil, err
	}

	versionTags = []GitTag{}
	for _, tag := range allTags {
		isVersionTag, err := tag.IsVersionTag()
		if err != nil {
			return nil, err
		}

		if isVersionTag {
			versionTags = append(versionTags, tag)
		}
	}

	return versionTags, nil
}

func (g *GitRepositoryBase) MustAddFilesByPath(pathsToAdd []string, verbose bool) {
	err := g.AddFilesByPath(pathsToAdd, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustBranchByNameExists(branchName string, verbose bool) (branchExists bool) {
	branchExists, err := g.BranchByNameExists(branchName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return branchExists
}

func (g *GitRepositoryBase) MustCheckHasNoUncommittedChanges(verbose bool) {
	err := g.CheckHasNoUncommittedChanges(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCheckIsGolangApplication(verbose bool) {
	err := g.CheckIsGolangApplication(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCheckIsGolangPackage(verbose bool) {
	err := g.CheckIsGolangPackage(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCheckIsOnLocalhost(verbose bool) {
	err := g.CheckIsOnLocalhost(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCheckIsPreCommitRepository(verbose bool) {
	err := g.CheckIsPreCommitRepository(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := g.CommitAndPush(commitOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdCommit
}

func (g *GitRepositoryBase) MustCommitIfUncommittedChanges(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := g.CommitIfUncommittedChanges(commitOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdCommit
}

func (g *GitRepositoryBase) MustContainsGoSourceFileOfMainPackageWithMainFunction(verbose bool) (mainFound bool) {
	mainFound, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mainFound
}

func (g *GitRepositoryBase) MustCreateAndInit(createOptions *CreateRepositoryOptions) {
	err := g.CreateAndInit(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustDirectoryByPathExists(verbose bool, path ...string) (exists bool) {
	exists, err := g.DirectoryByPathExists(verbose, path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return exists
}

func (g *GitRepositoryBase) MustEnsureMainReadmeMdExists(verbose bool) {
	err := g.EnsureMainReadmeMdExists(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustGetCurrentCommitMessage(verbose bool) (currentCommitMessage string) {
	currentCommitMessage, err := g.GetCurrentCommitMessage(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return currentCommitMessage
}

func (g *GitRepositoryBase) MustGetCurrentCommitsNewestVersion(verbose bool) (newestVersion Version) {
	newestVersion, err := g.GetCurrentCommitsNewestVersion(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestVersion
}

func (g *GitRepositoryBase) MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion Version) {
	newestVersion, err := g.GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return newestVersion
}

func (g *GitRepositoryBase) MustGetFileByPath(path ...string) (file File) {
	file, err := g.GetFileByPath(path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return file
}

func (g *GitRepositoryBase) MustGetLatestTagVersion(verbose bool) (latestTagVersion Version) {
	latestTagVersion, err := g.GetLatestTagVersion(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestTagVersion
}

func (g *GitRepositoryBase) MustGetLatestTagVersionAsString(verbose bool) (latestTagVersion string) {
	latestTagVersion, err := g.GetLatestTagVersionAsString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestTagVersion
}

func (g *GitRepositoryBase) MustGetLatestTagVersionOrNilIfNotFound(verbose bool) (latestTagVersion Version) {
	latestTagVersion, err := g.GetLatestTagVersionOrNilIfNotFound(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return latestTagVersion
}

func (g *GitRepositoryBase) MustGetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository) {
	parentRepositoryForBaseClass, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentRepositoryForBaseClass
}

func (g *GitRepositoryBase) MustGetPathAndHostDescription() (path string, hostDescription string) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path, hostDescription
}

func (g *GitRepositoryBase) MustHasNoUncommittedChanges(verbose bool) (hasNoUncommittedChanges bool) {
	hasNoUncommittedChanges, err := g.HasNoUncommittedChanges(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasNoUncommittedChanges
}

func (g *GitRepositoryBase) MustIsGolangApplication(verbose bool) (isGolangApplication bool) {
	isGolangApplication, err := g.IsGolangApplication(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isGolangApplication
}

func (g *GitRepositoryBase) MustIsGolangPackage(verbose bool) (isGolangPackage bool) {
	isGolangPackage, err := g.IsGolangPackage(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isGolangPackage
}

func (g *GitRepositoryBase) MustIsOnLocalhost(verbose bool) (isOnLocalhost bool) {
	isOnLocalhost, err := g.IsOnLocalhost(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isOnLocalhost
}

func (g *GitRepositoryBase) MustIsPreCommitRepository(verbose bool) (isPreCommitRepository bool) {
	isPreCommitRepository, err := g.IsPreCommitRepository(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isPreCommitRepository
}

func (g *GitRepositoryBase) MustListVersionTags(verbose bool) (versionTags []GitTag) {
	versionTags, err := g.ListVersionTags(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitRepositoryBase) MustSetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) {
	err := g.SetParentRepositoryForBaseClass(parentRepositoryForBaseClass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustWriteStringToFile(content string, verbose bool, path ...string) (writtenFile File) {
	writtenFile, err := g.WriteStringToFile(content, verbose, path...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return writtenFile
}

func (g *GitRepositoryBase) SetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) (err error) {
	if parentRepositoryForBaseClass == nil {
		return tracederrors.TracedErrorf("parentRepositoryForBaseClass is nil")
	}

	g.parentRepositoryForBaseClass = parentRepositoryForBaseClass

	return nil
}

func (g *GitRepositoryBase) WriteStringToFile(content string, verbose bool, path ...string) (writtenFile File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	return parent.WriteBytesToFile([]byte(content), verbose, path...)
}
