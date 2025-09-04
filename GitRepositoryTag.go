package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GitRepositoryTag struct {
	GitTagBase
	name          string
	gitRepository gitinterfaces.GitRepository
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

func GetGitRepositoryTagByNameAndRepository(tagName string, gitRepository gitinterfaces.GitRepository) (g *GitRepositoryTag, err error) {
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

func MustGetGitRepositoryTagByNameAndRepository(tagName string, gitRepository gitinterfaces.GitRepository) (g *GitRepositoryTag) {
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

func (g *GitRepositoryTag) GetGitRepository() (gitRepository gitinterfaces.GitRepository, err error) {
	if g.gitRepository == nil {
		return nil, tracederrors.TracedErrorf("gitRepository not set")
	}

	return g.gitRepository, nil
}

func (g *GitRepositoryTag) GetHash(ctx context.Context) (hash string, err error) {
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

	version, err = versionutils.ReadFromString(name)
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

	return versionutils.IsVersionString(name), nil
}

func (g *GitRepositoryTag) SetGitRepository(gitRepository gitinterfaces.GitRepository) (err error) {
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
