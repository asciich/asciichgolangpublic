package certificatescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/certificatescmd/tlscmd"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/certificatescmd/truststorecmd"
	"os"
)

func NewCertificatesCmd() *cobra.Command {
	const short = "Commmands related to certificates."

	cmd := &cobra.Command{
		Use:   "certificates",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` certificates certificates`,
	}

	cmd.AddCommand(
		tlscmd.NewTlsCmd(),
		truststorecmd.NewTrustStoreCmd(),
	)

	return cmd
}
