package publicipscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/publicips"
)

func NewGetPublicIpCmd() *cobra.Command {
	short := "Get the current public IP address. When in a natted network this command can be used to get the public IP address of the router."

	cmd := &cobra.Command{
		Use: "get-public-ip",
		Short: short,
		Long: short + `

This function creates a webrequest to ` + publicips.GET_PUBLIC_IP_URL + ` which sends back the public client IP.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			ip := mustutils.Must(publicips.GetPublicIp(ctx))
			
			fmt.Println(ip)
		},
	}

	return cmd
}