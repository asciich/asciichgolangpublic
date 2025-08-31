package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateReleaseOptions struct {
	Name        string
	Description string
}

func NewGitlabCreateReleaseOptions() (g *GitlabCreateReleaseOptions) {
	return new(GitlabCreateReleaseOptions)
}

func (g *GitlabCreateReleaseOptions) GetDescription() (description string, err error) {
	if g.Description == "" {
		return "", tracederrors.TracedErrorf("Description not set")
	}

	return g.Description, nil
}

func (g *GitlabCreateReleaseOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return g.Name, nil
}


func (g *GitlabCreateReleaseOptions) SetDescription(description string) (err error) {
	if description == "" {
		return tracederrors.TracedErrorf("description is empty string")
	}

	g.Description = description

	return nil
}

func (g *GitlabCreateReleaseOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}
