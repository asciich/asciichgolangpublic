package aidercmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/aiutils/aiderutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
)

func NewRunAiderCmd() *cobra.Command {
	const short = "Runs Aider with a given message and files."

	cmd := &cobra.Command{
		Use:   "run-aider",
		Short: short,
		Long: short + `

Usage:
    ` + os.Args[0] + ` ai aider run-aider`,

		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			if len(args) == 0 {
				logging.LogFatal("Please provide a message as the first argument.")
			}

			message := args[0]
			files := args[1:]

			prompt := fmt.Sprintf("%s\nFiles:\n%s", message, strings.Join(files, "\n"))

			startTime := time.Now()

			mustutils.Must0(aiderutils.RunAider(ctx, prompt, files))

			elapsedTime := time.Since(startTime)

			logging.LogGoodByCtxf(ctx, "Running Aider with message and files finished in %s.", elapsedTime)
		},
	}

	return cmd
}
