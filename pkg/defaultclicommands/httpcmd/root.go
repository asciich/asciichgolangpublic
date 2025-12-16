package httpcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/clientcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/testwebservercmd"
)

func NewHttpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "HTTP/ Web server and client related commands.",
	}

	cmd.AddCommand(
		clientcmd.NewClientCmd(),
		testwebservercmd.NewTestWebServerCmd(),
	)

	return cmd
}
