package nativegitoo

import (
	"context"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type NativeGitCommit struct {
	hash       string
	repository *NativeGitRepository
}

func NewNativeGitCommit(hash string, repository *NativeGitRepository) *NativeGitCommit {
	return &NativeGitCommit{
		hash:       hash,
		repository: repository,
	}
}

func (c *NativeGitCommit) getCommitObject() (commitObject *object.Commit, err error) {
	path, err := c.repository.GetPath()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	hash := plumbing.NewHash(c.hash)
	commitObject, err = repo.CommitObject(hash)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get commit object for hash '%s' failed: %w", c.hash, err)
	}

	return commitObject, nil
}

func (c *NativeGitCommit) GetHash(ctx context.Context) (hash string, err error) {
	if c.hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}
	return c.hash, nil
}

func (c *NativeGitCommit) GetCommitMessage(ctx context.Context) (message string, err error) {
	commitObject, err := c.getCommitObject()
	if err != nil {
		return "", err
	}

	return commitObject.Message, nil
}

func (c *NativeGitCommit) GetAuthorString(ctx context.Context) (authorString string, err error) {
	commitObject, err := c.getCommitObject()
	if err != nil {
		return "", err
	}

	return commitObject.Author.String(), nil
}

func (c *NativeGitCommit) GetAuthorEmail(ctx context.Context) (authorEmail string, err error) {
	commitObject, err := c.getCommitObject()
	if err != nil {
		return "", err
	}

	return commitObject.Author.Email, nil
}

func (c *NativeGitCommit) GetAgeSeconds(ctx context.Context) (ageSeconds float64, err error) {
	commitObject, err := c.getCommitObject()
	if err != nil {
		return 0, err
	}

	age := time.Since(commitObject.Author.When)
	return age.Seconds(), nil
}

func (c *NativeGitCommit) GetParentCommits(ctx context.Context, options *parameteroptions.GitCommitGetParentsOptions) (parents []gitinterfaces.GitCommit, err error) {
	commitObject, err := c.getCommitObject()
	if err != nil {
		return nil, err
	}

	parents = []gitinterfaces.GitCommit{}
	for _, parentHash := range commitObject.ParentHashes {
		parents = append(parents, NewNativeGitCommit(parentHash.String(), c.repository))
	}

	return parents, nil
}

func (c *NativeGitCommit) HasParentCommit(ctx context.Context) (hasParent bool, err error) {
	commitObject, err := c.getCommitObject()
	if err != nil {
		return false, err
	}

	return len(commitObject.ParentHashes) > 0, nil
}

func (c *NativeGitCommit) CreateTag(ctx context.Context, options *gitparameteroptions.GitRepositoryCreateTagOptions) (gitinterfaces.GitTag, error) {
	return c.repository.CreateTag(ctx, options)
}

func (c *NativeGitCommit) ListTags(ctx context.Context) (tags []gitinterfaces.GitTag, err error) {
	return c.repository.ListTagsForCommitHash(ctx, c.hash)
}

func (c *NativeGitCommit) ListTagNames(ctx context.Context) (tagNames []string, err error) {
	tags, err := c.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	tagNames = []string{}
	for _, tag := range tags {
		name, err := tag.GetName()
		if err != nil {
			return nil, err
		}
		tagNames = append(tagNames, name)
	}

	return tagNames, nil
}

func (c *NativeGitCommit) ListVersionTagNames(ctx context.Context) (versionTagNames []string, err error) {
	tagNames, err := c.ListTagNames(ctx)
	if err != nil {
		return nil, err
	}

	versionTagNames = []string{}
	for _, tagName := range tagNames {
		_, err := versionutils.NewFromString(tagName)
		if err == nil {
			versionTagNames = append(versionTagNames, tagName)
		}
	}

	return versionTagNames, nil
}

func (c *NativeGitCommit) HasVersionTag(ctx context.Context) (hasVersionTag bool, err error) {
	versionTagNames, err := c.ListVersionTagNames(ctx)
	if err != nil {
		return false, err
	}

	return len(versionTagNames) > 0, nil
}

func (c *NativeGitCommit) GetNewestTagVersion(ctx context.Context) (newestVersion versionutils.Version, err error) {
	newestVersion, err = c.GetNewestTagVersionOrNilIfUnset(ctx)
	if err != nil {
		return nil, err
	}

	if newestVersion == nil {
		return nil, tracederrors.TracedErrorf("No version tag found for commit '%s'", c.hash)
	}

	return newestVersion, nil
}

func (c *NativeGitCommit) GetNewestTagVersionOrNilIfUnset(ctx context.Context) (newestVersion versionutils.Version, err error) {
	versionTagNames, err := c.ListVersionTagNames(ctx)
	if err != nil {
		return nil, err
	}

	if len(versionTagNames) == 0 {
		return nil, nil
	}

	for _, tagName := range versionTagNames {
		version, err := versionutils.NewFromString(tagName)
		if err != nil {
			return nil, err
		}

		isNewerThan, err := version.IsNewerThan(newestVersion)
		if err != nil {
			return nil, err
		}

		if newestVersion == nil || isNewerThan {
			newestVersion = version
		}
	}

	return newestVersion, nil
}

func (c *NativeGitCommit) GetNewestTagVersionString(ctx context.Context) (versionString string, err error) {
	newestVersion, err := c.GetNewestTagVersion(ctx)
	if err != nil {
		return "", err
	}

	return newestVersion.String(), nil
}
