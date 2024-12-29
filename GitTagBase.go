package asciichgolangpublic

type GitTagBase struct {
	parentGitTagForBaseClass GitTag
}

func NewGitTagBase() (g *GitTagBase) {
	return new(GitTagBase)
}

func (g *GitTagBase) GetParentGitTagForBaseClass() (parentGitTagForBaseClass GitTag, err error) {

	return g.parentGitTagForBaseClass, nil
}

func (g *GitTagBase) GetVersion() (version Version, err error) {
	parent, err := g.GetParentGitTagForBaseClass()
	if err != nil {
		return nil, err
	}

	name, err := parent.GetName()
	if err != nil {
		return nil, err
	}

	return Versions().GetNewVersionByString(name)
}

func (g *GitTagBase) MustGetParentGitTagForBaseClass() (parentGitTagForBaseClass GitTag) {
	parentGitTagForBaseClass, err := g.GetParentGitTagForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentGitTagForBaseClass
}

func (g *GitTagBase) MustGetVersion() (version Version) {
	version, err := g.GetVersion()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return version
}

func (g *GitTagBase) MustSetParentGitTagForBaseClass(parentGitTagForBaseClass GitTag) {
	err := g.SetParentGitTagForBaseClass(parentGitTagForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitTagBase) SetParentGitTagForBaseClass(parentGitTagForBaseClass GitTag) (err error) {
	g.parentGitTagForBaseClass = parentGitTagForBaseClass

	return nil
}
