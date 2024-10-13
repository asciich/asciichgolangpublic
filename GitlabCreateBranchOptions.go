package asciichgolangpublic

type GitlabCreateBranchOptions struct {
	SourceBranchName    string
	BranchName          string
	Verbose             bool
	FailIfAlreadyExists bool
}

func NewGitlabCreateBranchOptions() (g *GitlabCreateBranchOptions) {
	return new(GitlabCreateBranchOptions)
}

func (g *GitlabCreateBranchOptions) GetBranchName() (branchName string, err error) {
	if g.BranchName == "" {
		return "", TracedErrorf("BranchName not set")
	}

	return g.BranchName, nil
}

func (g *GitlabCreateBranchOptions) GetFailIfAlreadyExists() (failIfAlreadyExists bool) {

	return g.FailIfAlreadyExists
}

func (g *GitlabCreateBranchOptions) GetSourceBranchName() (sourceBranchName string, err error) {
	if g.SourceBranchName == "" {
		return "", TracedErrorf("SourceBranchName not set")
	}

	return g.SourceBranchName, nil
}

func (g *GitlabCreateBranchOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabCreateBranchOptions) MustGetBranchName() (branchName string) {
	branchName, err := g.GetBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return branchName
}

func (g *GitlabCreateBranchOptions) MustGetSourceBranchName() (sourceBranchName string) {
	sourceBranchName, err := g.GetSourceBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sourceBranchName
}

func (g *GitlabCreateBranchOptions) MustSetBranchName(branchName string) {
	err := g.SetBranchName(branchName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateBranchOptions) MustSetSourceBranchName(sourceBranchName string) {
	err := g.SetSourceBranchName(sourceBranchName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateBranchOptions) SetBranchName(branchName string) (err error) {
	if branchName == "" {
		return TracedErrorf("branchName is empty string")
	}

	g.BranchName = branchName

	return nil
}

func (g *GitlabCreateBranchOptions) SetFailIfAlreadyExists(failIfAlreadyExists bool) {
	g.FailIfAlreadyExists = failIfAlreadyExists
}

func (g *GitlabCreateBranchOptions) SetSourceBranchName(sourceBranchName string) (err error) {
	if sourceBranchName == "" {
		return TracedErrorf("sourceBranchName is empty string")
	}

	g.SourceBranchName = sourceBranchName

	return nil
}

func (g *GitlabCreateBranchOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
