package files

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/changesummary"
	"github.com/asciich/asciichgolangpublic/checksums"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/datetime"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/signalmessenger"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var ErrFileBaseParentNotSet = errors.New("parent is not set")

// This is the base for `File` providing most convenience functions for file operations.
type FileBase struct {
	parentFileForBaseClass File
}

func NewFileBase() (f *FileBase) {
	return new(FileBase)
}

// Returns true if the file is a file on the local host.
//
// If a file can return a local path the assumption is it is a local file.
func (f *FileBase) IsLocalFile(verbose bool) (isLocalFile bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	localPath, err := parent.GetLocalPath()
	if err != nil {
		return false, tracederrors.TracedErrorf(
			"Not implemented for %w",
			err,
		)
	}

	if verbose {
		logging.LogInfof(
			"'%s' is a local file.",
			localPath,
		)
	}

	return true, nil
}

func (f *FileBase) AppendLine(line string, verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	toWrite := stringsutils.TrimAllLeadingAndTailingNewLines(line)
	toWrite = stringsutils.EnsureEndsWithExactlyOneLineBreak(toWrite)

	isEmptyFile, err := f.IsEmptyFile()
	if err != nil {
		return err
	}

	if !isEmptyFile {
		err = parent.EnsureEndsWithLineBreak(verbose)
		if err != nil {
			return err
		}
	}

	err = parent.AppendString(toWrite, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) CheckIsLocalFile(verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	isLocalFile, err := parent.IsLocalFile(verbose)
	if err != nil {
		return err
	}

	if !isLocalFile {
		return tracederrors.TracedError("Not a local file")
	}

	return nil
}

func (f *FileBase) ContainsLine(line string) (containsLine bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return false, err
	}

	return stringsutils.ContainsLine(content, line), nil
}

func (f *FileBase) CreateParentDirectory(verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	parentDir, err := parent.GetParentDirectory()
	if err != nil {
		return err
	}

	err = parentDir.Create(verbose)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) EnsureEndsWithLineBreak(verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	filePath, err := parent.GetLocalPath()
	if err != nil {
		return err
	}

	err = parent.Create(verbose)
	if err != nil {
		return err
	}

	isEmptyFile, err := parent.IsEmptyFile()
	if err != nil {
		return err
	}

	if isEmptyFile {
		err = parent.WriteString("\n", verbose)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf("Added newline to empty file '%s' to ensure ends with line break.", filePath)
		}

		return nil
	}

	lastChar, err := parent.ReadLastCharAsString()
	if err != nil {
		return err
	}

	if lastChar == "\n" {
		if verbose {
			logging.LogInfof("File '%s' already ends with a line break.", filePath)
		}
	} else {
		err = parent.AppendString("\n", verbose)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf("Added line break at end of '%s'.", filePath)
		}
	}

	return nil
}

func (f *FileBase) EnsureLineInFile(line string, verbose bool) (err error) {
	line = stringsutils.TrimAllLeadingAndTailingNewLines(line)

	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	err = parent.Create(verbose) // ensure the file is created if not existent.
	if err != nil {
		return err
	}

	lines, err := parent.ReadAsLines()
	if err != nil {
		return err
	}

	localPath, err := parent.GetLocalPath()
	if err != nil {
		return err
	}

	if slices.Contains(lines, line) {
		logging.LogInfof("Line '%s' already present in '%s'.", line, localPath)
	} else {
		err := parent.AppendLine(line, verbose)
		if err != nil {
			return err
		}
		logging.LogChangedf("Wrote line '%s' into '%s'.", line, localPath)
	}

	return nil
}

func (f *FileBase) GetFileTypeDescription(verbose bool) (fileTypeDescription string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	path, err := parent.GetLocalPath()
	if err != nil {
		return "", err
	}

	stdoutLines, err := commandexecutor.Bash().RunCommandAndGetStdoutAsLines(
		&parameteroptions.RunCommandOptions{
			Command: []string{"file", path},
			Verbose: verbose,
		},
	)

	if err != nil {
		return "", err
	}

	stdoutLines = slicesutils.RemoveEmptyStrings(stdoutLines)
	if len(stdoutLines) != 1 {
		return "", tracederrors.TracedErrorf("Expected exactly one line left bug got: '%v'", stdoutLines)
	}

	line := stdoutLines[0]
	splitted := strings.Split(line, ":")
	if len(splitted) != 2 {
		return "", tracederrors.TracedErrorf("Unexpected amount of splitted: '%v'", splitted)
	}

	fileTypeDescription = strings.TrimSpace(splitted[1])
	return fileTypeDescription, nil
}

