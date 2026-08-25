package bashcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/shellcmd/bashcmd/historycmd"
	"os"
)

func NewBashCmd() *cobra.Command {
	const short = "Bash related commands."

	cmd := &cobra.Command{
		Use:   "bash",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` shell bash bash`,
	}

	cmd.AddCommand(
		NewDefaultScriptStructureCmd(),

		historycmd.NewHistoryCmd(),
	)

	return cmd
}
