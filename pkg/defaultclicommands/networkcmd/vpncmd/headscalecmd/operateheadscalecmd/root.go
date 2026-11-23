package operateheadscalecmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleapikeyscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscalenodescmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleuserscmd"
)

func NewOperateCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	if options == nil {
		panic("options is nil")
	}

	if options.GetHeadScale == nil {
		panic("options.GetHeadScale is nil")
	}

	const short = "Operate the HeadScale running on the localhost."

	cmd := &cobra.Command{
		Use:   options.GetRootCmdUse(),
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network vpn headscale operate`,
	}

	cmd.AddCommand(
		operateheadscaleapikeyscmd.NewApiKeysCmd(options),
		operateheadscalenodescmd.NewNodesCmd(options),
		operateheadscaleuserscmd.NewUsersCmd(options),
	)

	return cmd
}
