package userscmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/storagecmd/s3cmd/miniocmd/miniocmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
)

func NewCreateUserCmd(options *miniocmdoptions.MinioCmdOptions) *cobra.Command {
	var password string
	var readOnly bool
	var keepCurrentPasswordIfUserExists bool

	const short = "Create a new user."

	cmd := &cobra.Command{
		Use:   "create",
		Short: short,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			adminClient := options.GetAdminClient(ctx, cmd)

			userName := args[0]

			createOptions := &s3options.CreateUserOptions{
				KeepCurrentPasswordIfUserExists: keepCurrentPasswordIfUserExists,
				ReadOnly:                        readOnly,
			}

			mustutils.Must0(nativeminioclient.CreateUser(ctx, adminClient, userName, password, createOptions))

			logging.LogInfoByCtxf(ctx, "Create user '%s' finished.", userName)
		},
	}

	cmd.Flags().StringVar(&password, "password", "", "Password for the new user")
	cmd.Flags().BoolVar(&readOnly, "read-only", false, "Give the user read only permissions")
	cmd.Flags().BoolVar(&keepCurrentPasswordIfUserExists, "keep-current-password-if-user-exists", false, "Keep the current password if the user already exists")

	return cmd
}
