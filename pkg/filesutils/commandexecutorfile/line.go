package commandexecutorfile

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func EndsWithLineBreak(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	// Read only the last byte of the file.
	// "--" makes sure a file path starting with "-" isn't read as an option.
	lastByte, err := commandExecutor.RunCommandAndGetStdoutAsBytes(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"tail", "-c", "1", "--", filePath,
			},
		},
	)
	if err != nil {
		return false, err
	}

	return len(lastByte) == 1 && lastByte[0] == '\n', nil
}

func AppendLine(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, line string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if line == "" {
		return tracederrors.TracedErrorEmptyString("line")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	err = CreateFile(ctx, commandExecutor, filePath, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	toWrite := stringsutils.TrimAllLeadingAndTailingNewLines(line)
	toWrite = stringsutils.EnsureEndsWithExactlyOneLineBreak(toWrite)

	emptyFile, err := IsEmptyFile(ctx, commandExecutor, filePath)
	if err != nil {
		return err
	}

	if !emptyFile {
		endsWithLineBreak, err := EndsWithLineBreak(ctx, commandExecutor, filePath)
		if err != nil {
			return err
		}

		if !endsWithLineBreak {
			toWrite = "\n" + toWrite
		}
	}

	err = AppendString(ctx, commandExecutor, filePath, toWrite)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Appended line to '%s' on '%s'.", filePath, hostDescription)

	return nil
}

func ReadAsLines(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string) ([]string, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("filePath")
	}

	content, err := ReadAsString(commandExecutor, filePath)
	if err != nil {
		return nil, err
	}

	return stringsutils.SplitLines(content, false), nil
}

func EnsureLineInFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, line string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if line == "" {
		return tracederrors.TracedErrorEmptyString("line")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	err = CreateFile(ctx, commandExecutor, filePath, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	lines, err := ReadAsLines(ctx, commandExecutor, filePath)
	if err != nil {
		return err
	}

	if slices.Contains(lines, line) {
		logging.LogInfof("Line '%s' already present in '%s' on '%s'.", line, filePath, hostDescription)
	} else {
		err = AppendLine(ctx, commandExecutor, filePath, line)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Wrote line '%s' into '%s' on '%s'.", line, filePath, hostDescription)
	}

	return nil
}
