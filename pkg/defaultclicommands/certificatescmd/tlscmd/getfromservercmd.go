package tlscmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/nativex509utils"
)

func NewGetFromWebserverCmd() *cobra.Command {
	const short = "Get TLS certificate from server by --hostname and --port."

	cmd := &cobra.Command{
		Use:   "get-from-server",
		Short: short,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			hostname, err := cmd.Flags().GetString("hostname")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if hostname == "" {
				logging.LogFatal("Please specify --hostname")
			}

			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if port <= 0 {
				logging.LogFatalf("Invalid port '%d'", port)
			}

			certs := mustutils.Must(nativex509utils.GetServersCertificateChain(ctx, hostname, port))

			for i := range len(certs) {
				fmt.Println(genericx509utils.GetInfoString(certs[len(certs)-1-i]))
			}
		},
	}

	cmd.Flags().String("hostname", "", "The hostname to get the certificate from.")
	cmd.Flags().Int("port", 443, "The port to establish the connection to and get the certificate. Default is 443/ SSL")

	return cmd
}
