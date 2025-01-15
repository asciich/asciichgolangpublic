package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabAddRunnerOptions struct {
	Name       string
	RunnerTags []string
	Verbose    bool
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

func (g *GitlabAddRunnerOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabAddRunnerOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabAddRunnerOptions) MustGetRunnerName() (runnerName string) {
	runnerName, err := g.GetRunnerName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runnerName
}

func (g *GitlabAddRunnerOptions) MustGetRunnerTags() (runnerTags []string) {
	runnerTags, err := g.GetRunnerTags()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runnerTags
}

func (g *GitlabAddRunnerOptions) MustGetTags() (runnerTags []string) {
	runnerTags, err := g.GetTags()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return runnerTags
}

func (g *GitlabAddRunnerOptions) MustGetTagsCommaSeparated() (tagsCommaSeperated string) {
	tagsCommaSeperated, err := g.GetTagsCommaSeparated()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tagsCommaSeperated
}

func (g *GitlabAddRunnerOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabAddRunnerOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabAddRunnerOptions) MustSetRunnerTags(runnerTags []string) {
	err := g.SetRunnerTags(runnerTags)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabAddRunnerOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (g *GitlabAddRunnerOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

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
