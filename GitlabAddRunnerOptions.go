package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabAddRunnerOptions struct {
	Name       string
	RunnerTags []string
}

func NewGitlabAddRunnerOptions() (g *GitlabAddRunnerOptions) {
	return new(GitlabAddRunnerOptions)
}

func (g *GitlabAddRunnerOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitlabAddRunnerOptions) GetRunnerTags() (runnerTags []string, err error) {
	if g.RunnerTags == nil {
		return nil, tracederrors.TracedErrorf("RunnerTags not set")
	}

	if len(g.RunnerTags) <= 0 {
		return nil, tracederrors.TracedErrorf("RunnerTags has no elements")
	}

	return g.RunnerTags, nil
}

func (g *GitlabAddRunnerOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabAddRunnerOptions) SetRunnerTags(runnerTags []string) (err error) {
	if runnerTags == nil {
		return tracederrors.TracedErrorf("runnerTags is nil")
	}

	if len(runnerTags) <= 0 {
		return tracederrors.TracedErrorf("runnerTags has no elements")
	}

	g.RunnerTags = runnerTags

	return nil
}

func (o *GitlabAddRunnerOptions) GetRunnerName() (runnerName string, err error) {
	if len(o.Name) <= 0 {
		return "", tracederrors.TracedError("Name not set")
	}

	return o.Name, nil
}

func (o *GitlabAddRunnerOptions) GetTags() (runnerTags []string, err error) {
	if len(o.RunnerTags) <= 0 {
		return nil, tracederrors.TracedError("RunnerTags not set")
	}

	return o.RunnerTags, nil
}

func (o *GitlabAddRunnerOptions) GetTagsCommaSeparated() (tagsCommaSeperated string, err error) {
	tags, err := o.GetTags()
	if err != nil {
		return "", err
	}

	tagsCommaSeperated = strings.Join(tags, ",")
	return tagsCommaSeperated, nil
}
