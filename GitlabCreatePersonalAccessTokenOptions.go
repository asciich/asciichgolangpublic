package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreatePersonalAccessTokenOptions struct {
	Name string
}

func NewGitlabCreatePersonalAccessTokenOptions() (g *GitlabCreatePersonalAccessTokenOptions) {
	return new(GitlabCreatePersonalAccessTokenOptions)
}

func (g *GitlabCreatePersonalAccessTokenOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}
func (o *GitlabCreatePersonalAccessTokenOptions) GetName() (name string, err error) {
	if len(o.Name) <= 0 {
		return "", tracederrors.TracedError("name not set")
	}

	return o.Name, nil
}
