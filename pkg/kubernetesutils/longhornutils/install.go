package longhornutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

const longhornctlVersion = "v1.13.0"

func InstallLonghornctl(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Install longhornctl started.")
	err := installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			SrcUrl:          "https://github.com/longhorn/cli/releases/download/" + longhornctlVersion + "/longhornctl-linux-amd64",
			Mode:            "u=rwx,g=rx,o=rx",
			InstallPath:     "/bin/longhornctl",
			UseSudo:         true,
			ReplaceExisting: true,
			Sha256Sum:       "2354ec1b9ff058fa4b7e178451d3445959cd4b6406317d7778b637f0314e1051",
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install longhornctl finished.")
	return nil
}
