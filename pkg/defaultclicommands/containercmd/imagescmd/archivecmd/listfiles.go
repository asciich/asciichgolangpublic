package archivecmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewListFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-files",
		Short: "List all files in an image archive.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one archive path.")
			}

			for _, f := range mustutils.Must(containerimagehandler.ListFilesInArchive(ctx, args[0])) {
				fmt.Println(f)
			}
		},
	}

	return cmd
}
