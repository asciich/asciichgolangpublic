package operateheadscalenodescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
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
