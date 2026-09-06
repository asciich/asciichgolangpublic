package kubeadmutils

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
// installoptions.InstallOptions used to install kubeadm.
func prepareInstallOptions(options *InstallKubeadmOptions) (*installoptions.InstallOptions, error) {
	if options == nil {
		options = DefaultInstallKubeadmOptions()
	}

	// Set defaults for any unset options
	if options.InstallPath == "" {
		options.InstallPath = "/bin/kubeadm"
	}
	if options.Version == "" {
		options.Version = "v1.36.2"
	}
	// UseSudo defaults to true if not set

	sha256Sum := getSha256SumForVersion(options.Version)
	if sha256Sum == "" {
		return nil, fmt.Errorf("unsupported kubeadm version: %s", options.Version)
	}

	srcUrl := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubeadm", options.Version)

	return &installoptions.InstallOptions{
		SrcUrl:          srcUrl,
		InstallPath:     options.InstallPath,
		Mode:            "u=rwx,g=rx,o=rx",
		ReplaceExisting: true,
		UseSudo:         options.UseSudo,
		Sha256Sum:       sha256Sum,
	}, nil
}

func InstallKubeadm(ctx context.Context, options *InstallKubeadmOptions) error {
	installOptions, err := prepareInstallOptions(options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubeadm started.")

	err = installutils.Install(ctx, installOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubeadm finished.")

	return nil
}

func InstallKubeadmUsingCommandExecutor(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *InstallKubeadmOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	installOptions, err := prepareInstallOptions(options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubeadm started.")

	err = commandexecutorinstall.Install(ctx, commandExecutor, installOptions)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubeadm finished.")

	return nil
}

// getSha256SumForVersion returns the SHA256 checksum for a specific kubeadm version.
func getSha256SumForVersion(version string) string {
	checksums := map[string]string{
		"v1.36.2": "f95d1b345f97c673ce1c1347eb96bf2587850975155ae701cdfc629a5b575b70",
	}
	return checksums[version]
}
