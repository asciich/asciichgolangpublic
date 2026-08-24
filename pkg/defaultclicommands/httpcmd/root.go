package httpcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/httpclientcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/testwebservercmd"
	"os"
)

func NewHttpCmd() *cobra.Command {
	const short = "HTTP/ Web server and client related commands."

	cmd := &cobra.Command{
		Use:   "http",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` http http`,
	}

	cmd.AddCommand(
		httpclientcmd.NewClientCmd(nil),
		testwebservercmd.NewTestWebServerCmd(),
	)

	return cmd
}
