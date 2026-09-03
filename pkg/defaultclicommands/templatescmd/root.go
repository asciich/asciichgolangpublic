package templatescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/templatescmd/gotemplatecmd"
)

func NewTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "handle templates",
	}

	cmd.AddCommand(
		gotemplatecmd.NewGoTemplateCmd(),
	)

	return cmd
}
