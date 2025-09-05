package nativefiles

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Contains(ctx context.Context, filePath string, searchString string) (bool, error) {
	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	content, err := ReadAsString(contextutils.WithSilent(ctx), filePath)
	if err != nil {
		return false, err
	}

	contains := strings.Contains(content, searchString)

	if contains {
		logging.LogInfoByCtxf(ctx, "File '%s' does contain search string '%s'.", filePath, searchString)
	} else {
		logging.LogInfoByCtxf(ctx, "File '%s' does not contain search string '%s'.", filePath, searchString)
	}

	return contains, nil
}
