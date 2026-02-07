package operateheadscalecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

type OperateOptions struct {
	// override the root cmd cobra.Command fields. If unset the default values are taken.
	RootCmdUse   string
	RootCmdShort string

	// Get the HeadScale to operate:
	GetHeadScale func(cmd *cobra.Command) headscaleinterfaces.HeadScale
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
