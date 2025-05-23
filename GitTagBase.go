package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GitTagBase struct {
	parentGitTagForBaseClass GitTag
}

func NewGitTagBase() (g *GitTagBase) {
	return new(GitTagBase)
}

func (g *GitTagBase) GetParentGitTagForBaseClass() (parentGitTagForBaseClass GitTag, err error) {

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

	return versionutils.Versions().GetNewVersionByString(name)
}

func (g *GitTagBase) MustGetParentGitTagForBaseClass() (parentGitTagForBaseClass GitTag) {
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

func (g *GitTagBase) MustSetParentGitTagForBaseClass(parentGitTagForBaseClass GitTag) {
	err := g.SetParentGitTagForBaseClass(parentGitTagForBaseClass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitTagBase) SetParentGitTagForBaseClass(parentGitTagForBaseClass GitTag) (err error) {
	g.parentGitTagForBaseClass = parentGitTagForBaseClass

	return nil
}
