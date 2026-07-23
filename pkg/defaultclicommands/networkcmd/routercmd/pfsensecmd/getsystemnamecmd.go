package pfsensecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/routerutils/pfsenseutils"
)

func NewGetSystemNameCmd() *cobra.Command {
	const short = "Login to the pfSense router and get the system name"

	cmd := &cobra.Command{
		Use:   "get-system-name",
		Short: short,
		Long: short + `

Usage:
  PFSENSE_PASSWORD="<Your PASSWORD>" ` + os.Args[0] + ` network router pfsense get-system-name --url="https://192.168.1.1" --username=admin --verbose`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			username, err := cmd.Flags().GetString("username")
			if err != nil {
				logging.LogFatalWithTrace(err)
			}

			router := pfsenseutils.Router{
				Url:      url,
				UserName: username,
			}

			systemName := mustutils.Must(router.GetSystemName(ctx))
			fmt.Println(systemName)
		},
	}

	return cmd
}
