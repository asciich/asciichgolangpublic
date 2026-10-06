package hostnamecmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/nativehost"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewSetHostnameCmd() *cobra.Command {
	const short = "Set the hostname persistently"

	cmd := &cobra.Command{
		Use:   "set-hostname",
		Short: short,
		Long: short + `

Sets the system hostname persistently using native Go (no external commands).
Both the static hostname (persisted in /etc/hostname) and the transient
hostname (kernel, via sethostname(2)) are set, matching the behavior of
"hostnamectl set-hostname". The operation is idempotent - it will skip if
both hostnames are already set to the requested value.

Root privileges are required to set the hostname.

Usage:
    ` + os.Args[0] + ` hostname set-hostname <hostname>

Usage example:
    ` + os.Args[0] + ` hostname set-hostname --verbose myserver.example.com
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one hostname argument.")
			}

			hostname := args[0]

			if hostname == "" {
				logging.LogFatal("Hostname cannot be empty.")
			}

			sudo, err := cmd.Flags().GetBool("sudo")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			host := nativehost.NewNativeHost()

			mustutils.Must0(
				host.SetHostName(ctx, hostname, &hostsutilsoptions.SetHostnameOptions{
					UseSudo: sudo,
				}),
			)

			logging.LogGoodByCtxf(ctx, "Hostname set to '%s'.", hostname)
		},
	}

	cmd.Flags().Bool("sudo", false, "Use sudo for the hostnamectl set-hostname command")

	return cmd
}
