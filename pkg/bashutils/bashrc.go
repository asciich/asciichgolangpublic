package bashutils

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

func GetBashRcFileOfCurrentUser(ctx context.Context) (bashRcFile filesinterfaces.File, err error) {
	rcfile, err := userutils.GetFileInHomeDirectory(".bashrc")
	if err != nil {
		return nil, err
	}

	path, err := rcfile.GetPath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Bashrc file of current user is '%s'.", path)

	return rcfile, nil
}

func EnableImmediateHistoryReadAndWriteForCurrentUser(ctx context.Context) (err error) {
	logging.LogInfoByCtxf(ctx, "Enable immediate history read and write for current user started.")

	bashRcFile, err := GetBashRcFileOfCurrentUser(ctx)
	if err != nil {
		return err
	}

	bashRcPath, err := bashRcFile.GetLocalPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Use bashrc file '%s' to enable immediate history read and write.", bashRcPath)

	linesToAdd := []string{
		"export PROMPT_COMMAND=\"history -a;history -r\"",
		"export HISTCONTROL=ignorespace:erasedups",
	}
	for _, line := range linesToAdd {
		err = bashRcFile.EnsureLineInFile(line, contextutils.GetVerboseFromContext(ctx))
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Enable immediate history read and write for current user finished.")

	return nil
}

func SetBashHistorySizeOfCurrentUser(ctx context.Context, newBashHistorySize int) (err error) {
	logging.LogInfoByCtxf(ctx, "Set size of bash history for current user to '%d' entries started.", newBashHistorySize)

	if newBashHistorySize <= 0 {
		return tracederrors.TracedErrorf("Invalid bash history size: '%d'", newBashHistorySize)
	}

	bashRcfile, err := GetBashRcFileOfCurrentUser(ctx)
	if err != nil {
		return err
	}

	linesToAdd := []string{
		fmt.Sprintf("export HISTSIZE=%d", newBashHistorySize),
		fmt.Sprintf("export HISTFILESIZE=%d", newBashHistorySize),
	}

	for _, l := range linesToAdd {
		err = bashRcfile.EnsureLineInFile(l, contextutils.GetVerboseFromContext(ctx))
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Set size of bash history for current user to '%d' entries started.", newBashHistorySize)

	return nil
}