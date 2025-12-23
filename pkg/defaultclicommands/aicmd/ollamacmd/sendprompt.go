package ollamacmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/ollamautils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewSendPromptCmd() *cobra.Command {
	const short = "Send a prompt to a running ollama instance and return the response. The prompt is read from stdin."

	cmd := &cobra.Command{
		Use:   "send-prompt",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			prompt, err := io.ReadAll(os.Stdin)
			if err != nil {
				logging.LogFatalWithTracef("Failed to read prompt from stdin: %v", err)
			}

			if len(prompt) == 0 {
				logging.LogFatal("Empty prompt read from stdin")
			}

			model, err := cmd.Flags().GetString("model")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			imagePaths, err := cmd.Flags().GetStringArray("image")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			response := mustutils.Must(ollamautils.SendPrompt(
				ctx, string(prompt),
				&ollamautils.PromptOptions{
					ModelName:  model,
					ImagePaths: imagePaths,
				}),
			)
			fmt.Println(response)

			logging.LogGoodByCtxf(ctx, "Sending prompt to ollama and print response finished.")
		},
	}

	cmd.Flags().String("model", "", fmt.Sprintf(
		"Name of the LLM model to use. For a fast response use '%s'. For a good balanced response in quality and speed use '%s'.",
		ollamautils.GetFastModelName(),
		ollamautils.GetModerateSpeedModelName(),
	))

	cmd.Flags().StringArray("image", []string{}, "Images to send as part of the prompt. This implies a model able to process images is used like '"+ollamautils.GetImageProcessingModelName()+"'.")

	return cmd
}
