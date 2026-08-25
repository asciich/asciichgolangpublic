package publicipscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/publicips"
)

func NewGetPublicIpCmd() *cobra.Command {
	const short = "Get the current public IP address."

	cmd := &cobra.Command{
		Use:   "get-public-ip",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network publicips get-public-ip`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			ip := mustutils.Must(publicips.GetPublicIp(ctx))

			fmt.Println(ip)
		},
	}

	return cmd
}
