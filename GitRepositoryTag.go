package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitRepositoryTag struct {
	GitTagBase
	name          string
	gitRepository GitRepository
}

func GetGitRepositoryTagByName(tagName string) (g *GitRepositoryTag, err error) {
	if tagName == "" {
		return nil, tracederrors.TracedErrorEmptyString("tagName")
	}

	g = NewGitRepositoryTag()

	err = g.SetName(tagName)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func GetGitRepositoryTagByNameAndRepository(tagName string, gitRepository GitRepository) (g *GitRepositoryTag, err error) {
	if tagName == "" {
		return nil, tracederrors.TracedErrorEmptyString("tagName")
	}

	if gitRepository == nil {
		return nil, tracederrors.TracedErrorNil("gitRepository")
	}

	g, err = GetGitRepositoryTagByName(tagName)
	if err != nil {
		return nil, err
	}

	err = g.SetGitRepository(gitRepository)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func MustGetGitRepositoryTagByName(tagName string) (g *GitRepositoryTag) {
	g, err := GetGitRepositoryTagByName(tagName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return g
}

func MustGetGitRepositoryTagByNameAndRepository(tagName string, gitRepository GitRepository) (g *GitRepositoryTag) {
	g, err := GetGitRepositoryTagByNameAndRepository(tagName, gitRepository)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return g
}

func NewGitRepositoryTag() (g *GitRepositoryTag) {
	g = new(GitRepositoryTag)

	g.MustSetParentGitTagForBaseClass(g)

	return g
}

func (g *GitRepositoryTag) GetGitRepository() (gitRepository GitRepository, err error) {
	if g.gitRepository == nil {
		return nil, tracederrors.TracedErrorf("gitRepository not set")
	}

	return g.gitRepository, nil
}

func (g *GitRepositoryTag) GetHash() (hash string, err error) {
	repo, err := g.GetGitRepository()
	if err != nil {
		return "", err
	}

	name, err := g.GetName()
	if err != nil {
		return "", err
	}

	return repo.GetHashByTagName(name)
}

func (g *GitRepositoryTag) GetName() (name string, err error) {
	if g.name == "" {
		return "", tracederrors.TracedErrorf("name not set")
	}

	return g.name, nil
}

func (g *GitRepositoryTag) GetVersion() (version versionutils.Version, err error) {
	name, err := g.GetName()
	if err != nil {
		return nil, err
	}

	version, err = versionutils.GetVersionByString(name)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func (g *GitRepositoryTag) IsVersionTag() (isVersionTag bool, err error) {
	name, err := g.GetName()
	if err != nil {
		return false, err
	}

	return versionutils.Versions().IsVersionString(name), nil
}

func (g *GitRepositoryTag) MustGetGitRepository() (gitRepository GitRepository) {
	gitRepository, err := g.GetGitRepository()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitRepository
}

func (g *GitRepositoryTag) MustGetHash() (hash string) {
	hash, err := g.GetHash()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hash
}

func (g *GitRepositoryTag) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitRepositoryTag) MustGetVersion() (version versionutils.Version) {
	version, err := g.GetVersion()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return version
}

func (g *GitRepositoryTag) MustIsVersionTag() (isVersionTag bool) {
	isVersionTag, err := g.IsVersionTag()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isVersionTag
}

func (g *GitRepositoryTag) MustSetGitRepository(gitRepository GitRepository) {
	err := g.SetGitRepository(gitRepository)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryTag) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryTag) SetGitRepository(gitRepository GitRepository) (err error) {
	if gitRepository == nil {
		return tracederrors.TracedErrorf("gitRepository is nil")
	}

	g.gitRepository = gitRepository

	return nil
}

func (g *GitRepositoryTag) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.name = name

	return nil
}
