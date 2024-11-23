package asciichgolangpublic

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
		return "", TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitlabCreateTagOptions) GetRef() (ref string, err error) {
	if g.Ref == "" {
		return "", TracedErrorf("Ref not set")
	}

	return g.Ref, nil
}

func (g *GitlabCreateTagOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitlabCreateTagOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabCreateTagOptions) MustGetRef() (ref string) {
	ref, err := g.GetRef()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return ref
}

func (g *GitlabCreateTagOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateTagOptions) MustSetRef(ref string) {
	err := g.SetRef(ref)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateTagOptions) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateTagOptions) SetRef(ref string) (err error) {
	if ref == "" {
		return TracedErrorf("ref is empty string")
	}

	g.Ref = ref

	return nil
}

func (g *GitlabCreateTagOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
