package raspbianimagecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/embeddedsystems/raspberrypi/raspbian"
)

func NewDownloadCmd() *cobra.Command {
	const short = "Download the Raspberry Pi OS Lite image to --output-path."

	cmd := &cobra.Command{
		Use:   "download",
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			outputPath, err := cmd.Flags().GetString("output-path")
			if err != nil {
				return err
			}

			if outputPath == "" {
				return fmt.Errorf("flag '--output-path' is required but empty")
			}

			return raspbian.DownloadRaspbianImage(ctx, outputPath)
		},
	}

	cmd.Flags().String("output-path", "", "Path to save the downloaded Raspbian Lite image")

	return cmd
}
