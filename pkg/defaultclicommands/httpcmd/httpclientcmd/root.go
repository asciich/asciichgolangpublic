package httpclientcmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/httpcmd/httpclientcmd/httpclientcmdoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpnativeclientoo"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
)

func NewClientCmd(options *httpclientcmdoptions.HttpClientCmdOptions) *cobra.Command {
	if options == nil {
		options = &httpclientcmdoptions.HttpClientCmdOptions{}
	}

	if options.GetClient == nil {
		options.GetClient = defaultGetClient
	}

	const short = "HTTP client functions"

	cmd := &cobra.Command{
		Use:   "client",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` http httpclient client`,
	}

	cmd.AddCommand(
		NewGetCmd(options),
		NewPerformRequestCmd(options),
	)

	return cmd
}

func defaultGetClient() httputilsinterfaces.Client {
	return httpnativeclientoo.NewNativeClient()
}

func normalizeUrl(url string) string {
	if !strings.Contains(url, "://") {
		return "https://" + url
	}
	return url
}
