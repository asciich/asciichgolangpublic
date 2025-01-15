package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabCreatePersonalAccessTokenOptions struct {
	Name    string
	Verbose bool
}

func NewGitlabCreatePersonalAccessTokenOptions() (g *GitlabCreatePersonalAccessTokenOptions) {
	return new(GitlabCreatePersonalAccessTokenOptions)
}

func (g *GitlabCreatePersonalAccessTokenOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreatePersonalAccessTokenOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabCreatePersonalAccessTokenOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreatePersonalAccessTokenOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreatePersonalAccessTokenOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreatePersonalAccessTokenOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreatePersonalAccessTokenOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GitlabCreatePersonalAccessTokenOptions) GetName() (name string, err error) {
	if len(o.Name) <= 0 {
		return "", tracederrors.TracedError("name not set")
	}

	return o.Name, nil
}
