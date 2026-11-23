package operateheadscaleapikeyscmd

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
)

func NewApiKeysCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "Operate headscale API keys."

	cmd := &cobra.Command{
		Use:   "apikeys",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate apikeys <list> [flags]`,
	}

	cmd.AddCommand(
		NewListCmd(options),
	)

	return cmd
}
