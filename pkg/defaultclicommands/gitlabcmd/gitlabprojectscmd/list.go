package gitlabprojectscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabcmdgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitlabutils/gitlaboptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListCmd(options gitlabcmdgeneric.GitlabCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List gitlab projects",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			gitlab := mustutils.Must(options.GetAuthenticatedGitlabInstance(ctx, cmd, args))
			fqdn := mustutils.Must(gitlab.GetFqdn())

			list := mustutils.Must(gitlab.ListProjectPaths(ctx, &gitlaboptions.GitlabgetProjectListOptions{}))
			for _, p := range list {
				fmt.Println(p)
			}

			logging.LogGoodByCtxf(ctx, "List projects in Gitlab '%s' finished, found %d projects.", fqdn, len(list))
		},
	}

	return cmd
}
