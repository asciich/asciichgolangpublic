package gitlabcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabcmdgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewShowFqdnCmd(options gitlabcmdgeneric.GitlabCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-fqdn",
		Short: "Show the FQDN if the specified gitlab instance.",
		Run: func(cmd *cobra.Command, args []string) {
			gitlab := mustutils.Must(options.GetGitlabInstance(cmd, args))
			fmt.Println(mustutils.Must(gitlab.GetFqdn()))
		},
	}

	return cmd
}
