package documentationcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/cobrautils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewGenerateMarkdownCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-markdown",
		Short: "Generate the markdown documentation for this binary. The output will be a single page markdown text document.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(
				mustutils.Must(cobrautils.GenerateMarkdownDocumentation(rootCmd.Use, rootCmd)),
			)
		},
	}

	return cmd
}
