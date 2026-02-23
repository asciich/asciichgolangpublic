package gitgeneric

import (
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GitTagBase struct {
	parentGitTagForBaseClass gitinterfaces.GitTag
}

func NewGitTagBase() (g *GitTagBase) {
	return new(GitTagBase)
}

func (g *GitTagBase) GetParentGitTagForBaseClass() (parentGitTagForBaseClass gitinterfaces.GitTag, err error) {

	return g.parentGitTagForBaseClass, nil
}

func (g *GitTagBase) GetVersion() (version versionutils.Version, err error) {
	parent, err := g.GetParentGitTagForBaseClass()
	if err != nil {
		return nil, err
	}

	name, err := parent.GetName()
	if err != nil {
		return nil, err
	}

	return versionutils.NewFromString(name)
}

func (g *GitTagBase) MustGetParentGitTagForBaseClass() (parentGitTagForBaseClass gitinterfaces.GitTag) {
	parentGitTagForBaseClass, err := g.GetParentGitTagForBaseClass()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentGitTagForBaseClass
}

func (g *GitTagBase) MustGetVersion() (version versionutils.Version) {
	version, err := g.GetVersion()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func (g *GitTagBase) MustSetParentGitTagForBaseClass(parentGitTagForBaseClass gitinterfaces.GitTag) {
	err := g.SetParentGitTagForBaseClass(parentGitTagForBaseClass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitTagBase) SetParentGitTagForBaseClass(parentGitTagForBaseClass gitinterfaces.GitTag) (err error) {
	g.parentGitTagForBaseClass = parentGitTagForBaseClass

	return nil
}
