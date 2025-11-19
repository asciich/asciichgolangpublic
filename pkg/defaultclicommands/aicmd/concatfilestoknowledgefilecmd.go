package aicmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewConcatFilesToKnowledgeFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "concat-files-to-knowledge-file",
		Short: "Concat files in a directory to one knowledge file.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one directory to include the files.")
			}

			sourcePath := args[0]

			fmt.Println(mustutils.Must(aiutils.ConcatFilesToKnowledgeFile(ctx, sourcePath)))
		},
	}

	return cmd
}
