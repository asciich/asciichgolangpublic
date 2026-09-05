package kubeletutils

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

// prepareInstallOptions applies defaults to the given options and builds the
// installoptions.InstallOptions used to install kubelet.
func prepareInstallOptions(options *InstallKubeletOptions) (*installoptions.InstallOptions, error) {
	if options == nil {
		options = DefaultInstallKubeletOptions()
	}

	// Set defaults for any unset options
	if options.InstallPath == "" {
		options.InstallPath = "/bin/kubelet"
	}
	if options.Version == "" {
		options.Version = "v1.36.2"
	}
	// UseSudo defaults to true if not set

	sha256Sum := getSha256SumForVersion(options.Version)
	if sha256Sum == "" {
		return nil, fmt.Errorf("unsupported kubelet version: %s", options.Version)
	}

	srcUrl := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubelet", options.Version)

	return &installoptions.InstallOptions{
		SrcUrl:          srcUrl,
		InstallPath:     options.InstallPath,
		Mode:            "u=rwx,g=rx,o=rx",
		ReplaceExisting: true,
		UseSudo:         options.UseSudo,
		Sha256Sum:       sha256Sum,
	}, nil
}

func InstallKubelet(ctx context.Context, options *InstallKubeletOptions) error {
	installOptions, err := prepareInstallOptions(options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubelet started.")

	err = installutils.Install(ctx, installOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubelet finished.")

	return nil
}

func InstallKubeletUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InstallKubeletOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	installOptions, err := prepareInstallOptions(options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubelet started.")

	err = commandexecutorinstall.Install(ctx, commandExecutor, installOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubelet finished.")

	return nil
}

// getSha256SumForVersion returns the SHA256 checksum for a specific kubelet version.
func getSha256SumForVersion(version string) string {
	checksums := map[string]string{
		"v1.36.2": "3ba68cb4f0906053c08ab843131d5cec559cf8c718dbd679950a078b306abbb4",
	}
	return checksums[version]
}
