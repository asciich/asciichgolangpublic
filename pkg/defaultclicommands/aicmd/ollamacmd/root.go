package ollamacmd

import (
	"os"

	"github.com/spf13/cobra"
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
		NewRunProxyCmd(),
		NewSendPromptCmd(),
	)

	return cmd
}
