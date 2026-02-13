package httpclientcmd

import (
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

	cmd := &cobra.Command{
		Use:   "client",
		Short: "HTTP client functions",
	}

	cmd.AddCommand(
		NewGetCmd(options),
		NewPerformRequestCmd(options),
	)

	return cmd
}

func defaultGetClient() httputilsinterfaces.Client {
	return httpnativeclientoo.GetNativeClient()
}
