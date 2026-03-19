package filesgeneric

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/changesummary"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"

	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/datetime"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/signalmessenger"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var ErrFileBaseParentNotSet = errors.New("parent is not set")

// This is the base for `File` providing most convenience functions for file operations.
type FileBase struct {
	parentFileForBaseClass filesinterfaces.File
}

func NewFileBase() (f *FileBase) {
	return new(FileBase)
}

// Returns true if the file is a file on the local host.
//
// If a file can return a local path the assumption is it is a local file.
func (f *FileBase) IsLocalFile(ctx context.Context) (isLocalFile bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	hostDescription, err := parent.GetHostDescription()
	if err != nil {
		return false, err
	}

	return hostDescription == "localhost", nil
}

func (f *FileBase) AppendLine(ctx context.Context, line string) (err error) {
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
		err = parent.EnsureEndsWithLineBreak(ctx)
		if err != nil {
			return err
		}
	}

	err = parent.AppendString(ctx, toWrite)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) CheckIsLocalFile(ctx context.Context) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	isLocalFile, err := parent.IsLocalFile(ctx)
	if err != nil {
		return err
	}

	if !isLocalFile {
		return tracederrors.TracedError("Not a local file")
	}

	return nil
}

func (f *FileBase) ContainsLine(ctx context.Context, line string) (containsLine bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return false, err
	}

	return stringsutils.ContainsLine(content, line), nil
}

func (f *FileBase) CreateParentDirectory(ctx context.Context) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	parentDir, err := parent.GetParentDirectory(ctx)
	if err != nil {
		return err
	}

	err = parentDir.Create(ctx, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) EnsureEndsWithLineBreak(ctx context.Context) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	filePath, err := parent.GetLocalPath()
	if err != nil {
		return err
	}

	err = parent.Create(ctx, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	isEmptyFile, err := parent.IsEmptyFile()
	if err != nil {
		return err
	}

	if isEmptyFile {
		err = parent.WriteString(ctx, "\n", &filesoptions.WriteOptions{})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Added newline to empty file '%s' to ensure ends with line break.", filePath)

		return nil
	}

	lastChar, err := parent.ReadLastCharAsString(ctx)
	if err != nil {
		return err
	}

	if lastChar == "\n" {
		logging.LogInfoByCtxf(ctx, "File '%s' already ends with a line break.", filePath)
	} else {
		err = parent.AppendString(ctx, "\n")
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Added line break at end of '%s'.", filePath)
	}

	return nil
}

func (f *FileBase) EnsureLineInFile(ctx context.Context, line string) (err error) {
	line = stringsutils.TrimAllLeadingAndTailingNewLines(line)

	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	err = parent.Create(ctx, &filesoptions.CreateOptions{}) // ensure the file is created if not existent.
	if err != nil {
		return err
	}

	lines, err := parent.ReadAsLines(ctx)
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
		err := parent.AppendLine(ctx, line)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Wrote line '%s' into '%s'.", line, localPath)
	}

	return nil
}

func (f *FileBase) GetFileTypeDescription(ctx context.Context) (fileTypeDescription string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	path, err := parent.GetLocalPath()
	if err != nil {
		return "", err
	}

	stdoutLines, err := commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"file", path},
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

func (f *FileBase) GetMimeType(ctx context.Context) (mimeType string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	const atMostConsideredBytesByHttpDetectContentType int = 512
	beginningOfFile, err := parent.ReadFirstNBytes(ctx, atMostConsideredBytesByHttpDetectContentType)
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

				logging.LogInfoByCtxf(ctx, "Adjusted mimeType of '%s' to '%s' to make it compliant to output of unix 'file' command.", path, mimeType)
			}
		} else if len(beginningOfFile) <= 0 {
			mimeType = "inode/x-empty"

			logging.LogInfoByCtxf(ctx, "Adjusted mimeType of '%s' to '%s' to make it compliant to output of unix 'file' command.", path, mimeType)
		}
	}

	return mimeType, nil
}

func (f *FileBase) GetNumberOfLinesWithPrefix(ctx context.Context, prefix string, trimLines bool) (nLines int, err error) {
	contentString, err := f.ReadAsString(ctx)
	if err != nil {
		return -1, err
	}

	nLines = stringsutils.GetNumberOfLinesWithPrefix(contentString, prefix, trimLines)

	return nLines, nil
}

func (f *FileBase) GetNumberOfNonEmptyLines(ctx context.Context) (nLines int, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return -1, err
	}

	lines, err := parent.ReadAsLines(ctx)
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

func (f *FileBase) GetParentDirectoryPath(ctx context.Context) (parentDirectoryPath string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	parentDir, err := parent.GetParentDirectory(ctx)
	if err != nil {
		return "", err
	}

	parentDirectoryPath, err = parentDir.GetLocalPath()
	if err != nil {
		return "", err
	}

	return parentDirectoryPath, nil
}

