package kindutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

// Installs the `kind` binary as documented here:
//
//	https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries
func InstallKind(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Install kind started.")

	err := installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			SrcUrl:          "https://kind.sigs.k8s.io/dl/v0.32.0/kind-linux-amd64",
			InstallPath:     "/bin/kind",
			Mode:            "u=rwx,g=rx,o=rx",
			UseSudo:         true,
			ReplaceExisting: true,
			Sha256Sum:       "50030de23cf40a18505f20426f6a8506bedf13c6e509244bd1fa9463721b0f54",
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install finished started.")

	return nil
}
