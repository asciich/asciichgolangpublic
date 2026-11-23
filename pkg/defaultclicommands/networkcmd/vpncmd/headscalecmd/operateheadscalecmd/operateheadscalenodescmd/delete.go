package operateheadscalenodescmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewDeleteCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	const short = "Delete a headscale node by name or ID."

	cmd := &cobra.Command{
		Use:   "delete <node-name|node-id>",
		Short: short,
		Long: short + `

The argument is interpreted as node ID if it is a number, otherwise as node name.
Deleting a node which does not exist is not an error.

Usage:
    ` + os.Args[0] + ` network vpn headscale operate nodes delete <node-name|node-id>

Usage examples:
    ` + os.Args[0] + ` network vpn headscale operate nodes delete my-node
    ` + os.Args[0] + ` network vpn headscale operate nodes delete 42`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one argument: the node name or the node ID.")
			}

			nodeIdentifier := args[0]
			if nodeIdentifier == "" {
				logging.LogFatal("The node name or node ID must not be empty.")
			}

			headscale := options.GetHeadScale(ctx, cmd)

			id, err := strconv.ParseUint(nodeIdentifier, 10, 64)
			if err == nil {
				mustutils.Must0(headscale.DeleteNodeById(ctx, id))
				logging.LogGoodByCtxf(ctx, "Headscale node with ID '%d' deleted.", id)
			} else {
				mustutils.Must0(headscale.DeleteNodeByName(ctx, nodeIdentifier))
				logging.LogGoodByCtxf(ctx, "Headscale node '%s' deleted.", nodeIdentifier)
			}
		},
	}

	return cmd
}
