package asciichgolangpublic

type GitCommitOptions struct {
	Message    string
	AllowEmpty bool
	Verbose    bool
}

func NewGitCommitOptions() (g *GitCommitOptions) {
	return new(GitCommitOptions)
}

func (g *GitCommitOptions) GetAllowEmpty() (allowEmpty bool) {

	return g.AllowEmpty
}

func (g *GitCommitOptions) GetMessage() (message string, err error) {
	if g.Message == "" {
		return "", TracedErrorf("Message not set")
	}

	return g.Message, nil
}

func (g *GitCommitOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitCommitOptions) MustGetMessage() (message string) {
	message, err := g.GetMessage()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return message
}

func (g *GitCommitOptions) MustSetMessage(message string) {
	err := g.SetMessage(message)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitCommitOptions) SetAllowEmpty(allowEmpty bool) {
	g.AllowEmpty = allowEmpty
}

func (g *GitCommitOptions) SetMessage(message string) (err error) {
	if message == "" {
		return TracedErrorf("message is empty string")
	}

	g.Message = message

	return nil
}

func (g *GitCommitOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
