package operateheadscaleuserscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/networkcmd/vpncmd/headscalecmd/operateheadscalecmd/operateheadscaleoptions"
)

func NewUsersCmd(options *operateheadscaleoptions.OperateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Operate headscale users",
	}

	cmd.AddCommand(
		NewCreateUserCmd(options),
		NewGetUserIdCmd(options),
		NewListCmd(options),
		NewCreatePreauthKeyCmd(options),
	)

	return cmd
}
