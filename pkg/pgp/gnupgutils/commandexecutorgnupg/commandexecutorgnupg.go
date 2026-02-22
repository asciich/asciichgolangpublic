package commandexecutorgnupg

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/gnupgoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckSignatureValid(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, signatureFile filesinterfaces.File) (err error) {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if signatureFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	path, hostDescription, err := signatureFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	hostDescriptionCommandExecutor, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != hostDescriptionCommandExecutor {
		return tracederrors.TracedErrorf("Mismatching hostDescriptions: CommandExecutor is on hostdescription='%s' while the file to check is on '%s'.", hostDescriptionCommandExecutor, hostDescription)
	}

	return CheckSingnatureByPathValid(ctx, commandExecutor, path)
}

func CheckSingnatureByPathValid(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, signaturePath string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if signaturePath == "" {
		return tracederrors.TracedErrorEmptyString("signaturePath")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Validate GnuPG signature from '%s' on host '%s' started.", signaturePath, hostDescription)

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"gpg", "--verify", signaturePath},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "GnuPG signature from '%s' on host '%s' validated.", signaturePath, hostDescription)

	return nil
}

func SignFileByPath(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, toSignPath string, options *gnupgoptions.SignOption) error {
	if toSignPath == "" {
		return tracederrors.TracedErrorEmptyString("toSignPath")
	}

	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		options = &gnupgoptions.SignOption{}
	}

	logging.LogInfoByCtxf(ctx, "Sign '%s' using gnupg started.", toSignPath)

	if !options.AsciiArmor {
		return tracederrors.TracedError("Only implemented for asciiArmor at the moment")
	}

	if !options.DetachedSign {
		return tracederrors.TracedError("Only implemented for DetachedSign at the moment")
	}

	signaturePath := toSignPath + ".asc"
	err := commandexecutorfile.Delete(ctx, commandExecutor, signaturePath, &filesoptions.DeleteOptions{})
	if err != nil {
		return err
	}

	signCommand := []string{
		"gpg",
		"--armor",
		"--detach-sig",
		"--output",
		signaturePath,
		toSignPath,
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
		},
	)
	if err != nil {
		return err
	}

	signatureFileExists, err := commandexecutorfile.Exists(ctx, commandExecutor, signaturePath)
	if err != nil {
		return err
	}

	if !signatureFileExists {
		return tracederrors.TracedErrorf(
			"Signing '%s' failed. Expected signature file '%s' does not exits.",
			toSignPath,
			signaturePath,
		)
	}

	logging.LogInfoByCtxf(ctx, "Sign '%s' using gnupg finished.", toSignPath)

	return nil
}

func SignFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, fileToSign filesinterfaces.File, options *gnupgoptions.SignOption) (err error) {
	if fileToSign == nil {
		return tracederrors.TracedError("fileToSign is nil")
	}

	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	toSignPath, err := fileToSign.GetPath()
	if err != nil {
		return err
	}

	return SignFileByPath(ctx, commandExecutor, toSignPath, options)
}
