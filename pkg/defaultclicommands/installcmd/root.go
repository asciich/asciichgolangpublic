package installcmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbash"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hosts"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func NewInstallCmd() (cmd *cobra.Command) {
	const short = "Install this binary on the current system."

	cmd = &cobra.Command{
		Use:   "install",
		Short: short,
		Long: short + `

It's recommended to explicitly specify the --binary-name during the installation to avoid version numbers in the binary name after downloading it.

Usage:
    asciichgolangpublic install --verbose --binary-name=asciichgolangpublic.
`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := contextutils.GetVerbosityContextByCobraCmd(cmd)

			binaryName, err := cmd.Flags().GetString("binary-name")
			if err != nil {
				logging.LogGoErrorFatal(err)
			}

			cliInstall(ctx, binaryName)
		},
	}

	cmd.PersistentFlags().String("binary-name", "", "If explicitly specified this tool will be installed as '--binary-name'. Otherwise the basename is taken as to perform the installation.")

	return cmd
}

func cliInstall(ctx context.Context, binaryName string) {
	srcPath := os.Args[0]

	if srcPath == "" {
		logging.LogFatalWithTrace("binaryName is empty string after evaluation.")
	}

	if binaryName == "" {
		binaryName = filepath.Base(srcPath)
		logging.LogInfoByCtxf(ctx, "Binary name '%s' calculated by the srcPath='%s' is taken to install this binary.", binaryName, srcPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Explicit binary name '%s' is taken to install this binary", binaryName)
	}

	if binaryName == "" {
		logging.LogFatalWithTrace("binaryName is empty string after evaluation.")
	}

	mustutils.Must(hosts.MustGetLocalHost().InstallBinary(
		ctx,
		&parameteroptions.InstallOptions{
			SourcePath:       srcPath,
			InstallationPath: filepath.Join("/bin", binaryName),
			UseSudoToInstall: true,
			BinaryName:       binaryName,
		},
	))

	mustutils.Must(commandexecutorbash.RunOneLiner(ctx, fmt.Sprintf("%s completion bash | sudo tee /etc/bash_completion.d/%s > /dev/null", binaryName, binaryName)))
	logging.LogChangedByCtx(ctx, "Installed bash completion")

	logging.LogGoodByCtxf(ctx, "Sucessfully installed '%s'", binaryName)
}
