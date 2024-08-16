package asciichgolangpublic

import (
	"errors"
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
