package asciichgolangpublic

type GitRepositoryBase struct {
	parentRepositoryForBaseClass GitRepository
}

func NewGitRepositoryBase() (g *GitRepositoryBase) {
	return new(GitRepositoryBase)
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

		return TracedErrorf(
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

	return TracedErrorf(
		"git repository '%s' on host '%s' is not a golang application",
		path,
		hostDescription,
	)
}

func (g *GitRepositoryBase) CommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit, err error) {
	if commitOptions == nil {
		return nil, TracedErrorNil("commitOptions")
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

func (g *GitRepositoryBase) ContainsGoSourceFileOfMainPackageWithMainFunction(verbose bool) (mainFound bool, err error) {
	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return false, err
	}

	goFiles, err := parent.ListFiles(
		&ListFileOptions{
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
				LogInfof(
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
		LogInfof(
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
		return TracedErrorNil("createOptions")
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
		return nil, TracedError("path has no elements")
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
	versionTags, err := g.ListVersionTags(verbose)
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

func (g *GitRepositoryBase) GetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository, err error) {
	if g.parentRepositoryForBaseClass == nil {
		return nil, TracedErrorf("parentRepositoryForBaseClass not set")
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
			LogInfof(
				"Git repository '%s' on host '%s' has no uncommitted changes.",
				path,
				hostDescription,
			)
		} else {
			LogInfof(
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
			LogInfof(
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
			LogInfof("'%s' contains a go application since 'main.go' was found.", repoPath)
		}
	}

	isMainFuncPresent, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(verbose)
	if err != nil {
		return false, err
	}

	if isMainFuncPresent {
		if verbose {
			LogInfof(
				"'%s' contains a go application since 'main' function was found.",
				repoPath,
			)
		}

		return true, nil
	}

	if verbose {
		LogInfof(
			"'%s' does not contain a go application.",
			repoPath,
		)
	}

	return false, nil
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

func (g *GitRepositoryBase) MustCheckHasNoUncommittedChanges(verbose bool) {
	err := g.CheckHasNoUncommittedChanges(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCheckIsGolangApplication(verbose bool) {
	err := g.CheckIsGolangApplication(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustCommitAndPush(commitOptions *GitCommitOptions) (createdCommit *GitCommit) {
	createdCommit, err := g.CommitAndPush(commitOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return createdCommit
}

func (g *GitRepositoryBase) MustContainsGoSourceFileOfMainPackageWithMainFunction(verbose bool) (mainFound bool) {
	mainFound, err := g.ContainsGoSourceFileOfMainPackageWithMainFunction(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return mainFound
}

func (g *GitRepositoryBase) MustCreateAndInit(createOptions *CreateRepositoryOptions) {
	err := g.CreateAndInit(createOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustGetCurrentCommitsNewestVersion(verbose bool) (newestVersion Version) {
	newestVersion, err := g.GetCurrentCommitsNewestVersion(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newestVersion
}

func (g *GitRepositoryBase) MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose bool) (newestVersion Version) {
	newestVersion, err := g.GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return newestVersion
}

func (g *GitRepositoryBase) MustGetFileByPath(path ...string) (file File) {
	file, err := g.GetFileByPath(path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return file
}

func (g *GitRepositoryBase) MustGetLatestTagVersion(verbose bool) (latestTagVersion Version) {
	latestTagVersion, err := g.GetLatestTagVersion(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return latestTagVersion
}

func (g *GitRepositoryBase) MustGetLatestTagVersionAsString(verbose bool) (latestTagVersion string) {
	latestTagVersion, err := g.GetLatestTagVersionAsString(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return latestTagVersion
}

func (g *GitRepositoryBase) MustGetParentRepositoryForBaseClass() (parentRepositoryForBaseClass GitRepository) {
	parentRepositoryForBaseClass, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentRepositoryForBaseClass
}

func (g *GitRepositoryBase) MustGetPathAndHostDescription() (path string, hostDescription string) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path, hostDescription
}

func (g *GitRepositoryBase) MustHasNoUncommittedChanges(verbose bool) (hasNoUncommittedChanges bool) {
	hasNoUncommittedChanges, err := g.HasNoUncommittedChanges(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return hasNoUncommittedChanges
}

func (g *GitRepositoryBase) MustIsGolangApplication(verbose bool) (isGolangApplication bool) {
	isGolangApplication, err := g.IsGolangApplication(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isGolangApplication
}

func (g *GitRepositoryBase) MustListVersionTags(verbose bool) (versionTags []GitTag) {
	versionTags, err := g.ListVersionTags(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return versionTags
}

func (g *GitRepositoryBase) MustSetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) {
	err := g.SetParentRepositoryForBaseClass(parentRepositoryForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryBase) MustWriteStringToFile(content string, verbose bool, path ...string) (writtenFile File) {
	writtenFile, err := g.WriteStringToFile(content, verbose, path...)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return writtenFile
}

func (g *GitRepositoryBase) SetParentRepositoryForBaseClass(parentRepositoryForBaseClass GitRepository) (err error) {
	if parentRepositoryForBaseClass == nil {
		return TracedErrorf("parentRepositoryForBaseClass is nil")
	}

	g.parentRepositoryForBaseClass = parentRepositoryForBaseClass

	return nil
}

func (g *GitRepositoryBase) WriteStringToFile(content string, verbose bool, path ...string) (writtenFile File, err error) {
	if len(path) <= 0 {
		return nil, TracedError("path has no elements")
	}

	parent, err := g.GetParentRepositoryForBaseClass()
	if err != nil {
		return nil, err
	}

	return parent.WriteBytesToFile([]byte(content), verbose, path...)
}
