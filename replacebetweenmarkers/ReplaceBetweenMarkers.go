package asciichgolangpublic

import (
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetContentToInsertDefinedInStartLineAsLines(line string, options *ReplaceBetweenMarkersOptions) (lines []string, err error) {
	if line == "" {
		return nil, tracederrors.TracedError("line is empty string")
	}

	if options == nil {
		return nil, tracederrors.TracedError("options is nil")
	}

	sourceFile, err := GetSourceFile(line, options)
	if err != nil {
		return nil, err
	}

	lines, err = sourceFile.ReadAsLines()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func GetSourceFile(line string, options *ReplaceBetweenMarkersOptions) (sourceFile files.File, err error) {
	if line == "" {
		return nil, tracederrors.TracedError("line is empty string")
	}

	if options == nil {
		return nil, tracederrors.TracedError("options is nil")
	}

	sourcePath, err := GetSourcePath(line, options)
	if err != nil {
		return nil, err
	}

	sourceFile, err = files.GetLocalFileByPath(sourcePath)
	if err != nil {
		return nil, err
	}

	return sourceFile, nil
}

func GetSourcePath(line string, options *ReplaceBetweenMarkersOptions) (sourcePath string, err error) {
	if line == "" {
		return "", tracederrors.TracedError("line is empty string")
	}

	if options == nil {
		return "", tracederrors.TracedError("options is nil")
	}

	if !IsReplaceBetweenMarkerStart(line) {
		return "", tracederrors.TracedErrorf("Unable to get source path of unknown line '%s'", line)
	}

	splitted := strings.Split(line, "source=")
	if len(splitted) != 2 {
		return "", tracederrors.TracedErrorf("Unexpected split: '%v'", splitted)
	}

	sourcePath = splitted[1]
	sourcePath = strings.Split(sourcePath, " ")[0]
	sourcePath = strings.TrimSpace(sourcePath)
	sourcePath = stringsutils.RemoveSurroundingQuotationMarks(sourcePath)
	sourcePath = strings.TrimSpace(sourcePath)

	if sourcePath == "" {
		return "", tracederrors.TracedErrorf("sourcePath is empty string after evaluationg source path in line '%s'", line)
	}

	if pathsutils.IsRelativePath(sourcePath) {
		workingDirectory, err := options.GetWorkingDirPath()
		if err != nil {
			return "", err
		}

		sourcePath = filepath.Join(workingDirectory, sourcePath)
	}

	if options.Verbose {
		logging.LogInfof(
			"Source path found in replace between markers line '%s' is '%s'",
			line,
			sourcePath,
		)
	}

	return sourcePath, nil
}

func IsReplaceBetweenMarkerEnd(line string) (isReplaceBetweenMarkerEnd bool) {
	commentContent := stringsutils.RemoveCommentMarkersAndTrimSpace(line)
	return strings.HasPrefix(commentContent, "REPLACE_BETWEEN_MARKERS END")
}

func IsReplaceBetweenMarkerStart(line string) (isReplaceBetweenMarkerStart bool) {
	commentContent := stringsutils.RemoveCommentMarkersAndTrimSpace(line)
	return strings.HasPrefix(commentContent, "REPLACE_BETWEEN_MARKERS START")
}

func MustGetContentToInsertDefinedInStartLineAsLines(line string, options *ReplaceBetweenMarkersOptions) (lines []string) {
	lines, err := GetContentToInsertDefinedInStartLineAsLines(line, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return lines
}

func MustGetSourceFile(line string, options *ReplaceBetweenMarkersOptions) (sourceFile files.File) {
	sourceFile, err := GetSourceFile(line, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourceFile
}

func MustGetSourcePath(line string, options *ReplaceBetweenMarkersOptions) (sourcePath string) {
	sourcePath, err := GetSourcePath(line, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourcePath
}

func MustReplaceBySourcesInString(input string, options *ReplaceBetweenMarkersOptions) (replaced string) {
	replaced, err := ReplaceBySourcesInString(input, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return replaced
}

func ReplaceBySourcesInString(input string, options *ReplaceBetweenMarkersOptions) (replaced string, err error) {
	if options == nil {
		return "", tracederrors.TracedError("options is nil")
	}

	outLines := []string{}

	var startMarkerFound bool = false
	for _, line := range stringsutils.SplitLines(input, false) {
		if stringsutils.IsComment(line) {
			if IsReplaceBetweenMarkerStart(line) {
				startMarkerFound = true
				outLines = append(outLines, line)

				additionalLines, err := GetContentToInsertDefinedInStartLineAsLines(line, options)
				if err != nil {
					return "", err
				}

				outLines = append(outLines, additionalLines...)

				continue
			}

			if IsReplaceBetweenMarkerEnd(line) {
				startMarkerFound = false
				outLines = append(outLines, line)
				continue
			}
		}

		if startMarkerFound {
			continue
		}

		outLines = append(outLines, line)
	}

	replaced = strings.Join(outLines, "\n")

	return replaced, nil
}
