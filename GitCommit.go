package asciichgolangpublic

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type GitCommit interface {
	CreateTag(ctx context.Context, options *gitparameteroptions.GitRepositoryCreateTagOptions) (GitTag, error)
	GetAgeSeconds(ctx context.Context) (float64, error)
	GetAuthorString(ctx context.Context) (string, error)
	GetAuthorEmail(ctx context.Context) (string, error)
	GetCommitMessage(ctx context.Context) (string, error)
	GetHash(ctx context.Context) (string, error)
	GetNewestTagVersion(ctx context.Context) (versionutils.Version, error)
	GetNewestTagVersionOrNilIfUnset(ctx context.Context) (versionutils.Version, error)
	GetNewestTagVersionString(ctx context.Context) (string, error)
	GetParentCommits(ctx context.Context, options *parameteroptions.GitCommitGetParentsOptions) ([]GitCommit, error)
	HasVersionTag(ctx context.Context) (bool, error)
	HasParentCommit(ctx context.Context) (bool, error)
	ListTagNames(ctx context.Context) ([]string, error)
	ListTags(ctx context.Context) ([]GitTag, error)
	ListVersionTagNames(ctx context.Context) ([]string, error)
}
