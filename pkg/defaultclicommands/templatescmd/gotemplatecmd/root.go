package gotemplatecmd

import "github.com/spf13/cobra"

func NewGoTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gotemplate",
		Short: "gotemplate related commands",
	}

	cmd.AddCommand(
		NewRenderCmd(),
	)

	return cmd
}
