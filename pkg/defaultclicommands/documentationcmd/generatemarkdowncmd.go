package documentationcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/cobrautils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewGenerateMarkdownCmd(rootCmd *cobra.Command) *cobra.Command {
	const short = "Generate the markdown documentation for this binary. The output will be a single page markdown text document."

	cmd := &cobra.Command{
		Use:   "generate-markdown",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` documentation generate-markdown`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(
				mustutils.Must(cobrautils.GenerateMarkdownDocumentation(rootCmd.Use, rootCmd)),
			)
		},
	}

	return cmd
}
