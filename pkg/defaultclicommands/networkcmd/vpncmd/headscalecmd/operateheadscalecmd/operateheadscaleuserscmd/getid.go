package operateheadscaleuserscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func NewGetUserIdCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "Get the user ID of the given user name."

	cmd := &cobra.Command{
		Use:   "get-id",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate users get-id`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one user name to get the headscale ID.")
			}

			userName := args[0]

			id := mustutils.Must(options.GetHeadScale(ctx, cmd).GetUserId(ctx, userName))
			fmt.Printf("%d\n", id)

			logging.LogGoodByCtxf(ctx, "Got user id '%d' for headscale user '%s'.", id, userName)
		},
	}

	return cmd
}
