package aiutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Concat all files beside a info.json to a single file.
// The generated output can be given to an AI as knowledge.
func ConcatFilesToKnowledgeFile(ctx context.Context, sourcePath string) (string, error) {
	if sourcePath == "" {
		return "", tracederrors.TracedErrorEmptyString("sourcePath")
	}

	logging.LogInfoByCtxf(ctx, "Concat files to knowledge file from '%s' started.", sourcePath)

	files, err := nativefiles.ListFiles(ctx, sourcePath, &parameteroptions.ListFileOptions{
		MatchBasenamePattern: []string{"^info.json$"},
	})
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Found %d info.json files.", len(files))

	builder := new(strings.Builder)
	var delimiter = "\n\n" + strings.Repeat("=", 20) + "\n\n"
	var fileinfoDelimiter = "\n" + strings.Repeat("-", 10) + "\n"

	for _, f := range files {
		infoContent, err := os.ReadFile(f)
		if err != nil {
			return "", tracederrors.TracedErrorf("Failed to read info file '%s': %w", f, err)
		}

		baseName, err := jsonutils.RunJqAgainstJsonStringAsString(string(infoContent), ".basename")
		if err != nil {
			return "", tracederrors.TracedErrorf("Unable to extract basename from '%s': %w", f, err)
		}

		contentPath := filepath.Join(filepath.Dir(f), baseName)

		content, err := os.ReadFile(contentPath)
		if err != nil {
			return "", tracederrors.TracedErrorf("Unable to read the content file '%s'.", contentPath)
		}

		fmt.Fprint(builder, delimiter)
		fmt.Fprint(builder, "Fileinfo: \n")
		fmt.Fprint(builder, fileinfoDelimiter)
		fmt.Fprint(builder, string(infoContent))
		fmt.Fprint(builder, fileinfoDelimiter)
		fmt.Fprint(builder, "Content: \n")
		fmt.Fprint(builder, fileinfoDelimiter)
		fmt.Fprint(builder, string(content))
		fmt.Fprint(builder, fileinfoDelimiter)
		fmt.Fprint(builder, delimiter)
	}

	logging.LogInfoByCtxf(ctx, "Concat files to knowledge file from '%s' finished.", sourcePath)

	return builder.String(), nil
}