func (f *FileBase) GetMimeType(verbose bool) (mimeType string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	const atMostConsideredBytesByHttpDetectContentType int = 512
	beginningOfFile, err := parent.ReadFirstNBytes(atMostConsideredBytesByHttpDetectContentType)
	if err != nil {
		return "", err
	}

	mimeType = http.DetectContentType(beginningOfFile)
	if mimeType == "" {
		return "", tracederrors.TracedError("Mimetype is empty string after evaluation")
	}

	path, err := parent.GetLocalPath()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(mimeType, "text/plain;") {
		const beginPgpMessage string = "-----BEGIN PGP MESSAGE-----\n"
		if len(beginningOfFile) > len(beginPgpMessage) {
			if strings.HasPrefix(string(beginningOfFile), beginPgpMessage) {
				mimeType = "application/pgp-encrypted"

				if verbose {
					logging.LogInfof(
						"Adjusted mimeType of '%s' to '%s' to make it compliant to output of unix 'file' command.",
						path,
						mimeType,
					)
				}
			}
		} else if len(beginningOfFile) <= 0 {
			mimeType = "inode/x-empty"

			if verbose {
				logging.LogInfof(
					"Adjusted mimeType of '%s' to '%s' to make it compliant to output of unix 'file' command.",
					path,
					mimeType,
				)
			}

		}
	}

	return mimeType, nil
}

func (f *FileBase) GetNumberOfLinesWithPrefix(prefix string, trimLines bool) (nLines int, err error) {
	contentString, err := f.ReadAsString()
	if err != nil {
		return -1, err
	}

	nLines = stringsutils.GetNumberOfLinesWithPrefix(contentString, prefix, trimLines)

	return nLines, nil
}

func (f *FileBase) GetNumberOfNonEmptyLines() (nLines int, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return -1, err
	}

	lines, err := parent.ReadAsLines()
	if err != nil {
		return -1, err
	}

	nLines = 0
	for _, l := range lines {
		if l != "" {
			nLines += 1
		}
	}

	return nLines, nil
}

func (f *FileBase) GetParentDirectoryPath() (parentDirectoryPath string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	parentDir, err := parent.GetParentDirectory()
	if err != nil {
		return "", err
	}

	parentDirectoryPath, err = parentDir.GetLocalPath()
	if err != nil {
		return "", err
	}

	return parentDirectoryPath, nil
}

func (f *FileBase) GetParentFileForBaseClass() (parentFileForBaseClass File, err error) {
	if f.parentFileForBaseClass == nil {
		return nil, tracederrors.TracedErrorf("%w", ErrFileBaseParentNotSet)
	}
	return f.parentFileForBaseClass, nil
}

func (f *FileBase) GetPathAndHostDescription() (path string, hostDescription string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", "", err
	}

	path, err = parent.GetPath()
	if err != nil {
		return "", "", err
	}

	hostDescription, err = parent.GetHostDescription()
	if err != nil {
		return "", "", err
	}

	return path, hostDescription, nil
}

func (f *FileBase) GetSha256Sum() (sha256sum string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return "", err
	}

	sha256sum = checksums.GetSha256SumFromString(content)

	return sha256sum, nil
}

