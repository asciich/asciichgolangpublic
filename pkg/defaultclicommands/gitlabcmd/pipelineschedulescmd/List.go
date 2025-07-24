package pipelineschedulescmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/parameteroptions/authenticationoptions"
)

func NewListCommand() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "list",
		Short: "List scheduled pipelines for given project url",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly 1 url to list")
			}

			url := args[0]
			if url == "" {
				logging.LogFatal("Given --url is empty string")
			}

			listCmd(ctx, url)
		},
	}

	return cmd
}

func listCmd(ctx context.Context, url string) {
	access := []authenticationoptions.AuthenticationOption{
		&asciichgolangpublic.GitlabAuthenticationOptions{
			AccessToken: os.Getenv("GITLAB_ACCESS_TOKEN"),
			GitlabUrl:   url,
		},
	}
	project, err := asciichgolangpublic.GetGitlabProjectByUrlFromString(url, access, contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return
	}

	names := project.MustListScheduledPipelineNames(ctx)
	for _, n := range names {
		fmt.Println(n)
	}

	logging.LogGoodByCtxf(ctx, "Found '%d' scheduled pipelines in gitlab project %s", len(names), url)
}
