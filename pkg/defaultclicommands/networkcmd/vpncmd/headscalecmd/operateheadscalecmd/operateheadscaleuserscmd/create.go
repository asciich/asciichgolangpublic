package operateheadscaleuserscmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func NewCreateUserCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "Create a headscale user."

	cmd := &cobra.Command{
		Use:   "create",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate users create`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one username to create.")
			}

			userName := args[0]

			mustutils.Must0(options.GetHeadScale(ctx, cmd).CreateUser(ctx, userName))

			logging.LogGoodByCtxf(ctx, "Created headscale user '%s'.", userName)
		},
	}

	return cmd
}
