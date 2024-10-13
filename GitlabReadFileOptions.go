package asciichgolangpublic

type GitlabReadFileOptions struct {
	Path       string
	BranchName string
	Verbose    bool
}

func NewGitlabReadFileOptions() (g *GitlabReadFileOptions) {
	return new(GitlabReadFileOptions)
}

func (g *GitlabReadFileOptions) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", TracedErrorf("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabReadFileOptions) GetGitlabGetRepositoryFileOptions() (getOptions *GitlabGetRepositoryFileOptions, err error) {
	getOptions = NewGitlabGetRepositoryFileOptions()
	getOptions.Path = g.Path
	getOptions.BranchName = g.BranchName
	getOptions.Verbose = g.Verbose
	return getOptions, nil
}

func (g *GitlabReadFileOptions) GetPath() (path string, err error) {
	if g.Path == "" {
		return "", TracedErrorf("Path not set")
	}

	return g.Path, nil
}

func (g *GitlabReadFileOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabReadFileOptions) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitlabReadFileOptions) MustGetGitlabGetRepositoryFileOptions() (getOptions *GitlabGetRepositoryFileOptions) {
	getOptions, err := g.GetGitlabGetRepositoryFileOptions()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return getOptions
}

func (g *GitlabReadFileOptions) MustGetPath() (path string) {
	path, err := g.GetPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return path
}

func (g *GitlabReadFileOptions) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabReadFileOptions) MustSetPath(path string) {
	err := g.SetPath(path)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabReadFileOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabReadFileOptions) SetPath(path string) (err error) {
	if path == "" {
		return TracedErrorf("path is empty string")
	}

	g.Path = path

	return nil
}

func (g *GitlabReadFileOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}

func (g *GitlabReadFileOptions) GetDeepCopy() (deepCopy *GitlabReadFileOptions) {
	deepCopy = NewGitlabReadFileOptions()

	*deepCopy = *g

	return deepCopy
}
