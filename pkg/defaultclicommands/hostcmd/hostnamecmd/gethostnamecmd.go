package hostnamecmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/nativehost"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func NewGetHostnameCmd() *cobra.Command {
	const short = "Get the current hostname"

	cmd := &cobra.Command{
		Use:   "get-hostname",
		Short: short,
		Long: short + `

Retrieves and outputs the current static hostname of the system.

Usage:
    ` + os.Args[0] + ` hostname get-hostname

Usage example:
    ` + os.Args[0] + ` hostname get-hostname --verbose
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 0 {
				logging.LogFatal("This command does not accept any arguments.")
			}

			host := nativehost.NewNativeHost()

			hostname, err := host.RunCommandAndGetStdoutAsString(
				ctx,
				&parameteroptions.RunCommandOptions{
					Command: []string{"hostnamectl", "--static", "hostname"},
				},
			)
			mustutils.Must0(err)

			fmt.Println(strings.TrimSpace(hostname))
		},
	}

	return cmd
}
