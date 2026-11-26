package gitlabcmdgeneric

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils/gitlabinterfaces"
)

type GitlabCmdOptions interface {
	GetGitlabInstance(ctx context.Context, cmd *cobra.Command, args []string) (gitlabinterfaces.GitlabInstance, error)
	GetAuthenticatedGitlabInstance(ctx context.Context, cmd *cobra.Command, args []string) (gitlabinterfaces.GitlabInstance, error)
}
