package commandexecutorfile

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ReadFirstNBytes(commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, numberOfBytesToRead int) (firstBytes []byte, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExectuor")
	}

	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	if numberOfBytesToRead < 0 {
		return nil, tracederrors.TracedErrorf("Invalid number of bytes to read: %d", numberOfBytesToRead)
	}

	firstBytes, err = commandExecutor.RunCommandAndGetStdoutAsBytes(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"head",
				fmt.Sprintf(
					"--bytes=%d",
					numberOfBytesToRead,
				),
				filePath,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return firstBytes, nil
}

func ReadAsBytes(commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string) ([]byte, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExectuor")
	}

	if filePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	content, err := commandExecutor.RunCommandAndGetStdoutAsBytes(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"cat", filePath},
		},
	)
	if err != nil {
		return nil, err
	}

	return content, nil
}
