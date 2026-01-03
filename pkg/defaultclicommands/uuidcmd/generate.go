package uuidcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/uuidutils"
)

func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "generate",
		Short: "Generate UUID",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			fmt.Println(uuidutils.Generate(ctx))
		},
	}

	return cmd
}