package aicmd

import "github.com/spf13/cobra"

func NewAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ai",
		Short: "Artificial inteligence related commands.",
	}

	cmd.AddCommand(
		NewConcatFilesToKnowledgeFileCmd(),
	)

	return cmd
}