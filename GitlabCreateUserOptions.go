package asciichgolangpublic

type GitlabCreateUserOptions struct {
	Name     string
	Username string
	Password string
	Email    string
	Verbose  bool
}

func NewGitlabCreateUserOptions() (g *GitlabCreateUserOptions) {
	return new(GitlabCreateUserOptions)
}

func (g *GitlabCreateUserOptions) GetEmail() (email string, err error) {
	if g.Email == "" {
		return "", TracedErrorf("Email not set")
	}

	return g.Email, nil
}

func (g *GitlabCreateUserOptions) GetName() (name string, err error) {
	if g.Name == "" {
		return "", TracedErrorf("Name not set")
	}

	return g.Name, nil
}

func (g *GitlabCreateUserOptions) GetPassword() (password string, err error) {
	if g.Password == "" {
		return "", TracedErrorf("Password not set")
	}

	return g.Password, nil
}

func (g *GitlabCreateUserOptions) GetUsername() (username string, err error) {
	if g.Username == "" {
		return "", TracedErrorf("Username not set")
	}

	return g.Username, nil
}

func (g *GitlabCreateUserOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreateUserOptions) MustGetEmail() (email string) {
	email, err := g.GetEmail()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return email
}

func (g *GitlabCreateUserOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabCreateUserOptions) MustGetPassword() (password string) {
	password, err := g.GetPassword()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return password
}

func (g *GitlabCreateUserOptions) MustGetUsername() (username string) {
	username, err := g.GetUsername()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return username
}

func (g *GitlabCreateUserOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreateUserOptions) MustSetEmail(email string) {
	err := g.SetEmail(email)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateUserOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateUserOptions) MustSetPassword(password string) {
	err := g.SetPassword(password)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateUserOptions) MustSetUsername(username string) {
	err := g.SetUsername(username)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateUserOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateUserOptions) SetEmail(email string) (err error) {
	if email == "" {
		return TracedErrorf("email is empty string")
	}

	g.Email = email

	return nil
}

func (g *GitlabCreateUserOptions) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateUserOptions) SetPassword(password string) (err error) {
	if password == "" {
		return TracedErrorf("password is empty string")
	}

	g.Password = password

	return nil
}

func (g *GitlabCreateUserOptions) SetUsername(username string) (err error) {
	if username == "" {
		return TracedErrorf("username is empty string")
	}

	g.Username = username

	return nil
}

func (g *GitlabCreateUserOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}
