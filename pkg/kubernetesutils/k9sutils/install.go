package k9sutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func InstallK9s(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Install k9s started.")
	err := installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			SrcUrl:          "https://github.com/derailed/k9s/releases/download/v0.51.0/k9s_Linux_amd64.tar.gz",
			SrcArchivePath:  "k9s",
			Mode:            "u=rwx,g=rx,o=rx",
			InstallPath:     "/bin/k9s",
			UseSudo:         true,
			ReplaceExisting: true,
			Sha256Sum:       "da2a7b4844204e5f31da0e57536e85a848f4d2260cd58737402bbb04ce4c99a2",
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install k9s finished.")
	return nil
}
