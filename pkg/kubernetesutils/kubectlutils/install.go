package kubectlutils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func InstallKubectl(ctx context.Context, options *InstallKubectlOptions) error {
	if options == nil {
		options = DefaultInstallKubectlOptions()
	}

	// Set defaults for any unset options
	if options.InstallPath == "" {
		options.InstallPath = "/bin/kubectl"
	}
	if options.Version == "" {
		options.Version = "v1.36.2"
	}
	// UseSudo defaults to true if not set

	logging.LogInfoByCtxf(ctx, "Install kubectl started.")

	sha256Sum := getSha256SumForVersion(options.Version)
	if sha256Sum == "" {
		return fmt.Errorf("unsupported kubectl version: %s", options.Version)
	}

	srcUrl := fmt.Sprintf("https://dl.k8s.io/release/%s/bin/linux/amd64/kubectl", options.Version)

	err := installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			SrcUrl:          srcUrl,
			InstallPath:     options.InstallPath,
			Mode:            "u=rwx,g=rx,o=rx",
			ReplaceExisting: true,
			UseSudo:         options.UseSudo,
			Sha256Sum:       sha256Sum,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubectl finished.")

	return nil
}

// getSha256SumForVersion returns the SHA256 checksum for a specific kubectl version.
func getSha256SumForVersion(version string) string {
	checksums := map[string]string{
		"v1.36.2": "1e9045ec32bea85da43de85f0065358529ea7c7a152eca78154fba5b58c27d82",
	}
	return checksums[version]
}
