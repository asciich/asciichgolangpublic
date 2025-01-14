package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type GitConfigSetOptions struct {
	Name    string
	Email   string
	Verbose bool
}

func NewGitConfigSetOptions() (g *GitConfigSetOptions) {
	return new(GitConfigSetOptions)
}

func (g *GitConfigSetOptions) GetEmail() (email string, err error) {
	if g.Email == "" {
		return "", errors.TracedErrorf("Email not set")
	}

	return g.Email, nil
}

func (g *GitConfigSetOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", errors.TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitConfigSetOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GitConfigSetOptions) IsEmailSet() (isSet bool) {
	return g.Email != ""
}

func (g *GitConfigSetOptions) IsNameSet() (isSet bool) {
	return g.Name != ""
}

func (g *GitConfigSetOptions) MustGetEmail() (email string) {
	email, err := g.GetEmail()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return email
}

func (g *GitConfigSetOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitConfigSetOptions) MustSetEmail(email string) {
	err := g.SetEmail(email)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitConfigSetOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitConfigSetOptions) SetEmail(email string) (err error) {
	if email == "" {
		return errors.TracedErrorf("email is empty string")
	}

	g.Email = email

	return nil
}

func (g *GitConfigSetOptions) SetName(name string) (err error) {
	if name == "" {
		return errors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitConfigSetOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
