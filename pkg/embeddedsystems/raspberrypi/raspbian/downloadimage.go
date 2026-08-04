package raspbian

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/compressionutils/xzutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	raspbianImageURL         = "https://downloads.raspberrypi.org/raspios_lite_armhf/images/raspios_lite_armhf-2026-04-21/2026-04-21-raspios-trixie-armhf-lite.img.xz"
	raspbianArchiveSha256Sum = "f393b8bc3fc49aef49ddc5d5af124333002f34e4b23ede439789145e5280d210"
)

// DownloadRaspbianImage downloads and extracts the Raspberry Pi OS Lite image to 'outputPath'.
//
// This uses the Lite variant which is a minimal image without a desktop environment.
// The .img.xz archive is downloaded, verified against the expected sha256 checksum,
// decompressed to 'outputPath', and the archive is removed afterwards.
func DownloadRaspbianImage(ctx context.Context, outputPath string) error {
	if outputPath == "" {
		return tracederrors.TracedErrorEmptyString("outputPath")
	}

	logging.LogInfoByCtxf(ctx, "Download Raspbian Lite image to '%s' started.", outputPath)
	logging.LogInfoByCtxf(ctx, "Downloading from '%s'.", raspbianImageURL)
	logging.LogInfoByCtxf(ctx, "Expected sha256sum: '%s'.", raspbianArchiveSha256Sum)

	archivePath := outputPath + ".xz"

	_, err := httputils.DownloadAsFile(
		httpgeneric.WithDownloadProgressEveryNBytes(ctx, 10*1024*1024),
		&httpoptions.DownloadAsFileOptions{
			RequestOptions: &httpoptions.RequestOptions{
				Url: raspbianImageURL,
			},
			OutputPath:        archivePath,
			OverwriteExisting: true,
			Sha256Sum:         raspbianArchiveSha256Sum,
		},
	)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to download Raspbian Lite archive to '%s': %w", archivePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Download complete. Decompressing archive '%s' to '%s'.", archivePath, outputPath)

	err = xzutils.Decompress(ctx, archivePath, outputPath)
	if err != nil {
		os.Remove(archivePath)
		return tracederrors.TracedErrorf("Failed to decompress Raspbian Lite archive '%s' to '%s': %w", archivePath, outputPath, err)
	}

	err = os.Remove(archivePath)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to remove archive '%s' after decompression: %w", archivePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Archive '%s' removed after successful decompression.", archivePath)
	logging.LogInfoByCtxf(ctx, "Download Raspbian Lite image to '%s' finished.", outputPath)

	return nil
}
