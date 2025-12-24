package ollamacmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewOcrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ocr",
		Short: "Uses an llm to detect the characters shown in the image. OCR means 'Optical Character Recognition'",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			model, err := cmd.Flags().GetString("model")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one image path to detect characters.")
			}

			imagePath := args[0]

			fmt.Println(mustutils.Must(ollamautils.OpticalCharacterRecognition(ctx, imagePath, &ollamautils.PromptOptions{
				ModelName: model,
			})))

			logging.LogGoodByCtxf(ctx, "Detecting characters in image '%s' finished.", imagePath)
		},
	}

	cmd.Flags().String("model", ollamautils.GetOCRModelName(), "LLM model used to detect the characters in the image.")

	return cmd
}
