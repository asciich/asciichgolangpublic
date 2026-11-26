package gitlabprojectscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/gitlabcmd/gitlabcmdgeneric"
)

func NewProjectsCmd(options gitlabcmdgeneric.GitlabCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "Handle gitlab projects.",
	}

	cmd.AddCommand(
		NewListCmd(options),
	)

	return cmd
}
