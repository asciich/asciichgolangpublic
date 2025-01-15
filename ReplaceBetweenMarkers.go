package asciichgolangpublic

import (
	"path/filepath"
	"strings"

	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type ReplaceBetweenMarkersService struct{}

func NewReplaceBetweenMarkersService() (r *ReplaceBetweenMarkersService) {
	return new(ReplaceBetweenMarkersService)
}

func ReplaceBetweenMarkers() (r *ReplaceBetweenMarkersService) {
	return NewReplaceBetweenMarkersService()
}

func (r *ReplaceBetweenMarkersService) GetContentToInsertDefinedInStartLineAsLines(line string, options *ReplaceBetweenMarkersOptions) (lines []string, err error) {
	if line == "" {
		return nil, tracederrors.TracedError("line is empty string")
	}

	if options == nil {
		return nil, tracederrors.TracedError("options is nil")
	}

	sourceFile, err := r.GetSourceFile(line, options)
	if err != nil {
		return nil, err
	}

	lines, err = sourceFile.ReadAsLines()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func (r *ReplaceBetweenMarkersService) GetSourceFile(line string, options *ReplaceBetweenMarkersOptions) (sourceFile File, err error) {
	if line == "" {
		return nil, tracederrors.TracedError("line is empty string")
	}

	if options == nil {
		return nil, tracederrors.TracedError("options is nil")
	}

	sourcePath, err := r.GetSourcePath(line, options)
	if err != nil {
		return nil, err
	}

	sourceFile, err = GetLocalFileByPath(sourcePath)
	if err != nil {
		return nil, err
	}

	return sourceFile, nil
}

func (r *ReplaceBetweenMarkersService) GetSourcePath(line string, options *ReplaceBetweenMarkersOptions) (sourcePath string, err error) {
	if line == "" {
		return "", tracederrors.TracedError("line is empty string")
	}

	if options == nil {
		return "", tracederrors.TracedError("options is nil")
	}

	if !r.IsReplaceBetweenMarkerStart(line) {
		return "", tracederrors.TracedErrorf("Unable to get source path of unknown line '%s'", line)
	}

	splitted := strings.Split(line, "source=")
	if len(splitted) != 2 {
		return "", tracederrors.TracedErrorf("Unexpected split: '%v'", splitted)
	}

	sourcePath = splitted[1]
	sourcePath = strings.Split(sourcePath, " ")[0]
	sourcePath = strings.TrimSpace(sourcePath)
	sourcePath = astrings.RemoveSurroundingQuotationMarks(sourcePath)
	sourcePath = strings.TrimSpace(sourcePath)

	if sourcePath == "" {
		return "", tracederrors.TracedErrorf("sourcePath is empty string after evaluationg source path in line '%s'", line)
	}

	if Paths().IsRelativePath(sourcePath) {
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

func (r *ReplaceBetweenMarkersService) IsReplaceBetweenMarkerEnd(line string) (isReplaceBetweenMarkerEnd bool) {
	commentContent := astrings.RemoveCommentMarkersAndTrimSpace(line)
	return strings.HasPrefix(commentContent, "REPLACE_BETWEEN_MARKERS END")
}

func (r *ReplaceBetweenMarkersService) IsReplaceBetweenMarkerStart(line string) (isReplaceBetweenMarkerStart bool) {
	commentContent := astrings.RemoveCommentMarkersAndTrimSpace(line)
	return strings.HasPrefix(commentContent, "REPLACE_BETWEEN_MARKERS START")
}

func (r *ReplaceBetweenMarkersService) MustGetContentToInsertDefinedInStartLineAsLines(line string, options *ReplaceBetweenMarkersOptions) (lines []string) {
	lines, err := r.GetContentToInsertDefinedInStartLineAsLines(line, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return lines
}

func (r *ReplaceBetweenMarkersService) MustGetSourceFile(line string, options *ReplaceBetweenMarkersOptions) (sourceFile File) {
	sourceFile, err := r.GetSourceFile(line, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourceFile
}

func (r *ReplaceBetweenMarkersService) MustGetSourcePath(line string, options *ReplaceBetweenMarkersOptions) (sourcePath string) {
	sourcePath, err := r.GetSourcePath(line, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sourcePath
}

func (r *ReplaceBetweenMarkersService) MustReplaceBySourcesInString(input string, options *ReplaceBetweenMarkersOptions) (replaced string) {
	replaced, err := r.ReplaceBySourcesInString(input, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return replaced
}

func (r *ReplaceBetweenMarkersService) ReplaceBySourcesInString(input string, options *ReplaceBetweenMarkersOptions) (replaced string, err error) {
	if options == nil {
		return "", tracederrors.TracedError("options is nil")
	}

	outLines := []string{}

	var startMarkerFound bool = false
	for _, line := range astrings.SplitLines(input, false) {
		if astrings.IsComment(line) {
			if r.IsReplaceBetweenMarkerStart(line) {
				startMarkerFound = true
				outLines = append(outLines, line)

				additionalLines, err := r.GetContentToInsertDefinedInStartLineAsLines(line, options)
				if err != nil {
					return "", err
				}

				outLines = append(outLines, additionalLines...)

				continue
			}

			if r.IsReplaceBetweenMarkerEnd(line) {
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
