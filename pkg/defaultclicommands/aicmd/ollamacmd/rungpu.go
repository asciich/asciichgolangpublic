package ollamacmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunGpuCmd() *cobra.Command {
	const short = "Start ollama in a docker container with GPU support."

	cmd := &cobra.Command{
		Use:   "run-gpu",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(ollamautils.RunGPU(ctx))

			logging.LogGoodByCtxf(ctx, "Ollama with GPU support started.")
		},
	}

	return cmd
}
