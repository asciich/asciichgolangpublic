package asciichgolangpublic

type GitlabTag struct {
	gitlabProject *GitlabProject
	name          string
}

func NewGitlabTag() (g *GitlabTag) {
	return new(GitlabTag)
}

func (g *GitlabTag) GetGitlabProject() (gitlabProject *GitlabProject, err error) {
	if g.gitlabProject == nil {
		return nil, TracedErrorf("gitlabProject not set")
	}

	return g.gitlabProject, nil
}

func (g *GitlabTag) GetName() (name string, err error) {
	if g.name == "" {
		return "", TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GitlabTag) IsVersionTag() (isVersionTag bool, err error) {
	tagName, err := g.GetName()
	if err != nil {
		return false, err
	}

	return Versions().IsVersionString(tagName), nil
}

func (g *GitlabTag) MustGetGitlabProject() (gitlabProject *GitlabProject) {
	gitlabProject, err := g.GetGitlabProject()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitlabProject
}

func (g *GitlabTag) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabTag) MustIsVersionTag() (isVersionTag bool) {
	isVersionTag, err := g.IsVersionTag()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isVersionTag
}

func (g *GitlabTag) MustSetGitlabProject(gitlabProject *GitlabProject) {
	err := g.SetGitlabProject(gitlabProject)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabTag) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabTag) SetGitlabProject(gitlabProject *GitlabProject) (err error) {
	if gitlabProject == nil {
		return TracedErrorf("gitlabProject is nil")
	}

	g.gitlabProject = gitlabProject

	return nil
}

func (g *GitlabTag) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
