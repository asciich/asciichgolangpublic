package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GnuPGService struct {
}

func GnuPG() (gnuPG *GnuPGService) {
	return new(GnuPGService)
}

func NewGnuPGService() (g *GnuPGService) {
	return new(GnuPGService)
}

func (g *GnuPGService) CheckSignatureValid(signatureFile files.File, verbose bool) (err error) {
	if signatureFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	err = signatureFile.CheckIsLocalFile(verbose)
	if err != nil {
		return tracederrors.TracedErrorf("Only implemented for local available files: %w", err)
	}

	path, hostDescription, err := signatureFile.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Validate GnuPG signature from '%s' on host '%s' started.",
			path,
			hostDescription,
		)
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: []string{"gpg", "--verify", path},
			Verbose: verbose,
		},
	)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"GnuPG signature from '%s' on host '%s' validated.",
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GnuPGService) MustCheckSignatureValid(signatureFile files.File, verbose bool) {
	err := g.CheckSignatureValid(signatureFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GnuPGService) MustSignFile(fileToSign files.File, options *GnuPGSignOptions) {
	err := g.SignFile(fileToSign, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GnuPGService) SignFile(fileToSign files.File, options *GnuPGSignOptions) (err error) {
	if fileToSign == nil {
		return tracederrors.TracedError("fileToSign is nil")
	}

	if options == nil {
		return tracederrors.TracedError("options is nil")
	}

	err = fileToSign.CheckIsLocalFile(options.Verbose)
	if err != nil {
		return tracederrors.TracedErrorf("Only implemented for local available files: %w", err)
	}

	path, err := fileToSign.GetPath()
	if err != nil {
		return err
	}

	if options.Verbose {
		logging.LogInfof("Sign '%s' using gnupg started.", path)
	}

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

	if err = signatureFile.Delete(options.Verbose); err != nil {
		return err
	}

	signCommand := []string{
		"gpg",
		"--armor",
		"--detach-sig",
		path,
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: signCommand,
			Verbose: options.Verbose,
		},
	)
	if err != nil {
		return err
	}

	signatureFileExists, err := signatureFile.Exists(false)
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

	if options.Verbose {
		logging.LogInfof("Sign '%s' using gnupg finished.", path)
	}

	return nil
}
