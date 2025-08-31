package asciichgolangpublic

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GitRepositoryBase struct {
	parentRepositoryForBaseClass GitRepository
}

func NewGitRepositoryBase() (g *GitRepositoryBase) {
	return new(GitRepositoryBase)
}

func (g *GitRepositoryBase) CheckExists(ctx context.Context) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	exists, err := parent.Exists(ctx)
	if err != nil {
		return err
	}

	path, err := parent.GetPath()
	if err != nil {
		return err
	}

	if !exists {
		return tracederrors.TracedErrorf("Local git repository '%s' does not exist.", path)
	}

	logging.LogInfoByCtxf(ctx, "Local git repository '%s' exists.", path)

	return nil
}

func (g *GitRepositoryBase) AddFilesByPath(ctx context.Context, pathsToAdd []string) (err error) {
	if len(pathsToAdd) <= 0 {
		return tracederrors.TracedError("pathToAdd has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	for _, p := range pathsToAdd {
		err = parent.AddFileByPath(ctx, p)
		if err != nil {
			return err
		}
	}

	path, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Added '%d' files to git repository '%s' on host '%s'.", len(pathsToAdd), path, hostDescription)

	return nil
}

func (g *GitRepositoryBase) BranchByNameExists(ctx context.Context, branchName string) (branchExists bool, err error) {
	if branchName == "" {
		return false, tracederrors.TracedErrorEmptyString("branchName")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	branchNames, err := parent.ListBranchNames(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	branchExists = slices.Contains(branchNames, branchName)

	path, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	if branchExists {
		logging.LogInfoByCtxf(ctx, "Branch '%s' in git repository '%s' on host '%s' exists.", branchName, path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Branch '%s' in git repository '%s' on host '%s' does not exist.", branchName, path, hostDescription)
	}

	return branchExists, nil
}

func (g *GitRepositoryBase) CheckHasNoUncommittedChanges(ctx context.Context) (err error) {
	hasNoUncommitedChanges, err := g.HasNoUncommittedChanges(ctx)
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

func (g *GitRepositoryBase) CheckIsGolangApplication(ctx context.Context) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	ok, err := parent.IsGolangApplication(ctx)
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

func (g *GitRepositoryBase) CheckIsGolangPackage(ctx context.Context) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	isGolangPackage, err := parent.IsGolangPackage(ctx)
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

func (g *GitRepositoryBase) CheckIsOnLocalhost(ctx context.Context) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	isOnLocalhost, err := parent.IsOnLocalhost(ctx)
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

func (g *GitRepositoryBase) CheckIsPreCommitRepository(ctx context.Context) (err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	isPreCommitRepository, err := parent.IsPreCommitRepository(ctx)
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

func (g *GitRepositoryBase) CommitAndPush(ctx context.Context, commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	createdCommit, err = parent.Commit(ctx, commitOptions)
	if err != nil {
		return nil, err
	}

	err = parent.Push(ctx)
	if err != nil {
		return nil, err
	}

	return createdCommit, nil
}

func (g *GitRepositoryBase) CommitIfUncommittedChanges(ctx context.Context, commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	hasUncommitedChanges, err := parent.HasUncommittedChanges(ctx)
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

		createdCommit, err = parent.Commit(ctx, optionsToUse)
		if err != nil {
			return nil, err
		}

		createdHash, err := createdCommit.GetHash(ctx)
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "Commited all uncommited changes in git repository '%s' on host '%s' as commit '%s'.", path, hostDescription, createdHash)
	} else {
		createdCommit, err = parent.GetCurrentCommit(contextutils.WithSilent(ctx))
		if err != nil {
			return nil, err
		}

		logging.LogInfoByCtxf(ctx, "No uncommited changes to commit in git repository '%s' on host '%s'.", path, hostDescription)
	}

	return createdCommit, nil
}

func (g *GitRepositoryBase) ContainsGoSourceFileOfMainPackageWithMainFunction(ctx context.Context) (mainFound bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	goFiles, err := parent.ListFiles(
		ctx,
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
			logging.LogInfoByCtxf(ctx, "Found '%s' and '%s' in '%s'", packageMainString, funcMainString, filePath)
			return true, nil
		}
	}

	logging.LogInfoByCtxf(ctx, "No file containing '%s' and '%s' found in '%s'.", packageMainString, funcMainString, path)

	return false, nil
}

func (g *GitRepositoryBase) CreateAndInit(ctx context.Context, createOptions *parameteroptions.CreateRepositoryOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	err = parent.Create(ctx, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	err = parent.Init(ctx, createOptions)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitRepositoryBase) DirectoryByPathExists(ctx context.Context, path ...string) (exists bool, err error) {
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

	exists, err = subDir.Exists(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	repoPath, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	subDirPath, err := subDir.GetPath()
	if err != nil {
		return false, err
	}

	relativeSubDirPath, err := pathsutils.GetRelativePathTo(subDirPath, repoPath)
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Directory '%s' in git repository '%s' on host '%s' exists.", relativeSubDirPath, path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Directory '%s' in git repository '%s' on host '%s' does not exist.", relativeSubDirPath, path, hostDescription)
	}

	return exists, err
}

func (g *GitRepositoryBase) EnsureMainReadmeMdExists(ctx context.Context) (err error) {
	const fileName string = "README.md"

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return err
	}

	_, err = parent.CreateFileInDirectory(ctx, "README.md", &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	err = parent.AddFileByPath(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (g *GitRepositoryBase) GetCurrentCommitMessage(ctx context.Context) (currentCommitMessage string, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return "", err
	}

	currentCommit, err := parent.GetCurrentCommit(ctx)
	if err != nil {
		return "", err
	}

	return currentCommit.GetCommitMessage(ctx)
}

func (g *GitRepositoryBase) GetCurrentCommitsNewestVersion(ctx context.Context) (newestVersion versionutils.Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	currentCommit, err := parent.GetCurrentCommit(ctx)
	if err != nil {
		return nil, err
	}

	return currentCommit.GetNewestTagVersion(ctx)
}

func (g *GitRepositoryBase) GetCurrentCommitsNewestVersionOrNilIfNotPresent(ctx context.Context) (newestVersion versionutils.Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	currentCommit, err := parent.GetCurrentCommit(ctx)
	if err != nil {
		return nil, err
	}

	return currentCommit.GetNewestTagVersionOrNilIfUnset(ctx)
}

func (g *GitRepositoryBase) GetFileByPath(path ...string) (file filesinterfaces.File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	rootDir, err := parent.GetRootDirectory(contextutils.ContextSilent())
	if err != nil {
		return nil, err
	}

	return rootDir.GetFileInDirectory(path...)
}

func (g *GitRepositoryBase) GetLatestTagVersion(ctx context.Context) (latestTagVersion versionutils.Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	latestTagVersion, err = parent.GetLatestTagVersionOrNilIfNotFound(ctx)
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

func (g *GitRepositoryBase) GetLatestTagVersionAsString(ctx context.Context) (latestTagVersion string, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return "", err
	}

	version, err := parent.GetLatestTagVersion(ctx)
	if err != nil {
		return "", err
	}

	return version.GetAsString()
}

func (g *GitRepositoryBase) GetLatestTagVersionOrNilIfNotFound(ctx context.Context) (latestTagVersion versionutils.Version, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	versionTags, err := parent.ListVersionTags(ctx)
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

		latestTagVersion, err = versionutils.ReturnNewerVersion(latestTagVersion, toCheck)
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

func (g *GitRepositoryBase) HasNoUncommittedChanges(ctx context.Context) (hasNoUncommittedChanges bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	hasUncommittedChanges, err := parent.HasUncommittedChanges(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	path, hostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	hasNoUncommittedChanges = !hasUncommittedChanges

	if hasNoUncommittedChanges {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' has no uncommitted changes.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' has uncommitted changes.", path, hostDescription)
	}

	return hasNoUncommittedChanges, nil
}

func (g *GitRepositoryBase) IsGolangApplication(ctx context.Context) (isGolangApplication bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	repoPath, err := parent.GetPath()
	if err != nil {
		return false, err
	}

	goModExists, err := parent.FileByPathExists(ctx, "go.mod")
	if err != nil {
		return false, err
	}

	if !goModExists {
		logging.LogInfoByCtxf(ctx, "'%s' has no 'go.mod' present and is therefore not a golang application.", repoPath)
		return false, nil
	}

	mainGoExists, err := parent.FileByPathExists(ctx, "main.go")
	if err != nil {
		return false, err
	}

	if mainGoExists {
		logging.LogInfoByCtxf(ctx, "'%s' contains a go application since 'main.go' was found.", repoPath)
	}

	isMainFuncPresent, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(ctx)
	if err != nil {
		return false, err
	}

	if isMainFuncPresent {
		logging.LogInfoByCtxf(ctx, "'%s' contains a go application since 'main' function was found.", repoPath)

		return true, nil
	}

	logging.LogInfoByCtxf(ctx, "'%s' does not contain a go application.", repoPath)

	return false, nil
}

func (g *GitRepositoryBase) IsGolangPackage(ctx context.Context) (isGolangPackage bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	repoPath, err := parent.GetPath()
	if err != nil {
		return false, err
	}

	goModExists, err := parent.FileByPathExists(ctx, "go.mod")
	if err != nil {
		return false, err
	}

	if !goModExists {
		logging.LogInfoByCtxf(ctx, "'%s' has no 'go.mod' present is not a golang package.", repoPath)
		return false, nil
	}

	isMainFuncPresent, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(ctx)
	if err != nil {
		return false, err
	}

	if isMainFuncPresent {
		logging.LogInfoByCtxf(ctx, "'%s' contains not a go package since 'main' function was found.", repoPath)

		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "'%s' contains a go package.", repoPath)

	return true, nil
}

func (g *GitRepositoryBase) IsOnLocalhost(ctx context.Context) (isOnLocalhost bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	hostDescription, err := parent.GetHostDescription()
	if err != nil {
		return false, err
	}

	isOnLocalhost = hostDescription == "localhost"

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	if isOnLocalhost {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' is on localhost", path)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' is not on localhost. Gost is '%s'.", path, hostDescription)
	}

	return isOnLocalhost, nil
}

func (g *GitRepositoryBase) IsPreCommitRepository(ctx context.Context) (isPreCommitRepository bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	isPreCommitRepository, err = parent.DirectoryByPathExists(contextutils.WithSilent(ctx), "pre_commit_hooks")
	if err != nil {
		return false, err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	if isPreCommitRepository {
		logging.LogInfoByCtxf(ctx, "Git reposiotry '%s' on host '%s' is a pre-commit repository.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Git reposiotry '%s' on host '%s' is not a pre-commit repository.", path, hostDescription)
	}

	return isPreCommitRepository, nil
}

func (g *GitRepositoryBase) ListVersionTags(ctx context.Context) (versionTags []GitTag, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	allTags, err := parent.ListTags(ctx)
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

func (g *GitRepositoryBase) SetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) (err error) {
	if parentRepositoryForBaseClass == nil {
		return tracederrors.TracedErrorf("parentRepositoryForBaseClass is nil")
	}

	g.parentRepositoryForBaseClass = parentRepositoryForBaseClass

	return nil
}

func (g *GitRepositoryBase) WriteStringToFile(ctx context.Context, path string, content string, options *filesoptions.WriteOptions) (writtenFile filesinterfaces.File, err error) {
	if len(path) <= 0 {
		return nil, tracederrors.TracedError("path has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	return parent.WriteBytesToFile(ctx, path, []byte(content), options)
}
