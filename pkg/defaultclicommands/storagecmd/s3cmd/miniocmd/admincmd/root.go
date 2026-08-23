package admincmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
)

func NewAdminCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "Minio admin related commands.",
	}

	cmd.AddCommand(
		NewCheckClusterHealthCmd(options),
	)

	return cmd
}
