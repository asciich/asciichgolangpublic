package nativefiles

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SetBlockInFile(ctx context.Context, path string, blockName string, block string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if blockName == "" {
		return tracederrors.TracedErrorEmptyString("blockName")
	}

	logging.LogInfoByCtxf(ctx, "Set block '%s' in '%s' started.", path, blockName)

	content, err := ReadAsString(ctx, path, &filesoptions.ReadOptions{})
	if err != nil {
		return err
	}

	adjusted, err := stringsutils.BlockInString(ctx, content, blockName, block)
	if err != nil {
		return err
	}

	if adjusted != content {
		err = WriteString(ctx, path, adjusted)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Set block '%s' in '%s' finished.", path, blockName)

	return nil
}
