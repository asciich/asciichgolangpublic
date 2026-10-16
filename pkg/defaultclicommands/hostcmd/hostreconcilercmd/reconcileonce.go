package hostreconcilercmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostreconciler"
	"github.com/asciich/asciichgolangpublic/pkg/hostreconciler/hostreconcileroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewReconcileOnceCmd() *cobra.Command {
	const short = "Reconcile the current host once according to the given config directory"

	cmd := &cobra.Command{
		Use:   "reconcile-once",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` hostreconciler reconcile-once --config-dir <config-dir> [--docker-compose-force-recreate]

Flags:
    --docker-compose-force-recreate
        Recreate all docker compose containers even if their configuration
        and image haven't changed. Equivalent to 'docker compose up --force-recreate'.
        Note: With this flag the reconciliation is intentionally NOT idempotent.

Expected config directory structure:
    <config-dir>/
    └── host.yaml`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 0 {
				logging.LogFatal("No positional arguments expected. Please use '--config-dir' to specify the config directory.")
			}

			configDir, err := cmd.Flags().GetString("config-dir")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if configDir == "" {
				logging.LogFatal("Please specify the config directory using '--config-dir'.")
			}

			dockerComposeForceRecreate, err := cmd.Flags().GetBool("docker-compose-force-recreate")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			source := mustutils.Must(hostreconciler.NewDirectorySource(configDir))

			mustutils.Must0(hostreconciler.ReconcileOnce(ctx, source, &hostreconcileroptions.RunOptions{
				DockerComposeForceRecreate: dockerComposeForceRecreate,
			}))

			logging.LogGoodByCtxf(ctx, "Host reconciled according to config directory '%s'.", configDir)
		},
	}

	cmd.Flags().String("config-dir", "", "Path to the config directory containing the host.yaml")
	cmd.Flags().Bool("docker-compose-force-recreate", false, "Recreate all docker compose containers even if their configuration and image haven't changed (equivalent to 'docker compose up --force-recreate')")

	return cmd
}
