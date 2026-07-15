package shellcmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/errorutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/fzfutils"
)

func NewFzfCmd() *cobra.Command {
	const short = "Run embedded fzf (command line fuzzy finder)."

	cmd := &cobra.Command{
		Use:   "fzf",
		Short: short,
		Long: short + `

Reads lines from stdin or an input file and presents an interactive fuzzy finder.
The selected line(s) are written to stdout.

Examples:
  # Select a file to open
  ls | ` + os.Args[0] + ` shell fzf | xargs vim

  # Select a git branch to checkout
  git branch | ` + os.Args[0] + ` shell fzf | xargs git checkout

  # Select from an input file
  ` + os.Args[0] + ` shell fzf --input-file /path/to/list.txt | xargs echo
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			multi, err := cmd.Flags().GetBool("multi")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			inputFile, err := cmd.Flags().GetString("input-file")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			// Read items from input file or stdin
			var input []byte
			if inputFile == "" {
				input, err = io.ReadAll(os.Stdin)
				if err != nil {
					logging.LogFatalWithTrace(err)
				}
			} else {
				input, err = nativefiles.ReadAsBytes(ctx, inputFile)
				if err != nil {
					logging.LogFatalWithTrace(err)
				}
			}

			result, err := fzfutils.RunFuzzySearch(ctx, input, &fzfutils.SearchOptions{
				Multi: multi,
			})
			if err != nil {
				// Treat user abort (Ctrl+C) as a clean exit with code 130
				if errorutils.IsUserAbort(err) {
					os.Exit(130)
				}
				logging.LogGoErrorFatalWithTrace(err)
			}

			for _, r := range result {
				fmt.Println(r)
			}
		},
	}

	cmd.Flags().Bool("multi", false, "Enable multi-select (TAB to select multiple items)")
	cmd.Flags().String("input-file", "", "Read input items from a file instead of stdin")

	return cmd
}
