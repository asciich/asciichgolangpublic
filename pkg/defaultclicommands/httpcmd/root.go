package httpcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/httpclientcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/testwebservercmd"
)

func NewHttpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP/ Web server and client related commands.",
	}

	cmd.AddCommand(
		httpclientcmd.NewClientCmd(nil),
		testwebservercmd.NewTestWebServerCmd(),
	)

	return cmd
}
