package commandexecutortarutils

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// decompressIfGzip returns the raw tar bytes.
// If archiveBytes are gzip compressed (detected by magic bytes 0x1f 0x8b)
// they are decompressed first. This does not rely on a file extension since
// temp files may not have a .tar.gz suffix.
func decompressIfGzip(archiveBytes []byte) (tarBytes []byte, err error) {
	if archiveBytes == nil {
		return nil, tracederrors.TracedErrorNil("archiveBytes")
	}

	isGzip := len(archiveBytes) >= 2 && archiveBytes[0] == 0x1f && archiveBytes[1] == 0x8b
	if !isGzip {
		return archiveBytes, nil
	}

	gzReader, err := gzip.NewReader(bytes.NewReader(archiveBytes))
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarBytes, err = io.ReadAll(gzReader)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to decompress gzip archive: %w", err)
	}

	return tarBytes, nil
}

func ReadFileFromTarArchiveAsBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, archivePath string, fileName string) ([]byte, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if archivePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("archivePath")
	}

	if fileName == "" {
		return nil, tracederrors.TracedErrorEmptyString("fileName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Read '%s' from tar archive '%s' on '%s' started.", fileName, archivePath, hostDescription)

	archiveBytes, err := commandexecutorfile.ReadAsBytes(commandExecutor, archivePath)
	if err != nil {
		return nil, err
	}

	tarBytes, err := decompressIfGzip(archiveBytes)
	if err != nil {
		return nil, err
	}

	content, err := tarutils.ReadFileFromTarArchiveBytesAsBytes(tarBytes, fileName)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Read '%s' from tar archive '%s' on '%s' finished.", fileName, archivePath, hostDescription)

	return content, nil
}

func ReadFileFromTarArchiveAsString(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, archivePath string, fileName string) (string, error) {
	contentBytes, err := ReadFileFromTarArchiveAsBytes(ctx, commandExecutor, archivePath, fileName)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func ExtractFileFromTarArchive(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, archivePath string, fileName string, destPath string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if archivePath == "" {
		return tracederrors.TracedErrorEmptyString("archivePath")
	}

	if fileName == "" {
		return tracederrors.TracedErrorEmptyString("fileName")
	}

	if destPath == "" {
		return tracederrors.TracedErrorEmptyString("destPath")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Extract '%s' from tar archive '%s' to '%s' on '%s' started.", fileName, archivePath, destPath, hostDescription)

	content, err := ReadFileFromTarArchiveAsBytes(ctx, commandExecutor, archivePath, fileName)
	if err != nil {
		return err
	}

	err = commandexecutorfile.WriteBytes(ctx, commandExecutor, destPath, content, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Extract '%s' from tar archive '%s' to '%s' on '%s' finished.", fileName, archivePath, destPath, hostDescription)

	return nil
}
