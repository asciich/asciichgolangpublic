package publicipscmd

import (
	"github.com/spf13/cobra"
	"os"
)

func NewPublicIpsCmd() *cobra.Command {
	const short = "Commands related to public IPs."

	cmd := &cobra.Command{
		Use:   "public-ips",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network publicips public-ips`,
	}

	cmd.AddCommand(
		NewGetPublicIpCmd(),
	)

	return cmd
}
