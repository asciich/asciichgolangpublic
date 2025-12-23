package ollamacmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunCpuOnlyCmd() *cobra.Command {
	const short = "Start ollama in a docker container without GPU support."

	cmd := &cobra.Command{
		Use:   "run-cpu-only",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			mustutils.Must0(ollamautils.RunCpuOnly(ctx))

			logging.LogGoodByCtxf(ctx, "Ollama in CPU only mode started.")
		},
	}

	return cmd
}
