package gitlabinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils/gitlaboptions"
)

type GitlabInstance interface {
	Authenticate(ctx context.Context, options *gitlaboptions.GitlabAuthenticationOptions) error
	GetFqdn() (string, error)
	ListProjectPaths(ctx context.Context, options *gitlaboptions.GitlabgetProjectListOptions) ([]string, error)
}
