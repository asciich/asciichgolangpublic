package operateheadscalecmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewOperateCmd(options *OperateOptions) *cobra.Command {
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
		NewCreatePreauthKeyCmd(options),
		NewCreateUserCmd(options),
		NewGetUserIdCmd(options),
		NewListUsersCmd(options),
	)

	return cmd
}
