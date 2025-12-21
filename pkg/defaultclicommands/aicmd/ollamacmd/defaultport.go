package ollamacmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
)

func NewDefaultPortCmd() *cobra.Command {
	cmd :=&cobra.Command{
		Use: "default-port",
		Short: fmt.Sprintf("Outputs the default port (%d) used to serve ollama", ollamautils.GetDefaultPort()),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%d\n", ollamautils.GetDefaultPort())
		},
	}

	return cmd
}