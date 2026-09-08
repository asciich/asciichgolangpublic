package ciliumutils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/commandexecutorinstall"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// prepareCliInstallOptions applies defaults to the given options and builds the
// installoptions.InstallOptions used to install the cilium CLI.
func prepareCliInstallOptions(options *InstallCiliumOptions) (*installoptions.InstallOptions, error) {
	if options == nil {
		options = DefaultInstallCiliumOptions()
	}

	// Set defaults for any unset options.
	if options.CliInstallPath == "" {
		options.CliInstallPath = "/usr/local/bin/cilium"
	}
	if options.CliVersion == "" {
		options.CliVersion = "v0.20.0"
	}
	// UseSudo defaults to true if not set.

	sha256Sum := getCliSha256SumForVersion(options.CliVersion)
	if sha256Sum == "" {
		return nil, fmt.Errorf("unsupported cilium CLI version: %s", options.CliVersion)
	}

	srcUrl := fmt.Sprintf(
		"https://github.com/cilium/cilium-cli/releases/download/%s/cilium-linux-amd64.tar.gz",
		options.CliVersion,
	)

	return &installoptions.InstallOptions{
		SrcUrl:          srcUrl,
		InstallPath:     options.CliInstallPath,
		Mode:            "u=rwx,g=rx,o=rx",
		SrcArchivePath:  "cilium",
		ReplaceExisting: true,
		UseSudo:         options.UseSudo,
		Sha256Sum:       sha256Sum,
	}, nil

}

// InstallCiliumCli installs the cilium CLI on the local machine.
func InstallCiliumCli(ctx context.Context, options *InstallCiliumOptions) error {
	installOptions, err := prepareCliInstallOptions(options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install cilium CLI started.")

	err = installutils.Install(ctx, installOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install cilium CLI finished.")

	return nil
}

// InstallCiliumCliUsingCommandExecutor installs the cilium CLI using the given
// command executor.
func InstallCiliumCliUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InstallCiliumOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	installOptions, err := prepareCliInstallOptions(options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install cilium CLI started.")

	err = commandexecutorinstall.Install(ctx, commandExecutor, installOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install cilium CLI finished.")

	return nil
}

func getCliSha256SumForVersion(version string) string {
	checksums := map[string]string{
		"v0.20.0": "8336dd43466badff49099e352567068f406e75f666760b128b525f114b4f7456",
	}
	return checksums[version]
}
