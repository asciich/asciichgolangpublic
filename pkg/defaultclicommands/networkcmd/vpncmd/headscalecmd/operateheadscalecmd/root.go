package operateheadscalecmd

import (
	"github.com/spf13/cobra"
)

func NewOperateCmd(options *OperateOptions) *cobra.Command {
	if options == nil {
		panic("options is nil")
	}

	if options.GetHeadScale == nil {
		panic("options.GetHeadScale is nil")
	}

	cmd := &cobra.Command{
		Use:   options.GetRootCmdUse(),
		Short: options.GetRootCmdShort(),
	}

	cmd.AddCommand(
		NewCreatePreauthKeyCmd(options),
		NewCreateUserCmd(options),
		NewGetUserIdCmd(options),
		NewListUsersCmd(options),
	)

	return cmd
}
