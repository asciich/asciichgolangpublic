package asciichgolangpublic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrFileBaseParentNotSet = errors.New("parent is not set")

// This is the base for `File` providing most convenience functions for file operations.
type FileBase struct {
	parentFileForBaseClass File
}

func NewFileBase() (f *FileBase) {
	return new(FileBase)
}

func (f *FileBase) GetParentFileForBaseClass() (parentFileForBaseClass File, err error) {
	if f.parentFileForBaseClass == nil {
		return nil, TracedErrorf("%w", ErrFileBaseParentNotSet)
	}
	return f.parentFileForBaseClass, nil
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

	sha256sum = Checksums().GetSha256SumFromString(content)

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
					currentBlockWithoutComments := Strings().RemoveCommentsAndTrimSpace(blockToAdd)
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
		LogInfof("Splitted file into '%d' text blocks.", len(textBlocks))
	}

	return textBlocks, nil
}

func (f *FileBase) IsContentEqualByComparingSha256Sum(otherFile File, verbose bool) (isEqual bool, err error) {
	if otherFile == nil {
		return false, TracedErrorNil("otherFile")
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

func (f *FileBase) IsMatchingSha256Sum(sha256sum string) (isMatching bool, err error) {
	currentSum, err := f.GetSha256Sum()
	if err != nil {
		return false, err
	}

	isMatching = currentSum == sha256sum
	return isMatching, nil
}

func (f *FileBase) MustGetParentFileForBaseClass() (parentFileForBaseClass File) {
	parentFileForBaseClass, err := f.GetParentFileForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentFileForBaseClass
}

func (f *FileBase) MustGetSha256Sum() (sha256sum string) {
	sha256sum, err := f.GetSha256Sum()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sha256sum
}

func (f *FileBase) MustGetTextBlocks(verbose bool) (textBlocks []string) {
	textBlocks, err := f.GetTextBlocks(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return textBlocks
}

func (f *FileBase) MustIsContentEqualByComparingSha256Sum(otherFile File, verbose bool) (isEqual bool) {
	isEqual, err := f.IsContentEqualByComparingSha256Sum(otherFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isEqual
}

func (f *FileBase) MustIsMatchingSha256Sum(sha256sum string) (isMatching bool) {
	isMatching, err := f.IsMatchingSha256Sum(sha256sum)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isMatching
}

func (f *FileBase) MustReadAsBool() (boolValue bool) {
	boolValue, err := f.ReadAsBool()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return boolValue
}

func (f *FileBase) MustReadAsInt64() (readValue int64) {
	readValue, err := f.ReadAsInt64()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return readValue
}

func (f *FileBase) MustReadAsLines() (contentLines []string) {
	contentLines, err := f.ReadAsLines()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return contentLines
}

func (f *FileBase) MustReadAsLinesWithoutComments() (contentLines []string) {
	contentLines, err := f.ReadAsLinesWithoutComments()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return contentLines
}

func (f *FileBase) MustReadAsString() (content string) {
	content, err := f.ReadAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return content
}

func (f *FileBase) MustReadFirstLine() (firstLine string) {
	firstLine, err := f.ReadFirstLine()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return firstLine
}

func (f *FileBase) MustReadFirstLineAndTrimSpace() (firstLine string) {
	firstLine, err := f.ReadFirstLineAndTrimSpace()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return firstLine
}

func (f *FileBase) MustReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *ChangeSummary) {
	changeSummary, err := f.ReplaceLineAfterLine(lineToFind, replaceLineAfterWith, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return changeSummary
}

func (f *FileBase) MustSetParentFileForBaseClass(parentFileForBaseClass File) {
	err := f.SetParentFileForBaseClass(parentFileForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustSortBlocksInFile(verbose bool) {
	err := f.SortBlocksInFile(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteInt64(toWrite int64, verbose bool) {
	err := f.WriteInt64(toWrite, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteLines(linesToWrite []string, verbose bool) {
	err := f.WriteLines(linesToWrite, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteString(toWrite string, verbose bool) {
	err := f.WriteString(toWrite, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (f *FileBase) MustWriteTextBlocks(textBlocks []string, verbose bool) {
	err := f.WriteTextBlocks(textBlocks, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
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

	readValue, err = strconv.ParseInt(contentString, 10, 64)
	if err != nil {
		return 0, TracedErrorf(
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

	contentLines = Strings().SplitLines(content)

	return contentLines, nil
}

func (f *FileBase) ReadAsLinesWithoutComments() (contentLines []string, err error) {
	contentString, err := f.ReadAsString()
	if err != nil {
		return nil, err
	}

	contentString = Strings().RemoveComments(contentString)
	contentLines = Strings().SplitLines(contentString)

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

func (f *FileBase) ReadFirstLine() (firstLine string, err error) {
	content, err := f.ReadAsString()
	if err != nil {
		return "", err
	}

	firstLine = Strings().GetFirstLine(content)

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

func (f *FileBase) ReplaceLineAfterLine(lineToFind string, replaceLineAfterWith string, verbose bool) (changeSummary *ChangeSummary, err error) {
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
					LogInfof(
						"ReplaceLineAfterLine: No need to replace line '%d' in '%s' as already '%s'",
						lineNumber,
						path,
						replaceLineAfterWith,
					)
				} else {
					LogChangedf(
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
		LogChangedf(
			"ReplaceLineAfterLine: Appended line '%s' in '%s' since last read line was a match.",
			replaceLineAfterWith,
			path,
		)
		numberOfReplaces += 1
	}

	changeSummary = NewChangeSummary()
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
			LogChangedf(
				"ReplaceLineAfterLine: Replaced '%d' lines in '%s'.",
				numberOfReplaces,
				path,
			)
		}
	} else {
		if verbose {
			LogInfof(
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

	blocks = Slices().SortStringSlice(blocks)

	err = f.WriteTextBlocks(blocks, verbose)
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
		return TracedErrorNil("linesToWrite")
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
		blockToWrite = Strings().EnsureEndsWithExactlyOneLineBreak(blockToWrite)

		textToWrite += blockToWrite
	}

	err = f.WriteString(textToWrite, verbose)
	if err != nil {
		return nil
	}

	return nil

}
