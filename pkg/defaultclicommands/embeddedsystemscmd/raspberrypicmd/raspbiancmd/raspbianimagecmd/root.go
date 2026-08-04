package raspbianimagecmd

import "github.com/spf13/cobra"

func NewRaspbianImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raspbian-image",
		Short: "Raspbian image related commands",
	}

	cmd.AddCommand(
		NewDownloadCmd(),
	)

	return cmd
}