func (f *FileBase) GetTextBlocks(verbose bool) (textBlocks []string, err error) {
	lines, err := f.ReadAsLines()
	if err != nil {
		return nil, err
	}

	var blockToAdd string = ""
	textBlocks = []string{}

	if len(lines) >= 1 {
		if lines[0] == "---" {
			textBlocks = append(textBlocks, "---")
			lines = lines[1:]
		}
	}

	insideBlock := false
	braceEndMarker := ""
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if insideBlock {
			if line == braceEndMarker {
				if len(line) > 0 {
					blockToAdd += "\n" + line
				}
				textBlocks = append(textBlocks, blockToAdd)
				insideBlock = false
				braceEndMarker = ""
			} else {
				if !strings.HasPrefix(trimmedLine, "//") {
					currentBlockWithoutComments := stringsutils.RemoveCommentsAndTrimSpace(blockToAdd)
					if currentBlockWithoutComments == "" {
						if strings.HasSuffix(trimmedLine, "(") {
							braceEndMarker = ")"
						}
						if strings.HasSuffix(trimmedLine, "{") {
							braceEndMarker = "}"
						}
					}
				}
				blockToAdd += "\n" + line
			}
		} else {
			if trimmedLine == "" {
				continue
			} else {
				blockToAdd = line
				insideBlock = true
				if strings.HasSuffix(trimmedLine, "(") {
					braceEndMarker = ")"
				}
				if strings.HasSuffix(trimmedLine, "{") {
					braceEndMarker = "}"
				}
			}
		}
	}

	if insideBlock {
		textBlocks = append(textBlocks, blockToAdd)
		insideBlock = false
	}

	if verbose {
		logging.LogInfof("Splitted file into '%d' text blocks.", len(textBlocks))
	}

	return textBlocks, nil
}

func (f *FileBase) GetValueAsInt(key string) (value int, err error) {
	if key == "" {
		return -1, tracederrors.TracedErrorEmptyString("key")
	}

	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return -1, err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return -1, err
	}

	return stringsutils.GetValueAsInt(content, key)
}

func (f *FileBase) GetValueAsString(key string) (value string, err error) {
	if key == "" {
		return "", tracederrors.TracedErrorEmptyString("key")
	}

	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return "", err
	}

	return stringsutils.GetValueAsString(content, key)
}

func (f *FileBase) IsContentEqualByComparingSha256Sum(otherFile File, verbose bool) (isEqual bool, err error) {
	if otherFile == nil {
		return false, tracederrors.TracedErrorNil("otherFile")
	}

	thisChecksum, err := f.GetSha256Sum()
	if err != nil {
		return false, err
	}

	otherChecksum, err := otherFile.GetSha256Sum()
	if err != nil {
		return false, err
	}

	isEqual = thisChecksum == otherChecksum

	return isEqual, nil
}

func (f *FileBase) IsEmptyFile() (isEmtpyFile bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	size, err := parent.GetSizeBytes()
	if err != nil {
		return false, err
	}

	isEmtpyFile = size == 0

	return isEmtpyFile, nil
}

func (f *FileBase) IsMatchingSha256Sum(sha256sum string) (isMatching bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	exists, err := parent.Exists(false)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	currentSum, err := parent.GetSha256Sum()
	if err != nil {
		return false, err
	}

	isMatching = currentSum == sha256sum
	return isMatching, nil
}

func (f *FileBase) IsPgpEncrypted(verbose bool) (isEncrypted bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	mimeType, err := parent.GetMimeType(verbose)
	if err != nil {
		return false, err
	}

	if mimeType == "application/pgp-encrypted" {
		return true, nil
	}

	fileDescription, err := parent.GetFileTypeDescription(verbose)
	if err != nil {
		return false, err
	}

	if stringsutils.ContainsAtLeastOneSubstringIgnoreCase(
		fileDescription,
		[]string{"gpg", "pgp"},
	) {
		return true, nil
	}

	return false, nil
}

