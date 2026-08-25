package truststorecmd

import (
	"os"
	"slices"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
)

func NewListCmd() *cobra.Command {
	const short = "list currently installed certificates"

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` certificates truststore list`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			certs := mustutils.Must(truststoreutils.ListCaCertificates(ctx))

			out := make([]string, 0, len(certs))
			for _, c := range certs {
				out = append(out, stringsutils.EnsureEndsWithExactlyOneLineBreak(mustutils.Must(x509utils.GetInfoSring(c))))
			}

			slices.Sort(out)

			for _, o := range out {
				print(o)
			}

			logging.LogGoodByCtxf(ctx, "Listed %d certificates.", len(certs))
		},
	}

	return cmd
}
