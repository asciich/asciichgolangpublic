package croncommandcmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/croncommand"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
)

func NewRunCmd() *cobra.Command {
	const short = "Runs the given command repeatedly based on a crontab schedule expression."

	cmd := &cobra.Command{
		Use:   "run",
		Short: short,
		Long: short + `

The cron expression can be either:
  - A standard 5-field crontab entry:  '*/5 * * * *'   (every 5 minutes)
  - A 6-field entry with seconds:      '*/10 * * * * *' (every 10 seconds)

Examples:
  cron-command --name my-job --cron '*/5 * * * *' --command 'echo hello' --verbose
  cron-command --name my-job --cron '*/5 * * * *' --command 'echo hello' --verbose --metrics-port=9123
  cron-command --name my-job --cron '0 9 * * 1-5' --command '/usr/local/bin/my-script.sh --flag value' --verbose
  cron-command --name my-job --cron '*/10 * * * * *' --command 'curl https://example.com/healthcheck' --verbose

The application blocks until the command exits non-successfully or the process is interrupted.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			cron, err := cmd.Flags().GetString("cron")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}
			commandStr, err := cmd.Flags().GetString("command")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			metricsPort, err := cmd.Flags().GetInt("metrics-port")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			// Split the command string into a slice to support arguments
			command := mustutils.Must(shelllinehandler.Split(commandStr))

			mustutils.Must0(croncommand.RunCronCommand(ctx, name, cron, command, metricsPort))
		},
	}

	cmd.Flags().String(
		"name",
		"",
		"Name of the cron job, used for logging and identification (e.g. 'my-backup-job')",
	)
	cmd.Flags().String(
		"cron",
		"",
		"Crontab expression defining the schedule (e.g. '*/5 * * * *' or '*/10 * * * * *' with seconds precision)",
	)
	cmd.Flags().String(
		"command",
		"",
		"Command to run periodically, including any arguments (e.g. '/usr/bin/my-script.sh --flag value')",
	)
	cmd.Flags().Int(
		"metrics-port",
		0,
		"If defined prometheus metrics are exposed on given port.",
	)

	mustutils.Must0(cmd.MarkFlagRequired("name"))
	mustutils.Must0(cmd.MarkFlagRequired("cron"))
	mustutils.Must0(cmd.MarkFlagRequired("command"))

	return cmd
}
