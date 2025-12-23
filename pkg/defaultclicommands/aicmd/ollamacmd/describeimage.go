package ollamacmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewDescribeImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe-image",
		Short: "Uses an llm to describe whats shown on the given image.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			model, err := cmd.Flags().GetString("model")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one image path to describe.")
			}

			imagePath := args[0]

			fmt.Println(mustutils.Must(ollamautils.DescribeImage(ctx, imagePath, &ollamautils.PromptOptions{
				ModelName: model,
			})))

			logging.LogGoodByCtxf(ctx, "Describing image '%s' finished.", imagePath)
		},
	}

	cmd.Flags().String("model", ollamautils.GetImageProcessingModelName(), "LLM model used to describe the image.")

	return cmd
}
