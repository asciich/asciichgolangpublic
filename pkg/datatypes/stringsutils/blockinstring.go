package stringsutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func BlockInString(ctx context.Context, content string, blockName string, block string) (string, error) {
	if blockName == "" {
		return "", tracederrors.TracedErrorEmptyString("blockName")
	}

	beginMarker := fmt.Sprintf("# BEGIN %s", blockName)
	endMarker := fmt.Sprintf("# END %s", blockName)

	markedBlock := beginMarker + "\n" + block + "\n" + endMarker

	beginIdx := strings.Index(content, beginMarker)
	endIdx := strings.Index(content, endMarker)

	var modifiedContent string

	if beginIdx == -1 && endIdx == -1 {
		if content != "" && !strings.HasSuffix(content, "\n") {
			modifiedContent = content + "\n" + markedBlock + "\n"
		} else {
			modifiedContent = content + markedBlock + "\n"
		}
		contextutils.SetChangeIndicator(ctx, true)
	} else if beginIdx == -1 || endIdx == -1 || beginIdx > endIdx {
		return "", tracederrors.TracedErrorf(
			"Malformed block markers for '%s': begin or end marker missing or in wrong order",
			blockName,
		)
	} else {
		endMarkerEnd := endIdx + len(endMarker)
		existingBlock := content[beginIdx:endMarkerEnd]

		if existingBlock == markedBlock {
			modifiedContent = content
		} else {
			modifiedContent = content[:beginIdx] + markedBlock + content[endMarkerEnd:]
			contextutils.SetChangeIndicator(ctx, true)
		}
	}

	return modifiedContent, nil
}
