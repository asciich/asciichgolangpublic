package tailscalecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/nativetailscalehttpclient"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/tailscaleutils/tailscaleoptions"
)

func NewHttpRequestCmd() *cobra.Command {
	const short = "Connect to tailscale and perform a HTTP request."

	cmd := &cobra.Command{
		Use:   "http-request",
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

			url, err := cmd.Flags().GetString("url")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if url == "" {
				logging.LogFatal("Please specify --url to request.")
			}

			response, _, cancel := mustutils.Must3(nativetailscalehttpclient.SendRequest(
				ctx,
				&tailscaleoptions.ConnectOptions{
					HostName:   hostname,
					PreAuthKey: preauthKey,
					ControlURL: controlUrl,
				},
				&httpoptions.RequestOptions{
					Url: url,
				},
			))
			defer cancel()
			fmt.Print(mustutils.Must(response.GetBodyAsString()))
			logging.LogGoodByCtxf(ctx, "Perform HTTP request to %s over tailscale network finished.", url)
		},
	}

	cmd.Flags().String("hostname", "", "The hostname to use for the client.")
	cmd.Flags().String("control-url", "", "Url of the control plane. Can include the port. E.g: https://headscale.example.com.1234")
	cmd.Flags().String("url", "", "Url to request.")

	return cmd
}
