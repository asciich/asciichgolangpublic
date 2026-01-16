package gitgeneric
import (
	"context"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GenericGitCommit struct {
	gitRepo gitinterfaces.GitRepository
	hash    string
}

func NewGitCommit() (g *GenericGitCommit) {
	return new(GenericGitCommit)
}

func (g *GenericGitCommit) CreateTag(ctx context.Context, options *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag gitinterfaces.GitTag, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	hash, err := g.GetHash(ctx)
	if err != nil {
		return nil, err
	}

	optionsToUse := options.GetDeepCopy()

	err = optionsToUse.SetCommitHash(hash)
	if err != nil {
		return nil, err
	}

	return repo.CreateTag(ctx, optionsToUse)
}

func (g *GenericGitCommit) GetAgeSeconds(ctx context.Context) (age float64, err error) {
	hash, err := g.GetHash(ctx)
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

func (g *GenericGitCommit) GetAuthorEmail(ctx context.Context) (authorEmail string, err error) {
	hash, err := g.GetHash(ctx)
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

func (g *GenericGitCommit) GetAuthorString(ctx context.Context) (authorString string, err error) {
	hash, err := g.GetHash(ctx)
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

func (g *GenericGitCommit) GetCommitMessage(ctx context.Context) (commitMessage string, err error) {
	hash, err := g.GetHash(ctx)
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

func (g *GenericGitCommit) GetGitRepo() (gitRepo gitinterfaces.GitRepository, err error) {
	if g.gitRepo == nil {
		return nil, tracederrors.TracedError("gitRepo not set")
	}
	return g.gitRepo, nil
}

func (g *GenericGitCommit) GetHash(ctx context.Context) (hash string, err error) {
	if g.hash == "" {
		return "", tracederrors.TracedErrorf("hash not set")
	}

	return g.hash, nil
}

func (g *GenericGitCommit) GetNewestTagVersionString(ctx context.Context) (string, error) {
	version, err := g.GetNewestTagVersion(ctx)
	if err != nil {
		return "", err
	}

	return version.GetAsString()
}

func (g *GenericGitCommit) GetNewestTagVersion(ctx context.Context) (newestVersion versionutils.Version, err error) {
	newestVersion, err = g.GetNewestTagVersionOrNilIfUnset(ctx)
	if err != nil {
		return nil, err
	}

	if newestVersion == nil {
		hash, err := g.GetHash(ctx)
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

func (g *GenericGitCommit) GetNewestTagVersionOrNilIfUnset(ctx context.Context) (newestVersion versionutils.Version, err error) {
	versions, err := g.ListVersionTagVersions(ctx)
	if err != nil {
		return nil, err
	}

	if len(versions) <= 0 {
		return nil, err
	}

	return versionutils.GetLatestVersionFromSlice(versions)
}

func (g *GenericGitCommit) GetParentCommits(ctx context.Context, options *parameteroptions.GitCommitGetParentsOptions) ([]gitinterfaces.GitCommit, error) {
	hash, err := g.GetHash(ctx)
	if err != nil {
		return nil, err
	}

	repo, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	parentCommit, err := repo.GetCommitParentsByCommitHash(ctx, hash, options)
	if err != nil {
		return nil, err
	}

	return parentCommit, nil
}

func (g *GenericGitCommit) GetRepoRootPathAndHostDescription() (repoRootPath string, hostDescription string, err error) {
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

func (g *GenericGitCommit) HasParentCommit(ctx context.Context) (hasParentCommit bool, err error) {
	hash, err := g.GetHash(ctx)
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

func (g *GenericGitCommit) HasVersionTag(ctx context.Context) (hasVersionTag bool, err error) {
	tags, err := g.ListVersionTags(ctx)
	if err != nil {
		return false, err
	}

	return len(tags) > 0, nil
}

func (g *GenericGitCommit) ListTagNames(ctx context.Context) (tagNames []string, err error) {
	tags, err := g.ListTags(ctx)
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

func (g *GenericGitCommit) ListTags(ctx context.Context) (tags []gitinterfaces.GitTag, err error) {
	repository, err := g.GetGitRepo()
	if err != nil {
		return nil, err
	}

	hash, err := g.GetHash(ctx)
	if err != nil {
		return nil, err
	}

	return repository.ListTagsForCommitHash(ctx, hash)
}

func (g *GenericGitCommit) ListVersionTagNames(ctx context.Context) (tagNames []string, err error) {
	tags, err := g.ListVersionTags(ctx)
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

func (g *GenericGitCommit) ListVersionTagVersions(ctx context.Context) (versions []versionutils.Version, err error) {
	versionTags, err := g.ListVersionTags(ctx)
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

func (g *GenericGitCommit) ListVersionTags(ctx context.Context) (tags []gitinterfaces.GitTag, err error) {
	allTags, err := g.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	tags = []gitinterfaces.GitTag{}
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

func (g *GenericGitCommit) SetGitRepo(gitRepo gitinterfaces.GitRepository) (err error) {
	if gitRepo == nil {
		return tracederrors.TracedErrorNil("gitRepo")
	}

	g.gitRepo = gitRepo

	return nil
}

func (g *GenericGitCommit) SetHash(hash string) (err error) {
	if hash == "" {
		return tracederrors.TracedErrorf("hash is empty string")
	}

	g.hash = hash

	return nil
}
