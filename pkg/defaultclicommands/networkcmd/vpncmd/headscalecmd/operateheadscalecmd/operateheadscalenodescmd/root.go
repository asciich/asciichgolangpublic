package operateheadscalenodescmd

import (
	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
)

func NewNodesCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Operate headscale nodes",
	}

	cmd.AddCommand(
		NewDeleteCmd(options),
		NewListCmd(options),
	)

	return cmd
}
