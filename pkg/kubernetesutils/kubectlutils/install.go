package kubectlutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func InstallKubectl(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Install kubectl started.")

	err := installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			SrcUrl:          "https://dl.k8s.io/release/v1.36.2/bin/linux/amd64/kubectl",
			InstallPath:     "/bin/kubectl",
			Mode:            "u=rwx,g=rx,o=rx",
			ReplaceExisting: true,
			UseSudo:         true,
			Sha256Sum:       "1e9045ec32bea85da43de85f0065358529ea7c7a152eca78154fba5b58c27d82",
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kubectl finished.")

	return nil
}
