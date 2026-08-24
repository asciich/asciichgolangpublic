package ollamacmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewOllamaCmd() *cobra.Command {
	const short = "ollama related commands"

	cmd := &cobra.Command{
		Use:   "ollama",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai ollama ollama`,
	}

	cmd.AddCommand(
		NewDefaultPortCmd(),
		NewDescribeImageCmd(),
		NewOcrCmd(),
		NewRunCpuOnlyCmd(),
		NewRunGpuCmd(),
		NewRunMcpAgentCmd(),
		NewSendPromptCmd(),
	)

	return cmd
}
