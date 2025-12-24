package ollamacmd

import "github.com/spf13/cobra"

func NewOllamaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ollama",
		Short: "ollama related commands",
	}

	cmd.AddCommand(
		NewDefaultPortCmd(),
		NewDescribeImageCmd(),
		NewOcrCmd(),
		NewRunCpuOnlyCmd(),
		NewSendPromptCmd(),
	)

	return cmd
}
