package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateUserOptions struct {
	Name     string
	Username string
	Password string
	Email    string
}

func NewGitlabCreateUserOptions() (g *GitlabCreateUserOptions) {
	return new(GitlabCreateUserOptions)
}

func (g *GitlabCreateUserOptions) GetEmail() (email string, err error) {
	if g.Email == "" {
		return "", tracederrors.TracedErrorf("Email not set")
	}

	return g.Email, nil
}

func (g *GitlabCreateUserOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitlabCreateUserOptions) GetPassword() (password string, err error) {
	if g.Password == "" {
		return "", tracederrors.TracedErrorf("Password not set")
	}

	return g.Password, nil
}

func (g *GitlabCreateUserOptions) GetUsername() (username string, err error) {
	if g.Username == "" {
		return "", tracederrors.TracedErrorf("Username not set")
	}

	return g.Username, nil
}

func (g *GitlabCreateUserOptions) SetEmail(email string) (err error) {
	if email == "" {
		return tracederrors.TracedErrorf("email is empty string")
	}

	g.Email = email

	return nil
}

func (g *GitlabCreateUserOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateUserOptions) SetPassword(password string) (err error) {
	if password == "" {
		return tracederrors.TracedErrorf("password is empty string")
	}

	g.Password = password

	return nil
}

func (g *GitlabCreateUserOptions) SetUsername(username string) (err error) {
	if username == "" {
		return tracederrors.TracedErrorf("username is empty string")
	}

	g.Username = username

	return nil
}
