package commandexecutorfileoo

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (f *File) WriteBytes(ctx context.Context, toWrite []byte, options *filesoptions.WriteOptions) (err error) {
	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	commandExecutor, filePath, err := f.GetCommandExecutorAndFilePath()
	if err != nil {
		return err
	}

	hostDescription, err := f.GetHostDescription()
	if err != nil {
		return err
	}

	command := []string{"bash", "-c", fmt.Sprintf("cat > '%s'", filePath)}

	if options != nil && options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	_, err = commandExecutor.RunCommand(
		contextutils.WithSilent(ctx),
		&parameteroptions.RunCommandOptions{
			Command:     command,
			StdinString: string(toWrite),
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Wrote '%d' bytes to file '%s' on '%s'", len(toWrite), filePath, hostDescription)

	return nil
}
