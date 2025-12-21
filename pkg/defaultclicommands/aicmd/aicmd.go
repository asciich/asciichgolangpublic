package aicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/ollamacmd"
)

func NewAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ai",
		Short: "Artificial inteligence related commands.",
	}

	cmd.AddCommand(
		NewConcatFilesToKnowledgeFileCmd(),

		ollamacmd.NewOllamaCmd(),
	)

	return cmd
}