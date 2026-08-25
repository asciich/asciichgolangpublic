package shellcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/shellcmd/bashcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/shellcmd/croncommandcmd"
	"os"
)

func NewShellCmd() *cobra.Command {
	const short = "Shell related commands."

	cmd := &cobra.Command{
		Use:   "shell",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` shell shell`,
	}

	cmd.AddCommand(
		NewFzfCmd(),

		bashcmd.NewBashCmd(),
		croncommandcmd.NewCronCommandCmd(),
	)

	return cmd
}
