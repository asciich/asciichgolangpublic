package certificatescmd

import (
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/defaultclicommands/certificatescmd/truststorecmd"
)

func NewCertificatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificates",
		Short: "Commmands related to certificates.",
	}

	cmd.AddCommand(
		truststorecmd.NewTrustStoreCmd(),
	)

	return cmd
}
