package asciichgolangpublic

import (
	"os"
)

// A File represents any kind of file regardless if a local file or a remote file.
type File interface {
	AppendBytes(toWrite []byte, verbose bool) (err error)
	AppendString(toWrite string, verbose bool) (err error)
	CopyToFile(destFile File, verbose bool) (err error)
	Create(verbose bool) (err error)
	Delete(verbose bool) (err error)
	Exists() (exists bool, err error)
	GetBaseName() (baseName string, err error)
	GetDeepCopy() (deepCopy File)
	GetLocalPath() (localPath string, err error)
	GetLocalPathOrEmptyStringIfUnset() (localPath string)
	GetParentDirectory() (parentDirectory Directory, err error)
	GetUriAsString() (uri string, err error)
	MustAppendBytes(toWrtie []byte, verbose bool)
	MustAppendString(toWrtie string, verbose bool)
	MustCopyToFile(destFile File, verbose bool)
	MustCreate(verbose bool)
	MustDelete(verbose bool)
	MustExists() (exists bool)
	MustGetBaseName() (baseName string)
	MustGetLocalPath() (localPath string)
	MustGetParentDirectory() (parentDirectory Directory)
	MustGetUriAsString() (uri string)
	MustPrintContentOnStdout()
	MustReadAsBytes() (content []byte)
	MustSecurelyDelete(verbose bool)
	MustWriteBytes(toWrite []byte, verbose bool)
	PrintContentOnStdout() (err error)
	ReadAsBytes() (content []byte, err error)
	SecurelyDelete(verbose bool) (err error)
	WriteBytes(toWrite []byte, verbose bool) (err error)

	// All methods below this line can be implemented by embedding the `FileBase` struct:
	GetSha256Sum() (sha256sum string, err error)
	GetTextBlocks(verbose bool) (textBlocks []string, err error)
	IsContentEqualByComparingSha256Sum(other File, verbose bool) (isMatching bool, err error)
	IsMatchingSha256Sum(sha256sum string) (isMatching bool, err error)
	MustGetSha256Sum() (sha256sum string)
	MustGetTextBlocks(verbose bool) (textBlocks []string)
	MustIsContentEqualByComparingSha256Sum(other File, verbose bool) (isMatching bool)
	MustIsMatchingSha256Sum(sha256sum string) (isMatching bool)
	MustReadAsBool() (content bool)
	MustReadAsInt64() (content int64)
	MustReadAsLines() (contentLines []string)
	MustReadAsLinesWithoutComments() (contentLines []string)
	MustReadAsString() (content string)
	MustReadFirstLine() (firstLine string)
	MustReadFirstLineAndTrimSpace() (firstLine string)
	MustReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *ChangeSummary)
	MustSortBlocksInFile(verbose bool)
	MustWriteInt64(toWrite int64, verbose bool)
	MustWriteLines(linesToWrite []string, verbose bool)
	MustWriteString(content string, verbose bool)
	MustWriteTextBlocks(textBlocks []string, verose bool)
	ReadAsBool() (content bool, err error)
	ReadAsInt64() (content int64, err error)
	ReadAsLines() (contentLines []string, err error)
	ReadAsLinesWithoutComments() (contentLines []string, err error)
	ReadAsString() (content string, err error)
	ReadFirstLine() (firstLine string, err error)
	ReadFirstLineAndTrimSpace() (firstLine string, err error)
	ReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *ChangeSummary, err error)
	SortBlocksInFile(verbose bool) (err error)
	WriteInt64(toWrite int64, verboe bool) (err error)
	WriteLines(linesToWrite []string, verbose bool) (err error)
	WriteString(content string, verbose bool) (err error)
	WriteTextBlocks(textBlocks []string, verbose bool) (err error)
}

func GetFileByOsFile(osFile *os.File) (file File, err error) {
	if osFile == nil {
		return nil, TracedError("osFile is nil")
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
		LogGoErrorFatal(err)
	}

	return file
}
