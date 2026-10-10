package commandexecutoriscsi

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/storage/iscsiutils/iscsigeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const initiatorNameFilePath = "/etc/iscsi/initiatorname.iscsi"

// EnsureInitiorNameConfig ensures that the /etc/iscsi/initiatorname.iscsi file exists.
// If it does not exist, it creates it with a proper, new generated InitiatorName.
func EnsureInitiorNameConfig(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *filesoptions.CreateOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		options = &filesoptions.CreateOptions{}
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Ensure ISCSI initiator name in '%s' on '%s' started.", initiatorNameFilePath, hostDescription)

	writeNeeded, err := commandexecutorfile.IsEmptyFile(ctx, commandExecutor, initiatorNameFilePath)
	if err != nil {
		if filesgeneric.IsErrFileNotFound(err) {
			writeNeeded = true
		} else {
			return err
		}
	}

	if writeNeeded {
		name, err := iscsigeneric.GenerateInitiatorName(ctx)
		if err != nil {
			return err
		}

		content := "InitiatorName=" + name + "\n"

		err = commandexecutorfile.WriteString(ctx, commandExecutor, initiatorNameFilePath, content, &filesoptions.WriteOptions{
			UseSudo: options.UseSudo,
		})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Wrote new ISCSI initiator name '%s' to '%s' on '%s'.", name, initiatorNameFilePath, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "ISCSI initiator name already present in '%s' on '%s'. Skip recreation.", initiatorNameFilePath, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Ensure ISCSI initiator name in '%s' on '%s' finished.", initiatorNameFilePath, hostDescription)

	return nil
}
