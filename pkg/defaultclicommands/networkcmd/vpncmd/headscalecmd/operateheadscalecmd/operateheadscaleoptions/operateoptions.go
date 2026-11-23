package operateheadscaleoptions

import (
	"context"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

type OperateOptions struct {
	// override the root cmd cobra.Command fields. If unset the default values are taken.
	RootCmdUse   string
	RootCmdShort string

	// Get the HeadScale to operate:
	GetHeadScale func(ctx context.Context, cmd *cobra.Command) headscaleinterfaces.Headscale
}

func (o *OperateOptions) GetRootCmdUse() string {
	if o.RootCmdUse == "" {
		return "operate"
	}

	return o.RootCmdUse
}

func (o *OperateOptions) GetRootCmdShort() string {
	if o.RootCmdUse == "" {
		return "Operate the HeadScale running on the localhost."
	}

	return o.RootCmdUse
}