func (f *FileBase) GetParentFileForBaseClass() (parentFileForBaseClass filesinterfaces.File, err error) {
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

func (f *FileBase) GetSha256Sum(ctx context.Context) (sha256sum string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return "", err
	}

	sha256sum = checksumutils.GetSha256SumFromString(content)

	return sha256sum, nil
}

func (f *FileBase) GetTextBlocks(ctx context.Context) (textBlocks []string, err error) {
	lines, err := f.ReadAsLines(ctx)
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

	logging.LogInfoByCtxf(ctx, "Splitted file into '%d' text blocks.", len(textBlocks))

	return textBlocks, nil
}

func (f *FileBase) GetValueAsInt(ctx context.Context, key string) (value int, err error) {
	if key == "" {
		return -1, tracederrors.TracedErrorEmptyString("key")
	}

	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return -1, err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return -1, err
	}

	return stringsutils.GetValueAsInt(content, key)
}

func (f *FileBase) GetValueAsString(ctx context.Context, key string) (value string, err error) {
	if key == "" {
		return "", tracederrors.TracedErrorEmptyString("key")
	}

	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return "", err
	}

	return stringsutils.GetValueAsString(content, key)
}

func (f *FileBase) IsContentEqualByComparingSha256Sum(ctx context.Context, otherFile filesinterfaces.File) (isEqual bool, err error) {
	if otherFile == nil {
		return false, tracederrors.TracedErrorNil("otherFile")
	}

	thisChecksum, err := f.GetSha256Sum(ctx)
	if err != nil {
		return false, err
	}

	otherChecksum, err := otherFile.GetSha256Sum(ctx)
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

	exists, err := parent.Exists(contextutils.ContextSilent())
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	currentSum, err := parent.GetSha256Sum(contextutils.ContextSilent())
	if err != nil {
		return false, err
	}

	isMatching = currentSum == sha256sum
	return isMatching, nil
}

func (f *FileBase) IsPgpEncrypted(ctx context.Context) (isEncrypted bool, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return false, err
	}

	mimeType, err := parent.GetMimeType(ctx)
	if err != nil {
		return false, err
	}

	if mimeType == "application/pgp-encrypted" {
		return true, nil
	}

	fileDescription, err := parent.GetFileTypeDescription(ctx)
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

func (f *FileBase) PrintContentOnStdout(ctx context.Context) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return err
	}

	fmt.Print(content)

	return nil
}

func (f *FileBase) ReadAsBool(ctx context.Context) (boolValue bool, err error) {
	contentString, err := f.ReadAsString(ctx)
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

func (f *FileBase) ReadAsFloat64(ctx context.Context) (content float64, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return -1, err
	}

	contentString, err := parent.ReadAsString(ctx)
	if err != nil {
		return -1, err
	}

	content, err = strconv.ParseFloat(contentString, 64)
	if err != nil {
		return -1, err
	}

	return content, nil
}

func (f *FileBase) ReadAsInt(ctx context.Context) (readValue int, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return 0, err
	}

	contentString, err := parent.ReadAsString(ctx)
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

func (f *FileBase) ReadAsInt64(ctx context.Context) (readValue int64, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return 0, err
	}

	contentString, err := parent.ReadAsString(ctx)
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

func (f *FileBase) ReadAsLines(ctx context.Context) (contentLines []string, err error) {
	content, err := f.ReadAsString(ctx)
	if err != nil {
		return nil, err
	}

	contentLines = stringsutils.SplitLines(content, false)

	return contentLines, nil
}

func (f *FileBase) ReadAsLinesWithoutComments(ctx context.Context) (contentLines []string, err error) {
	contentString, err := f.ReadAsString(ctx)
	if err != nil {
		return nil, err
	}

	contentString = stringsutils.RemoveComments(contentString)
	contentLines = stringsutils.SplitLines(contentString, false)

	return contentLines, nil
}

func (f *FileBase) ReadAsString(ctx context.Context) (content string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	contentBytes, err := parent.ReadAsBytes(ctx)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}

func (f *FileBase) ReadAsTimeTime(ctx context.Context) (date *time.Time, err error) {
	contentString, err := f.ReadAsString(ctx)
	if err != nil {
		return nil, err
	}

	date, err = datetime.Dates().ParseString(contentString)
	if err != nil {
		return nil, err
	}

	return date, nil
}

func (f *FileBase) ReadFirstLine(ctx context.Context) (firstLine string, err error) {
	content, err := f.ReadAsString(ctx)
	if err != nil {
		return "", err
	}

	firstLine = stringsutils.GetFirstLine(content)

	return firstLine, nil
}

func (f *FileBase) ReadFirstLineAndTrimSpace(ctx context.Context) (firstLine string, err error) {
	firstLine, err = f.ReadFirstLine(ctx)
	if err != nil {
		return "", err
	}

	firstLine = strings.TrimSpace(firstLine)

	return firstLine, nil
}

func (f *FileBase) ReadLastCharAsString(ctx context.Context) (lastChar string, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return "", err
	}

	content, err := parent.ReadAsString(ctx)
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

