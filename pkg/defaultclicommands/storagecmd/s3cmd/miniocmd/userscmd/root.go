package userscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"os"
)

func NewUsersCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	const short = "Users related commands."

	cmd := &cobra.Command{
		Use:   "users",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` storage s3 minio users`,
	}

	cmd.AddCommand(
		NewCreateUserCmd(options),
		NewDeleteUserCmd(options),
		NewListUsersCmd(options),
	)

	return cmd
}
