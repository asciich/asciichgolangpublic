package gnupg

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckSignatureValid(ctx context.Context, signatureFile filesinterfaces.File) (err error) {
	if signatureFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	err = signatureFile.CheckIsLocalFile(ctx)
	if err != nil {
		return tracederrors.TracedErrorf("Only implemented for local available files: %w", err)
	}

	path, hostDescription, err := signatureFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Validate GnuPG signature from '%s' on host '%s' started.", path, hostDescription)

	_, err = commandexecutorbashoo.Bash().RunCommand(ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"gpg", "--verify", path},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "GnuPG signature from '%s' on host '%s' validated.", path, hostDescription)

	return nil
}

func SignFile(ctx context.Context, fileToSign filesinterfaces.File, options *GnuPGSignOptions) (err error) {
	if fileToSign == nil {
		return tracederrors.TracedError("fileToSign is nil")
	}

	if options == nil {
		return tracederrors.TracedError("options is nil")
	}

	err = fileToSign.CheckIsLocalFile(ctx)
	if err != nil {
		return tracederrors.TracedErrorf("Only implemented for local available files: %w", err)
	}

	path, err := fileToSign.GetPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Sign '%s' using gnupg started.", path)

	if !options.AsciiArmor {
		return tracederrors.TracedError("Only implemented for asciiArmor at the moment")
	}

	if !options.DetachedSign {
		return tracederrors.TracedError("Only implemented for DetachedSign at the moment")
	}

	signaturePath := path + ".asc"
	signatureFile, err := files.GetLocalFileByPath(signaturePath)
	if err != nil {
		return err
	}

	if err = signatureFile.Delete(ctx, &filesoptions.DeleteOptions{}); err != nil {
		return err
	}

	signCommand := []string{
		"gpg",
		"--armor",
		"--detach-sig",
		path,
	}

	_, err = commandexecutorbashoo.Bash().RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
		},
	)
	if err != nil {
		return err
	}

	signatureFileExists, err := signatureFile.Exists(contextutils.ContextSilent())
	if err != nil {
		return err
	}

	if !signatureFileExists {
		return tracederrors.TracedErrorf(
			"Signing '%s' failed. Expected signature file '%s' does not exits.",
			path,
			signaturePath,
		)
	}

	logging.LogInfoByCtxf(ctx, "Sign '%s' using gnupg finished.", path)

	return nil
}
