package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateTagOptions struct {
	Name    string
	Verbose bool
	Ref     string
}

func NewGitlabCreateTagOptions() (g *GitlabCreateTagOptions) {
	return new(GitlabCreateTagOptions)
}

func (g *GitlabCreateTagOptions) GetDeepCopy() (deepCopy *GitlabCreateTagOptions) {
	deepCopy = NewGitlabCreateTagOptions()

	*deepCopy = *g

	return deepCopy
}

func (g *GitlabCreateTagOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitlabCreateTagOptions) GetRef() (ref string, err error) {
	if g.Ref == "" {
		return "", tracederrors.TracedErrorf("Ref not set")
	}

	return g.Ref, nil
}

func (g *GitlabCreateTagOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabCreateTagOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateTagOptions) SetRef(ref string) (err error) {
	if ref == "" {
		return tracederrors.TracedErrorf("ref is empty string")
	}

	g.Ref = ref

	return nil
}

func (g *GitlabCreateTagOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
