package collectioncmd

import (
	"strings"

	"github.com/spf13/cobra"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/ansibleutils/ansiblegalaxyutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func NewCreateFileStructureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-file-structure",
		Short: "Create collection file structure in given directory path. Useful as a starting point to write your own collection.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			namespace, err := cmd.Flags().GetString("namespace")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if namespace == "" {
				logging.LogFatal("Please specify --namespace.")
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if name == "" {
				logging.LogFatal("Please specify --name.")
			}

			author, err := cmd.Flags().GetString("author")
			if err != nil {
				logging.LogGoErrorFatalWithTrace(err)
			}

			if author == "" {
				logging.LogFatal("Please specify --author.")
			}

			if len(args) != 1 {
				logging.LogFatal("Please specify exaclty one output directory path.")
			}

			path := strings.TrimSpace(args[0])
			if path == "" {
				logging.LogFatalf("'%s' is not a valid path", path)
			}

			mustutils.Must0(ansiblegalaxyutils.CreateFileStructure(
				ctx,
				path,
				&ansiblegalaxyutils.CreateCollectionFileStructureOptions{
					Namespace: namespace,
					Name:      name,
					Version:   "v0.1.0",
					Authors:   []string{author},
				},
			))

			logging.LogGoodByCtxf(ctx, "Created ansible galaxy collection file structure in '%s'.", path)
		},
	}

	cmd.PersistentFlags().String("namespace", "", "Namespace to use for the collection.")
	cmd.PersistentFlags().String("name", "", "Name to use for the collection.")
	cmd.PersistentFlags().String("author", "", "Author to use for the collection.")

	return cmd
}
