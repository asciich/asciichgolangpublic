package operateheadscaleuserscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/datatypes/stringsutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func NewCreatePreauthKeyCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "Create a preauth key for a the specified headscale user."

	cmd := &cobra.Command{
		Use:   "create-preauth-key",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate users create-preauth-key`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one user name to create a preauth key.")
			}

			userName := args[0]

			preauthKey := mustutils.Must(options.GetHeadScale(ctx, cmd).GeneratePreauthKeyForUser(ctx, userName))
			fmt.Print(stringsutils.EnsureEndsWithExactlyOneLineBreak(preauthKey))

			logging.LogGoodByCtxf(ctx, "Generated preauth key for user '%s'.", userName)
		},
	}

	return cmd
}
