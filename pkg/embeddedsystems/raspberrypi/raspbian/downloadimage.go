package raspbian

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	raspbianImageURL  = "https://downloads.raspberrypi.org/raspios_lite_armhf/images/raspios_lite_armhf-2026-04-21/2026-04-21-raspios-trixie-armhf-lite.img.xz"
	raspbianSha256Sum = "f393b8bc3fc49aef49ddc5d5af124333002f34e4b23ede439789145e5280d210"
)

// DownloadRaspbianImage downloads the Raspberry Pi OS Lite image to 'outputPath'.
//
// This uses the Lite variant which is a minimal image without a desktop environment.
// The downloaded file is in .img.xz format and needs to be decompressed before use.
// The download is verified against the expected sha256 checksum.
func DownloadRaspbianImage(ctx context.Context, outputPath string) error {
	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Download Raspbian Lite image to '%s' started.", outputPath)
	logging.LogInfoByCtxf(ctx, "Downloading from '%s'.", raspbianImageURL)
	logging.LogInfoByCtxf(ctx, "Expected sha256sum: '%s'.", raspbianSha256Sum)

	_, err := httputils.DownloadAsFile(
		httpgeneric.WithDownloadProgressEveryNBytes(ctx, 10*1024*1024),
		&httpoptions.DownloadAsFileOptions{
			RequestOptions: &httpoptions.RequestOptions{
				Url: raspbianImageURL,
			},
			OutputPath:        outputPath,
			OverwriteExisting: true,
			Sha256Sum:         raspbianSha256Sum,
		},
	)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to download Raspbian Lite image to '%s': %w", outputPath, err)
	}

	logging.LogInfoByCtxf(ctx, "Download Raspbian Lite image to '%s' finished.", outputPath)

	return nil
}
