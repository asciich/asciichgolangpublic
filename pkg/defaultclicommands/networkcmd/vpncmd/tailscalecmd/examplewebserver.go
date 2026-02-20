package tailscalecmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscalehttpserver"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
	"github.com/asciich/asciichgolangpublic/pkg/userinteraction"
)

func NewExampleWebserverCmd() *cobra.Command {
	const short = "Connect a simple example webserver purely running in go to tailscale."

	cmd := &cobra.Command{
		Use:   "example-webserver",
		Short: short,
		Long: short + `

The webserver will be reachable under the current --hostname for other connected tailscale clients:
	http://hostname/hello

For the authentication the preauth key from the env var '` + tailscalegeneric.PREAUTH_KEY_ENV_VAR_NAME + `' is used.
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			hostname, err := cmd.Flags().GetString("hostname")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if hostname == "" {
				logging.LogFatal("Please specify --hostname")
			}

			logging.LogInfoByCtxf(ctx, "Going to connect to tailscale network as hostname '%s'.", hostname)

			controlUrl, err := cmd.Flags().GetString("control-url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if controlUrl == "" {
				logging.LogFatal("Please specify --control-url")
			}

			preauthKey := os.Getenv(tailscalegeneric.PREAUTH_KEY_ENV_VAR_NAME)
			if preauthKey == "" {
				logging.LogFatalf("Please export the PreAuth key as '%s'.", tailscalegeneric.PREAUTH_KEY_ENV_VAR_NAME)
			}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				logging.LogInfoByCtxf(ctx, "/ was requested.")
				fmt.Fprintf(w, "Hello from Tailscale example webserver server\n")
			})

			server := mustutils.Must(nativetailscalehttpserver.StartHttpServer(
				ctx,
				mux,
				80,
				&tailscaleoptions.ConnectOptions{
					HostName:   hostname,
					PreAuthKey: preauthKey,
					ControlURL: controlUrl,
				},
			))

			userinteraction.WaitUserAbortf("Tailscale webserver running and reachable as 'http://%s/' over tailscale.", hostname)

			mustutils.Must0(server.Close(ctx))
		},
	}

	cmd.Flags().String("hostname", "", "The hostname to use for the example webserver. The webserver will be reachalbe under this name by other tailscale connected clients.")
	cmd.Flags().String("control-url", "", "Url of the control plane. Can include the port. E.g: https://headscale.example.com.1234")

	return cmd
}
