package asciichgolangpublic

type GitlabCreateMergeRequestOptions struct {
	SourceBranchName string
	TargetBranchName string
	Title            string
	Verbose          bool
}

func NewGitlabCreateMergeRequestOptions() (g *GitlabCreateMergeRequestOptions) {
	return new(GitlabCreateMergeRequestOptions)
}

func (g *GitlabCreateMergeRequestOptions) GetDeepCopy() (copy *GitlabCreateMergeRequestOptions) {
	copy = NewGitlabCreateMergeRequestOptions()
	*copy = *g
	return copy
}

func (g *GitlabCreateMergeRequestOptions) GetSourceBranchName() (sourceBranchName string, err error) {
	if g.SourceBranchName == "" {
		return "", TracedErrorf("SourceBranchName not set")
	}

	return g.SourceBranchName, nil
}

func (g *GitlabCreateMergeRequestOptions) GetTargetBranchName() (targetBranchName string, err error) {
	if g.TargetBranchName == "" {
		return "", TracedErrorf("TargetBranchName not set")
	}

	return g.TargetBranchName, nil
}

func (g *GitlabCreateMergeRequestOptions) GetTitle() (title string, err error) {
	if g.Title == "" {
		return "", TracedErrorf("Title not set")
	}

	return g.Title, nil
}

func (g *GitlabCreateMergeRequestOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabCreateMergeRequestOptions) IsTargetBranchSet() (isSet bool) {
	return g.TargetBranchName != ""
}

func (g *GitlabCreateMergeRequestOptions) MustGetSourceBranchName() (sourceBranchName string) {
	sourceBranchName, err := g.GetSourceBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sourceBranchName
}

func (g *GitlabCreateMergeRequestOptions) MustGetTargetBranchName() (targetBranchName string) {
	targetBranchName, err := g.GetTargetBranchName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return targetBranchName
}

func (g *GitlabCreateMergeRequestOptions) MustGetTitle() (title string) {
	title, err := g.GetTitle()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return title
}

func (g *GitlabCreateMergeRequestOptions) MustSetSourceBranchName(sourceBranchName string) {
	err := g.SetSourceBranchName(sourceBranchName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) MustSetTargetBranchName(targetBranchName string) {
	err := g.SetTargetBranchName(targetBranchName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) MustSetTitle(title string) {
	err := g.SetTitle(title)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateMergeRequestOptions) SetSourceBranchName(sourceBranchName string) (err error) {
	if sourceBranchName == "" {
		return TracedErrorf("sourceBranchName is empty string")
	}

	g.SourceBranchName = sourceBranchName

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetTargetBranchName(targetBranchName string) (err error) {
	if targetBranchName == "" {
		return TracedErrorf("targetBranchName is empty string")
	}

	g.TargetBranchName = targetBranchName

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetTitle(title string) (err error) {
	if title == "" {
		return TracedErrorf("title is empty string")
	}

	g.Title = title

	return nil
}

func (g *GitlabCreateMergeRequestOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
