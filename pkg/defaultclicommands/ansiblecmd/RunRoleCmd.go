package ansiblecmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run-role",
		Short: "Run the ansible --role against the given --host.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			hostname, err := cmd.Flags().GetString("host")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if hostname == "" {
				logging.LogFatal("Please specify --host")
			}

			role, err := cmd.Flags().GetString("role")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if role == "" {
				logging.LogFatal("Please specify --role")
			}

			remoteUser, err := cmd.Flags().GetString("remote-user")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			vePath, err := cmd.Flags().GetString("virtualenv-path")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			keepTemporaryPlaybook, err := cmd.Flags().GetBool("keep-temporary-playbook")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			mustutils.Must0(ansibleutils.RunRoles(
				ctx,
				[]string{role},
				&ansibleparemeteroptions.RunOptions{
					Limit:                 hostname,
					AnsibleVirtualenvPath: vePath,
					KeepTemporaryPlaybook: keepTemporaryPlaybook,
					RemoteUser:            remoteUser,
				},
			))
		},
	}

	cmd.PersistentFlags().String("host", "", "Host to run the ansible role against.")
	cmd.PersistentFlags().String("role", "", "Name of the role to run.")
	cmd.PersistentFlags().String("virtualenv-path", "", "Path to the python virtualenv containing the ansible installation.")
	cmd.PersistentFlags().String("remote-user", "", "The remote user name to set in the playbook. E.g. '--remote-user=root' if you want ansible to connect to the --host as 'root' user.")
	cmd.PersistentFlags().Bool("keep-temporary-playbook", false, "Do not automatically delete the temporary used playbook. Useful for debugging.")

	return cmd
}
