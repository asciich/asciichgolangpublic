package admincmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"os"
)

func NewAdminCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "Minio admin related commands."

	cmd := &cobra.Command{
		Use:   "admin",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage s3 minio admin admin`,
	}

	cmd.AddCommand(
		NewCheckClusterHealthCmd(options),
	)

	return cmd
}
