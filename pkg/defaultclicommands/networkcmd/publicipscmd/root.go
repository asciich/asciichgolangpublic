package publicipscmd

import "github.com/spf13/cobra"

func NewPublicIpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "public-ips",
		Short: "Commands related to public IPs.",
	}

	cmd.AddCommand(
		NewGetPublicIpCmd(),
	)

	return cmd
}