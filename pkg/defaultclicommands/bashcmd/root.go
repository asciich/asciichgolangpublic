package bashcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/bashcmd/historycmd"
)

func NewBashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bash",
		Short: "Bash related commands.",
	}

	cmd.AddCommand(
		NewDefaultScriptStructureCmd(),

		historycmd.NewHistoryCmd(),
	)

	return cmd
}
