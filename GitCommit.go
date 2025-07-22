package asciichgolangpublic

import (
	"sort"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitCommit struct {
	gitRepo GitRepository
	hash    string
}

func NewGitCommit() (g *GitCommit) {
	return new(GitCommit)
}

func (g *GitCommit) CreateTag(options *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag GitTag, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	hash, err := g.GetHash()
	if err != nil {
		return nil, err
	}

	optionsToUse := options.GetDeepCopy()

	err = optionsToUse.SetCommitHash(hash)
	if err != nil {
		return nil, err
	}

	return repo.CreateTag(
		optionsToUse,
	)
}

func (g *GitCommit) GetAgeSeconds() (age float64, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return -1, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return -1, err
	}

	age, err = repo.GetCommitAgeSecondsByCommitHash(hash)
	if err != nil {
		return
	}

	return age, nil
}

func (g *GitCommit) GetAuthorEmail() (authorEmail string, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return "", err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return "", err
	}

	authorEmail, err = repo.GetAuthorEmailByCommitHash(hash)
	if err != nil {
		return
	}

	return authorEmail, nil
}

func (g *GitCommit) GetAuthorString() (authorString string, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return "", err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return "", err
	}

	authorString, err = repo.GetAuthorStringByCommitHash(hash)
	if err != nil {
		return
	}

	return authorString, nil
}

func (g *GitCommit) GetCommitMessage() (commitMessage string, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return "", err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return "", err
	}

	commitMessage, err = repo.GetCommitMessageByCommitHash(hash)
	if err != nil {
		return
	}

	return commitMessage, nil
}

func (g *GitCommit) GetGitRepo() (gitRepo GitRepository, err error) {
	if g.gitRepo == nil {
		return nil, tracederrors.TracedError("gitRepo not set")
	}
	return g.gitRepo, nil
}

func (g *GitCommit) GetHash() (hash string, err error) {
	if g.hash == "" {
		return "", tracederrors.TracedErrorf("hash not set")
	}

	return g.hash, nil
}

func (g *GitCommit) GetNewestTagVersionString(verbose bool) (string, error) {
	version, err := g.GetNewestTagVersion(verbose)
	if err != nil {
		return "", err
	}

	return version.GetAsString()
}

func (g *GitCommit) GetNewestTagVersion(verbose bool) (newestVersion versionutils.Version, err error) {
	newestVersion, err = g.GetNewestTagVersionOrNilIfUnset(verbose)
	if err != nil {
		return nil, err
	}

	if newestVersion == nil {
		hash, err := g.GetHash()
		if err != nil {
			return nil, err
		}

		path, hostDescription, err := g.GetRepoRootPathAndHostDescription()
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"no version tag found for commit '%s' in repository '%s' on host '%s'",
			hash,
			path,
			hostDescription,
		)
	}

	return newestVersion, nil
}

func (g *GitCommit) GetNewestTagVersionOrNilIfUnset(verbose bool) (newestVersion versionutils.Version, err error) {
	versions, err := g.ListVersionTagVersions(verbose)
	if err != nil {
		return nil, err
	}

	if len(versions) <= 0 {
		return nil, err
	}

	return versionutils.GetLatestVersionFromSlice(versions)
}

func (g *GitCommit) GetParentCommits(options *parameteroptions.GitCommitGetParentsOptions) (parentCommit []*GitCommit, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return nil, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	parentCommit, err = repo.GetCommitParentsByCommitHash(hash, options)
	if err != nil {
		return
	}

	return parentCommit, nil
}

func (g *GitCommit) GetRepoRootPathAndHostDescription() (repoRootPath string, hostDescription string, err error) {
	repo, err := g.GetGitRepo()
	if err != nil {
		return "", "", err
	}

	repoRootPath, err = repo.GetRootDirectoryPath(contextutils.ContextSilent())
	if err != nil {
		return "", "", err
	}

	hostDescription, err = repo.GetHostDescription()
	if err != nil {
		return "", "", err
	}

	return repoRootPath, hostDescription, nil
}

func (g *GitCommit) HasParentCommit() (hasParentCommit bool, err error) {
	hash, err := g.GetHash()
	if err != nil {
		return false, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return false, err
	}

	hasParentCommit, err = repo.CommitHasParentCommitByCommitHash(hash)
	if err != nil {
		return false, err
	}

	return hasParentCommit, nil
}

func (g *GitCommit) HasVersionTag(verbose bool) (hasVersionTag bool, err error) {
	tags, err := g.ListVersionTags(verbose)
	if err != nil {
		return false, err
	}

	return len(tags) > 0, nil
}

func (g *GitCommit) ListTagNames(verbose bool) (tagNames []string, err error) {
	tags, err := g.ListTags(verbose)
	if err != nil {
		return nil, err
	}

	tagNames = []string{}
	for _, t := range tags {
		toAdd, err := t.GetName()
		if err != nil {
			return nil, err
		}

		tagNames = append(tagNames, toAdd)
	}

	sort.Strings(tagNames)

	return tagNames, nil
}

func (g *GitCommit) ListTags(verbose bool) (tags []GitTag, err error) {
	repository, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	hash, err := g.GetHash()
	if err != nil {
		return nil, err
	}

	return repository.ListTagsForCommitHash(hash, verbose)
}

func (g *GitCommit) ListVersionTagNames(verbose bool) (tagNames []string, err error) {
	tags, err := g.ListVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	tagNames = []string{}
	for _, t := range tags {
		toAdd, err := t.GetName()
		if err != nil {
			return nil, err
		}

		tagNames = append(tagNames, toAdd)
	}

	tagNames, err = versionutils.SortStringSlice(tagNames)
	if err != nil {
		return nil, err
	}

	return tagNames, nil
}

func (g *GitCommit) ListVersionTagVersions(verbose bool) (versions []versionutils.Version, err error) {
	versionTags, err := g.ListVersionTags(verbose)
	if err != nil {
		return nil, err
	}

	versions = []versionutils.Version{}
	for _, v := range versionTags {
		toAdd, err := v.GetVersion()
		if err != nil {
			return nil, err
		}

		versions = append(versions, toAdd)
	}

	return versions, nil
}

func (g *GitCommit) ListVersionTags(verbose bool) (tags []GitTag, err error) {
	allTags, err := g.ListTags(verbose)
	if err != nil {
		return nil, err
	}

	tags = []GitTag{}
	for _, t := range allTags {
		isVersionTag, err := t.IsVersionTag()
		if err != nil {
			return nil, err
		}

		if isVersionTag {
			tags = append(tags, t)
		}
	}

	return tags, nil
}

func (g *GitCommit) SetGitRepo(gitRepo GitRepository) (err error) {
	if gitRepo == nil {
		return tracederrors.TracedErrorNil("gitRepo")
	}

	g.gitRepo = gitRepo

	return nil
}

func (g *GitCommit) SetHash(hash string) (err error) {
	if hash == "" {
		return tracederrors.TracedErrorf("hash is empty string")
	}

	g.hash = hash

	return nil
}
