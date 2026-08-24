package aicmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/aidercmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/copilotcmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/ollamacmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/openhandscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/aicmd/vectordatabasecmd"
	"os"
)

func NewAICmd() *cobra.Command {
	const short = "Artificial inteligence related commands."

	cmd := &cobra.Command{
		Use:   "ai",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai ai`,
	}

	cmd.AddCommand(
		NewConcatFilesToKnowledgeFileCmd(),

		aidercmd.NewAiderCmd(),
		copilotcmd.NewCopilotCmd(),
		ollamacmd.NewOllamaCmd(),
		openhandscmd.NewOpenHandsCmd(),
		vectordatabasecmd.NewVectorDatabaseCmd(),
	)

	return cmd
}
