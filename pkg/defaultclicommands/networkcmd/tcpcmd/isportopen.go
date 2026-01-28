package tcpcmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils"
)

func NewIsPortOpenCmd() *cobra.Command {
	const short = "Check if a TCP --port on the given --host is open. Exit 0 means the port is open, 1 means not open."
	cmd := &cobra.Command{
		Use:   "is-port-open",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` network tcp is-port-open --host="localhost" --port=22
    `,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			hostname, err := cmd.Flags().GetString("host")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if hostname == "" {
				logging.LogFatal("Please specify --host.")
			}

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			if port <= 0 {
				logging.LogFatal("Please specify a valid --port.")
			}

			isOpen := mustutils.Must(netutils.IsTcpPortOpen(ctx, hostname, port))

			if isOpen {
				logging.LogGoodByCtxf(ctx, "TCP port %d on host '%s' is open.", port, hostname)
				os.Exit(0)
			}

			logging.LogErrorByCtxf(ctx, "TCP port %d on host '%s' is not open.", port, hostname)
			os.Exit(1)
		},
	}

	cmd.Flags().String("host", "", "The host to check.")
	cmd.Flags().Int("port", 0, "The port to check.")

	return cmd
}