func (f *FileBase) IsYYYYmmdd_HHMMSSPrefix() (hasDatePrefix bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	basename, err := parent.GetBaseName()
	if err != nil {
		return false, err
	}

	layoutString := datetime.Dates().LayoutStringYYYYmmdd_HHMMSS()

	if len(basename) < len(layoutString) {
		return false, nil
	}

	toParse := basename[:len(layoutString)]
	_, err = datetime.Dates().ParseStringWithGivenLayout(toParse, layoutString)
	if err != nil {
		if strings.Contains(err.Error(), "Unable to parse as date") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f *FileBase) MustAppendLine(line string, verbose bool) {
	err := f.AppendLine(line, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustCheckIsLocalFile(verbose bool) {
	err := f.CheckIsLocalFile(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustContainsLine(line string) (containsLine bool) {
	containsLine, err := f.ContainsLine(line)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return containsLine
}

func (f *FileBase) MustCreateParentDirectory(verbose bool) {
	err := f.CreateParentDirectory(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustEnsureEndsWithLineBreak(verbose bool) {
	err := f.EnsureEndsWithLineBreak(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustEnsureLineInFile(line string, verbose bool) {
	err := f.EnsureLineInFile(line, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustGetCreationDateByFileName(verbose bool) (creationDate *time.Time) {
	creationDate, err := f.GetCreationDateByFileName(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return creationDate
}

func (f *FileBase) MustGetFileTypeDescription(verbose bool) (fileTypeDescription string) {
	fileTypeDescription, err := f.GetFileTypeDescription(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return fileTypeDescription
}

func (f *FileBase) MustGetMimeType(verbose bool) (mimeType string) {
	mimeType, err := f.GetMimeType(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mimeType
}

func (f *FileBase) MustGetNumberOfLinesWithPrefix(prefix string, trimLines bool) (nLines int) {
	nLines, err := f.GetNumberOfLinesWithPrefix(prefix, trimLines)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nLines
}

func (f *FileBase) MustGetNumberOfNonEmptyLines() (nLines int) {
	nLines, err := f.GetNumberOfNonEmptyLines()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return nLines
}

func (f *FileBase) MustGetParentDirectoryPath() (parentDirectoryPath string) {
	parentDirectoryPath, err := f.GetParentDirectoryPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentDirectoryPath
}

func (f *FileBase) MustGetParentFileForBaseClass() (parentFileForBaseClass File) {
	parentFileForBaseClass, err := f.GetParentFileForBaseClass()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentFileForBaseClass
}

func (f *FileBase) MustGetPathAndHostDescription() (path string, hostDescription string) {
	path, hostDescription, err := f.GetPathAndHostDescription()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return path, hostDescription
}

func (f *FileBase) MustGetSha256Sum() (sha256sum string) {
	sha256sum, err := f.GetSha256Sum()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sha256sum
}

func (f *FileBase) MustGetTextBlocks(verbose bool) (textBlocks []string) {
	textBlocks, err := f.GetTextBlocks(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return textBlocks
}

func (f *FileBase) MustGetValueAsInt(key string) (value int) {
	value, err := f.GetValueAsInt(key)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return value
}

func (f *FileBase) MustGetValueAsString(key string) (value string) {
	value, err := f.GetValueAsString(key)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return value
}

func (f *FileBase) MustIsContentEqualByComparingSha256Sum(otherFile File, verbose bool) (isEqual bool) {
	isEqual, err := f.IsContentEqualByComparingSha256Sum(otherFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isEqual
}

func (f *FileBase) MustIsEmptyFile() (isEmtpyFile bool) {
	isEmtpyFile, err := f.IsEmptyFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isEmtpyFile
}

func (f *FileBase) MustIsLocalFile(verbose bool) (isLocalFile bool) {
	isLocalFile, err := f.IsLocalFile(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isLocalFile
}

func (f *FileBase) MustIsMatchingSha256Sum(sha256sum string) (isMatching bool) {
	isMatching, err := f.IsMatchingSha256Sum(sha256sum)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isMatching
}

func (f *FileBase) MustIsPgpEncrypted(verbose bool) (isEncrypted bool) {
	isEncrypted, err := f.IsPgpEncrypted(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isEncrypted
}

func (f *FileBase) MustIsYYYYmmdd_HHMMSSPrefix() (hasDatePrefix bool) {
	hasDatePrefix, err := f.IsYYYYmmdd_HHMMSSPrefix()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hasDatePrefix
}

func (f *FileBase) MustPrintContentOnStdout() {
	err := f.PrintContentOnStdout()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustReadAsBool() (boolValue bool) {
	boolValue, err := f.ReadAsBool()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return boolValue
}

func (f *FileBase) MustReadAsFloat64() (content float64) {
	content, err := f.ReadAsFloat64()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (f *FileBase) MustReadAsInt() (readValue int) {
	readValue, err := f.ReadAsInt()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return readValue
}

func (f *FileBase) MustReadAsInt64() (readValue int64) {
	readValue, err := f.ReadAsInt64()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return readValue
}

func (f *FileBase) MustReadAsLines() (contentLines []string) {
	contentLines, err := f.ReadAsLines()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return contentLines
}

func (f *FileBase) MustReadAsLinesWithoutComments() (contentLines []string) {
	contentLines, err := f.ReadAsLinesWithoutComments()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return contentLines
}

func (f *FileBase) MustReadAsString() (content string) {
	content, err := f.ReadAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return content
}

func (f *FileBase) MustReadAsTimeTime() (date *time.Time) {
	date, err := f.ReadAsTimeTime()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return date
}

func (f *FileBase) MustReadFirstLine() (firstLine string) {
	firstLine, err := f.ReadFirstLine()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return firstLine
}

func (f *FileBase) MustReadFirstLineAndTrimSpace() (firstLine string) {
	firstLine, err := f.ReadFirstLineAndTrimSpace()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return firstLine
}

func (f *FileBase) MustReadLastCharAsString() (lastChar string) {
	lastChar, err := f.ReadLastCharAsString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return lastChar
}

func (f *FileBase) MustRemoveLinesWithPrefix(prefix string, verbose bool) {
	err := f.RemoveLinesWithPrefix(prefix, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

/* todo remove or move
func (f *FileBase) MustReplaceBetweenMarkers(verbose bool) {
	err := f.ReplaceBetweenMarkers(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}
*/

func (f *FileBase) MustReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *changesummary.ChangeSummary) {
	changeSummary, err := f.ReplaceLineAfterLine(lineToFind, replaceLineAfterWith, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return changeSummary
}

func (f *FileBase) MustSetParentFileForBaseClass(parentFileForBaseClass File) {
	err := f.SetParentFileForBaseClass(parentFileForBaseClass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustSortBlocksInFile(verbose bool) {
	err := f.SortBlocksInFile(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustTrimSpacesAtBeginningOfFile(verbose bool) {
	err := f.TrimSpacesAtBeginningOfFile(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteInt64(toWrite int64, verbose bool) {
	err := f.WriteInt64(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteLines(linesToWrite []string, verbose bool) {
	err := f.WriteLines(linesToWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteString(toWrite string, verbose bool) {
	err := f.WriteString(toWrite, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteTextBlocks(textBlocks []string, verbose bool) {
	err := f.WriteTextBlocks(textBlocks, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (f *FileBase) PrintContentOnStdout() (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return err
	}

	fmt.Print(content)

	return nil
}

func (f *FileBase) ReadAsBool() (boolValue bool, err error) {
	contentString, err := f.ReadAsString()
	if err != nil {
		return false, err
	}

	contentString = strings.TrimSpace(contentString)

	boolValue, err = strconv.ParseBool(contentString)
	if err != nil {
		return false, err
	}

	return boolValue, nil
}

func (f *FileBase) ReadAsFloat64() (content float64, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return -1, err
	}

	contentString, err := parent.ReadAsString()
	if err != nil {
		return -1, err
	}

	content, err = strconv.ParseFloat(contentString, 64)
	if err != nil {
		return -1, err
	}

	return content, nil
}

func (f *FileBase) ReadAsInt() (readValue int, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return 0, err
	}

	contentString, err := parent.ReadAsString()
	if err != nil {
		return 0, err
	}

	localPath, err := parent.GetLocalPath()
	if err != nil {
		return 0, err
	}

	contentString = strings.TrimSpace(contentString)

	readValue, err = strconv.Atoi(contentString)
	if err != nil {
		return 0, tracederrors.TracedErrorf(
			"Unable to parse file '%s' as int: '%w'",
			localPath,
			err,
		)
	}

	return readValue, nil
}

func (f *FileBase) ReadAsInt64() (readValue int64, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return 0, err
	}

	contentString, err := parent.ReadAsString()
	if err != nil {
		return 0, err
	}

	localPath, err := parent.GetLocalPath()
	if err != nil {
		return 0, err
	}

	contentString = strings.TrimSpace(contentString)

	readValue, err = strconv.ParseInt(contentString, 10, 64)
	if err != nil {
		return 0, tracederrors.TracedErrorf(
			"Unable to parse file '%s' as int64: '%w'",
			localPath,
			err,
		)
	}

	return readValue, nil
}

func (f *FileBase) ReadAsLines() (contentLines []string, err error) {
	content, err := f.ReadAsString()
	if err != nil {
		return nil, err
	}

	contentLines = stringsutils.SplitLines(content, false)

	return contentLines, nil
}

func (f *FileBase) ReadAsLinesWithoutComments() (contentLines []string, err error) {
	contentString, err := f.ReadAsString()
	if err != nil {
		return nil, err
	}

	contentString = stringsutils.RemoveComments(contentString)
	contentLines = stringsutils.SplitLines(contentString, false)

	return contentLines, nil
}

func (f *FileBase) ReadAsString() (content string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	contentBytes, err := parent.ReadAsBytes()
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func (f *FileBase) ReadAsTimeTime() (date *time.Time, err error) {
	contentString, err := f.ReadAsString()
	if err != nil {
		return nil, err
	}

	date, err = datetime.Dates().ParseString(contentString)
	if err != nil {
		return nil, err
	}

	return date, nil
}

func (f *FileBase) ReadFirstLine() (firstLine string, err error) {
	content, err := f.ReadAsString()
	if err != nil {
		return "", err
	}

	firstLine = stringsutils.GetFirstLine(content)

	return firstLine, nil
}

func (f *FileBase) ReadFirstLineAndTrimSpace() (firstLine string, err error) {
	firstLine, err = f.ReadFirstLine()
	if err != nil {
		return "", err
	}

	firstLine = strings.TrimSpace(firstLine)

	return firstLine, nil
}

func (f *FileBase) ReadLastCharAsString() (lastChar string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return "", err
	}

	localPath, err := parent.GetLocalPath()
	if err != nil {
		return "", err
	}

	if len(content) <= 0 {
		return "", tracederrors.TracedErrorf("Get last char failed, '%s' is empty file", localPath)
	}

	lastChar = content[len(content)-1:]

	return lastChar, nil
}

func (f *FileBase) RemoveLinesWithPrefix(prefix string, verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return err
	}

	replaced := stringsutils.RemoveLinesWithPrefix(content, prefix)

	path, err := parent.GetPath()
	if err != nil {
		return err
	}

	if content == replaced {
		if verbose {
			logging.LogInfof(
				"No lines with prefix '%s' to remove in '%s'.",
				prefix,
				path,
			)
		}
	} else {
		err = parent.WriteString(
			replaced,
			false,
		)
		if err != nil {
			return err
		}

		if verbose {
			logging.LogChangedf(
				"Replaced all lines with prefix '%s' in '%s'.",
				prefix,
				path,
			)
		}
	}

	return nil
}

/* TODO remove or move
func (f *FileBase) ReplaceBetweenMarkers(verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return err
	}

	workingDirPath, err := parent.GetParentDirectoryPath()
	if err != nil {
		return err
	}

	content, err = ReplaceBetweenMarkers().ReplaceBySourcesInString(
		content,
		&ReplaceBetweenMarkersOptions{
			WorkingDirPath: workingDirPath,
			Verbose:        verbose,
		},
	)
	if err != nil {
		return err
	}

	err = parent.WriteString(content, verbose)
	if err != nil {
		return err
	}

	path, err := parent.GetLocalPath()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Replace between markers finished in '%s'.", path)
	}

	return nil
}
*/

func (f *FileBase) ReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *changesummary.ChangeSummary, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return nil, err
	}

	lines, err := parent.ReadAsLines()
	if err != nil {
		return nil, err
	}

	path, err := parent.GetLocalPath()
	if err != nil {
		return nil, err
	}

	matchFound := false
	linesToWrite := []string{}

	numberOfReplaces := 0

	for i, line := range lines {

		if matchFound {
			lineNumber := i + 1

			if verbose {
				if line == replaceLineAfterWith {
					logging.LogInfof(
						"ReplaceLineAfterLine: No need to replace line '%d' in '%s' as already '%s'",
						lineNumber,
						path,
						replaceLineAfterWith,
					)
				} else {
					logging.LogChangedf(
						"ReplaceLineAfterLine: Replace line '%d' in '%s' by '%s' (was '%s')",
						lineNumber,
						path,
						replaceLineAfterWith,
						line,
					)
					numberOfReplaces += 1
				}
			}

			linesToWrite = append(linesToWrite, replaceLineAfterWith)

			matchFound = false
		} else {
			if line == lineToFind {
				matchFound = true
			}
			linesToWrite = append(linesToWrite, line)
		}
	}

	if matchFound {
		linesToWrite = append(linesToWrite, replaceLineAfterWith)
		logging.LogChangedf(
			"ReplaceLineAfterLine: Appended line '%s' in '%s' since last read line was a match.",
			replaceLineAfterWith,
			path,
		)
		numberOfReplaces += 1
	}

	changeSummary = changesummary.NewChangeSummary()
	err = changeSummary.SetNumberOfChanges(numberOfReplaces)
	if err != nil {
		return nil, err
	}

	if changeSummary.IsChanged() {
		err = parent.WriteLines(linesToWrite, verbose)
		if err != nil {
			return nil, err
		}

		if verbose {
			logging.LogChangedf(
				"ReplaceLineAfterLine: Replaced '%d' lines in '%s'.",
				numberOfReplaces,
				path,
			)
		}
	} else {
		if verbose {
			logging.LogInfof(
				"ReplaceLineAfterLine: No replaces in '%s' made since no matches were found.",
				path,
			)
		}
	}

	return changeSummary, nil
}

func (f *FileBase) SetParentFileForBaseClass(parentFileForBaseClass File) (err error) {
	f.parentFileForBaseClass = parentFileForBaseClass

	return nil
}

func (f *FileBase) SortBlocksInFile(verbose bool) (err error) {
	blocks, err := f.GetTextBlocks(verbose)
	if err != nil {
		return err
	}

	sort.Strings(blocks)

	err = f.WriteTextBlocks(blocks, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) TrimSpacesAtBeginningOfFile(verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString()
	if err != nil {
		return err
	}

	content = stringsutils.TrimSpacesLeft(content)

	err = f.WriteString(content, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) WriteInt64(toWrite int64, verbose bool) (err error) {
	stringRepresentation := fmt.Sprintf("%d", toWrite)

	err = f.WriteString(stringRepresentation, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) WriteLines(linesToWrite []string, verbose bool) (err error) {
	if linesToWrite == nil {
		return tracederrors.TracedErrorNil("linesToWrite")
	}

	contentToWrite := strings.Join(linesToWrite, "\n")

	err = f.WriteString(contentToWrite, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) WriteString(toWrite string, verbose bool) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	return parent.WriteBytes([]byte(toWrite), verbose)
}

func (f *FileBase) WriteTextBlocks(textBlocks []string, verbose bool) (err error) {
	textToWrite := ""

	for i, blockToWrite := range textBlocks {
		if i > 0 {
			blockToWrite = "\n" + blockToWrite
		}
		blockToWrite = stringsutils.EnsureEndsWithExactlyOneLineBreak(blockToWrite)

		textToWrite += blockToWrite
	}

	err = f.WriteString(textToWrite, verbose)
	if err != nil {
		return err
	}

	return nil

}

func (xxx *FileBase) GetCreationDateByFileName(verbose bool) (creationDate *time.Time, err error) {
	parent, err := xxx.GetParentFileForBaseClass()
	if err != nil {
		return nil, err
	}

	basename, err := parent.GetBaseName()
	if err != nil {
		return nil, err
	}

	creationDate, err = datetime.Dates().ParseStringPrefixAsDate(basename)
	if err != nil {
		if strings.Contains(err.Error(), "Unable to parse prefix ") {
			err = nil
		} else if strings.Contains(err.Error(), "Unable to parse date ") {
			err = nil
		} else {
			return nil, err
		}
	}

	if creationDate == nil {
		creationDate, err = signalmessenger.ParseCreationDateFromSignalPictureBaseName(basename)
		if err != nil {
			if strings.Contains(err.Error(), "Unable to parse date ") {
				err = nil
			} else if strings.Contains(err.Error(), "is not a singal picture base name") {
				err = nil
			} else {
				return nil, err
			}
		}
	}

	if creationDate == nil {
		return nil, tracederrors.TracedError("All attempts failed to extract creationDate")
	}

	return creationDate, nil
}
