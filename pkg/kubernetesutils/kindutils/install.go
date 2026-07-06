package kindutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

// Installs the `kind` binary as documented here:
//   https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries
func InstallKind(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Install kind started.")

	binaryFile, err := httputils.DownloadAsFile(ctx,
		&httpoptions.DownloadAsFileOptions{
			RequestOptions: &httpoptions.RequestOptions{
				Url: "https://kind.sigs.k8s.io/dl/v0.32.0/kind-linux-amd64",
			},
			OutputPath:        "/bin/kind",
			OverwriteExisting: true,
			UseSudo:           true,
		})
	if err != nil {
		return err
	}

	err = binaryFile.Chmod(ctx, &filesoptions.ChmodOptions{
		PermissionsString: "u=rwx,g=rx,o=rx",
		UseSudo: true,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install finished started.")

	return nil
}
