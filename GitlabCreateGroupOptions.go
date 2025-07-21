package asciichgolangpublic

import "github.com/asciich/asciichgolangpublic/pkg/logging"

type GitlabCreateGroupOptions struct {
	Verbose bool
}

func NewGitlabCreateGroupOptions() (createOptions *GitlabCreateGroupOptions) {
	return new(GitlabCreateGroupOptions)
}

func (g *GitlabCreateGroupOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreateGroupOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreateGroupOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateGroupOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabCreateGroupOptions) GetDeepCopy() (copy *GitlabCreateGroupOptions) {
	copy = NewGitlabCreateGroupOptions()

	*copy = *o

	return copy
}
