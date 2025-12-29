package historycmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/bashutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func NewIncreaseSizeCmd() *cobra.Command {
	const increasedBashHistorySize int = 20000
	shortDescription := fmt.Sprintf("Enble bigger bash history size to '%d' entries", increasedBashHistorySize)

	cmd := &cobra.Command{
		Use:   "increase-size",
		Short: shortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)
			bashutils.SetBashHistorySizeOfCurrentUser(ctx, increasedBashHistorySize)

			logging.LogGoodByCtxf(ctx, "Increased bash history size to '%d'.", increasedBashHistorySize)
		},
	}

	return cmd
}
