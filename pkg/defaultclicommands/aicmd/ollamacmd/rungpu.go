package ollamacmd

import (
	"os"

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
		Long: short + `

The GPU vendor is autodetected and the matching docker configuration is used:
  - AMD GPUs (ROCm) are preferred if detected. The dedicated 'ollama/ollama:rocm'
    image is used and the required device nodes '/dev/kfd' and '/dev/dri' are
    passed through to the container.
  - NVIDIA GPUs are used as a second option via the default 'ollama/ollama'
    image and docker's '--gpus all' flag.

If no supported GPU is detected the command fails with an error. Falling back to
CPU only mode is intentionally not done; use the dedicated CPU only command for
that instead.

The container is started in an idempotent way: if an ollama container is already
running nothing is changed. All run modes share the same volume mount
'ollama:/root/.ollama' so downloaded models are reused regardless of the mode.

The context length (OLLAMA_CONTEXT_LENGTH) can be adjusted with the
'--context-length' flag. When not set (or set to 0) the ollamautils default
of 32768 is used.

Usage:
    ` + os.Args[0] + ` ai ollama run-gpu
    ` + os.Args[0] + ` ai ollama run-gpu --context-length 8192`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			contextLength := mustutils.Must(cmd.Flags().GetInt("context-length"))

			mustutils.Must0(ollamautils.RunGPU(ctx, &ollamautils.RunOptions{
				ContextLength: contextLength,
			}))

			logging.LogGoodByCtxf(ctx, "Ollama with GPU support started.")
		},
	}

	cmd.Flags().Int(
		"context-length",
		0,
		"Context length (OLLAMA_CONTEXT_LENGTH) to use. When 0 the default of 32768 is used.",
	)

	return cmd
}
