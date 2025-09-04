package filescmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all files in the given path",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one path to list the files.")
			}

			path := args[0]

			for _, f := range mustutils.Must(nativefiles.ListFiles(ctx, path, &parameteroptions.ListFileOptions{})) {
				fmt.Println(f)
			}

			logging.LogGoodByCtxf(ctx, "List files in '%s' finished.", path)
		},
	}

	return cmd
}
