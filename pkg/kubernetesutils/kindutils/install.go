package kindutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const DEFAULT_KIND_VERSION = "v0.33.0"

// SHA256 checksums of the `kind-linux-amd64` binary per released version.
// Add new versions here after verifying the checksum against the official release.
var kindLinuxAmd64Sha256Sums = map[string]string{
	"v0.32.0": "50030de23cf40a18505f20426f6a8506bedf13c6e509244bd1fa9463721b0f54",
	"v0.33.0": "aee6151561422756b764a4ae28e7f44cda5af5a9eead3cc9985112b1de8d8e0d",
}

// Returns the kind version with a leading "v", e.g. "0.33.0" becomes "v0.33.0".
func NormalizeKindVersion(kindVersion string) (string, error) {
	ret := strings.TrimSpace(kindVersion)
	if ret == "" {
		return "", tracederrors.TracedErrorEmptyString("kindVersion")
	}

	if !strings.HasPrefix(ret, "v") {
		ret = "v" + ret
	}

	return ret, nil
}

// Returns the expected SHA256 checksum of the `kind-linux-amd64` binary for the given version.
func GetKindLinuxAmd64Sha256Sum(kindVersion string) (string, error) {
	if kindVersion == "" {
		return "", tracederrors.TracedErrorEmptyString("kindVersion")
	}

	ret, ok := kindLinuxAmd64Sha256Sums[kindVersion]
	if !ok {
		return "", tracederrors.TracedErrorf("No known sha256 checksum for kind version '%s'.", kindVersion)
	}

	return ret, nil
}

// Installs the `kind` binary as documented here:
//
//	https://kind.sigs.k8s.io/docs/user/quick-start/#installing-from-release-binaries
//
// If kindVersion is empty, DEFAULT_KIND_VERSION is used.
// The binary is installed to /bin/kind using sudo.
func InstallKind(ctx context.Context, kindVersion string) error {
	if kindVersion == "" {
		kindVersion = DEFAULT_KIND_VERSION
	}

	kindVersion, err := NormalizeKindVersion(kindVersion)
	if err != nil {
		return err
	}

	sha256Sum, err := GetKindLinuxAmd64Sha256Sum(kindVersion)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kind version '%s' started.", kindVersion)

	srcUrl := "https://kind.sigs.k8s.io/dl/" + kindVersion + "/kind-linux-amd64"

	err = installutils.Install(
		ctx,
		&installoptions.InstallOptions{
			SrcUrl:          srcUrl,
			InstallPath:     "/bin/kind",
			Mode:            "u=rwx,g=rx,o=rx",
			UseSudo:         true,
			ReplaceExisting: true,
			Sha256Sum:       sha256Sum,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install kind version '%s' finished.", kindVersion)

	return nil
}