func (f *FileBase) RemoveLinesWithPrefix(ctx context.Context, prefix string) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return err
	}

	replaced := stringsutils.RemoveLinesWithPrefix(content, prefix)

	path, err := parent.GetPath()
	if err != nil {
		return err
	}

	if content == replaced {
		logging.LogInfoByCtxf(ctx, "No lines with prefix '%s' to remove in '%s'.", prefix, path)
	} else {
		err = parent.WriteString(contextutils.WithSilent(ctx), replaced, &filesoptions.WriteOptions{})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Replaced all lines with prefix '%s' in '%s'.", prefix, path)
	}

	return nil
}

func (f *FileBase) ReplaceLineAfterLine(ctx context.Context, lineToFind string, replaceLineAfterWith string) (changeSummary *changesummary.ChangeSummary, err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return nil, err
	}

	lines, err := parent.ReadAsLines(ctx)
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

			if line == replaceLineAfterWith {
				logging.LogInfoByCtxf(ctx, "ReplaceLineAfterLine: No need to replace line '%d' in '%s' as already '%s'", lineNumber, path, replaceLineAfterWith)
			} else {
				logging.LogInfoByCtxf(ctx, "ReplaceLineAfterLine: Replace line '%d' in '%s' by '%s' (was '%s')", lineNumber, path, replaceLineAfterWith, line)
				numberOfReplaces += 1
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
		err = parent.WriteLines(ctx, linesToWrite)
		if err != nil {
			return nil, err
		}

		logging.LogChangedByCtxf(ctx, "ReplaceLineAfterLine: Replaced '%d' lines in '%s'.", numberOfReplaces, path)
	} else {
		logging.LogInfoByCtxf(ctx, "ReplaceLineAfterLine: No replaces in '%s' made since no matches were found.", path)
	}

	return changeSummary, nil
}

func (f *FileBase) SetParentFileForBaseClass(parentFileForBaseClass filesinterfaces.File) (err error) {
	f.parentFileForBaseClass = parentFileForBaseClass

	return nil
}

func (f *FileBase) SortBlocksInFile(ctx context.Context) (err error) {
	blocks, err := f.GetTextBlocks(ctx)
	if err != nil {
		return err
	}

	sort.Strings(blocks)

	err = f.WriteTextBlocks(ctx, blocks)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) TrimSpacesAtBeginningOfFile(ctx context.Context) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	content, err := parent.ReadAsString(ctx)
	if err != nil {
		return err
	}

	content = stringsutils.TrimSpacesLeft(content)

	err = f.WriteString(ctx, content, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) WriteInt64(ctx context.Context, toWrite int64) (err error) {
	stringRepresentation := fmt.Sprintf("%d", toWrite)

	err = f.WriteString(ctx, stringRepresentation, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) WriteLines(ctx context.Context, linesToWrite []string) (err error) {
	if linesToWrite == nil {
		return tracederrors.TracedErrorNil("linesToWrite")
	}

	contentToWrite := strings.Join(linesToWrite, "\n")

	err = f.WriteString(ctx, contentToWrite, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (f *FileBase) WriteString(ctx context.Context, toWrite string, options *filesoptions.WriteOptions) (err error) {
	parent, err := f.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	return parent.WriteBytes(ctx, []byte(toWrite), options)
}

func (f *FileBase) WriteTextBlocks(ctx context.Context, textBlocks []string) (err error) {
	textToWrite := ""

	for i, blockToWrite := range textBlocks {
		if i > 0 {
			blockToWrite = "\n" + blockToWrite
		}
		blockToWrite = stringsutils.EnsureEndsWithExactlyOneLineBreak(blockToWrite)

		textToWrite += blockToWrite
	}

	err = f.WriteString(ctx, textToWrite, &filesoptions.WriteOptions{})
	if err != nil {
		return err
	}

	return nil

}

func (xxx *FileBase) GetCreationDateByFileName(ctx context.Context) (creationDate *time.Time, err error) {
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

func (xxx *FileBase) CopyToFile(ctx context.Context, destFile filesinterfaces.File, options *filesoptions.CopyOptions) error {
	if destFile == nil {
		return tracederrors.TracedErrorNil("destFile")
	}

	if options == nil {
		options = &filesoptions.CopyOptions{}
	}

	parent, err := xxx.GetParentFileForBaseClass()
	if err != nil {
		return err
	}

	srcPath, srcHostDescription, err := parent.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	destPath, destHostDescription, err := destFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Copy '%s' on '%s' to '%s' on '%s' started.", srcPath, srcHostDescription, destPath, destHostDescription)

	src, err := parent.OpenAsReadCloser(ctx)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := destFile.OpenAsWriteCloser(ctx, &filesoptions.WriteOptions{UseSudo: options.UseSudo})
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to copy file '%s' on '%s' to '%s' on '%s': %w", srcPath, srcHostDescription, destPath, destHostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Copy '%s' on '%s' to '%s' on '%s' finished.", srcPath, srcHostDescription, destPath, destHostDescription)

	return nil
}
