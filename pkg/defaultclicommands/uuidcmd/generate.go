package uuidcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
)

func NewGenerateCmd() *cobra.Command {
	const short = "Generate UUID"

	cmd := &cobra.Command{
		Use:   "generate",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` uuid generate`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			fmt.Println(uuidutils.Generate(ctx))
		},
	}

	return cmd
}
