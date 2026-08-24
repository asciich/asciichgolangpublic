package ollamacmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
)

func NewDefaultPortCmd() *cobra.Command {
	const short = "Outputs the default port used to serve ollama"

	cmd := &cobra.Command{
		Use:   "default-port",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai ollama default-port`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%d\n", ollamautils.GetDefaultPort())
		},
	}

	return cmd
}
