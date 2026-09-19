package filescmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewMkdirCmd() *cobra.Command {
	const short = "Ensure a directory exists"

	cmd := &cobra.Command{
		Use:   "mkdir",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` files mkdir <path>`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one path for the directory to create.")
			}

			path := args[0]

			mustutils.Must0( // Use the mustutils to fail on error ...
				nativefiles.CreateDirectory(ctx, path, &filesoptions.CreateOptions{}), // ... and reuse already existing functions.
			)

			logging.LogGoodByCtxf(ctx, "Directory '%s' ensured to exist.", path)
		},
	}

	return cmd
}
