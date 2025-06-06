package gitparameteroptions

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitCommitOptions struct {
	// Message of the commit:
	Message string

	// Allow empty commit:
	AllowEmpty bool

	// Commit all changes, not only the added ones:
	CommitAllChanges bool

	// Enable verbose output
	Verbose bool
}

func NewGitCommitOptions() (g *GitCommitOptions) {
	return new(GitCommitOptions)
}

func (g *GitCommitOptions) GetAllowEmpty() (allowEmpty bool) {

	return g.AllowEmpty
}

func (g *GitCommitOptions) GetCommitAllChanges() (commitAllChanges bool) {

	return g.CommitAllChanges
}

func (g *GitCommitOptions) GetDeepCopy() (deepCopy *GitCommitOptions) {
	deepCopy = NewGitCommitOptions()

	*deepCopy = *g

	return deepCopy
}

func (g *GitCommitOptions) GetMessage() (message string, err error) {
	if g.Message == "" {
		return "", tracederrors.TracedErrorf("Message not set")
	}

	return g.Message, nil
}

func (g *GitCommitOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitCommitOptions) MustGetMessage() (message string) {
	message, err := g.GetMessage()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return message
}

func (g *GitCommitOptions) MustSetMessage(message string) {
	err := g.SetMessage(message)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitCommitOptions) SetAllowEmpty(allowEmpty bool) {
	g.AllowEmpty = allowEmpty
}

func (g *GitCommitOptions) SetCommitAllChanges(commitAllChanges bool) {
	g.CommitAllChanges = commitAllChanges
}

func (g *GitCommitOptions) SetMessage(message string) (err error) {
	if message == "" {
		return tracederrors.TracedErrorf("message is empty string")
	}

	g.Message = message

	return nil
}

func (g *GitCommitOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
