package gotemplatecmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/templateutils/gotemplateutils"
)

func NewRenderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render gotemplate.",
		Long:  "Render gotemplates.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.ContextVerbose()

			if len(args) != 1 {
				logging.LogFatal("Please specify exactly one template to to render")
			}

			filePathToRender := args[0]

			variables := map[string]interface{}{}
			vars, err := cmd.Flags().GetStringArray("var")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			for _, v := range vars {
				splitted := strings.SplitN(v, "=", 2)
				if len(splitted) != 2 {
					logging.LogFatalf("Unable to parse var='%s'", v)
				}

				variables[splitted[0]] = splitted[1]
			}

			rendered := mustutils.Must(gotemplateutils.RenderTemplateFromFilePathAsString(ctx, filePathToRender, variables))

			fmt.Print(rendered)

			logging.LogGoodByCtxf(ctx, "Rendered gotemplate '%s'", filePathToRender)
		},
	}

	return cmd
}
