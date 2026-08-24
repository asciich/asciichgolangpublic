package commandexecutorgnupg

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/gnupgoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckSignatureValid(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, signatureFile filesinterfaces.File, toValidateFile filesinterfaces.File) (err error) {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if signatureFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	if toValidateFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	hostDescriptionCommandExecutor, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	signaturePath, hostDescription, err := signatureFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != hostDescriptionCommandExecutor {
		return tracederrors.TracedErrorf("Mismatching hostDescriptions: CommandExecutor is on hostdescription='%s' while the signature file is on '%s'.", hostDescriptionCommandExecutor, hostDescription)
	}

	toValidatePath, hostDescription, err := toValidateFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if hostDescription != hostDescriptionCommandExecutor {
		return tracederrors.TracedErrorf("Mismatching hostDescriptions: CommandExecutor is on hostdescription='%s' while the file to validate is on '%s'.", hostDescriptionCommandExecutor, hostDescription)
	}

	return CheckSingnatureByPathValid(ctx, commandExecutor, signaturePath, toValidatePath)
}

func CheckSingnatureByPathValid(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, signaturePath string, toValidatePath string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if signaturePath == "" {
		return tracederrors.TracedErrorEmptyString("signaturePath")
	}

	if toValidatePath == "" {
		return tracederrors.TracedErrorEmptyString("toValidatePath")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Validate GnuPG signature '%s' for '%s' on host '%s' started.", signaturePath, toValidatePath, hostDescription)

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"gpg", "--verify", signaturePath, toValidatePath},
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Validate GnuPG signature '%s' for '%s' on host '%s' finished.", signaturePath, toValidatePath, hostDescription)

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

// SignContentString signs the given content string and returns the signature as a string.
func SignContentString(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, content string, options *gnupgoptions.SignOption) (signature string, err error) {
	sigBytes, err := SignContentBytes(ctx, commandExecutor, []byte(content), options)
	if err != nil {
		return "", err
	}
	return string(sigBytes), nil
}

// SignContentBytes signs the given content bytes and returns the signature as bytes.
func SignContentBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, content []byte, options *gnupgoptions.SignOption) (signature []byte, err error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	if options == nil {
		options = &gnupgoptions.SignOption{}
	}

	logging.LogInfoByCtxf(ctx, "Sign content using gnupg started.")

	if !options.AsciiArmor {
		return nil, tracederrors.TracedError("Only implemented for asciiArmor at the moment")
	}

	// Create temporary file for content
	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	tempContentPath, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = commandexecutorfile.Delete(ctx, commandExecutor, tempContentPath, &filesoptions.DeleteOptions{})
	}()

	// Write content to temp file
	err = commandexecutorfile.WriteBytes(ctx, commandExecutor, tempContentPath, content, nil)
	if err != nil {
		return nil, err
	}

	// Create temporary file for signature
	tempSigPath, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = commandexecutorfile.Delete(ctx, commandExecutor, tempSigPath, &filesoptions.DeleteOptions{})
	}()

	// Sign the content file
	signCommand := []string{
		"gpg",
		"--yes",
		"--armor",
		"--detach-sig",
		"--output",
		tempSigPath,
		tempContentPath,
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
		},
	)
	if err != nil {
		return nil, err
	}

	// Read the signature
	signature, err = commandexecutorfile.ReadAsBytes(commandExecutor, tempSigPath)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Sign content using gnupg finished on host '%s'.", hostDescription)

	return signature, nil
}

// CheckSignatureValidForContentString verifies that the signature is valid for the given content string.
func CheckSignatureValidForContentString(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, signature string, content string) error {
	return CheckSignatureValidForContentBytes(ctx, commandExecutor, []byte(signature), []byte(content))
}

// CheckSignatureValidForContentBytes verifies that the signature is valid for the given content bytes.
func CheckSignatureValidForContentBytes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, signature []byte, content []byte) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if signature == nil {
		return tracederrors.TracedErrorNil("signature")
	}

	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	logging.LogInfoByCtxf(ctx, "Validate content signature using gnupg started.")

	// Create temporary file for content
	tempContentPath, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = commandexecutorfile.Delete(ctx, commandExecutor, tempContentPath, &filesoptions.DeleteOptions{})
	}()

	// Write content to temp file
	err = commandexecutorfile.WriteBytes(ctx, commandExecutor, tempContentPath, content, nil)
	if err != nil {
		return err
	}

	// Create temporary file for signature
	tempSigPath, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = commandexecutorfile.Delete(ctx, commandExecutor, tempSigPath, &filesoptions.DeleteOptions{})
	}()

	// Write signature to temp file
	err = commandexecutorfile.WriteBytes(ctx, commandExecutor, tempSigPath, signature, nil)
	if err != nil {
		return err
	}

	// Verify the signature
	return CheckSingnatureByPathValid(ctx, commandExecutor, tempSigPath, tempContentPath)
}
