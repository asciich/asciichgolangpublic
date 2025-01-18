package files

import (
	"os"
	"time"

	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A File represents any kind of file regardless if a local file or a remote file.
type File interface {
	AppendBytes(toWrite []byte, verbose bool) (err error)
	AppendString(toWrite string, verbose bool) (err error)
	Chmod(options *parameteroptions.ChmodOptions) (err error)
	Chown(options *parameteroptions.ChownOptions) (err error)
	CopyToFile(destFile File, verbose bool) (err error)
	Create(verbose bool) (err error)
	Delete(verbose bool) (err error)
	Exists(verbose bool) (exists bool, err error)
	GetBaseName() (baseName string, err error)
	GetDeepCopy() (deepCopy File)
	GetHostDescription() (hostDescription string, err error)
	GetLocalPath() (localPath string, err error)
	GetLocalPathOrEmptyStringIfUnset() (localPath string, err error)
	GetParentDirectory() (parentDirectory Directory, err error)
	GetPath() (path string, err error)
	GetSizeBytes() (fileSize int64, err error)
	GetUriAsString() (uri string, err error)
	MoveToPath(destPath string, useSudo bool, verbose bool) (movedFile File, err error)
	MustAppendBytes(toWrtie []byte, verbose bool)
	MustAppendString(toWrtie string, verbose bool)
	MustChmod(options *parameteroptions.ChmodOptions)
	MustChown(options *parameteroptions.ChownOptions)
	MustCopyToFile(destFile File, verbose bool)
	MustCreate(verbose bool)
	MustDelete(verbose bool)
	MustExists(verbose bool) (exists bool)
	MustGetBaseName() (baseName string)
	MustGetHostDescription() (hostDescription string)
	MustGetLocalPath() (localPath string)
	MustGetLocalPathOrEmptyStringIfUnset() (localPath string)
	MustGetPath() (path string)
	MustGetParentDirectory() (parentDirectory Directory)
	MustGetSizeBytes() (fileSize int64)
	MustGetUriAsString() (uri string)
	MustMoveToPath(destPath string, useSudo bool, verbose bool) (movedFile File)
	MustReadAsBytes() (content []byte)
	MustSecurelyDelete(verbose bool)
	MustTruncate(newSizeBytes int64, verbose bool)
	MustWriteBytes(toWrite []byte, verbose bool)
	ReadAsBytes() (content []byte, err error)
	SecurelyDelete(verbose bool) (err error)
	Truncate(newSizeBytes int64, verbose bool) (err error)
	WriteBytes(toWrite []byte, verbose bool) (err error)

	// All methods below this line can be implemented by embedding the `FileBase` struct:
	AppendLine(line string, verbose bool) (err error)
	CheckIsLocalFile(verbose bool) (err error)
	ContainsLine(line string) (containsLine bool, err error)
	CreateParentDirectory(verbose bool) (err error)
	EnsureLineInFile(line string, verbose bool) (err error)
	EnsureEndsWithLineBreak(verbose bool) (err error)
	GetCreationDateByFileName(verbose bool) (creationDate *time.Time, err error)
	GetFileTypeDescription(verbose bool) (fileTypeDescription string, err error)
	GetMimeType(verbose bool) (mimeType string, err error)
	GetNumberOfLinesWithPrefix(prefix string, trimLines bool) (nLines int, err error)
	GetNumberOfNonEmptyLines() (nLines int, err error)
	GetParentDirectoryPath() (parentDirectoryPath string, err error)
	GetPathAndHostDescription() (path string, hostDescription string, err error)
	GetSha256Sum() (sha256sum string, err error)
	GetTextBlocks(verbose bool) (textBlocks []string, err error)
	GetValueAsInt(key string) (value int, err error)
	GetValueAsString(key string) (value string, err error)
	IsContentEqualByComparingSha256Sum(other File, verbose bool) (isMatching bool, err error)
	IsEmptyFile() (isEmpty bool, err error)
	IsLocalFile(verbose bool) (isLocalFile bool, err error)
	IsMatchingSha256Sum(sha256sum string) (isMatching bool, err error)
	IsPgpEncrypted(verbose bool) (isPgpEncrypted bool, err error)
	IsYYYYmmdd_HHMMSSPrefix() (hasDatePrefix bool, err error)
	MustAppendLine(line string, verbose bool)
	MustCheckIsLocalFile(verbose bool)
	MustContainsLine(line string) (containsLine bool)
	MustCreateParentDirectory(verbose bool)
	MustEnsureEndsWithLineBreak(verbose bool)
	MustEnsureLineInFile(line string, verbose bool)
	MustGetCreationDateByFileName(verbose bool) (creationDate *time.Time)
	MustGetFileTypeDescription(verbose bool) (fileTypeDescription string)
	MustGetMimeType(verbose bool) (mimeType string)
	MustGetNumberOfLinesWithPrefix(prefix string, trimLines bool) (nLines int)
	MustGetNumberOfNonEmptyLines() (nLines int)
	MustGetParentDirectoryPath() (parentDirectoryPath string)
	MustGetPathAndHostDescription() (path string, hostDescription string)
	MustGetSha256Sum() (sha256sum string)
	MustGetTextBlocks(verbose bool) (textBlocks []string)
	MustGetValueAsInt(key string) (value int)
	MustGetValueAsString(key string) (value string)
	MustIsContentEqualByComparingSha256Sum(other File, verbose bool) (isMatching bool)
	MustIsEmptyFile() (isEmpty bool)
	MustIsLocalFile(verbose bool) (isLocalFile bool)
	MustIsMatchingSha256Sum(sha256sum string) (isMatching bool)
	MustIsPgpEncrypted(verbose bool) (isPgpEncrypted bool)
	MustIsYYYYmmdd_HHMMSSPrefix() (hasDatePrefix bool)
	MustPrintContentOnStdout()
	MustReadAsBool() (content bool)
	MustReadAsFloat64() (content float64)
	MustReadAsInt() (content int)
	MustReadAsInt64() (content int64)
	MustReadAsLines() (contentLines []string)
	MustReadAsLinesWithoutComments() (contentLines []string)
	MustReadAsTimeTime() (time *time.Time)
	MustReadAsString() (content string)
	MustReadFirstLine() (firstLine string)
	MustReadFirstLineAndTrimSpace() (firstLine string)
	MustReadLastCharAsString() (lastChar string)
	MustReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte)
	MustRemoveLinesWithPrefix(prefix string, verbose bool)
	MustReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *changesummary.ChangeSummary)
	MustSortBlocksInFile(verbose bool)
	MustTrimSpacesAtBeginningOfFile(verbose bool)
	MustWriteInt64(toWrite int64, verbose bool)
	MustWriteLines(linesToWrite []string, verbose bool)
	MustWriteString(content string, verbose bool)
	MustWriteTextBlocks(textBlocks []string, verose bool)
	PrintContentOnStdout() (err error)
	ReadAsBool() (content bool, err error)
	ReadAsFloat64() (content float64, err error)
	ReadAsInt() (content int, err error)
	ReadAsInt64() (content int64, err error)
	ReadAsLines() (contentLines []string, err error)
	ReadAsLinesWithoutComments() (contentLines []string, err error)
	ReadAsString() (content string, err error)
	ReadAsTimeTime() (time *time.Time, err error)
	ReadFirstLine() (firstLine string, err error)
	ReadFirstNBytes(numberOfBytesToRead int) (firstBytes []byte, err error)
	ReadFirstLineAndTrimSpace() (firstLine string, err error)
	ReadLastCharAsString() (lastChar string, err error)
	RemoveLinesWithPrefix(prefix string, verbose bool) (err error)
	ReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *changesummary.ChangeSummary, err error)
	SortBlocksInFile(verbose bool) (err error)
	TrimSpacesAtBeginningOfFile(verbose bool) (err error)
	WriteInt64(toWrite int64, verboe bool) (err error)
	WriteLines(linesToWrite []string, verbose bool) (err error)
	WriteString(content string, verbose bool) (err error)
	WriteTextBlocks(textBlocks []string, verbose bool) (err error)
}

func GetFileByOsFile(osFile *os.File) (file File, err error) {
	if osFile == nil {
		return nil, tracederrors.TracedError("osFile is nil")
	}

	file, err = NewLocalFileByPath(osFile.Name())
	if err != nil {
		return nil, err
	}

	return file, nil
}

func MustGetFileByOsFile(osFile *os.File) (file File) {
	file, err := GetFileByOsFile(osFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return file
}