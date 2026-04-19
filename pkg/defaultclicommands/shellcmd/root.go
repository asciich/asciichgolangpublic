package shellcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/shellcmd/bashcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/shellcmd/croncommandcmd"
)

func NewShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "Shell related commands.",
	}

	cmd.AddCommand(
		bashcmd.NewBashCmd(),
		croncommandcmd.NewCronCommandCmd(),
	)

	return cmd
}
