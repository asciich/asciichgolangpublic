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

func NewShowInfoCmd() *cobra.Command {
	const short = "Show certificate info from a local file or stdin."

	cmd := &cobra.Command{
		Use:   "show-info",
		Short: short,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			filePath := args[0]

			if filePath == "" {
				logging.LogFatal("Please specify a file path or '-' for stdin")
			}

			cert := mustutils.Must(nativex509utils.ReadCertificateFromFileOrStdin(ctx, filePath))

			fmt.Println(genericx509utils.GetInfoString(cert))
		},
	}

	return cmd
}
